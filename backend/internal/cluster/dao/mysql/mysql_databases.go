package mysql

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/ast"
	_ "github.com/pingcap/tidb/pkg/types/parser_driver"

	"github.com/keepsty/go_rds/internal/cluster/models"
)

type visitor struct {
	tableNames []string
}

func (v *visitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	if name, ok := in.(*ast.TableName); ok {
		v.tableNames = append(v.tableNames, name.Name.O)
	}
	return in, false
}

func (v *visitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	return in, true
}

func getTableName(slices []string) string {
	var res string
	if len(slices) < 1 {
		return ""
	}
	for _, v := range slices {
		if res == "" {
			res = fmt.Sprintf("%s", v)
		} else {
			res = fmt.Sprintf("%s,%s", res, v)
		}
	}
	return res
}

func quoteIdentifier(name string) string {
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}

func GetDNSByDBID(id int64) (dsn *models.MySQLDsn, err error) {
	sqlGetDBDsnStr := `select IFNULL(h.ip,'') as ip,IFNULL(di.port,0) as port,IFNULL(db.name,'') as dbname from database_database as db left join database_servicegroup as sg on sg.id=db.servicegroup_id left join database_instance as di on di.servicegroup_id=sg.id left join hosts as h on h.id=di.host_id where di.role=0 and db.id=? order by di.id limit 1`
	dsn = new(models.MySQLDsn)
	dsn.User = "root"
	dsn.Password = "123123"
	if err = db.Get(dsn, sqlGetDBDsnStr, id); err != nil {
		return nil, err
	}
	return
}

func GetClusterDBByUsername(username string) (data []*models.SgDBOptions, err error) {
	sqlDBTotalStr := `select count(db.name) from database_database as db,database_servicegroup as sg where db.servicegroup_id=sg.id and db.rd_user like CONCAT('%,', ?, ',%')`
	sqlSGTotalStr := `select count(db.name) from database_database as db,database_servicegroup as sg where db.servicegroup_id=sg.id and db.rd_user like CONCAT('%,', ?, ',%') group by sg.id`
	var dbCnt int64
	var sgCnt int64
	if err = db.Get(&dbCnt, sqlDBTotalStr, username); err != nil {
		return nil, err
	}
	if err := db.Get(&sgCnt, sqlSGTotalStr, username); err != nil {
		return nil, err
	}

	data = make([]*models.SgDBOptions, 0, sgCnt)
	dbdata := make([]*models.SgDBOptionsChild, 0, dbCnt)
	sqlSGStr := `select sg.id sg_id,IFNULL(sg.name,'') sg_name from database_database as db,database_servicegroup as sg where db.servicegroup_id=sg.id and db.rd_user like CONCAT('%,', ?, ',%') group by sg.id`
	sqlDBStr := `select db.id db_id,db.name db_name,sg.id as sg_id from database_database as db,database_servicegroup as sg where db.servicegroup_id=sg.id and db.rd_user like CONCAT('%,', ?, ',%')`
	if err = db.Select(&data, sqlSGStr, username); err != nil {
		return nil, err
	}
	if err = db.Select(&dbdata, sqlDBStr, username); err != nil {
		return nil, err
	}
	for _, rows := range data {
		for _, dbRows := range dbdata {
			if dbRows.SGID == rows.Value {
				rows.Children = append(rows.Children, dbRows)
			}
		}
	}
	return
}

func GetUsernameHistoryQuery(rqUser *models.RequestGetUserHistoryQueryHandler) (data *models.RpsQueryHistoryData, err error) {
	sqlUsernameTotalStr := "select count(1) from database_dms_query_record where username=?"
	sqlDBStr := "select host,db_name,tables,execute_query,query_consume_time,affected_rows,create_time from database_dms_query_record where username=?"
	if rqUser.Table != "" {
		sqlUsernameTotalStr = sqlUsernameTotalStr + " and (db_name=? or tables=?)"
		sqlDBStr = sqlDBStr + " and (db_name=? or tables=?)"
	}
	var userQueryCount int64
	if err := db.Get(&userQueryCount, sqlUsernameTotalStr, rqUser.Username, rqUser.Table, rqUser.Table); err != nil {
		return nil, err
	}
	if userQueryCount < 1 {
		err = ErrorNoRows
		return nil, err
	}
	rows, err := db.Query(sqlDBStr, rqUser.Username, rqUser.Table, rqUser.Table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data = new(models.RpsQueryHistoryData)
	col, _ := rows.Columns()
	data.Columns = make([]*models.RpsQueryHistoryColumn, 0, len(col))

	for _, v := range col {
		cf := new(models.RpsQueryHistoryColumn)
		cf.Title = v
		cf.Key = v
		cf.DataIndex = v
		cf.Escape = true
		data.Columns = append(data.Columns, cf)
	}
	values := make([][]byte, len(col))
	scans := make([]interface{}, len(col))
	for i := range values {
		scans[i] = &values[i]
	}
	res := make([]map[string]string, 0)
	for rows.Next() {
		_ = rows.Scan(scans...)
		row := make(map[string]string)
		for k, v := range values {
			key := col[k]
			row[key] = string(v)
		}
		res = append(res, row)
	}
	data.Data = res
	return
}

func ClusterGetDBTableInfoHandler(tbInfo *models.RequestGetDBTableInfo) (data *models.TableStructureData, err error) {
	mysqlDsn, err := GetDNSByDBID(tbInfo.DBID)
	if err != nil {
		return
	}
	itemDsn, err := mysqlDsn.InitMySQL()
	if err != nil {
		return
	}
	defer mysqlDsn.CloseMySQL(itemDsn)

	data = new(models.TableStructureData)
	if tbInfo.Type == "table_structure" {
		sqlDBStr := fmt.Sprintf("show create table %s.%s ;", quoteIdentifier(mysqlDsn.DBName), quoteIdentifier(tbInfo.Table))
		if err = itemDsn.Select(&data.TS, sqlDBStr); err != nil {
			return nil, err
		}
	} else if tbInfo.Type == "table_base" {
		sqlDBStr := `select TABLE_NAME as '表名', TABLE_TYPE as '表类型', ENGINE as '引擎', ROW_FORMAT as '行格式', TABLE_ROWS as '表行数(估算值)', round(DATA_LENGTH/1024, 2) as '数据大小(KB)', round(INDEX_LENGTH/1024, 2) as '索引大小(KB)', TABLE_COLLATION as '字符集校验规则', TABLE_COMMENT as '备注', CREATE_TIME as '创建时间' from information_schema.tables where table_schema=? and table_name=?`
		if itemDsn.Select(&data.TB, sqlDBStr, mysqlDsn.DBName, tbInfo.Table) != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("查询类型异常")
	}
	return data, nil
}

func GetClusterTableBySGDBID(sgID, dbID int64) (data []*models.GetTableDBData, err error) {
	mysqlDsn, err := GetDNSByDBID(dbID)
	if err != nil {
		return
	}
	itemDsn, err := mysqlDsn.InitMySQL()
	if err != nil {
		return
	}
	defer mysqlDsn.CloseMySQL(itemDsn)

	var tbCount int64
	sqlGetTableCountStr := `select count(1) from information_schema.columns where table_schema=? and table_name not regexp '^_(.*)[_ghc|_gho|_del]$' group by table_schema, table_name order by table_name`
	if err = itemDsn.Get(&tbCount, sqlGetTableCountStr, mysqlDsn.DBName); err != nil {
		return nil, err
	}

	data = make([]*models.GetTableDBData, 0, tbCount)
	sqlStr := `select table_schema as 'database', table_name as 'table', group_concat(concat(column_name, ' ', column_type) SEPARATOR '#') as 'join_columns', group_concat(column_name) as 'columns' from information_schema.columns where table_schema=? and table_name not regexp '^_(.*)[_ghc|_gho|_del]$' group by table_schema, table_name order by table_name`
	if err = itemDsn.Select(&data, sqlStr, mysqlDsn.DBName); err != nil {
		return nil, err
	}
	return
}

func ClusterDBReadQueryExecuteHandler(rqData *models.DBReadQueryExecute) (data *models.ReadQueryData, err error) {
	mysqlDsn, err := GetDNSByDBID(rqData.DBId)
	if err != nil {
		return
	}
	sqlStart := time.Now()
	itemDsn, err := mysqlDsn.InitMySQL()
	elapsedTime := time.Since(sqlStart)
	if err != nil {
		return nil, err
	}
	defer mysqlDsn.CloseMySQL(itemDsn)

	p := parser.New()
	stmt, err := p.ParseOneStmt(rqData.Sql, "", "")
	if err != nil {
		return nil, err
	}
	v := new(visitor)
	stmt.Accept(v)
	tableStr := getTableName(v.tableNames)

	execErrStr := ""
	rows, err := itemDsn.Query(rqData.Sql)
	if err != nil {
		execErrStr = fmt.Sprintf("Fail,%s", err.Error())
		data = new(models.ReadQueryData)
		data.Columns = make([]*models.RpsColumn, 0)
	} else {
		execErrStr = "Success"
		defer rows.Close()
		data = new(models.ReadQueryData)
		col, _ := rows.Columns()
		data.Columns = make([]*models.RpsColumn, 0, len(col))
		cf := new(models.RpsColumn)
		cf.Field = "state"
		cf.Checkbox = true
		data.Columns = append(data.Columns, cf)

		for _, v := range col {
			c := new(models.RpsColumn)
			c.Field = v
			c.Title = v
			c.Escape = true
			data.Columns = append(data.Columns, c)
		}
		values := make([][]byte, len(col))
		scans := make([]interface{}, len(col))
		for i := range values {
			scans[i] = &values[i]
		}
		res := make([]map[string]string, 0)
		for rows.Next() {
			_ = rows.Scan(scans...)
			row := make(map[string]string)
			for k, v := range values {
				key := col[k]
				row[key] = string(v)
			}
			res = append(res, row)
		}
		data.Data = res
	}

	queryRecord := `insert into database_dms_query_record(username,host,db_name,tables,execute_query,query_consume_time,query_status,affected_rows) values (?,?,?,?,?,?,?,?)`
	if _, err = db.Exec(queryRecord, rqData.Username, mysqlDsn.IP, mysqlDsn.DBName, tableStr, strings.Replace(rqData.Sql, "\n", " ", -1), elapsedTime.Seconds(), execErrStr, len(data.Data)); err != nil {
		return nil, err
	}
	return
}

func ClusterDBDataDictHandler(dbData *models.DBReadQueryExecute) (data []*models.DataDictJson, dbName string, err error) {
	mysqlDsn, err := GetDNSByDBID(dbData.DBId)
	if err != nil {
		return
	}
	itemDsn, err := mysqlDsn.InitMySQL()
	if err != nil {
		return
	}
	defer mysqlDsn.CloseMySQL(itemDsn)

	var tableCnt int64
	tableCntSql := "select count(1) from information_schema.TABLES where TABLE_SCHEMA=?"
	if err = itemDsn.Get(&tableCnt, tableCntSql, mysqlDsn.DBName); err != nil {
		return nil, "", err
	}
	data = make([]*models.DataDictJson, 0, tableCnt)
	sql := `select t.table_name,if(t.TABLE_COMMENT!='',t.TABLE_COMMENT,'None') as table_comment,t.create_time, group_concat(distinct concat_ws('<b>', c.COLUMN_NAME,c.COLUMN_TYPE,if(c.IS_NULLABLE='NO','NOT NULL','NULL'),ifnull(c.COLUMN_DEFAULT, ''),ifnull(c.CHARACTER_SET_NAME,''), ifnull(c.COLLATION_NAME,''),ifnull(c.COLUMN_COMMENT, '')) separator '<a>') as columns_info, group_concat(distinct concat_ws('<b>', s.INDEX_NAME,if(s.NON_UNIQUE=0,'唯一','不唯一'),s.Cardinality, s.INDEX_TYPE,s.COLUMN_NAME) separator '<a>') as index_info from information_schema.COLUMNS c join information_schema.TABLES t on c.TABLE_SCHEMA = t.TABLE_SCHEMA and c.TABLE_NAME = t.TABLE_NAME left join information_schema.STATISTICS s on c.TABLE_SCHEMA = s.TABLE_SCHEMA and c.TABLE_NAME = s.TABLE_NAME where t.TABLE_SCHEMA=? group by t.TABLE_NAME,t.TABLE_COMMENT,t.CREATE_TIME`
	if err = itemDsn.Select(&data, sql, mysqlDsn.DBName); err != nil {
		return nil, "", err
	}
	return data, mysqlDsn.DBName, nil
}
