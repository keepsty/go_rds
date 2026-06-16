package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/keepsty/go_rds/internal/global"
	"github.com/keepsty/go_rds/middleware"
)

func Routers(r *gin.Engine) {
	apiV1 := r.Group("/api/v1/das")
	apiV1.Use(global.App.JWT.MiddlewareFunc())
	{
		RegisterDasRoutes(apiV1)
	}

	adminV1 := r.Group("/api/v1/admin/das")
	adminV1.Use(global.App.JWT.MiddlewareFunc(), middleware.HasAdminPermission())
	{
		RegisterAdminRoutes(adminV1)
	}
}
