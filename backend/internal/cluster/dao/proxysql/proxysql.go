package proxysql

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"

	"github.com/keepsty/go_rds/internal/cluster/models"
	"github.com/keepsty/go_rds/internal/global"
)

func SetProxyConfig(job <-chan *models.ProxyHosts, wg *sync.WaitGroup, status int64) {
	defer wg.Done()
	conf := global.App.Config.ProxySQL
	for i := range job {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", conf.User, conf.Password, i.ProxyIP, i.AdminPort, conf.DBName)
		itemDB, err := sqlx.Connect("mysql", dsn)
		if err != nil {
			return
		}
		queryStr := ""
		switch status {
		case 0:
			queryStr = fmt.Sprintf("UPDATE mysql_servers SET status='OFFLINE_SOFT' WHERE hostname='%s'", i.MysqlIP)
		case 1:
			queryStr = fmt.Sprintf("UPDATE mysql_servers SET status='ONLINE' WHERE hostname='%s'", i.MysqlIP)
		}
		if queryStr == "" {
			itemDB.Close()
			return
		}
		_, err = itemDB.Exec(queryStr)
		if err != nil {
			itemDB.Close()
			return
		}
		_, _ = itemDB.Exec("load mysql servers to runtime;")
		_, _ = itemDB.Exec("save mysql servers to disk;")
		itemDB.Close()
	}
}
