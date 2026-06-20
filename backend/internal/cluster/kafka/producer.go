package kafka

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Producer Kafka 消息生产者
type Producer struct {
	brokers []string
	topic   string
	logDir  string
}

// NewProducer 创建 Kafka 生产者
// 当 brokers 为空时，消息写入 logDir 下的日志文件（开发模式）
func NewProducer(brokers []string, topic, logDir string) *Producer {
	return &Producer{brokers: brokers, topic: topic, logDir: logDir}
}

// Send 发送备份任务消息到 Kafka
func (p *Producer) Send(msg *BackupTaskMessage) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	if len(p.brokers) == 0 || p.brokers[0] == "" {
		// 开发模式：写入日志文件
		return p.writeToFile(data)
	}

	// 生产模式：发送到 Kafka（请确保已引入 github.com/segmentio/kafka-go）
	// 当前使用日志模式，连接 Kafka 时取消下方注释
	/*
		writer := &kafka.Writer{
			Addr:     kafka.TCP(p.brokers...),
			Topic:    p.topic,
			Balancer: &kafka.LeastBytes{},
		}
		defer writer.Close()
		return writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(fmt.Sprintf("%d", msg.TaskID)),
				Value: data,
			},
		)
	*/
	return p.writeToFile(data)
}

func (p *Producer) writeToFile(data []byte) error {
	logPath := p.logDir + "/kafka_backup_messages.log"
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), string(data)))
	return err
}
