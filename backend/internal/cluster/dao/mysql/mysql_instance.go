package mysql

import (
	"github.com/keepsty/go_rds/internal/cluster/models"
)

func GetProxysqlInfoBYSGID(id int64, job chan<- *models.ProxyHosts) {
	defer close(job)
	proxyHost := make([]*models.ProxyHosts, 0, 100)
	proxySQL := `select dp.id as proxy_id,dp.admin_port,h.ip as proxy_ip,dh.ip as mysql_ip from database_proxysql as dp,hosts as h,hosts as dh,database_instance as di where di.servicegroup_id=dp.servicegroup_id and dp.host_id=h.id and dh.id=di.host_id and di.id=?`
	err := db.Select(&proxyHost, proxySQL, id)
	if err != nil {
		return
	}
	for i := 0; i < len(proxyHost); i++ {
		job <- proxyHost[i]
	}
}

func SetInstanceStatusByID(id, status int64) (data int64, err error) {
	sqlStr := "update database_instance set status = ? where id =?"
	_, err = db.Exec(sqlStr, status, id)
	return
}
