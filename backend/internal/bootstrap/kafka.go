package bootstrap

import (
	"github.com/keepsty/go_rds/internal/cluster/kafka"
	"github.com/keepsty/go_rds/internal/global"
)

func InitializeKafka() {
	global.App.KafkaProducer = kafka.NewProducer(
		global.App.Config.Kafka.Brokers,
		global.App.Config.Kafka.Topic,
		global.App.Config.Kafka.LogDir,
	)
}
