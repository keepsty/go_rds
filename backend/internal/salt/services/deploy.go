package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/keepsty/go_rds/internal/salt/models"
	"github.com/keepsty/go_rds/internal/config"
	saltConfig "github.com/keepsty/go_rds/internal/salt/config"
	"github.com/keepsty/go_rds/internal/salt/sftp"
	"github.com/keepsty/go_rds/internal/global"
)

// ---------- MySQL 部署 ----------

func MySQLConfInit(hp []*models.SaltMysqlHostPost) (files []*models.SaltStateFiles, hosts []string, err error) {
	curr_dir, err := os.Getwd()
	if err != nil {
		return nil, nil, fmt.Errorf("获取工作目录失败: %w", err)
	}
	config_path := fmt.Sprintf("%s/salt/config/default", curr_dir)
	files = make([]*models.SaltStateFiles, 0, 20)
	hosts = make([]string, 0, len(hp))
	for _, v := range hp {
		if v != nil {
			hosts = append(hosts, v.Host)
			v_num, _ := strconv.ParseInt(v.Version, 10, 64)
			deploy_mysql_install_sls := new(models.SaltStateFiles)
			mycnf := new(models.SaltStateFiles)
			init_mysql := new(models.SaltStateFiles)

			if v_num > 8000 {
				mycnf.FileName = "mysql_80_cnf"
			} else {
				mycnf.FileName = "mysql_57_cnf"
			}
			mycnf.FilePath = config_path
			mycnf.TargetFileName = fmt.Sprintf("my_%d.cnf", v.Port)
			mycnf.TargetFilePath = fmt.Sprintf("%s/salt/config/%s/mysql_%d", curr_dir, v.Host, v.Port)
			files = append(files, mycnf)

			deploy_mysql_install_sls.FilePath = fmt.Sprintf("%s/salt/config/state_sls", curr_dir)
			deploy_mysql_install_sls.FileName = "deploy_mysql_instance.sls"
			deploy_mysql_install_sls.TargetFileName = "deploy_mysql_instance.sls"
			deploy_mysql_install_sls.TargetFilePath = fmt.Sprintf("%s/salt/config/%s/mysql_%d", curr_dir, v.Host, v.Port)
			files = append(files, deploy_mysql_install_sls)

			init_mysql.FilePath = fmt.Sprintf("%s/salt/config/default", curr_dir)
			init_mysql.FileName = "init_mysql.sh"
			init_mysql.TargetFileName = fmt.Sprintf("init_mysql_%d.sh", v.Port)
			init_mysql.TargetFilePath = fmt.Sprintf("%s/salt/config/%s/mysql_%d", curr_dir, v.Host, v.Port)
			files = append(files, init_mysql)
			err = saltConfig.FileTemplate(v, files)
			if err != nil {
				return nil, nil, err
			}
		}
	}
	return
}

func InstallMySQLHandler(env string, si *models.SaltMysqlServerInfo, saltConf *config.Salt) (data []interface{}, err error) {
	sftpCnf := global.App.Config.Sftp
	sftp_conn, err := sftp.Connect(sftpCnf.User, sftpCnf.Password, sftpCnf.Hostname, sftpCnf.Port)
	if err != nil {
		return nil, fmt.Errorf("SFTP 连接失败: %w", err)
	}
	curr_dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("获取工作目录失败: %w", err)
	}
	if si == nil {
		return nil, fmt.Errorf("SaltMysqlServerInfo 为空")
	}
	_, _, err = MySQLConfInit(si.HostPort)

	svc := NewSaltService()
	for _, v := range si.HostPort {
		sftp.UploadDirectory(sftp_conn, fmt.Sprintf("%s/salt/config/%s", curr_dir, v.Host), fmt.Sprintf("%s/%s/%s", saltConf.BaseDir, env, v.Host))
		state_sls := fmt.Sprintf("%s.mysql_%d.deploy_mysql_instance", v.Host, v.Port)
		resp, err := svc.RunState(v.Host, state_sls, false)
		if err != nil {
			zap.L().Error("InstallMySQLHandler failed, err: ", zap.Error(err))
			return nil, err
		}
		zap.L().Info("InstallMySQLHandler result", zap.Any("result", resp))
		b, _ := json.Marshal(resp)
		data = append(data, b)
	}
	return
}

func MysqlDependCheckHandler(env string, imd *models.SaltMysqlDep, hp []*models.SaltMysqlHostPost, conf *config.Salt) (data []interface{}, err error) {
	return MysqlDependConfInit(env, imd, hp, conf)
}

func MysqlDependConfInit(env string, imd *models.SaltMysqlDep, hp []*models.SaltMysqlHostPost, conf *config.Salt) (data []interface{}, err error) {
	sftpCnf := global.App.Config.Sftp
	sftp_conn, err := sftp.Connect(sftpCnf.User, sftpCnf.Password, sftpCnf.Hostname, sftpCnf.Port)
	if err != nil {
		return nil, fmt.Errorf("SFTP 连接失败: %w", err)
	}
	curr_dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("获取工作目录失败: %w", err)
	}
	files := make([]*models.SaltStateFiles, 0, 20)
	for _, v := range hp {
		if v == nil {
			continue
		}
		deploy_mysql_install_sls := new(models.SaltStateFiles)
		deploy_mysql_install_sls.FilePath = fmt.Sprintf("%s/salt/config/state_sls", curr_dir)
		deploy_mysql_install_sls.FileName = "deploy_mysql_dep.sls"
		deploy_mysql_install_sls.TargetFileName = "deploy_mysql_dep.sls"
		deploy_mysql_install_sls.TargetFilePath = fmt.Sprintf("%s/salt/config/%s/mysql_%d", curr_dir, v.Host, v.Port)
		files = append(files, deploy_mysql_install_sls)
	}

	err = saltConfig.FileTemplate(imd, files)
	if err != nil {
		return nil, err
	}

	svc := NewSaltService()
	for _, v := range hp {
		sftp.UploadDirectory(sftp_conn, fmt.Sprintf("%s/salt/config/%s", curr_dir, v.Host), fmt.Sprintf("%s/%s/%s", conf.BaseDir, env, v.Host))
		state_sls := fmt.Sprintf("%s.mysql_%d.deploy_mysql_dep", v.Host, v.Port)
		resp, err := svc.RunState(v.Host, state_sls, false)
		if err != nil {
			zap.L().Error("MysqlDependConfInit failed, err: ", zap.Error(err))
			return nil, err
		}
		zap.L().Info("MysqlDependConfInit result", zap.Any("result", resp))
		b, _ := json.Marshal(resp)
		data = append(data, b)
	}
	return
}

// ---------- ProxySQL 部署 ----------

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

	svc := NewSaltService()
	for _, v := range php.HostIP {
		state_sls := fmt.Sprintf("%s.proxysql_%d.deploy_proxysql", v.Hostname, php.AdminPort)
		resp, err := svc.RunState(v.Hostname, state_sls, false)
		if err != nil {
			zap.L().Error("InstallProxysqlHandler failed, err: ", zap.Error(err))
			return nil, err
		}
		zap.L().Info("InstallProxysqlHandler result", zap.Any("result", resp))
		b, _ := json.Marshal(resp)
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
	files := make([]*models.SaltStateFiles, 0, 1)
	curr_dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %w", err)
	}
	config_path := fmt.Sprintf("%s/salt/config/default", curr_dir)
	sftpCnf := global.App.Config.Sftp
	sftp_conn, err := sftp.Connect(sftpCnf.User, sftpCnf.Password, sftpCnf.Hostname, sftpCnf.Port)
	if err != nil {
		return fmt.Errorf("SFTP 连接失败: %w", err)
	}
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
