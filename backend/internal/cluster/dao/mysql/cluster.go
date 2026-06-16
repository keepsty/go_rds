package mysql

import (
	"github.com/keepsty/go_rds/internal/cluster/forms"
	"github.com/keepsty/go_rds/internal/cluster/models"
)

func GetClusters(pl *forms.ParamClusterList) (data []*models.ServiceGroupList, err error) {
	data = make([]*models.ServiceGroupList, 0, pl.Size)
	sqlStr := `select sg.id,
		IFNULL(pl.name,'') as prod_name,
		IFNULL(sg.name,'') as name,
		IFNULL(dba_user,'') as dba_user,
		IFNULL(rd_owner,'') as rd_owner,
		IFNULL(service_level,0) as service_level,
		IFNULL(environment,0) as environment,
		IFNULL(ha_type,0) as ha_type,
		IFNULL(middleware,0) as middleware,
		IFNULL(dns_name,'') as dns_name,
		IFNULL(vip,'') as vip,
		IFNULL(peak_time,'') as peak_time,
		sg.create_time
		from database_servicegroup sg left join production_line pl on sg.prod_id = pl.id
		order by sg.id desc limit ?,?`
	err = db.Select(&data, sqlStr, (pl.Page-1)*pl.Size, pl.Size)
	return
}

func GetClusterByID(id int64) (data []*models.ServiceGroupDetail, err error) {
	data = make([]*models.ServiceGroupDetail, 0, 1)
	sqlStr := `select sg.id,
		IFNULL(pl.name,'') as prod_name,
		IFNULL(sg.name,'') as sg_name,
		IFNULL(dba_user,'') as dba_user,
		IFNULL(rd_owner,'') as rd_owner,
		IFNULL(service_level,0) as service_level,
		IFNULL(environment,0) as environment,
		IFNULL(ha_type,0) as ha_type,
		IFNULL(middleware,0) as middleware,
		IFNULL(dns_name,'') as dns_name,
		IFNULL(vip,'') as vip,
		IFNULL(cluster_description,'') as cluster_description,
		IFNULL(peak_time,'') as peak_time,
		sg.create_time
		from database_servicegroup sg left join production_line pl on sg.prod_id = pl.id
		where sg.id=?`
	err = db.Select(&data, sqlStr, id)
	return
}

func GetClusterInsByID(id int64) (data []*models.ServiceGroupInsDetail, err error) {
	var insCount int64
	sqlTotalStr := `select count(1) from database_instance as di left join hosts as h on di.host_id = h.id left join database_servicegroup as dsg on di.servicegroup_id = dsg.id and dsg.id=?`
	if err := db.Get(&insCount, sqlTotalStr, id); err != nil {
		return nil, err
	}
	if insCount < 1 {
		err = ErrorNoRows
		return nil, err
	}
	data = make([]*models.ServiceGroupInsDetail, 0, insCount)
	sqlStr := `select di.id,
		IFNULL(dsg.name,'') as sg_name,
		IFNULL(h.name,'') as hostname,
		IFNULL(h.ip,'') as ip,
		IFNULL(concat(h.cpu_number,'C - ',h.memory,'G - ',h.disk_size,'G'),'') as host_info,
		IFNULL(mysql_version,'') as mysql_version,
		di.port,di.role,di.status,
		IFNULL(di.purpose,'') as purpose,
		di.create_time
		from database_instance as di
		left join hosts as h on di.host_id = h.id
		left join database_servicegroup as dsg on di.servicegroup_id = dsg.id
		where dsg.id=?`
	err = db.Select(&data, sqlStr, id)
	return
}

func GetDataCount(query string, param interface{}) (total int64, err error) {
	if err := db.Get(&total, query, param); err != nil {
		return 0, err
	}
	if total < 1 {
		err = ErrorNoRows
		return 0, err
	}
	return total, err
}

func GetClusterDBsByID(id int64, Total int64) (data []*models.ServiceGroupDBsDetail, err error) {
	data = make([]*models.ServiceGroupDBsDetail, 0, Total)
	sqlStr := `select id,
		IFNULL(name,'') as name,
		IFNULL(rd_user,'') as rd_user,
		servicegroup_id,
		IFNULL(database_size,0) as database_size,
		IFNULL(database_charset,'') as database_charset
		from database_database as db where db.servicegroup_id=?`
	err = db.Select(&data, sqlStr, id)
	return
}

func GetClusterDBDetailByID(id int64) (data []*models.ServiceGroupDBTableDetail, err error) {
	var dbsCount int64
	sqlTotalStr := `select count(1) from database_tables where db_id=?`
	if err := db.Get(&dbsCount, sqlTotalStr, id); err != nil {
		return nil, err
	}
	if dbsCount < 1 {
		err = ErrorNoRows
		return nil, err
	}
	data = make([]*models.ServiceGroupDBTableDetail, 0, dbsCount)
	sqlStr := `select id,
		IFNULL(table_name,'') as table_name,db_id,
		IFNULL(table_size,0) as table_size,
		IFNULL(free_size,0) as free_size,
		IFNULL(table_rows,0) as table_rows,
		IFNULL(auto_increase,0) as auto_increase,
		IFNULL(table_collation,'') as table_collation,
		IFNULL(tb_schema,'') as tb_schema
		from database_tables where db_id=?`
	err = db.Select(&data, sqlStr, id)
	return
}

func GetClusterTablesBySGID(id int64) (data []*models.ServiceGroupDBTableDetail, err error) {
	var dbsCount int64
	sqlTotalStr := `select count(1) from database_database as db,database_tables as tbs,database_servicegroup as dsg where db.servicegroup_id = dsg.id and db.id = tbs.db_id and dsg.id=?`
	if err := db.Get(&dbsCount, sqlTotalStr, id); err != nil {
		return nil, err
	}
	if dbsCount < 1 {
		err = ErrorNoRows
		return nil, err
	}
	data = make([]*models.ServiceGroupDBTableDetail, 0, dbsCount)
	sqlStr := `select tbs.id,
		IFNULL(tbs.table_name,'') as table_name,db_id,
		IFNULL(table_size,0) as table_size,
		IFNULL(free_size,0) as free_size,
		IFNULL(table_rows,0) as table_rows,
		IFNULL(tbs.auto_increase,0) as auto_increase,
		IFNULL(table_collation,'') as table_collation,
		IFNULL(tb_schema,'') as tb_schema
		from database_database as db,database_tables as tbs,database_servicegroup as dsg
		where db.servicegroup_id = dsg.id and db.id = tbs.db_id and dsg.id=?`
	err = db.Select(&data, sqlStr, id)
	return
}

func GetClusterProxyListByID(id int64, cnt int64) (data []*models.ServiceGroupProxyListDetail, err error) {
	if cnt < 1 {
		err = ErrorNoRows
		return nil, err
	}
	data = make([]*models.ServiceGroupProxyListDetail, 0, cnt)
	sqlStr := `select dp.id,
		IFNULL(h.name,'') as hostname,
		dp.admin_port,dp.app_port,
		IFNULL(dp.proxy_weight,0) as proxy_weight,
		IFNULL(dp.proxy_version,'') as proxy_version,
		IFNULL(dp.rule_info,'') as rule_info
		from database_proxysql as dp
		left join hosts as h on dp.host_id=h.id
		where dp.servicegroup_id=?`
	err = db.Select(&data, sqlStr, id)
	return
}

func GetClusterProxyDetailBySGID(id int64) (data []*models.ServiceGroupProxyDetail, err error) {
	sqlTotalStr := `select count(1) from database_proxysql as dp,database_proxysql_details as dpd,database_servicegroup as ds where dp.servicegroup_id=ds.id and dp.id=dpd.proxysql_id and ds.id=?`
	dbsCount, err := GetDataCount(sqlTotalStr, &id)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ServiceGroupProxyDetail, 0, dbsCount)
	sqlStr := `select dpd.id,dpd.proxysql_id,IFNULL(dpd.username,'') as username,IFNULL(dpd.flag_out,0) as flag_out,
		IFNULL(dpd.match_pattern,'') as match_pattern,IFNULL(dpd.destination_hostgroup,0) as destination_hostgroup,
		IFNULL(dpd.hostname,'') as hostname,IFNULL(dpd.port,0) as port,IFNULL(dpd.status,'') as status,dpd.update_time
		from database_proxysql as dp,database_proxysql_details as dpd,database_servicegroup as ds
		where dp.servicegroup_id=ds.id and dp.id=dpd.proxysql_id and ds.id=?`
	err = db.Select(&data, sqlStr, id)
	return
}

func GetClusterProxyHostInfoBySGID(id int64) (data []*models.ProxyHostInfoDetail, err error) {
	sqlTotalStr := `select count(1) from database_proxysql_runtime_mysql_server as dms,database_proxysql as dp,database_servicegroup as ds where dp.servicegroup_id=ds.id and dp.id=dms.proxysql_id and ds.id=?`
	dbsCount, err := GetDataCount(sqlTotalStr, &id)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ProxyHostInfoDetail, 0, dbsCount)
	sqlStr := `select dms.proxysql_id,dms.hostgroup_id,IFNULL(dms.mysql_hostname,'') as mysql_hostname,
		IFNULL(dms.port,0) as port,IFNULL(dms.status,'') as status,IFNULL(dms.weight,0) as weight,dms.update_time
		from database_proxysql_runtime_mysql_server as dms,database_proxysql as dp,database_servicegroup as ds
		where dp.servicegroup_id=ds.id and dp.id=dms.proxysql_id and ds.id=?`
	err = db.Select(&data, sqlStr, id)
	return
}
