package models

// BackupTemplate 备份模板
type BackupTemplate struct {
	ID          int64  `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	DbType      string `db:"db_type" json:"db_type"`
	Description string `db:"description" json:"description"`
	ConfigSchema string `db:"config_schema" json:"config_schema"`
	DefaultConfig string `db:"default_config" json:"default_config"`
	CreateTime  string `db:"create_time" json:"create_time"`
	UpdateTime  string `db:"update_time" json:"update_time"`
}

// BackupConfig 备份配置
type BackupConfig struct {
	ID            int64  `db:"id" json:"id"`
	Name          string `db:"name" json:"name"`
	DbType        string `db:"db_type" json:"db_type"`
	InstanceID    string `db:"instance_id" json:"instance_id"`
	BackupType    string `db:"backup_type" json:"backup_type"`
	ScheduleCron  string `db:"schedule_cron" json:"schedule_cron"`
	RetentionDays int    `db:"retention_days" json:"retention_days"`
	StoragePath   string `db:"storage_path" json:"storage_path"`
	Status        string `db:"status" json:"status"`
	Remark        string `db:"remark" json:"remark"`
	CreateTime    string `db:"create_time" json:"create_time"`
	UpdateTime    string `db:"update_time" json:"update_time"`
}

// BackupTask 备份任务
type BackupTask struct {
	ID         int64   `db:"id" json:"id"`
	ConfigID   int64   `db:"config_id" json:"config_id"`
	Name       string  `db:"name" json:"name"`
	DbType     string  `db:"db_type" json:"db_type"`
	BackupType string  `db:"backup_type" json:"backup_type"`
	InstanceID string  `db:"instance_id" json:"instance_id"`
	Status     string  `db:"status" json:"status"`
	StartedAt  *string `db:"started_at" json:"started_at"`
	FinishedAt *string `db:"finished_at" json:"finished_at"`
	CreatedBy  string  `db:"created_by" json:"created_by"`
	CreateTime string  `db:"create_time" json:"create_time"`
}

// BackupRecord 备份记录
type BackupRecord struct {
	ID         int64   `db:"id" json:"id"`
	TaskID     int64   `db:"task_id" json:"task_id"`
	FileName   *string `db:"file_name" json:"file_name"`
	FileSize   int64   `db:"file_size" json:"file_size"`
	FilePath   *string `db:"file_path" json:"file_path"`
	Status     string  `db:"status" json:"status"`
	Output     *string `db:"output" json:"output"`
	StartedAt  *string `db:"started_at" json:"started_at"`
	FinishedAt *string `db:"finished_at" json:"finished_at"`
	CreateTime string  `db:"create_time" json:"create_time"`
}