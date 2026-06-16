package forms

// SaltMySQLProxysqlMha SaltStack MySQL+ProxySQL+MHA 部署参数
type SaltMySQLProxysqlMha struct {
	ENV                      string               `json:"env"`
	SaltMysqlDepJson         *SaltMysqlDep         `json:"salt_mysql_dep_json"`
	SaltMysqlServerInfoJson  *SaltMysqlServerInfo  `json:"salt_mysql_server_info_json"`
	SaltProxySqlHostPostJson *SaltProxySqlHostPost `json:"salt_proxysql_host_post_json"`
}

// SaltMysqlDep MySQL 依赖包配置
type SaltMysqlDep struct {
	Version             string `json:"version"`
	PerconaToolkit      string `json:"percona_toolkit"`
	Xtrabackup24Rpm     string `json:"xtrabackup_24_rpm"`
	Xtrabackup24Tarball string `json:"xtrabackup_24_tarball"`
	Xtrabackup80Rpm     string `json:"xtrabackup_80_rpm"`
	Xtrabackup80Tarball string `json:"xtrabackup_80_tarball"`
}

// SaltMysqlServerInfo MySQL 服务器信息
type SaltMysqlServerInfo struct {
	User         string              `json:"user"`
	Password     string              `json:"password"`
	Database     string              `json:"database"`
	ServerName   string              `json:"server_name"`
	ServerMainer string              `json:"server_mainer"`
	ServerOps    string              `json:"server_ops"`
	HostPort     []*SaltMysqlHostPost `json:"host_port"`
}

// SaltMysqlHostPost MySQL 主机端口配置
type SaltMysqlHostPost struct {
	Port                      int64  `json:"port"`
	InstanceType              int64  `json:"instance_type"`
	ServerId                  int64  `json:"server_id"`
	InnodbBufferPoolInstances int64  `json:"innodb_buffer_pool_instances"`
	Host                      string `json:"host"`
	Version                   string `json:"version"`
	MysqlDir                  string `json:"mysql_dir"`
	Datadir                   string `json:"datadir"`
	BaseDir                   string `json:"base_dir"`
	InnodbBufferPoolSize      string `json:"innodb_buffer_pool_size"`
	MysqlIp                   string `json:"mysql_ip"`
}

// SaltProxySqlHostPost ProxySQL 主机配置
type SaltProxySqlHostPost struct {
	AdminPort       int64             `json:"admin_port"`
	HostIP          []*SaltProxyHostIP `json:"host_ip"`
	InstanceType    string             `json:"instance_type"`
	ProxysqlDir     string             `json:"proxysql_dir"`
	ProxysqlConfDir string             `json:"proxysql_conf_dir"`
	ProxysqlRpm     string             `json:"proxysql_rpm"`
	MysqlVersion    string             `json:"mysql_version"`
}

// SaltProxyHostIP ProxySQL 主机 IP
type SaltProxyHostIP struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}
