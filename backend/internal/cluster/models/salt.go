package models

type SaltPreCheckResult struct {
	Result bool   `json:"result"`
	Msg    string `json:"msg"`
}

type SaltMySQLProxysqlMha struct {
	ENV                      string                `json:"env"`
	SaltMysqlDepJson         *SaltMysqlDep         `json:"salt_mysql_dep_json"`
	SaltMysqlServerInfoJson  *SaltMysqlServerInfo  `json:"salt_mysql_server_info_json"`
	SaltProxySqlHostPostJson *SaltProxySqlHostPost `json:"salt_proxysql_host_post_json"`
}

type SaltMysqlDep struct {
	Version             string `json:"version"`
	PerconaToolkit      string `json:"percona_toolkit"`
	Xtrabackup24Rpm     string `json:"xtrabackup_24_rpm"`
	Xtrabackup24Tarball string `json:"xtrabackup_24_tarball"`
	Xtrabackup80Rpm     string `json:"xtrabackup_80_rpm"`
	Xtrabackup80Tarball string `json:"xtrabackup_80_tarball"`
}

type SaltMysqlServerInfo struct {
	User         string               `json:"user"`
	Password     string               `json:"password"`
	Database     string               `json:"database"`
	ServerName   string               `json:"server_name"`
	ServerMainer string               `json:"server_mainer"`
	ServerOps    string               `json:"server_ops"`
	HostPort     []*SaltMysqlHostPost `json:"host_port"`
}

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

type SaltProxySqlConfigJson struct {
	WriterHostgroup           int64  `json:"writer_hostgroup"`
	ReaderHostgroup           int64  `json:"reader_hostgroup"`
	BackupWriteHostGroup      int64  `json:"backup_write_host_group"`
	OfflineHostGroup          int64  `json:"offline_host_group"`
	AdminPort                 int64  `json:"admin_port"`
	ProxysqlDir               string `json:"proxysql_dir"`
	AdminPassword             string `json:"admin_password"`
	ClusterPassword           string `json:"cluster_password"`
	DbaPassword               string `json:"dba_password"`
	StatsPassword             string `json:"stats_password"`
	ChickAlicePassword        string `json:"chick_alice_password"`
	MonitorPassword           string `json:"monitor_password"`
	ProxysqlMysqlUserUsername string `json:"proxysql_mysql_user_username"`
	ProxysqlMysqlUserPassword string `json:"proxysql_mysql_user_password"`
	MonitorRPassword          string `json:"monitor_r_password"`
	MonitorRwPassword         string `json:"monitor_rw_password"`
	MhaPassword               string `json:"mha_password"`
	ProxysqlServersJson       string `json:"proxysql_servers_json"`
	ProxysqlMysqlServersJson  string `json:"proxysql_mysql_servers_json"`
}

type SaltProxySqlHostPost struct {
	AdminPort       int64              `json:"admin_port"`
	HostIP          []*SaltProxyHostIP `json:"host_ip"`
	InstanceType    string             `json:"instance_type"`
	ProxysqlDir     string             `json:"proxysql_dir"`
	ProxysqlConfDir string             `json:"proxysql_conf_dir"`
	ProxysqlRpm     string             `json:"proxysql_rpm"`
	MysqlVersion    string             `json:"mysql_version"`
}

type SaltProxyHostIP struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

type SaltStateFiles struct {
	FilePath       string `json:"file_path"`
	FileName       string `json:"file_name"`
	TargetFileName string `json:"target_file_name"`
	TargetFilePath string `json:"target_file_path"`
}
