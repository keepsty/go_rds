package models

import (
	"github.com/keepsty/go_rds/internal/common/models"
	"gorm.io/datatypes"
)

// ---------- 部署模板 ----------

// InsightSaltTemplates Salt 部署任务模板
type InsightSaltTemplates struct {
	*models.Model
	Name         string         `gorm:"type:varchar(64);not null;uniqueIndex:uniq_name;comment:模板标识" json:"name"`
	Title        string         `gorm:"type:varchar(128);not null;comment:模板标题" json:"title"`
	Description  string         `gorm:"type:varchar(512);not null;default:'';comment:模板描述" json:"description"`
	FieldsSchema datatypes.JSON `gorm:"type:json;comment:字段定义JSON" json:"fields_schema"`
	Defaults     datatypes.JSON `gorm:"type:json;comment:默认值JSON" json:"defaults"`
}

func (InsightSaltTemplates) TableName() string {
	return "insight_salt_templates"
}

// ---------- 主机配置 ----------

// SaltHostConfig 主机配置
type SaltHostConfig struct {
	*models.Model
	Name        string         `gorm:"type:varchar(128);not null;uniqueIndex:uniq_name;comment:配置名称" json:"name"`
	Hosts       datatypes.JSON `gorm:"type:json;comment:主机列表JSON" json:"hosts"`
	Description string         `gorm:"type:varchar(512);not null;default:'';comment:描述" json:"description"`
}

func (SaltHostConfig) TableName() string {
	return "insight_salt_host_configs"
}

// ---------- 部署任务 ----------

// SaltTask 部署任务
type SaltTask struct {
	*models.Model
	Name         string         `gorm:"type:varchar(128);not null;index;comment:任务名称" json:"name"`
	TemplateName string         `gorm:"type:varchar(64);not null;comment:模板标识" json:"template_name"`
	HostConfigID uint64         `gorm:"not null;comment:主机配置ID" json:"host_config_id"`
	ConfigParams datatypes.JSON `gorm:"type:json;comment:任务配置参数" json:"config_params"`
	Status       string         `gorm:"type:varchar(32);not null;default:'pending';comment:状态 pending/approved/running/success/failed" json:"status"`
	CreatedBy    string         `gorm:"type:varchar(64);not null;comment:创建人" json:"created_by"`
	ApprovedBy   string         `gorm:"type:varchar(64);not null;default:'';comment:审批人" json:"approved_by"`
	RunOutput    datatypes.JSON `gorm:"type:json;comment:执行输出" json:"run_output"`
	StartedAt    *models.LocalTime `gorm:"comment:开始时间" json:"started_at"`
	FinishedAt   *models.LocalTime `gorm:"comment:完成时间" json:"finished_at"`
}

func (SaltTask) TableName() string {
	return "insight_salt_tasks"
}

// ---------- 部署相关类型 ----------

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
