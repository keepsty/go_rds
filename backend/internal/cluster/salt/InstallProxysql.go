package salt

import (
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/keepsty/go_rds/internal/cluster/models"
	saltConfig "github.com/keepsty/go_rds/internal/cluster/salt/config"
	"github.com/keepsty/go_rds/internal/cluster/salt/sftp"
	"github.com/keepsty/go_rds/internal/config"
	"github.com/keepsty/go_rds/internal/global"
)

const (
	ProxysqlMaxReplicationLagSuper  = 0
	ProxysqlMaxReplicationLagNormal = 120
	ProxysqlWeightSuper             = 100
	ProxysqlWeightNormal            = 99
	ProxysqlMHAWriteHostGroup       = 10
	ProxysqlMHAReadHostGroup        = 11
	ProxysqlMGRWriteHostGroup       = 10
	ProxysqlMGRBackupWriteHostGroup = 11
	ProxysqlMGRReadHostGroup        = 12
	ProxysqlMGROfflineHostGroup     = 13
	ProxysqlMysqlWightSuper         = 1
	ProxysqlMysqlWightNormal        = 1000
	ProxysqlMysqlStatusOnline       = "ONLINE"
)

func InstallProxysqlHandler(env string, php *models.SaltProxySqlHostPost, si *models.SaltMysqlServerInfo, conf *config.ProxySQL, saltConf *config.Salt) (saltData []interface{}, err error) {
	data, err := ProxysqlConfHandler(php, si, conf)
	if err != nil {
		return
	}
	if SaltProxysqlConfInit(env, data, php, saltConf) != nil {
		return nil, err
	}

	for _, v := range php.HostIP {
		authBody := InitAuth(saltConf)
		state_sls := fmt.Sprintf(`{"client": "local", "tgt": "%s", "fun": "state.sls","kwarg":{"mods": "%s.proxysql_%d.deploy_proxysql","saltenv":"%s"}}`, v.Hostname, v.Hostname, php.AdminPort, env)
		body, err := DoQuery(authBody, state_sls, saltConf)
		if err != nil {
			zap.L().Error("InstallMySQLHandler failed, err: ", zap.Error(err))
			return nil, err
		}
		zap.L().Info(string(body))
		b, _ := json.Marshal(body)
		saltData = append(saltData, b)
	}
	return
}

func ProxysqlConfHandler(php *models.SaltProxySqlHostPost, si *models.SaltMysqlServerInfo, conf *config.ProxySQL) (data *models.SaltProxySqlConfigJson, err error) {
	data = &models.SaltProxySqlConfigJson{
		ReaderHostgroup:           ProxysqlMHAReadHostGroup,
		WriterHostgroup:           ProxysqlMHAWriteHostGroup,
		ProxysqlDir:               php.ProxysqlDir,
		AdminPort:                 php.AdminPort,
		ProxysqlMysqlUserUsername: si.User,
		ProxysqlMysqlUserPassword: si.Password,
		AdminPassword:             conf.Password,
		ClusterPassword:           conf.ClusterPassword,
		DbaPassword:               conf.DbaPassword,
		StatsPassword:             conf.StatsPassword,
		ChickAlicePassword:        conf.ChickAlicePassword,
		MonitorPassword:           conf.MonitorPassword,
		MonitorRPassword:          conf.MonitorRPassword,
		MonitorRwPassword:         conf.MonitorRwPassword,
		MhaPassword:               conf.MhaPassword,
	}
	for k, v := range php.HostIP {
		if k == 0 {
			data.ProxysqlServersJson = fmt.Sprintf("{hostname=\"%s\",port=%d,weight=%d}", v.IP, php.AdminPort, ProxysqlWeightSuper)
		} else {
			data.ProxysqlServersJson += fmt.Sprintf(",\n{hostname=\"%s\",port=%d,weight= %d}", v.IP, php.AdminPort, ProxysqlWeightNormal)
		}
	}
	for k, v := range si.HostPort {
		if k == 0 {
			data.ProxysqlMysqlServersJson = fmt.Sprintf("{\naddress=\"%s\"\nport=%d\nhostgroup=%d\nstatus=\"%s\"\nweight=%d\nmax_replication_lag=%d\n}", v.MysqlIp, v.Port, ProxysqlMHAWriteHostGroup, ProxysqlMysqlStatusOnline, ProxysqlMysqlWightSuper, ProxysqlMaxReplicationLagSuper)
			data.ProxysqlMysqlServersJson += fmt.Sprintf("\n,{\naddress=\"%s\"\nport=%d\nhostgroup=%d\nstatus=\"%s\"\nweight=%d\nmax_replication_lag=%d\n}", v.MysqlIp, v.Port, ProxysqlMHAReadHostGroup, ProxysqlMysqlStatusOnline, ProxysqlMysqlWightSuper, ProxysqlMaxReplicationLagSuper)
		} else {
			data.ProxysqlMysqlServersJson += fmt.Sprintf("\n,{\naddress=\"%s\"\nport=%d\nhostgroup=%d\nstatus=\"%s\"\nweight=%d\nmax_replication_lag=%d\n}", v.MysqlIp, v.Port, ProxysqlMHAReadHostGroup, ProxysqlMysqlStatusOnline, ProxysqlMysqlWightNormal, ProxysqlMaxReplicationLagNormal)
		}
	}
	return
}

func SaltProxysqlConfInit(env string, data *models.SaltProxySqlConfigJson, php *models.SaltProxySqlHostPost, saltConf *config.Salt) (err error) {
	files := make([]*models.SaltStateFiles, 1)
	curr_dir, _ := os.Getwd()
	config_path := fmt.Sprintf("%s/salt/config/default", curr_dir)
	sftpCnf := global.App.Config.Sftp
	sftp_conn, _ := sftp.Connect(sftpCnf.User, sftpCnf.Password, sftpCnf.Hostname, sftpCnf.Port)
	for _, h := range php.HostIP {
		proxycnf := new(models.SaltStateFiles)
		init_proxysql_user_sh := new(models.SaltStateFiles)
		if php.InstanceType == "mha" || php.InstanceType == "mgr" {
			proxycnf.FileName = fmt.Sprintf("%s_proxysql_conf", php.InstanceType)
		} else {
			err = fmt.Errorf("InstanceType '%s' is error.", php.InstanceType)
			return
		}
		proxycnf.FilePath = config_path
		proxycnf.TargetFileName = "proxysql.cnf"
		proxycnf.TargetFilePath = fmt.Sprintf("%s/salt/config/%s/proxysql_%d", curr_dir, h.Hostname, php.AdminPort)
		files = append(files, proxycnf)
		init_proxysql_user_sh.FileName = "init_proxysql_user.sh"
		init_proxysql_user_sh.FilePath = config_path
		init_proxysql_user_sh.TargetFileName = "init_proxysql_user.sh"
		init_proxysql_user_sh.TargetFilePath = fmt.Sprintf("%s/salt/config/%s/proxysql_%d", curr_dir, h.Hostname, php.AdminPort)
		files = append(files, init_proxysql_user_sh)
		err := saltConfig.FileTemplate(data, files)
		if err != nil {
			return err
		}

		proxysql_sls := new(models.SaltStateFiles)
		//proxysql_sls_path :=fmt.Sprintf("%s/salt/config/state_sls",curr_dir)
		php.ProxysqlConfDir = fmt.Sprintf("%s/proxysql_%d", h.Hostname, php.AdminPort)
		proxysql_sls.FilePath = fmt.Sprintf("%s/salt/config/state_sls", curr_dir)
		proxysql_sls.FileName = "deploy_proxysql.sls"
		proxysql_sls.TargetFileName = "deploy_proxysql.sls"
		proxysql_sls.TargetFilePath = fmt.Sprintf("%s/salt/config/%s/proxysql_%d", curr_dir, h.Hostname, php.AdminPort)
		f := make([]*models.SaltStateFiles, 1)
		f = append(f, proxysql_sls)
		err = saltConfig.FileTemplate(php, f)
		if err != nil {
			return err
		}

		sftp.UploadDirectory(sftp_conn, fmt.Sprintf("%s/salt/config/%s", curr_dir, h.Hostname), fmt.Sprintf("%s/%s/%s", saltConf.BaseDir, env, h.Hostname))
	}
	return err
}
