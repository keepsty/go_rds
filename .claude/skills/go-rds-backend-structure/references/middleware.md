# 中间件栈

## 中间件总览

| 中间件 | 文件 | 说明 |
|--------|------|------|
| `middleware.InitAuthMiddleware()` | `middleware/jwt.go` | JWT 认证，基于 `appleboy/gin-jwt/v2` |
| `global.App.JWT.MiddlewareFunc()` | — | 保护需要登录的路由 |
| `middleware.HasAdminPermission()` | `middleware/permissions.go` | 检查 `is_superuser`，非管理员返回 403 |
| `middleware.OTPMiddleware()` | `middleware/otp.go` | 登录前验证用户名/密码 + TOTP 双因子 |
| `GetUserNameFromJWT(c)` | `middleware/claims.go` | 从 JWT claims 提取用户名 `(string, bool)` |
| `LoggerRequestToFile(logger)` | `middleware/log.go` | 请求日志（logrus + lumberjack 轮转） |

## 挂载位置

- **全局中间件**：在 `api/api.go` 的 `Init()` 中挂载（CORS、request-id、请求日志）
- **路由级中间件**：在模块的 `routers/routers.go` 中通过 `r.Group().Use()` 挂载

## 使用示例

```go
// 路由分组保护
apiV1 := r.Group("/api/v1/{模块}")
apiV1.Use(global.App.JWT.MiddlewareFunc())

adminV1 := r.Group("/api/v1/admin/{模块}")
adminV1.Use(global.App.JWT.MiddlewareFunc(), middleware.HasAdminPermission())

// 登录链：OTP → JWT
v1.POST("/user/login", middleware.OTPMiddleware(), global.App.JWT.LoginHandler)
```

## 从 JWT 提取用户名（claims.go）

```go
func GetUserNameFromJWT(c *gin.Context) (string, bool) {
    claims := jwt.ExtractClaims(c)
    raw, ok := claims["id"]
    if !ok { return "", false }
    username, ok := raw.(string)
    return username, ok
}
```
