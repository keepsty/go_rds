package config

type App struct {
	ListenAddress string `mapstructure:"listen_address" json:"listen_address" yaml:"listen_address"`
	Environment   string `mapstructure:"environment" json:"environment" yaml:"environment"`
	SECRET_KEY    string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
}

type Crontab struct {
	SyncDBMetas string `mapstructure:"sync_db_metas" json:"sync_db_metas" yaml:"sync_db_metas"`
}

type Log struct {
	Level   string `mapstructure:"level" json:"level" yaml:"level"`
	RootDir string `mapstructure:"root_dir" json:"root_dir" yaml:"root_dir"`
}

type Database struct {
	Driver          string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Host            string `mapstructure:"host" json:"host" yaml:"host"`
	Port            int    `mapstructure:"port" json:"port" yaml:"port"`
	Database        string `mapstructure:"database" json:"database" yaml:"database"`
	UserName        string `mapstructure:"username" json:"username" yaml:"username"`
	Password        string `mapstructure:"password" json:"password" yaml:"password"`
	Charset         string `mapstructure:"charset" json:"charset" yaml:"charset"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxIdleTime int    `mapstructure:"conn_max_idle_time" json:"conn_max_idle_time" yaml:"conn_max_idle_time"`
	ConnMaxLifetime int    `mapstructure:"conn_max_life_time" json:"conn_max_life_time" yaml:"conn_max_life_time"`
}

type Redis struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

type Das struct {
	MaxExecutionTime  uint64 `mapstructure:"max_execution_time" json:"max_execution_time" yaml:"max_execution_time"`
	DefaultReturnRows uint64 `mapstructure:"default_return_rows" json:"default_return_rows" yaml:"default_return_rows"`
	MaxReturnRows     uint64 `mapstructure:"max_return_rows" json:"max_return_rows" yaml:"max_return_rows"`
}

type Ghost struct {
	Path string   `mapstructure:"path" json:"path" yaml:"path"`
	Args []string `mapstructure:"args" json:"args" yaml:"args"`
}

type Notify struct {
	NoticeURL string `mapstructure:"notice_url" json:"notice_url" yaml:"notice_url"`
	Wechat    struct {
		Enable  bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
		Webhook string `mapstructure:"webhook" json:"webhook" yaml:"webhook"`
	}
	Mail struct {
		Enable   bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
		Username string `mapstructure:"username" json:"username" yaml:"username"`
		Password string `mapstructure:"password" json:"password" yaml:"password"`
		Host     string `mapstructure:"host" json:"host" yaml:"host"`
		Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	}
	DingTalk struct {
		Enable   bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
		Webhook  string `mapstructure:"webhook" json:"webhook" yaml:"webhook"`
		Keywords string `mapstructure:"keywords" json:"keywords" yaml:"keywords"`
	}
}

// ProxySQL 配置
type ProxySQL struct {
	User               string `mapstructure:"user" json:"user" yaml:"user"`
	Password           string `mapstructure:"password" json:"password" yaml:"password"`
	DBName             string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
	ClusterPassword    string `mapstructure:"cluster_password" json:"cluster_password" yaml:"cluster_password"`
	DbaPassword        string `mapstructure:"dba_password" json:"dba_password" yaml:"dba_password"`
	StatsPassword      string `mapstructure:"stats_password" json:"stats_password" yaml:"stats_password"`
	ChickAlicePassword string `mapstructure:"chick_alice_password" json:"chick_alice_password" yaml:"chick_alice_password"`
	MonitorPassword    string `mapstructure:"monitor_password" json:"monitor_password" yaml:"monitor_password"`
	MonitorRPassword   string `mapstructure:"monitor_r_password" json:"monitor_r_password" yaml:"monitor_r_password"`
	MonitorRwPassword  string `mapstructure:"monitor_rw_password" json:"monitor_rw_password" yaml:"monitor_rw_password"`
	MhaPassword        string `mapstructure:"mha_password" json:"mha_password" yaml:"mha_password"`
}

// Salt 配置
type Salt struct {
	URL      string `mapstructure:"url" json:"url" yaml:"url"`
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Timeout  uint64 `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	BaseDir  string `mapstructure:"basedir" json:"basedir" yaml:"basedir"`
}

// SFTP 配置
type Sftp struct {
	User     string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Hostname string `mapstructure:"hostname" json:"hostname" yaml:"hostname"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
}
// Kafka 配置
type Kafka struct {
	Brokers []string `mapstructure:"brokers" json:"brokers" yaml:"brokers"`
	Topic   string   `mapstructure:"topic" json:"topic" yaml:"topic"`
	LogDir  string   `mapstructure:"log_dir" json:"log_dir" yaml:"log_dir"`
}


type Configuration struct {
	App      App      `mapstructure:"app" json:"app" yaml:"app"`
	Crontab  Crontab  `mapstructure:"crontab" json:"crontab" yaml:"crontab"`
	Log      Log      `mapstructure:"log" json:"log" yaml:"log"`
	Database Database `mapstructure:"database" json:"database" yaml:"database"`
	Redis    Redis    `mapstructure:"redis" json:"redis" yaml:"redis"`
	Das      Das      `mapstructure:"das" json:"das" yaml:"das"`
	Ghost    Ghost    `mapstructure:"ghost" json:"ghost" yaml:"ghost"`
	Notify   Notify   `mapstructure:"-" json:"notify" yaml:"-"`
	ProxySQL ProxySQL `mapstructure:"proxysql" json:"proxysql" yaml:"proxysql"`
	Salt     Salt     `mapstructure:"salt" json:"salt" yaml:"salt"`
	Sftp     Sftp     `mapstructure:"sftp" json:"sftp" yaml:"sftp"`
	Kafka    Kafka    `mapstructure:"kafka" json:"kafka" yaml:"kafka"`
}
