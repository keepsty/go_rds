package bootstrap

import (
	"github.com/keepsty/go_rds/middleware"

	"github.com/keepsty/go_rds/internal/global"
)

func InitializeLog() {
	global.App.Log = middleware.InitLogger("app.log")
}
