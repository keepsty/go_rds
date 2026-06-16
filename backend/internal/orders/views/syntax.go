package views

import (
	"github.com/keepsty/go_rds/pkg/response"

	"github.com/keepsty/go_rds/internal/orders/forms"
	"github.com/keepsty/go_rds/internal/orders/services"

	"github.com/gin-gonic/gin"
)

// 语法审核
func InspectOrderSyntaxView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form *forms.InspectOrderSyntaxForm = &forms.InspectOrderSyntaxForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.InspectOrderSyntaxService{
			InspectOrderSyntaxForm: form,
			C:                      c,
			Username:               username,
		}
		returnData, err := service.Run()
		if err != nil {
			response.Fail(c, err.Error())
		} else {
			response.Success(c, returnData, "success")
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}
