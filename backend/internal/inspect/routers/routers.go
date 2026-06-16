package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/keepsty/go_rds/internal/global"
	"github.com/keepsty/go_rds/middleware"
)

func Routers(r *gin.Engine) {
	adminV1 := r.Group("/api/v1/admin/inspect")
	adminV1.Use(global.App.JWT.MiddlewareFunc(), middleware.HasAdminPermission())
	{
		RegisterAdminRoutes(adminV1)
	}
}
