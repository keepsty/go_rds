package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/keepsty/go_rds/internal/global"
	"github.com/keepsty/go_rds/internal/salt/views"
	"github.com/keepsty/go_rds/middleware"
)

func Routers(r *gin.Engine) {
	// Salt 部署 API（需 JWT 认证，普通用户可用）
	apiV1 := r.Group("/api/v1/salt")
	apiV1.Use(global.App.JWT.MiddlewareFunc())
	{
		RegisterApiRoutes(apiV1)
	}

	// Salt 管理 API（需 JWT + 管理员权限）
	adminV1 := r.Group("/api/v1/admin/salt")
	adminV1.Use(global.App.JWT.MiddlewareFunc(), middleware.HasAdminPermission())
	{
		RegisterAdminRoutes(adminV1)
	}
}

// RegisterApiRoutes 普通用户可访问的 salt API
func RegisterApiRoutes(v1 *gin.RouterGroup) {
	v1.GET("/templates", views.GetTemplatesView)
	v1.POST("/templates/deploy", views.DeployTemplateView)
	v1.GET("/minions", views.ListMinionsView)
	// 主机配置（普通用户只读）
	v1.GET("/host-configs", views.GetHostConfigsView)
	// 部署任务
	v1.GET("/tasks", views.GetTasksView)
	v1.POST("/tasks", views.CreateTaskView)
	v1.POST("/tasks/:id/run", views.RunTaskView)
}

// RegisterAdminRoutes 管理员可访问的 salt API
func RegisterAdminRoutes(admin *gin.RouterGroup) {
	// 模版管理
	admin.POST("/templates", views.CreateTemplateView)
	admin.PUT("/templates/:id", views.UpdateTemplateView)
	admin.DELETE("/templates/:id", views.DeleteTemplateView)
	// 主机配置管理
	admin.POST("/host-configs", views.CreateHostConfigView)
	admin.PUT("/host-configs/:id", views.UpdateHostConfigView)
	admin.DELETE("/host-configs/:id", views.DeleteHostConfigView)
	// 任务审批
	admin.PUT("/tasks/:id/approve", views.ApproveTaskView)
	// MySQL 部署
	admin.POST("/deploy/mysql-cluster", views.AddMySQLClusterHandler)
}