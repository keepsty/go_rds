package views

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/keepsty/go_rds/internal/cluster/forms"
	"github.com/keepsty/go_rds/internal/cluster/models"
	"github.com/keepsty/go_rds/internal/cluster/services"
	"github.com/keepsty/go_rds/internal/global"
	"github.com/keepsty/go_rds/pkg/response"
)

// AddMySQLClusterHandler 通过SaltStack自动部署MySQL+ProxySQL集群
func AddMySQLClusterHandler(c *gin.Context) {
	mha := new(forms.SaltMySQLProxysqlMha)
	if err := c.ShouldBindJSON(mha); err != nil {
		if _, ok := err.(validator.ValidationErrors); !ok {
			response.ValidateFail(c, "参数错误")
			return
		}
		response.ValidateFail(c, err.Error())
		return
	}

	version := mha.SaltMysqlDepJson.Version
	versionStr := strings.ReplaceAll(version, ".", "")
	for _, v := range mha.SaltMysqlServerInfoJson.HostPort {
		v.BaseDir = fmt.Sprintf("/usr/local/mysql_%s", versionStr)
		v.Version = versionStr
		v.MysqlDir = fmt.Sprintf("/data/mysql_%d", v.Port)
		v.Datadir = fmt.Sprintf("%s/data", v.MysqlDir)
		s := strings.Split(v.MysqlIp, ".")
		if len(s) < 4 {
			response.Fail(c, "IP格式错误")
			return
		}
		v.ServerId, _ = strconv.ParseInt(fmt.Sprintf("%s%s%d", s[2], s[3], v.Port), 10, 64)
	}

	mha.SaltProxySqlHostPostJson.ProxysqlDir = fmt.Sprintf("/data/proxysql_%d", mha.SaltProxySqlHostPostJson.AdminPort)
	mha.SaltMysqlDepJson.Version = versionStr

	saltConf := global.App.Config.Salt
	proxySQLConf := global.App.Config.ProxySQL

	// 将 forms 类型转换为 models 类型
	b, _ := json.Marshal(mha.SaltProxySqlHostPostJson)
	var phpModel models.SaltProxySqlHostPost
	json.Unmarshal(b, &phpModel)

	b2, _ := json.Marshal(mha.SaltMysqlServerInfoJson)
	var siModel models.SaltMysqlServerInfo
	json.Unmarshal(b2, &siModel)

	_, err := services.SaltInstallProxysqlHandler(
		mha.ENV,
		&phpModel,
		&siModel,
		&proxySQLConf,
		&saltConf,
	)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "部署任务已提交")
}
