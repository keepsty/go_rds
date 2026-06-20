package forms

// CreateBackupConfigForm 创建备份配置
type CreateBackupConfigForm struct {
	Name          string `json:"name" binding:"required,min=2,max=128"`
	DbType        string `json:"db_type" binding:"required,oneof=MySQL Redis TiDB"`
	InstanceID    string `json:"instance_id" binding:"required"`
	BackupType    string `json:"backup_type" binding:"required,oneof=full incremental"`
	ScheduleCron  string `json:"schedule_cron"`
	RetentionDays int    `json:"retention_days"`
	StoragePath   string `json:"storage_path"`
	Remark        string `json:"remark"`
}

// UpdateBackupConfigForm 更新备份配置
type UpdateBackupConfigForm struct {
	Name          string `json:"name" binding:"required,min=2,max=128"`
	DbType        string `json:"db_type" binding:"required,oneof=MySQL Redis TiDB"`
	InstanceID    string `json:"instance_id" binding:"required"`
	BackupType    string `json:"backup_type" binding:"required,oneof=full incremental"`
	ScheduleCron  string `json:"schedule_cron"`
	RetentionDays int    `json:"retention_days"`
	StoragePath   string `json:"storage_path"`
	Status        string `json:"status" binding:"required,oneof=enabled disabled"`
	Remark        string `json:"remark"`
}

// CreateBackupTaskForm 创建备份任务
type CreateBackupTaskForm struct {
	ConfigID   int64  `json:"config_id" binding:"required"`
	Name       string `json:"name" binding:"required,min=2,max=128"`
	BackupType string `json:"backup_type" binding:"required,oneof=full incremental"`
	ConfigParams map[string]interface{} `json:"config_params"`
}

// UpdateTaskStatusForm 更新任务状态
type UpdateTaskStatusForm struct {
	Status string `json:"status" binding:"required,oneof=running success failed"`
}
// 备份模板
type CreateBackupTemplateForm struct {
	Name         string `json:"name" binding:"required,min=2,max=128"`
	DbType       string `json:"db_type" binding:"required,oneof=MySQL Redis TiDB"`
	Description  string `json:"description"`
	ConfigSchema string `json:"config_schema" binding:"required"`
	DefaultConfig string `json:"default_config"`
}

type UpdateBackupTemplateForm struct {
	Name         string `json:"name" binding:"required,min=2,max=128"`
	DbType       string `json:"db_type" binding:"required,oneof=MySQL Redis TiDB"`
	Description  string `json:"description"`
	ConfigSchema string `json:"config_schema" binding:"required"`
	DefaultConfig string `json:"default_config"`
}
