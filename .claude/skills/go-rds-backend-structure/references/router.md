# 路由注册规范（routers/）

模块 `routers/` 目录负责定义路由分组和中间件挂载。

## 主入口

每个模块必须导出 `func Routers(r *gin.Engine)`，由 `api.Include()` 统一收集。

## 路由层级

```
/api/v1/              → 无认证（如登录、刷新 token）
/api/v1/{模块}        → JWT 认证
/api/v1/admin/{模块}  → JWT + 管理员权限
```

## Routers 模板

```go
package routers

import (
    "github.com/gin-gonic/gin"
    "github.com/keepsty/go_rds/internal/global"
    "github.com/keepsty/go_rds/middleware"
    "github.com/keepsty/go_rds/internal/{模块}/views"
)

func Routers(r *gin.Engine) {
    // 用户侧 API（需 JWT）
    apiV1 := r.Group("/api/v1/{模块}")
    apiV1.Use(global.App.JWT.MiddlewareFunc())
    { RegisterApiRoutes(apiV1) }

    // 管理侧 API（需 JWT + 管理员）
    adminV1 := r.Group("/api/v1/admin/{模块}")
    adminV1.Use(global.App.JWT.MiddlewareFunc(), middleware.HasAdminPermission())
    { RegisterAdminRoutes(adminV1) }
}

func RegisterApiRoutes(v1 *gin.RouterGroup) {
    v1.GET("resource", views.GetResourceView)
    v1.POST("resource", views.CreateResourceView)
    v1.PUT("resource/:id", views.UpdateResourceView)
    v1.DELETE("resource/:id", views.DeleteResourceView)
}

func RegisterAdminRoutes(admin *gin.RouterGroup) {
    admin.GET("/config", views.AdminGetConfigView)
    admin.POST("/config", views.AdminCreateConfigView)
}
```

## 有公共路由 + 无需 admin 前缀的场景

如果模块有自己独立的 admin 前缀且不是 `/api/v1/admin/`，直接在 `Routers()` 内分层：

```go
func Routers(r *gin.Engine) {
    // 无需认证的公共路由
    r.GET("/ws/:channel", views.WebSocketHandler)

    apiV1 := r.Group("/api/v1/orders")
    apiV1.Use(global.App.JWT.MiddlewareFunc())
    { RegisterApiRoutes(apiV1) }

    adminV1 := r.Group("/api/v1/admin/approval-flows")
    adminV1.Use(global.App.JWT.MiddlewareFunc(), middleware.HasAdminPermission())
    { RegisterAdminRoutes(adminV1) }
}
```

## 路由注册到 app.go

在 `internal/app/app.go` 中加入 `api.Include()`：

```go
api.Include(
    userRouter.Routers,
    commonRouter.Routers,
    yourModuleRouter.Routers,
)
```
