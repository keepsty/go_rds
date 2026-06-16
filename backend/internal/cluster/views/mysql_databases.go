package views

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/keepsty/go_rds/internal/cluster/forms"
	"github.com/keepsty/go_rds/internal/cluster/models"
	"github.com/keepsty/go_rds/internal/cluster/services"
	"github.com/keepsty/go_rds/internal/global"
	"github.com/keepsty/go_rds/pkg/response"
)

// GetDBsByClusterNameHandler 根据用户名获取数据库选项
func GetDBsByClusterNameHandler(c *gin.Context) {
	username := c.Param("username")
	data, err := services.GetClusterDBByUsername(username)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, data, "success")
}

// GetDBsHandler 获取数据库列表
func GetDBsHandler(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	data, err := services.GetClusterProxyDetailBySGID(id)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, data, "success")
}

// GetTbsByClusterNameHandler 根据SGID和DBID获取表列表（树形结构）
func GetTbsByClusterNameHandler(c *gin.Context) {
	sgID, ok := parseInt64Param(c, "sg_id")
	if !ok {
		return
	}
	dbID, ok := parseInt64Param(c, "db_id")
	if !ok {
		return
	}
	data, err := services.GetClusterTableBySGDBID(sgID, dbID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, data, "success")
}

// ClusterDBReadQueryExecuteHandler 执行SQL查询
func ClusterDBReadQueryExecuteHandler(c *gin.Context) {
	readQueryData := new(forms.DBReadQueryExecute)
	if err := c.ShouldBindJSON(readQueryData); err != nil {
		if _, ok := err.(validator.ValidationErrors); !ok {
			response.ValidateFail(c, "参数错误")
			return
		}
		response.ValidateFail(c, err.Error())
		return
	}
	if readQueryData.DBId == 0 {
		response.ValidateFail(c, "db_id不能为空")
		return
	}
	username, ok := getUsername(c)
	if !ok {
		return
	}
	readQueryData.Username = username

	dasConfig := global.App.Config.Das
	data, err := services.ClusterDBReadQueryExecuteHandler((*models.DBReadQueryExecute)(readQueryData), &dasConfig)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, data, "success")
}

// ClusterGetDBDataDictHandler 获取数据字典
func ClusterGetDBDataDictHandler(c *gin.Context) {
	dbData := new(forms.DBReadQueryExecute)
	if err := c.ShouldBindJSON(dbData); err != nil {
		if _, ok := err.(validator.ValidationErrors); !ok {
			response.ValidateFail(c, "参数错误")
			return
		}
		response.ValidateFail(c, err.Error())
		return
	}
	if dbData.DBId == 0 {
		response.ValidateFail(c, "db_id不能为空")
		return
	}
	data, dbName, err := services.ClusterDBDataDictHandler((*models.DBReadQueryExecute)(dbData))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	res := new(models.DataDictResponse)
	res.Data = data
	res.DBName = dbName
	response.Success(c, res, "success")
}

// ClusterGetDBTableInfoHandler 获取表结构/基础信息
func ClusterGetDBTableInfoHandler(c *gin.Context) {
	dbData := new(forms.RequestGetDBTableInfo)
	if err := c.ShouldBind(dbData); err != nil {
		if _, ok := err.(validator.ValidationErrors); !ok {
			response.ValidateFail(c, "参数错误")
			return
		}
		response.ValidateFail(c, err.Error())
		return
	}
	if dbData.DBID == 0 {
		response.ValidateFail(c, "db_id不能为空")
		return
	}
	data, err := services.ClusterGetDBTableInfoHandler((*models.RequestGetDBTableInfo)(dbData))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, data, "success")
}

// GetUserHistorySqlHandler 获取用户查询历史
func GetUserHistorySqlHandler(c *gin.Context) {
	his := new(models.RequestGetUserHistoryQueryHandler)
	if err := c.ShouldBind(his); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	// 从 JWT 获取用户名
	username, ok := getUsername(c)
	if !ok {
		return
	}
	his.Username = username

	data, err := services.GetUsernameHistoryQuery(his)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, data, "success")
}
