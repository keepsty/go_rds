package app

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/keepsty/go_rds/api"
	"github.com/keepsty/go_rds/internal/bootstrap"
	clusterMysql "github.com/keepsty/go_rds/internal/cluster/dao/mysql"
	"github.com/keepsty/go_rds/internal/global"
	"github.com/keepsty/go_rds/middleware"
	"github.com/keepsty/go_rds/web"

	clusterRouter "github.com/keepsty/go_rds/internal/cluster/routers"
	commonRouter "github.com/keepsty/go_rds/internal/common/routers"
	dasRouter "github.com/keepsty/go_rds/internal/das/routers"
	inspectRouter "github.com/keepsty/go_rds/internal/inspect/routers"
	ordersRouter "github.com/keepsty/go_rds/internal/orders/routers"
	saltRouter "github.com/keepsty/go_rds/internal/salt/routers"
	userRouter "github.com/keepsty/go_rds/internal/users/routers"

	"github.com/gin-gonic/gin"
)

const mediaDir = "./media"

func setupStaticFiles(r *gin.Engine) error {
	// Embedded file system - 映射整个 dist 目录
	distFS, err := fs.Sub(web.StaticFS, "dist")
	if err != nil {
		return fmt.Errorf("error accessing embedded filesystem: %w", err)
	}

	// 映射 assets 目录
	assetsFS, err := fs.Sub(distFS, "assets")
	if err != nil {
		return fmt.Errorf("error accessing assets directory: %w", err)
	}
	r.StaticFS("/assets", http.FS(assetsFS))

	// public 目录下的文件（avatar.png / favicon.ico）
	r.StaticFileFS("/avatar.png", "avatar.png", http.FS(distFS))
	r.StaticFileFS("/favicon.ico", "favicon.ico", http.FS(distFS))

	// 业务上传文件目录
	if _, err := os.Stat(mediaDir); os.IsNotExist(err) {
		if err := os.MkdirAll(mediaDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create media directory: %w", err)
		}
	}
	r.Static("/media", mediaDir)

	return nil
}

func setupNoRoute(r *gin.Engine) {
	// Fix 404 issue on page refresh
	r.NoRoute(func(c *gin.Context) {
		if strings.Contains(c.Request.Header.Get("Accept"), "text/html") {
			if content, err := web.StaticFS.ReadFile("dist/index.html"); err == nil {
				c.Header("Content-Type", "text/html; charset=utf-8")
				c.Data(http.StatusOK, "text/html; charset=utf-8", content)
				return
			}
		}
		c.String(http.StatusNotFound, "Not Found")
	})
}

func setupRootRoute(r *gin.Engine) {
	// Root route
	r.GET("/", func(c *gin.Context) {
		if data, err := web.StaticFS.ReadFile("dist/index.html"); err == nil {
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
		} else {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
	})
}

func RunServer() {
	// Production mode
	if global.App.Config.App.Environment == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize authentication middleware
	var err error
	if global.App.JWT, err = middleware.InitAuthMiddleware(); err != nil {
		fmt.Println("Failed to initialize authentication middleware:", err)
		return
	}

	// Load route configs for multiple APPs
	api.Include(
		userRouter.Routers,
		commonRouter.Routers,
		inspectRouter.Routers,
		dasRouter.Routers,
		ordersRouter.Routers,
		clusterRouter.Routers,
		saltRouter.Routers,
	)

	// Initialize router
	r := api.Init()

	// Static files and routes
	if err := setupStaticFiles(r); err != nil {
		fmt.Println(err)
		return
	}
	setupNoRoute(r)
	setupRootRoute(r)

	// Error handling
	r.Use(gin.Recovery())

	// Start server
	if err := r.Run(global.App.Config.App.ListenAddress); err != nil {
		fmt.Println("Failed to start server: ", err.Error())
	}
}

func Run(configFile string) {
	bootstrap.InitializeConfig(configFile)
	bootstrap.InitializeLog()
	global.App.DB = bootstrap.InitializeDB()
	global.App.Redis = bootstrap.InitializeRedis()

	// 初始化 cluster 模块的独立 sqlx 数据库连接
	if err := clusterMysql.Init(&global.App.Config.Database); err != nil {
		fmt.Println("Failed to initialize cluster database connection:", err)
		return
	}

	// 自动创建 cluster 模块所需的数据库表
	if err := clusterMysql.AutoMigrate(); err != nil {
		fmt.Println("Failed to auto migrate cluster tables:", err)
		return
	}

	// cluster 模块初始化种子数据
	if err := clusterMysql.InitializeSeedData(); err != nil {
		fmt.Println("Failed to initialize cluster seed data:", err)
		return
	}
	defer func() {
		if global.App.DB != nil {
			db, _ := global.App.DB.DB()
			db.Close()
		}
	}()
	bootstrap.InitializeCron()
	bootstrap.InitializeKafka()
	RunServer()
}
