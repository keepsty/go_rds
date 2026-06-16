package services

import (
	"github.com/keepsty/go_rds/internal/cluster/dao/mysql"
	"github.com/keepsty/go_rds/internal/cluster/forms"
	"github.com/keepsty/go_rds/internal/cluster/models"
)

func GetClusters(pl *forms.ParamClusterList) (data []*models.ServiceGroupList, err error) {
	data, err = mysql.GetClusters(pl)
	return
}

func GetClusterByID(id int64) (data []*models.ServiceGroupDetail, err error) {
	data, err = mysql.GetClusterByID(id)
	return
}

func GetClusterInsByID(id int64) (data []*models.ServiceGroupInsDetail, err error) {
	data, err = mysql.GetClusterInsByID(id)
	return
}

func GetClusterDBSByID(id int64) (data []*models.ServiceGroupDBsDetail, err error) {
	queryDbCnt := "select count(1) from database_database as db where db.servicegroup_id=?"
	cnt, err := mysql.GetDataCount(queryDbCnt, &id)
	if err != nil {
		return nil, err
	}
	DBdata := make([]*models.ServiceGroupDBsDetail, 0, cnt)
	DBdata, err = mysql.GetClusterDBsByID(id, cnt)
	if err != nil {
		return nil, err
	}
	tables, err := GetClusterTablesBySGID(id)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(DBdata); i++ {
		DBdata[i].Tables = make([]*models.ServiceGroupDBTableDetail, 0, len(tables))
		for j := 0; j < len(tables); j++ {
			if DBdata[i].ID == tables[j].DbID {
				DBdata[i].Tables = append(DBdata[i].Tables, tables[j])
			}
		}
	}
	data = DBdata
	return
}

func GetClusterTablesBySGID(id int64) (data []*models.ServiceGroupDBTableDetail, err error) {
	data, err = mysql.GetClusterTablesBySGID(id)
	return
}

func GetClusterDBDetailByID(id int64) (data []*models.ServiceGroupDBTableDetail, err error) {
	data, err = mysql.GetClusterDBDetailByID(id)
	return
}

func GetClusterProxyListByID(id int64) (data []*models.ServiceGroupProxyListDetail, err error) {
	queryDbCnt := "select count(1) from database_proxysql as dp where  dp.servicegroup_id=?"
	cnt, err := mysql.GetDataCount(queryDbCnt, &id)
	if err != nil {
		return nil, err
	}
	ProxyList := make([]*models.ServiceGroupProxyListDetail, 0, cnt)
	ProxyList, err = mysql.GetClusterProxyListByID(id, cnt)
	if err != nil {
		return nil, err
	}
	ProxyMysqlServer, err := GetClusterProxyHostInfoBySGID(id)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(ProxyList); i++ {
		ProxyList[i].ProxyMysqlServer = make([]*models.ProxyHostInfoDetail, 0, len(ProxyMysqlServer))
		for k := 0; k < len(ProxyMysqlServer); k++ {
			if ProxyList[i].ID == ProxyMysqlServer[k].ProxysqlID {
				ProxyList[i].ProxyMysqlServer = append(ProxyList[i].ProxyMysqlServer, ProxyMysqlServer[k])
			}
		}
	}
	data = ProxyList
	return
}

func GetClusterProxyDetailBySGID(id int64) (data []*models.ServiceGroupProxyDetail, err error) {
	data, err = mysql.GetClusterProxyDetailBySGID(id)
	return
}

func GetClusterProxyHostInfoBySGID(id int64) (data []*models.ProxyHostInfoDetail, err error) {
	data, err = mysql.GetClusterProxyHostInfoBySGID(id)
	return
}
