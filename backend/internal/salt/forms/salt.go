package forms

import "github.com/keepsty/go_rds/internal/salt/models"

// SaltMySQLProxysqlMha SaltStack MySQL+ProxySQL+MHA 部署参数
type SaltMySQLProxysqlMha struct {
	ENV                      string                  `json:"env"`
	SaltMysqlDepJson         *models.SaltMysqlDep         `json:"salt_mysql_dep_json"`
	SaltMysqlServerInfoJson  *models.SaltMysqlServerInfo  `json:"salt_mysql_server_info_json"`
	SaltProxySqlHostPostJson *models.SaltProxySqlHostPost `json:"salt_proxysql_host_post_json"`
}
