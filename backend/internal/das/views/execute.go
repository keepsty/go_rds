package views

import (
	"github.com/keepsty/go_rds/internal/global"

	"github.com/keepsty/go_rds/pkg/response"

	"github.com/keepsty/go_rds/internal/das/forms"
	"github.com/keepsty/go_rds/internal/das/models"
	"github.com/keepsty/go_rds/internal/das/services"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// 执行MySQL/TiDB查询
func ExecuteMySQLQueryView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form *forms.ExecuteMySQLQueryForm = &forms.ExecuteMySQLQueryForm{}
	var RequestID string = requestid.Get(c)

	if err := c.ShouldBind(&form); err == nil {
		service := services.ExecuteMySQLQueryService{
			ExecuteMySQLQueryForm: form,
			C:                     c,
			Username:              username,
		}
		returnData, err := service.Run()
		if err != nil {
			// 更新数据库记录
			global.App.DB.Model(&models.InsightDASRecords{}).
				Where("request_id=? and username=?", RequestID, username).
				Updates(map[string]any{"error_msg": err.Error(), "is_finish": true})
			response.Fail(c, err.Error())
		} else {
			// 更新数据库记录
			global.App.DB.Model(&models.InsightDASRecords{}).
				Where("request_id=? and username=?", RequestID, username).
				Updates(map[string]any{"is_finish": true})
			response.Success(c, returnData, "success")
		}
	} else {
		// 有效性验证不通过，更新数据库记录
		global.App.DB.Model(&models.InsightDASRecords{}).
			Where("request_id=? and username=?", RequestID, username).
			Updates(map[string]any{"error_msg": err.Error(), "is_finish": true})
		response.ValidateFail(c, err.Error())
	}
}

// 执行ClickHouse查询
func ExecuteClickHouseQueryView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form *forms.ExecuteClickHouseQueryForm = &forms.ExecuteClickHouseQueryForm{}
	var RequestID string = requestid.Get(c)

	if err := c.ShouldBind(&form); err == nil {
		service := services.ExecuteClickHouseQueryService{
			ExecuteClickHouseQueryForm: form,
			C:                          c,
			Username:                   username,
		}
		returnData, err := service.Run()
		if err != nil {
			// 更新数据库记录
			global.App.DB.Model(&models.InsightDASRecords{}).
				Where("request_id=? and username=?", RequestID, username).
				Updates(map[string]any{"error_msg": err.Error(), "is_finish": true})
			response.Fail(c, err.Error())
		} else {
			// 更新数据库记录
			global.App.DB.Model(&models.InsightDASRecords{}).
				Where("request_id=? and username=?", RequestID, username).
				Updates(map[string]any{"is_finish": true})
			response.Success(c, returnData, "success")
		}
	} else {
		// 有效性验证不通过，更新数据库记录
		global.App.DB.Model(&models.InsightDASRecords{}).
			Where("request_id=? and username=?", RequestID, username).
			Updates(map[string]any{"error_msg": err.Error(), "is_finish": true})
		response.ValidateFail(c, err.Error())
	}
}
