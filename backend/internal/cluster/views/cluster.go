package views

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/keepsty/go_rds/internal/cluster/forms"
	"github.com/keepsty/go_rds/internal/cluster/models"
	"github.com/keepsty/go_rds/internal/cluster/services"
	"github.com/keepsty/go_rds/pkg/response"
)

// ClustersHandler 获取集群列表
func ClustersHandler(c *gin.Context) {
	cl := forms.ParamClusterList{
		Page: 1,
		Size: 10,
	}
	if err := c.ShouldBindQuery(&cl); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	data, err := services.GetClusters(&cl)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	resData := new(models.ServiceGroupResponse)
	resData.Total = len(data)
	resData.Page = cl.Page
	resData.Clusters = data
	response.Success(c, resData, "success")
}

// GetClusterByIDHandler 根据ID获取集群详情
func GetClusterByIDHandler(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	data, err := services.GetClusterByID(id)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, data, "success")
}

// GetClusterInsByIDHandler 根据ID获取集群实例
func GetClusterInsByIDHandler(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	data, err := services.GetClusterInsByID(id)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, data, "success")
}

// GetClusterDBsByIDHandler 根据ID获取集群数据库列表
func GetClusterDBsByIDHandler(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	data, err := services.GetClusterDBSByID(id)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, data, "success")
}

// GetClusterDBDetailByIDHandler 根据ID获取数据库表详情
func GetClusterDBDetailByIDHandler(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	data, err := services.GetClusterDBDetailByID(id)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, data, "success")
}

// GetClusterProxyListByIDHandler 根据ID获取ProxySQL列表
func GetClusterProxyListByIDHandler(c *gin.Context) {
	id, ok := parseInt64Param(c, "id")
	if !ok {
		return
	}
	data, err := services.GetClusterProxyListByID(id)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, data, "success")
}

// SetInstanceStatusByIDHandler 设置实例状态（上线/下线）
func SetInstanceStatusByIDHandler(c *gin.Context) {
	ins := new(forms.InstanceDetail)
	if err := c.ShouldBindJSON(ins); err != nil {
		if _, ok := err.(validator.ValidationErrors); !ok {
			response.ValidateFail(c, "参数错误")
			return
		}
		response.ValidateFail(c, err.Error())
		return
	}
	data, err := services.SetInstanceStatusByID(ins.ID, ins.Status)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, data, "success")
}
