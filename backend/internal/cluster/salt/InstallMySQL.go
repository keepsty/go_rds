package salt

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/keepsty/go_rds/internal/cluster/models"
	saltConfig "github.com/keepsty/go_rds/internal/cluster/salt/config"
	"github.com/keepsty/go_rds/internal/cluster/salt/sftp"
	"github.com/keepsty/go_rds/internal/config"
	"github.com/keepsty/go_rds/internal/global"
)

func MySQLConfInit(hp []*models.SaltMysqlHostPost) (files []*models.SaltStateFiles, hosts []string, err error) {
	curr_dir, err := os.Getwd()
	config_path := fmt.Sprintf("%s/salt/config/default", curr_dir)
	files = make([]*models.SaltStateFiles, 20)
	hosts = make([]string, len(hp))
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
	sftp_conn, _ := sftp.Connect(sftpCnf.User, sftpCnf.Password, sftpCnf.Hostname, sftpCnf.Port)
	curr_dir, err := os.Getwd()
	_, _, err = MySQLConfInit(si.HostPort)

	for _, v := range si.HostPort {
		sftp.UploadDirectory(sftp_conn, fmt.Sprintf("%s/salt/config/%s", curr_dir, v.Host), fmt.Sprintf("%s/%s/%s", saltConf.BaseDir, env, v.Host))
		authBody := InitAuth(saltConf)
		//state_sls := fmt.Sprintf(`{"client": "runner", "tgt": "%s", "fun": "state.sls","kwarg":{"mods": "mysql.files.%s.%d.mysql_install_instance","saltenv":"%s"}}`,v.Host,v.Host, v.Port, "prod")
		//state_sls := fmt.Sprintf(`{"client": "local", "tgt": "%s", "fun": "state.sls","kwarg":{"mods": "mysql.files.%s.%d.mysql_install_instance","saltenv":"%s"}}`,v.Host,v.Host, v.Port, "prod")
		state_sls := fmt.Sprintf(`{"client": "local", "tgt": "%s", "fun": "state.sls","kwarg":{"mods": "%s.mysql_%d.deploy_mysql_instance","saltenv":"%s"}}`, v.Host, v.Host, v.Port, env)
		body, err := DoQuery(authBody, state_sls, saltConf)
		if err != nil {
			zap.L().Error("InstallMySQLHandler failed, err: ", zap.Error(err))
			return nil, err
		}
		zap.L().Info(string(body))
		b, _ := json.Marshal(body)
		data = append(data, b)
	}
	return
}

func MysqlDependCheckHandler(env string, imd *models.SaltMysqlDep, hp []*models.SaltMysqlHostPost, conf *config.Salt) (data []interface{}, err error) {
	MysqlDependConfInit(env, imd, hp, conf)
	return
}

func MysqlDependConfInit(env string, imd *models.SaltMysqlDep, hp []*models.SaltMysqlHostPost, conf *config.Salt) (data []interface{}, err error) {
	sftpCnf := global.App.Config.Sftp
	sftp_conn, _ := sftp.Connect(sftpCnf.User, sftpCnf.Password, sftpCnf.Hostname, sftpCnf.Port)
	curr_dir, err := os.Getwd()
	files := make([]*models.SaltStateFiles, 20)
	for _, v := range hp {
		if v == nil {
			return
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

	for _, v := range hp {
		sftp.UploadDirectory(sftp_conn, fmt.Sprintf("%s/salt/config/%s", curr_dir, v.Host), fmt.Sprintf("%s/%s/%s", conf.BaseDir, env, v.Host))
		authBody := InitAuth(conf)
		state_sls := fmt.Sprintf(`{"client": "local", "tgt": "%s", "fun": "state.sls","kwarg":{"mods": "%s.mysql_%d.deploy_mysql_dep","saltenv":"%s"}}`, v.Host, v.Host, v.Port, env)
		body, err := DoQuery(authBody, state_sls, conf)
		if err != nil {
			zap.L().Error("InstallMySQLHandler failed, err: ", zap.Error(err))
			return nil, err
		}
		zap.L().Info(string(body))
		b, _ := json.Marshal(body)
		data = append(data, b)
	}
	return
}
