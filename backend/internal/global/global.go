package global

import (
	config "github.com/keepsty/go_rds/internal/config"
	kafkaPkg "github.com/keepsty/go_rds/internal/cluster/kafka"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Application struct {
	ConfigViper    *viper.Viper
	Config         config.Configuration
	JWT            *jwt.GinJWTMiddleware
	Log            *logrus.Logger
	DB             *gorm.DB
	Redis          *redis.Client
	Cron           *cron.Cron
	KafkaProducer  *kafkaPkg.Producer
}

var App = new(Application)
