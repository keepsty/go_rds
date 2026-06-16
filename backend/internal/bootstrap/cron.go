package bootstrap

import (
	"time"

	"github.com/keepsty/go_rds/internal/global"

	commonTasks "github.com/keepsty/go_rds/internal/common/tasks"
	dasTasks "github.com/keepsty/go_rds/internal/das/tasks"

	"github.com/robfig/cron/v3"
)

func InitializeCron() {
	global.App.Cron = cron.New()

	go func() {
		_, err := global.App.Cron.AddFunc(global.App.Config.Crontab.SyncDBMetas, func() {
			global.App.Log.Info("Run SyncDBMeta At:", time.Now())
			commonTasks.SyncDBMeta()
		})
		if err != nil {
			global.App.Log.Error(err)
		}

		_, err = global.App.Cron.AddFunc("*/5 * * * *", func() {
			global.App.Log.Info("Run KillTiDBQuery At:", time.Now())
			kill := dasTasks.KillTiDBQuery{}
			kill.Run()
		})
		if err != nil {
			global.App.Log.Error(err)
		}
		global.App.Cron.Start()
		defer global.App.Cron.Stop()
		select {}
	}()
}
