package mysql

import (
	"fmt"

	"github.com/keepsty/go_rds/internal/cluster/models"
)

// ---------- 备份配置 ----------

func GetBackupConfigs(dbType string) ([]models.BackupConfig, error) {
	var list []models.BackupConfig
	query := "SELECT id, name, db_type, instance_id, backup_type, schedule_cron, retention_days, storage_path, status, remark, create_time, update_time FROM database_backup_config"
	args := []interface{}{}
	if dbType != "" {
		query += " WHERE db_type = ? ORDER BY id DESC"
		args = append(args, dbType)
	} else {
		query += " ORDER BY id DESC"
	}
	err := db.Select(&list, query, args...)
	return list, err
}

func GetBackupConfigByID(id int64) (*models.BackupConfig, error) {
	var c models.BackupConfig
	err := db.Get(&c, `SELECT id, name, db_type, instance_id, backup_type, schedule_cron, retention_days, storage_path, status, remark, create_time, update_time FROM database_backup_config WHERE id=?`, id)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func CreateBackupConfig(c *models.BackupConfig) (int64, error) {
	result, err := db.Exec(
		`INSERT INTO database_backup_config (name, db_type, instance_id, backup_type, schedule_cron, retention_days, storage_path, status, remark)
		 VALUES (?, ?, ?, ?, ?, ?, ?, 'enabled', ?)`,
		c.Name, c.DbType, c.InstanceID, c.BackupType, c.ScheduleCron, c.RetentionDays, c.StoragePath, c.Remark,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取插入ID失败: %w", err)
	}
	return id, nil
}

func UpdateBackupConfig(id int64, c *models.BackupConfig) error {
	_, err := db.Exec(
		`UPDATE database_backup_config SET name=?, db_type=?, instance_id=?, backup_type=?, schedule_cron=?, retention_days=?, storage_path=?, status=?, remark=? WHERE id=?`,
		c.Name, c.DbType, c.InstanceID, c.BackupType, c.ScheduleCron, c.RetentionDays, c.StoragePath, c.Status, c.Remark, id,
	)
	return err
}

func DeleteBackupConfig(id int64) error {
	// 检查是否有未完成的任务
	var taskCount int
	err := db.Get(&taskCount, `SELECT COUNT(*) FROM database_backup_task WHERE config_id=? AND status NOT IN ('success','failed')`, id)
	if err != nil {
		return err
	}
	if taskCount > 0 {
		return fmt.Errorf("该配置下有 %d 个未完成的任务，请先处理后再删除", taskCount)
	}
	_, err = db.Exec("DELETE FROM database_backup_config WHERE id=?", id)
	return err
}

// ---------- 备份任务 ----------

func GetBackupTasks(status string) ([]models.BackupTask, error) {
	var list []models.BackupTask
	query := "SELECT id, config_id, name, db_type, backup_type, instance_id, status, started_at, finished_at, created_by, create_time FROM database_backup_task"
	args := []interface{}{}
	if status != "" {
		query += " WHERE status = ? ORDER BY id DESC"
		args = append(args, status)
	} else {
		query += " ORDER BY id DESC"
	}
	err := db.Select(&list, query, args...)
	return list, err
}

func GetBackupTaskByID(id int64) (*models.BackupTask, error) {
	var t models.BackupTask
	err := db.Get(&t, "SELECT id, config_id, name, db_type, backup_type, instance_id, status, started_at, finished_at, created_by, create_time FROM database_backup_task WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func CreateBackupTask(t *models.BackupTask) (int64, error) {
	result, err := db.Exec(
		`INSERT INTO database_backup_task (config_id, name, db_type, backup_type, instance_id, status, created_by)
		 VALUES (?, ?, ?, ?, ?, 'pending', ?)`,
		t.ConfigID, t.Name, t.DbType, t.BackupType, t.InstanceID, t.CreatedBy,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取插入ID失败: %w", err)
	}
	return id, nil
}

func UpdateBackupTaskStatus(id int64, status string) error {
	_, err := db.Exec(
		`UPDATE database_backup_task
		 SET started_at = CASE WHEN ? = 'running' THEN NOW() ELSE started_at END,
		     finished_at = CASE WHEN ? IN ('success','failed') THEN NOW() ELSE finished_at END,
		     status = ?
		 WHERE id=?`,
		status, status, status, id,
	)
	return err
}

// ---------- 备份记录 ----------

func GetBackupRecords(taskID int64) ([]models.BackupRecord, error) {
	var list []models.BackupRecord
	query := "SELECT id, task_id, file_name, file_size, file_path, status, output, started_at, finished_at, create_time FROM database_backup_record"
	args := []interface{}{}
	if taskID > 0 {
		query += " WHERE task_id = ? ORDER BY id DESC"
		args = append(args, taskID)
	} else {
		query += " ORDER BY id DESC LIMIT 100"
	}
	err := db.Select(&list, query, args...)
	return list, err
}

func CreateBackupRecord(r *models.BackupRecord) (int64, error) {
	result, err := db.Exec(
		`INSERT INTO database_backup_record (task_id, file_name, file_size, file_path, status, output)
		 VALUES (?, ?, ?, ?, 'running', '')`,
		r.TaskID, r.FileName, r.FileSize, r.FilePath,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取插入ID失败: %w", err)
	}
	return id, nil
}

// ---------- 备份模板 ----------

func GetBackupTemplates(dbType string) ([]models.BackupTemplate, error) {
	var list []models.BackupTemplate
	query := "SELECT id, name, db_type, description, config_schema, default_config, create_time, update_time FROM database_backup_template"
	args := []interface{}{}
	if dbType != "" {
		query += " WHERE db_type = ? ORDER BY id DESC"
		args = append(args, dbType)
	} else {
		query += " ORDER BY id DESC"
	}
	err := db.Select(&list, query, args...)
	return list, err
}

func GetBackupTemplateByID(id int64) (*models.BackupTemplate, error) {
	var t models.BackupTemplate
	err := db.Get(&t, "SELECT id, name, db_type, description, config_schema, default_config, create_time, update_time FROM database_backup_template WHERE id=?", id)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func CreateBackupTemplate(t *models.BackupTemplate) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO database_backup_template (name, db_type, description, config_schema, default_config) VALUES (?, ?, ?, ?, ?)",
		t.Name, t.DbType, t.Description, t.ConfigSchema, t.DefaultConfig,
	)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("获取插入ID失败: %w", err)
	}
	return id, nil
}

func UpdateBackupTemplate(id int64, t *models.BackupTemplate) error {
	_, err := db.Exec(
		"UPDATE database_backup_template SET name=?, db_type=?, description=?, config_schema=?, default_config=? WHERE id=?",
		t.Name, t.DbType, t.Description, t.ConfigSchema, t.DefaultConfig, id,
	)
	return err
}

func DeleteBackupTemplate(id int64) error {
	_, err := db.Exec("DELETE FROM database_backup_template WHERE id=?", id)
	return err
}
