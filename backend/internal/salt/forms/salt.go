package forms

import "github.com/keepsty/go_rds/internal/salt/models"

// SaltMySQLProxysqlMha SaltStack MySQL+ProxySQL+MHA 部署参数
type SaltMySQLProxysqlMha struct {
	ENV                      string                       `json:"env"`
	SaltMysqlDepJson         *models.SaltMysqlDep         `json:"salt_mysql_dep_json"`
	SaltMysqlServerInfoJson  *models.SaltMysqlServerInfo  `json:"salt_mysql_server_info_json"`
	SaltProxySqlHostPostJson *models.SaltProxySqlHostPost `json:"salt_proxysql_host_post_json"`
}

// DeployRequest Salt 通用部署请求
type DeployRequest struct {
	Template string                 `json:"template" binding:"required"`
	Targets  []string               `json:"targets" binding:"required,min=1"`
	Config   map[string]interface{} `json:"config"`
}

// CreateTemplateForm 创建模板
type CreateTemplateForm struct {
	Name         string `form:"name" json:"name" binding:"required,min=2,max=64"`
	Title        string `form:"title" json:"title" binding:"required,min=2,max=128"`
	Description  string `form:"description" json:"description"`
	FieldsSchema string `form:"fields_schema" json:"fields_schema" binding:"required"`
	Defaults     string `form:"defaults" json:"defaults"`
}

// UpdateTemplateForm 更新模板
type UpdateTemplateForm struct {
	Title        string `form:"title" json:"title" binding:"required,min=2,max=128"`
	Description  string `form:"description" json:"description"`
	FieldsSchema string `form:"fields_schema" json:"fields_schema" binding:"required"`
	Defaults     string `form:"defaults" json:"defaults"`
}

// CreateHostConfigForm 创建主机配置
type CreateHostConfigForm struct {
	Name        string `form:"name" json:"name" binding:"required,min=2,max=128"`
	Hosts       string `form:"hosts" json:"hosts" binding:"required"`
	Description string `form:"description" json:"description"`
}

// UpdateHostConfigForm 更新主机配置
type UpdateHostConfigForm struct {
	Name        string `form:"name" json:"name" binding:"required,min=2,max=128"`
	Hosts       string `form:"hosts" json:"hosts" binding:"required"`
	Description string `form:"description" json:"description"`
}

// CreateTaskForm 创建部署任务
type CreateTaskForm struct {
	Name         string                 `form:"name" json:"name" binding:"required,min=2,max=128"`
	TemplateName string                 `form:"template_name" json:"template_name" binding:"required"`
	HostConfigID uint64                 `form:"host_config_id" json:"host_config_id" binding:"required"`
	ConfigParams map[string]interface{} `form:"config_params" json:"config_params"`
}

// ApproveTaskForm 审批任务
type ApproveTaskForm struct {
	Action string `form:"action" json:"action" binding:"required,oneof=approve reject"`
}
