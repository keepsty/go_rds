package views

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/keepsty/go_rds/internal/global"
	"github.com/keepsty/go_rds/internal/salt/forms"
	"github.com/keepsty/go_rds/internal/salt/models"
	"github.com/keepsty/go_rds/internal/salt/services"
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

	_, err := services.InstallProxysqlHandler(
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

// GetTemplatesView 获取可用部署模板列表
func GetTemplatesView(c *gin.Context) {
	templates, err := services.GetTemplates()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, templates, "success")
}

// DeployTemplateView 执行模板部署（同步等待结果）
func DeployTemplateView(c *gin.Context) {
	var req forms.DeployRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	result, err := services.DeployFromTemplate(req.Template, req.Config, req.Targets)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, result, "部署完成")
}

// ListMinionsView 获取 minion 列表
func ListMinionsView(c *gin.Context) {
	svc := services.NewSaltService()
	list, err := svc.ListMinions()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, list, "success")
}

// CreateTemplateView 创建模板
func CreateTemplateView(c *gin.Context) {
	var form forms.CreateTemplateForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	if err := services.CreateTemplate(form.Name, form.Title, form.Description, form.FieldsSchema, form.Defaults); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "创建成功")
}

// UpdateTemplateView 更新模板
func UpdateTemplateView(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	var form forms.UpdateTemplateForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	if err := services.UpdateTemplate(id, form.Title, form.Description, form.FieldsSchema, form.Defaults); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "更新成功")
}

// DeleteTemplateView 删除模板
func DeleteTemplateView(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok {
		return
	}
	if err := services.DeleteTemplate(id); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "删除成功")
}

// ---------- 主机配置 ----------

func GetHostConfigsView(c *gin.Context) {
	list, err := services.GetHostConfigs()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, list, "success")
}

func CreateHostConfigView(c *gin.Context) {
	var form forms.CreateHostConfigForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	if err := services.CreateHostConfig(form.Name, form.Hosts, form.Description); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "创建成功")
}

func UpdateHostConfigView(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok { return }
	var form forms.UpdateHostConfigForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	if err := services.UpdateHostConfig(id, form.Name, form.Hosts, form.Description); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "更新成功")
}

func DeleteHostConfigView(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok { return }
	if err := services.DeleteHostConfig(id); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "删除成功")
}

// ---------- 部署任务 ----------

func GetTasksView(c *gin.Context) {
	status := c.Query("status")
	list, err := services.GetTasks(status)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, list, "success")
}

func CreateTaskView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok { return }
	var form forms.CreateTaskForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	id, err := services.CreateTask(form.Name, form.TemplateName, form.HostConfigID, form.ConfigParams, username)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id}, "创建成功")
}

func ApproveTaskView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok { return }
	id, ok := parseUint64Param(c, "id")
	if !ok { return }
	var form forms.ApproveTaskForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	if err := services.ApproveTask(id, form.Action, username); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "操作成功")
}

func RunTaskView(c *gin.Context) {
	id, ok := parseUint64Param(c, "id")
	if !ok { return }
	if err := services.RunTask(id); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "执行完成")
}
