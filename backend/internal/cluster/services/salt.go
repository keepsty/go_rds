package services

import (
	"github.com/keepsty/go_rds/internal/config"
	"github.com/keepsty/go_rds/internal/salt/models"
	saltServices "github.com/keepsty/go_rds/internal/salt/services"
)

func SaltAddClusterClusterPreCheckHandler() (data *models.SaltPreCheckResult, err error) {
	return
}

func SaltInstallMysqlDeployHandler(env string, imd *models.SaltMysqlDep, hp []*models.SaltMysqlHostPost, conf *config.Salt) (data []interface{}, err error) {
	saltServices.MysqlDependCheckHandler(env, imd, hp, conf)
	return
}

func SaltAddMysqlClusterHandler(env string, si *models.SaltMysqlServerInfo, conf *config.Salt) (data []interface{}, err error) {
	data, err = saltServices.InstallMySQLHandler(env, si, conf)
	if err != nil {
		return
	}
	return data, err
}

func SaltInstallProxysqlHandler(env string, php *models.SaltProxySqlHostPost, si *models.SaltMysqlServerInfo, conf *config.ProxySQL, saltConf *config.Salt) (data []interface{}, err error) {
	saltServices.InstallProxysqlHandler(env, php, si, conf, saltConf)
	return
}
