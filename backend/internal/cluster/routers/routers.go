package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/keepsty/go_rds/internal/global"
)

func Routers(r *gin.Engine) {
	apiV1 := r.Group("/api/v1/cluster")
	apiV1.Use(global.App.JWT.MiddlewareFunc())
	{
		RegisterApiRoutes(apiV1)
	}
}
