package services

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pingcap/tidb/pkg/parser"
	"github.com/pingcap/tidb/pkg/parser/ast"

	"github.com/keepsty/go_rds/internal/cluster/dao/mysql"
	"github.com/keepsty/go_rds/internal/cluster/models"
	"github.com/keepsty/go_rds/internal/config"
)

func GetClusterDBByUsername(username string) (data []*models.SgDBOptions, err error) {
	data, err = mysql.GetClusterDBByUsername(username)
	return
}

func GetUsernameHistoryQuery(rqUser *models.RequestGetUserHistoryQueryHandler) (data *models.RpsQueryHistoryData, err error) {
	data, err = mysql.GetUsernameHistoryQuery(rqUser)
	return
}

func ClusterGetDBTableInfoHandler(tbInfo *models.RequestGetDBTableInfo) (data *models.TableStructureData, err error) {
	data, err = mysql.ClusterGetDBTableInfoHandler(tbInfo)
	return
}

func GetClusterTableBySGDBID(sgID, dbID int64) (sgTableTree []*models.SGDBTableList, err error) {
	data, err := mysql.GetClusterTableBySGDBID(sgID, dbID)
	if err != nil {
		return nil, err
	}
	sgTableTree = make([]*models.SGDBTableList, 0, len(data))

	for _, val := range data {
		tmpTab := new(models.SGDBTableList)
		tmpTab.Columns = strings.Split(val.JoinColumn, "#")
		sgTableTreeChild := make([]*models.SGDBTableListChild, 0, len(tmpTab.Columns))
		for _, v := range tmpTab.Columns {
			tmpTabChild := new(models.SGDBTableListChild)
			tmpTabChild.Label = v
			tmpTabChild.Value = fmt.Sprintf("%s___%s___%s", val.Database, val.Table, v)
			sgTableTreeChild = append(sgTableTreeChild, tmpTabChild)
		}
		tmpTab.Label = fmt.Sprintf("%s___%s", val.Database, val.Table)
		tmpTab.Value = val.Table
		tmpTab.Children = sgTableTreeChild
		sgTableTree = append(sgTableTree, tmpTab)
	}
	return
}

func ReadQueryOptimize(node ast.StmtNode, sql string, st *config.Das) (ruleSql string, err error) {
	upperSQL := strings.ToUpper(sql)
	if len(strings.Split(upperSQL, "LIMIT")) > 2 {
		return "", errors.New("SQL语法异常，出现多个limit")
	}

	ruleLimitN, _ := regexp.Compile(`LIMIT([\s]*\d+[\s]*)$`)
	ruleLimitPoint, _ := regexp.Compile(`LIMIT([\s]*\d+[\s]*)(,)([\s]*\d+[\s]*)$`)
	ruleLimitOffset, _ := regexp.Compile(`LIMIT([\s]*\d+[\s]*)(OFFSET)([\s]*\d+[\s]*)$`)

	byteSQL := []byte(upperSQL)
	limitStr := ""
	limitNumSlice := make([]int64, 0, 2)

	switch {
	case ruleLimitN.Match(byteSQL):
		limitStr = string(ruleLimitN.Find(byteSQL))
	case ruleLimitPoint.Match(byteSQL):
		limitStr = string(ruleLimitPoint.Find(byteSQL))
	case ruleLimitOffset.Match(byteSQL):
		limitStr = string(ruleLimitOffset.Find(byteSQL))
		strings.Replace(limitStr, "OFFSET", ",", -1)
	default:
		tmpSQL := strings.Replace(sql, ";", "", -1)
		ruleSql = fmt.Sprintf("%s LIMIT %d", tmpSQL, st.DefaultReturnRows)
		return ruleSql, nil
	}
	splitSQL := strings.Split(upperSQL, limitStr)
	if len(splitSQL) > 2 {
		return "", errors.New("SQL语法异常，出现多个limit")
	}
	intStr, err := regexp.Compile(`\d+`)
	if err != nil {
		return "", err
	}

	for _, n := range intStr.FindAllString(limitStr, -1) {
		tmpN, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			return "", err
		}
		limitNumSlice = append(limitNumSlice, tmpN)
	}
	if len(limitNumSlice) > 1 && uint64(limitNumSlice[1]) > st.MaxReturnRows {
		limitStr = fmt.Sprintf("LIMIT %d, %d", limitNumSlice[0], st.MaxReturnRows)
	} else if len(limitNumSlice) == 1 && uint64(limitNumSlice[0]) > st.MaxReturnRows {
		limitNumSlice[0] = int64(st.MaxReturnRows)
		limitStr = fmt.Sprintf("LIMIT %d", st.MaxReturnRows)
	}
	ruleSql = fmt.Sprintf("%s %s", splitSQL[0], limitStr)
	return ruleSql, nil
}

func ClusterDBReadQueryExecuteHandler(rqData *models.DBReadQueryExecute, st *config.Das) (data *models.ReadQueryData, err error) {
	p := parser.New()
	stmt, err := p.ParseOneStmt(rqData.Sql, "", "")
	if err != nil {
		return nil, err
	}
	switch node := stmt.(type) {
	case *ast.SelectStmt:
		rqData.Sql, err = ReadQueryOptimize(node, rqData.Sql, st)
		if err != nil {
			return nil, err
		}
	case *ast.ExplainStmt:
		fmt.Println("desc / explain")
	case *ast.ShowStmt:
		fmt.Println("show create table")
	default:
		return nil, errors.New("sql类型异常")
	}
	data, err = mysql.ClusterDBReadQueryExecuteHandler(rqData)
	if err != nil {
		return
	}
	return data, err
}

func ClusterDBDataDictHandler(dbData *models.DBReadQueryExecute) (data []*models.DataDictJson, dbName string, err error) {
	return mysql.ClusterDBDataDictHandler(dbData)
}
