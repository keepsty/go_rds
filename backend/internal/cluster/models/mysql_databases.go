package models

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type SgDBOptions struct {
	Value    int64               `json:"value" db:"sg_id"`
	Label    string              `json:"label" db:"sg_name"`
	Children []*SgDBOptionsChild `json:"children"`
}

type SgDBOptionsChild struct {
	Value int64  `json:"value" db:"db_id"`
	Label string `json:"label" db:"db_name"`
	SGID  int64  `json:"sg_id" db:"sg_id"`
}

type MySQLDsn struct {
	User     string `json:"user" db:"user"`
	Password string `json:"password" db:"password"`
	IP       string `json:"ip" db:"ip"`
	DBName   string `json:"dbname" db:"dbname"`
	Port     int64  `json:"port" db:"port"`
}

func (m *MySQLDsn) InitMySQL() (itemDB *sqlx.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.User, m.Password, m.IP, m.Port, m.DBName)
	itemDB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return
}

func (m *MySQLDsn) CloseMySQL(db *sqlx.DB) {
	_ = db.Close()
}

type SGDBTableList struct {
	Label    string                `json:"label"`
	Value    string                `json:"value"`
	Children []*SGDBTableListChild `json:"children"`
	Columns  []string              `json:"tables"`
}

type GetTableDBData struct {
	Database   string `json:"database" db:"database"`
	Table      string `json:"table" db:"table"`
	JoinColumn string `json:"join_columns" db:"join_columns"`
	Columns    string `json:"columns" db:"columns"`
}

type SGDBTableListChild struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type DBReadQueryExecute struct {
	DBId      int64  `json:"db_id"`
	Sql       string `json:"sql"`
	Character string `json:"character"`
	QueryHash string `json:"query_hash"`
	Username  string `json:"username"`
}

type ReadQueryData struct {
	Columns []*RpsColumn        `json:"columns"`
	Data    []map[string]string `json:"data"`
}

type RpsColumn struct {
	Field    string `json:"field"`
	Title    string `json:"title"`
	Escape   bool   `json:"escape"`
	Checkbox bool   `json:"checkbox"`
}

type RpsQueryHistoryData struct {
	Columns []*RpsQueryHistoryColumn `json:"columns"`
	Data    []map[string]string      `json:"data"`
}

type RpsQueryHistoryColumn struct {
	Title     string `json:"title"`
	DataIndex string `json:"dataIndex"`
	Key       string `json:"key"`
	Escape    bool   `json:"escape"`
}

type DataDictResponse struct {
	Data   []*DataDictJson `json:"data"`
	DBName string          `json:"db_name"`
}

type DataDictJson struct {
	TableName    string `json:"table_name" db:"TABLE_NAME"`
	TableComment string `json:"table_comment" db:"table_comment"`
	CreateTime   string `json:"create_time" db:"CREATE_TIME"`
	ColumnsInfo  string `json:"columns_info" db:"columns_info"`
	IndexInfo    string `json:"index_info" db:"index_info"`
}

type DMSExecQueryRecordHandler struct {
	ID               int64     `json:"id" db:"id"`
	AffectedRows     int64     `json:"affected_rows" db:"affected_rows"`
	Username         string    `json:"username" db:"username"`
	Host             string    `json:"host" db:"host"`
	DBName           string    `json:"db_name" db:"db_name"`
	Tables           string    `json:"tables" db:"tables"`
	ExecuteQuery     string    `json:"execute_query" db:"execute_query"`
	QueryConsumeTime string    `json:"query_consume_time" db:"query_consume_time"`
	QueryStatus      string    `json:"query_status" db:"query_status"`
	CreateTime       time.Time `json:"create_time" db:"create_time"`
}

type RequestGetUserHistoryQueryHandler struct {
	PageSize int64  `form:"page_size"`
	Page     int64  `form:"page"`
	DBID     int64  `form:"db_id"`
	Username string `form:"username"`
	Table    string `form:"table"`
}

type RequestGetDBTableInfo struct {
	DBID  int64  `json:"db_id" form:"db_id"`
	Table string `json:"table" form:"table"`
	Type  string `json:"type" form:"type"`
}

type TableStructureData struct {
	TS []*TableStructure `json:"ts"`
	TB []*TableBase      `json:"tb"`
}

type TableStructure struct {
	Table       string `json:"table" db:"Table"`
	CreateTable string `json:"create_table" db:"Create Table"`
}

type TableBase struct {
	Table          string `json:"table" db:"表名"`
	TableType      string `json:"create_table" db:"表类型"`
	Engine         string `json:"engine" db:"引擎"`
	RowFormat      string `json:"row_format" db:"行格式"`
	TableRows      string `json:"table_rows" db:"表行数(估算值)"`
	DataLength     string `json:"data_length" db:"数据大小(KB)"`
	IndexLength    string `json:"index_length" db:"索引大小(KB)"`
	TableCollation string `json:"table_collation" db:"字符集校验规则"`
	TableComment   string `json:"table_comment" db:"备注"`
	CreateTime     string `json:"create_time" db:"创建时间"`
}
