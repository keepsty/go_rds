package kafka

// BackupTaskMessage Kafka 备份任务消息
type BackupTaskMessage struct {
	TaskID       int64                  `json:"task_id"`
	TaskName     string                 `json:"task_name"`
	ConfigName   string                 `json:"config_name"`
	DbType       string                 `json:"db_type"`
	BackupType   string                 `json:"backup_type"`
	InstanceID   string                 `json:"instance_id"`
	Config       map[string]interface{} `json:"config"`
	CreatedBy    string                 `json:"created_by"`
	Timestamp    string                 `json:"timestamp"`
}
