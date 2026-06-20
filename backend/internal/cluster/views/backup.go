package views

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/keepsty/go_rds/internal/cluster/forms"
	"github.com/keepsty/go_rds/internal/cluster/models"
	"github.com/keepsty/go_rds/internal/cluster/services"
	"github.com/keepsty/go_rds/pkg/response"
)

func handleValidationError(c *gin.Context, err error) {
	if _, ok := err.(validator.ValidationErrors); !ok {
		response.ValidateFail(c, "参数错误")
		return
	}
	response.ValidateFail(c, err.Error())
}

func GetBackupConfigsView(c *gin.Context) {
	dbType := c.Query("db_type")
	list, err := services.GetBackupConfigs(dbType)
	if err != nil { response.Fail(c, err.Error()); return }
	response.Success(c, list, "success")
}

func CreateBackupConfigView(c *gin.Context) {
	var form forms.CreateBackupConfigForm
	if err := c.ShouldBindJSON(&form); err != nil {
		handleValidationError(c, err)
		return
	}
	cfg := &models.BackupConfig{
		Name: form.Name, DbType: form.DbType, InstanceID: form.InstanceID,
		BackupType: form.BackupType, ScheduleCron: form.ScheduleCron,
		RetentionDays: form.RetentionDays, StoragePath: form.StoragePath, Remark: form.Remark,
	}
	id, err := services.CreateBackupConfig(cfg)
	if err != nil { response.Fail(c, err.Error()); return }
	response.Success(c, gin.H{"id": id}, "创建成功")
}

func UpdateBackupConfigView(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok { return }
	var form forms.UpdateBackupConfigForm
	if err := c.ShouldBindJSON(&form); err != nil {
		handleValidationError(c, err)
		return
	}
	cfg := &models.BackupConfig{
		Name: form.Name, DbType: form.DbType, InstanceID: form.InstanceID,
		BackupType: form.BackupType, ScheduleCron: form.ScheduleCron,
		RetentionDays: form.RetentionDays, StoragePath: form.StoragePath,
		Status: form.Status, Remark: form.Remark,
	}
	if err := services.UpdateBackupConfig(id, cfg); err != nil { response.Fail(c, err.Error()); return }
	response.Success(c, nil, "更新成功")
}

func DeleteBackupConfigView(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok { return }
	if err := services.DeleteBackupConfig(id); err != nil { response.Fail(c, err.Error()); return }
	response.Success(c, nil, "删除成功")
}

func GetBackupTasksView(c *gin.Context) {
	status := c.Query("status")
	list, err := services.GetBackupTasks(status)
	if err != nil { response.Fail(c, err.Error()); return }
	response.Success(c, list, "success")
}

func CreateBackupTaskView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok { return }
	var form forms.CreateBackupTaskForm
	if err := c.ShouldBindJSON(&form); err != nil {
		handleValidationError(c, err)
		return
	}
	id, err := services.CreateBackupTask(form.ConfigID, form.Name, form.BackupType, username, form.ConfigParams)
	if err != nil { response.Fail(c, err.Error()); return }
	response.Success(c, gin.H{"id": id}, "创建成功")
}

func UpdateBackupTaskStatusView(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok { return }
	var form forms.UpdateTaskStatusForm
	if err := c.ShouldBindJSON(&form); err != nil {
		handleValidationError(c, err)
		return
	}
	if err := services.UpdateBackupTaskStatus(id, form.Status); err != nil { response.Fail(c, err.Error()); return }
	response.Success(c, nil, "状态更新成功")
}

func GetBackupRecordsView(c *gin.Context) {
	taskIDStr := c.Query("task_id")
	var taskID int64
	if taskIDStr != "" {
		var err error
		taskID, err = strconv.ParseInt(taskIDStr, 10, 64)
		if err != nil { response.ValidateFail(c, "非法参数: task_id"); return }
	}
	list, err := services.GetBackupRecords(taskID)
	if err != nil { response.Fail(c, err.Error()); return }
	response.Success(c, list, "success")
}

func GetBackupTemplatesView(c *gin.Context) {
	dbType := c.Query("db_type")
	list, err := services.GetBackupTemplates(dbType)
	if err != nil { response.Fail(c, err.Error()); return }
	response.Success(c, list, "success")
}

func CreateBackupTemplateView(c *gin.Context) {
	var form forms.CreateBackupTemplateForm
	if err := c.ShouldBindJSON(&form); err != nil {
		handleValidationError(c, err)
		return
	}
	t := &models.BackupTemplate{
		Name: form.Name, DbType: form.DbType, Description: form.Description,
		ConfigSchema: form.ConfigSchema, DefaultConfig: form.DefaultConfig,
	}
	id, err := services.CreateBackupTemplate(t)
	if err != nil { response.Fail(c, err.Error()); return }
	response.Success(c, gin.H{"id": id}, "创建成功")
}

func UpdateBackupTemplateView(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok { return }
	var form forms.UpdateBackupTemplateForm
	if err := c.ShouldBindJSON(&form); err != nil {
		handleValidationError(c, err)
		return
	}
	t := &models.BackupTemplate{
		Name: form.Name, DbType: form.DbType, Description: form.Description,
		ConfigSchema: form.ConfigSchema, DefaultConfig: form.DefaultConfig,
	}
	if err := services.UpdateBackupTemplate(id, t); err != nil { response.Fail(c, err.Error()); return }
	response.Success(c, nil, "更新成功")
}

func DeleteBackupTemplateView(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok { return }
	if err := services.DeleteBackupTemplate(id); err != nil { response.Fail(c, err.Error()); return }
	response.Success(c, nil, "删除成功")
}
