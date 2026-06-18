# 启动流程 + 全局单例

## 启动链路

```
cmd/main.go 解析 -config 参数
  → app.Run(configFile)
    → bootstrap.InitializeConfig()    # Viper 读取 yaml → global.App.Config
    → bootstrap.InitializeLog()       # logrus 日志
    → bootstrap.InitializeDB()        # GORM 自动迁移 + 种子数据
    → bootstrap.InitializeRedis()     # Redis 客户端（可选）
    → 模块独立初始化（如 clusterMysql.Init/AutoMigrate/Seed）
    → bootstrap.InitializeCron()      # 定时任务
    → app.RunServer()
      → middleware.InitAuthMiddleware()       # JWT 中间件
      → api.Include(模块.Routers...)          # 收集所有模块路由
      → api.Init()                            # 创建 Gin 引擎 + CORS + request-id + 日志
      → setupStaticFiles() / setupNoRoute()   # 嵌入式前端 SPA
      → r.Run(listenAddress)
```

## 全局单例（internal/global/global.go）

```go
type Application struct {
    ConfigViper *viper.Viper
    Config      config.Configuration
    JWT         *jwt.GinJWTMiddleware
    Log         *logrus.Logger
    DB          *gorm.DB
    Redis       *redis.Client
    Cron        *cron.Cron
}

var App = new(Application)
```

任意位置通过 `global.App.DB`、`global.App.Config`、`global.App.Log` 访问。

## 模型自动迁移注册

所有模型集中在 `bootstrap/db.go` 的 `initializeTables()` 注册：

```go
func initializeTables(db *gorm.DB) {
    err := db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(
        &usersModels.InsightUsers{},
        &commonModels.InsightInstances{},
        // ...
        &{模块}Models.XxxModel{},   // ← 新增
    )
}
```

## 种子数据

在 `initializeMySQLGorm()` 中依次调用（按依赖顺序）：

```go
initializeTables(db)
initializeAllowedOperations(db)   // DAS 允许操作
initializeGlobalInspectParams(db) // 审核参数
initializeAdminUser(db)           // 管理员 admin 用户
initializeNotifySettings(db)      // 通知配置
```

## cmd/main.go 入口模板

```go
var (
    Version   string
    BuildTime string
    GitCommit string
    GitBranch string
)

func main() {
    var configFile string
    var showVersion bool
    flag.StringVar(&configFile, "config", "config.yaml", "config file path")
    flag.BoolVar(&showVersion, "version", false, "show version and exit")
    flag.Parse()
    if showVersion { printVersion(); return }
    app.Run(configFile)
}
```

## app.RunServer 路由注册

```go
api.Include(
    userRouter.Routers,
    commonRouter.Routers,
    inspectRouter.Routers,
    dasRouter.Routers,
    ordersRouter.Routers,
    clusterRouter.Routers,
    // 新增模块路由
)
```
