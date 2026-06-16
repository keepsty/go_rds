package models

import "time"

type ServiceGroup struct {
	ID           int       `json:"id" db:"id"`
	ProdID       int       `json:"prod_id" db:"prod_id"`
	ServiceLevel int8      `json:"service_level" db:"service_level"`
	Environment  int8      `json:"environment" db:"environment"`
	HAType       int8      `json:"ha_type" db:"ha_type"`
	Middleware   int8      `json:"middleware" db:"middleware"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
	DnsName      string    `json:"dns_name" db:"dns_name"`
	Vip          string    `json:"vip" db:"vip"`
	PeakTime     string    `json:"peak_time" db:"peak_time"`
	Name         string    `json:"name" db:"name"`
	ProdName     string    `json:"prod_name" db:"prod_name"`
	DBAUser      string    `json:"dba_user" db:"dba_user"`
	RDOwner      string    `json:"rd_owner" db:"rd_owner"`
}

type Instance struct {
	ID             int       `json:"id" db:"id"`
	HostID         int       `json:"host_id" db:"host_id"`
	ServiceGroupID int       `json:"servicegroup_id" db:"servicegroup_id"`
	Port           int       `json:"port" db:"port"`
	Status         int       `json:"status" db:"status"`
	Role           int       `json:"role" db:"role"`
	CreateTime     time.Time `json:"create_time" db:"create_time"`
	Purpose        string    `json:"purpose" db:"purpose"`
	MysqlVersion   string    `json:"mysql_version" db:"mysql_version"`
}

type Hosts struct {
	ID             int       `json:"id" db:"id"`
	DiskType       int       `json:"disk_type" db:"disk_type"`
	DiskSize       int       `json:"disk_size" db:"disk_size"`
	Memory         int       `json:"memory" db:"memory"`
	RaidType       int       `json:"raid_type" db:"raid_type"`
	CpuNumber      int       `json:"cpu_number" db:"cpu_number"`
	Name           string    `json:"name" db:"name"`
	IP             string    `json:"ip" db:"ip"`
	CpuPlatform    string    `json:"cpu_platform" db:"cpu_platform"`
	SystemType     string    `json:"system_type" db:"system_type"`
	SystemFilename string    `json:"system_filename" db:"system_filename"`
	IDCLocation    string    `json:"idc_location" db:"idc_location"`
	CreateTime     time.Time `json:"create_time" db:"create_time"`
}

type ServiceGroupList struct {
	ID           int       `json:"id" db:"id"`
	ServiceLevel int8      `json:"service_level" db:"service_level"`
	Environment  int8      `json:"environment" db:"environment"`
	HAType       int8      `json:"ha_type" db:"ha_type"`
	Middleware   int8      `json:"middleware" db:"middleware"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
	DnsName      string    `json:"dns_name" db:"dns_name"`
	Vip          string    `json:"vip" db:"vip"`
	PeakTime     string    `json:"peak_time" db:"peak_time"`
	ProdName     string    `json:"prod_name" db:"prod_name"`
	Name         string    `json:"name" db:"name"`
	DBAUser      string    `json:"dba_user" db:"dba_user"`
	RDOwner      string    `json:"rd_owner" db:"rd_owner"`
}

type ServiceGroupResponse struct {
	Total    int                 `json:"total"`
	Page     int                 `json:"page"`
	Clusters []*ServiceGroupList `json:"clusters"`
}

type ServiceGroupDetail struct {
	ID                 int       `json:"id" db:"id"`
	ServiceLevel       int8      `json:"service_level" db:"service_level"`
	Environment        int8      `json:"environment" db:"environment"`
	HAType             int8      `json:"ha_type" db:"ha_type"`
	Middleware         int8      `json:"middleware" db:"middleware"`
	CreateTime         time.Time `json:"create_time" db:"create_time"`
	SGName             string    `json:"sg_name" db:"sg_name"`
	ProdName           string    `json:"prod_name" db:"prod_name"`
	DnsName            string    `json:"dns_name" db:"dns_name"`
	Vip                string    `json:"vip" db:"vip"`
	ClusterDescription string    `json:"cluster_description" db:"cluster_description"`
	PeakTime           string    `json:"peak_time" db:"peak_time"`
	DBAUser            string    `json:"dba_user" db:"dba_user"`
	RDOwner            string    `json:"rd_owner" db:"rd_owner"`
}

type ServiceGroupInsDetail struct {
	ID           int64     `json:"id" db:"id"`
	Port         int64     `json:"port" db:"port"`
	Role         int64     `json:"role" db:"role"`
	Status       int64     `json:"status" db:"status"`
	CreateTime   time.Time `json:"create_time" db:"create_time"`
	SgName       string    `json:"sg_name" db:"sg_name"`
	Hostname     string    `json:"hostname" db:"hostname"`
	Purpose      string    `json:"purpose" db:"purpose"`
	HostInfo     string    `json:"host_info" db:"host_info"`
	MysqlVersion string    `json:"mysql_version" db:"mysql_version"`
	IP           string    `json:"ip" db:"ip"`
}

type InstanceDetail struct {
	ID             int64     `json:"id" db:"id"`
	HostId         int64     `json:"host_id" db:"host_id"`
	ServiceGroupID int64     `json:"servicegroup_id" db:"servicegroup_id"`
	Port           int64     `json:"port" db:"port"`
	Status         int64     `json:"status" db:"status"`
	Role           int64     `json:"role" db:"role"`
	CreateTime     time.Time `json:"create_time" db:"create_time"`
	Purpose        string    `json:"purpose" db:"purpose"`
	MysqlVersion   string    `json:"mysql_version" db:"mysql_version"`
}

type ServiceGroupDBsDetail struct {
	ID              int64                        `json:"id" db:"id"`
	ServiceGroupId  int64                        `json:"servicegroup_id" db:"servicegroup_id"`
	DatabaseSize    float64                      `json:"database_size" db:"database_size"`
	Name            string                       `json:"name" db:"name"`
	RDUser          string                       `json:"rd_user" db:"rd_user"`
	DBRdLeader      string                       `json:"db_rd_leader" db:"db_rd_leader"`
	DatabaseCharset string                       `json:"database_charset" db:"database_charset"`
	CreateTime      time.Time                    `json:"create_time" db:"create_time"`
	UpdateTime      time.Time                    `json:"update_time" db:"update_time"`
	Tables          []*ServiceGroupDBTableDetail `json:"tables" db:""`
}

type ServiceGroupDBTableDetail struct {
	ID             int64     `json:"id" db:"id"`
	DbID           int64     `json:"db_id" db:"db_id"`
	TableSize      float64   `json:"table_size" db:"table_size"`
	FreeSize       float64   `json:"free_size" db:"free_size"`
	TableRows      int64     `json:"table_rows" db:"table_rows"`
	AutoIncrease   int64     `json:"auto_increase" db:"auto_increase"`
	TableName      string    `json:"table_name" db:"table_name"`
	TableCollation string    `json:"table_collation" db:"table_collation"`
	TbSchema       string    `json:"tb_schema" db:"tb_schema"`
	CreateTime     time.Time `json:"create_time" db:"create_time"`
	UpdateTime     time.Time `json:"update_time" db:"update_time"`
}

type ServiceGroupProxyListDetail struct {
	ID             int64                      `json:"id" db:"id"`
	ServiceGroupId int64                      `json:"servicegroup_id" db:"servicegroup_id"`
	AdminPort      int64                      `json:"admin_port" db:"admin_port"`
	AppPort        int64                      `json:"app_port" db:"app_port"`
	ProxyWeight    int64                      `json:"proxy_weight" db:"proxy_weight"`
	Hostname       string                     `json:"hostname" db:"hostname"`
	ProxyVersion   string                     `json:"proxy_version" db:"proxy_version"`
	RuleInfo       string                     `json:"rule_info" db:"rule_info"`
	ProxyMysqlServer []*ProxyHostInfoDetail   `json:"proxy_mysql_server" db:""`
}

type ServiceGroupProxyDetail struct {
	ID                   int64     `json:"id" db:"id"`
	ProxysqlID           int64     `json:"proxysql_id" db:"proxysql_id"`
	FlagOut              int64     `json:"flag_out" db:"flag_out"`
	DestinationHostgroup int64     `json:"destination_hostgroup" db:"destination_hostgroup"`
	Port                 int64     `json:"port" db:"port"`
	MatchPattern         string    `json:"match_pattern" db:"match_pattern"`
	Hostname             string    `json:"hostname" db:"hostname"`
	Status               string    `json:"status" db:"status"`
	UpdateTime           time.Time `json:"update_time" db:"update_time"`
}

type ProxyHostInfoDetail struct {
	ProxysqlID    int64     `json:"proxysql_id" db:"proxysql_id"`
	HostgroupID   int64     `json:"hostgroup_id" db:"hostgroup_id"`
	Weight        int64     `json:"weight" db:"weight"`
	Port          int64     `json:"port" db:"port"`
	UpdateTime    time.Time `json:"update_time" db:"update_time"`
	MysqlHostname string    `json:"mysql_hostname" db:"mysql_hostname"`
	Status        string    `json:"status" db:"status"`
}
