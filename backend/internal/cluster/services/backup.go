package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/keepsty/go_rds/internal/cluster/kafka"
	"github.com/keepsty/go_rds/internal/cluster/dao/mysql"
	"github.com/keepsty/go_rds/internal/cluster/models"
	"github.com/keepsty/go_rds/internal/global"
)

// 备份任务状态机：合法的状态转换
var taskStatusTransitions = map[string][]string{
	"pending":  {"running", "failed"},
	"running":  {"success", "failed"},
	"success":  {},
	"failed":   {},
}

func validTaskStatusTransition(oldStatus, newStatus string) bool {
	allowed, ok := taskStatusTransitions[oldStatus]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == newStatus {
			return true
		}
	}
	return false
}

func sanitizeStoragePath(path string) string {
	// 移除潜在的路径穿越
	cleaned := strings.ReplaceAll(path, "..", "")
	return cleaned
}

func GetBackupConfigs(dbType string) ([]models.BackupConfig, error) {
	return mysql.GetBackupConfigs(dbType)
}

func CreateBackupConfig(c *models.BackupConfig) (int64, error) {
	if c.RetentionDays <= 0 {
		c.RetentionDays = 7
	}
	c.StoragePath = sanitizeStoragePath(c.StoragePath)
	return mysql.CreateBackupConfig(c)
}

func UpdateBackupConfig(id int64, c *models.BackupConfig) error {
	c.StoragePath = sanitizeStoragePath(c.StoragePath)
	return mysql.UpdateBackupConfig(id, c)
}

func DeleteBackupConfig(id int64) error {
	return mysql.DeleteBackupConfig(id)
}

func GetBackupTasks(status string) ([]models.BackupTask, error) {
	return mysql.GetBackupTasks(status)
}

func CreateBackupTask(configID int64, name, backupType, username string, configParams map[string]interface{}) (int64, error) {
	cfg, err := mysql.GetBackupConfigByID(configID)
	if err != nil {
		return 0, fmt.Errorf("备份配置不存在: %w", err)
	}
	task := &models.BackupTask{
		ConfigID:   configID,
		Name:       name,
		DbType:     cfg.DbType,
		BackupType: backupType,
		InstanceID: cfg.InstanceID,
		CreatedBy:  username,
	}
	taskID, err := mysql.CreateBackupTask(task)
	if err != nil {
		return 0, err
	}
	if configParams == nil {
		configParams = make(map[string]interface{})
	}
	msg := &kafka.BackupTaskMessage{
		TaskID:     taskID,
		TaskName:   name,
		ConfigName: cfg.Name,
		DbType:     cfg.DbType,
		BackupType: backupType,
		InstanceID: cfg.InstanceID,
		Config:     configParams,
		CreatedBy:  username,
		Timestamp:  time.Now().Format(time.RFC3339),
	}
	if global.App.KafkaProducer != nil {
		if err := global.App.KafkaProducer.Send(msg); err != nil {
			global.App.Log.Errorf("发送备份任务到Kafka失败: %v", err)
		}
	}
	return taskID, nil
}

func UpdateBackupTaskStatus(id int64, status string) error {
	task, err := mysql.GetBackupTaskByID(id)
	if err != nil {
		return fmt.Errorf("任务不存在: %w", err)
	}
	if !validTaskStatusTransition(task.Status, status) {
		return fmt.Errorf("非法状态转换: %s → %s", task.Status, status)
	}
	return mysql.UpdateBackupTaskStatus(id, status)
}

func GetBackupRecords(taskID int64) ([]models.BackupRecord, error) {
	return mysql.GetBackupRecords(taskID)
}

func GetBackupTemplates(dbType string) ([]models.BackupTemplate, error) {
	return mysql.GetBackupTemplates(dbType)
}

func CreateBackupTemplate(t *models.BackupTemplate) (int64, error) {
	return mysql.CreateBackupTemplate(t)
}

func UpdateBackupTemplate(id int64, t *models.BackupTemplate) error {
	return mysql.UpdateBackupTemplate(id, t)
}

func DeleteBackupTemplate(id int64) error {
	return mysql.DeleteBackupTemplate(id)
}
