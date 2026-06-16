package routers

import (
	"github.com/keepsty/go_rds/internal/global"
	"github.com/keepsty/go_rds/middleware"

	"github.com/keepsty/go_rds/internal/orders/views"

	"github.com/gin-gonic/gin"
)

func Routers(r *gin.Engine) {
	r.GET("/ws/:channel", views.WebSocketHandler)

	apiV1 := r.Group("/api/v1/orders")
	apiV1.Use(global.App.JWT.MiddlewareFunc())
	{
		RegisterApiRoutes(apiV1)
	}

	adminV1 := r.Group("/api/v1/admin/approval-flows")
	adminV1.Use(global.App.JWT.MiddlewareFunc(), middleware.HasAdminPermission())
	{
		RegisterAdminRoutes(adminV1)
	}
}
