package models

type ProxyHosts struct {
	ProxyID   int64  `json:"proxy_id" db:"proxy_id"`
	AdminPort int64  `json:"admin_port" db:"admin_port"`
	ProxyIP   string `json:"proxy_ip" db:"proxy_ip"`
	MysqlIP   string `json:"mysql_ip" db:"mysql_ip"`
}

type ProxysqlMysqlServers struct {
	Hostgroup         int64  `json:"hostgroup" db:"hostgroup"`
	Weight            int64  `json:"weight" db:"weight"`
	Port              int64  `json:"port" db:"port"`
	MaxReplicationLag int64  `json:"max_replication_lag"`
	Hostname          string `json:"hostname" db:"hostname"`
	Status            string `json:"status" db:"status"`
}

type ProxysqlServers struct {
	Port     int64  `json:"port"`
	Weight   int64  `json:"weight"`
	Hostname string `json:"hostname"`
}
