package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/keepsty/go_rds/internal/global"
	"github.com/keepsty/go_rds/internal/salt/views"
)

func Routers(r *gin.Engine) {
	// Salt 操作 API（需 JWT 认证）
	apiV1 := r.Group("/api/v1/salt")
	apiV1.Use(global.App.JWT.MiddlewareFunc())
	{
		RegisterApiRoutes(apiV1)
	}
}

func RegisterApiRoutes(v1 *gin.RouterGroup) {
	v1.POST("/deploy/mysql-cluster", views.AddMySQLClusterHandler)
}
