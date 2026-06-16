package services

import (
	"sync"

	"github.com/keepsty/go_rds/internal/cluster/dao/mysql"
	"github.com/keepsty/go_rds/internal/cluster/dao/proxysql"
	"github.com/keepsty/go_rds/internal/cluster/models"
)

func SetInstanceStatusByID(id int64, status int64) (data int64, err error) {
	jobChan := make(chan *models.ProxyHosts, 20)

	mysql.GetProxysqlInfoBYSGID(id, jobChan)
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go proxysql.SetProxyConfig(jobChan, wg, status)
	}
	wg.Wait()
	_, err = mysql.SetInstanceStatusByID(id, status)
	return
}
