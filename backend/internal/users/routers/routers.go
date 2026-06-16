package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/keepsty/go_rds/internal/global"
	"github.com/keepsty/go_rds/internal/users/views"
	"github.com/keepsty/go_rds/middleware"
)

func Routers(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// user auth
		v1.POST("/user/login", middleware.OTPMiddleware(), global.App.JWT.LoginHandler)
		v1.POST("/user/logout", global.App.JWT.LogoutHandler)
		v1.POST("/user/otp-auth-url", views.GetOTPAuthURLView)
		v1.POST("/user/otp-auth-callback", views.GetOTPAuthCallbackView)
		v1.GET("/user/refresh_token", global.App.JWT.RefreshHandler)
	}

	// 下面接口需要认证
	apiV1 := r.Group("/api/v1/profile")
	apiV1.Use(global.App.JWT.MiddlewareFunc())
	{
		RegisterApiRoutes(apiV1)
	}

	adminV1 := r.Group("/api/v1/admin")
	adminV1.Use(global.App.JWT.MiddlewareFunc(), middleware.HasAdminPermission())
	{
		RegisterAdminRoutes(adminV1)
	}
}
