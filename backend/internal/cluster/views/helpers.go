package views

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/keepsty/go_rds/pkg/response"
)

func getUsername(c *gin.Context) (string, bool) {
	username, ok := GetUserNameFromJWT(c)
	if !ok {
		response.Fail(c, "认证信息无效")
		return "", false
	}
	return username, true
}

func parseUintParam(c *gin.Context, name string) (uint64, bool) {
	raw := c.Param(name)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		response.ValidateFail(c, "非法参数: "+name)
		return 0, false
	}
	return id, true
}

func parseInt64Param(c *gin.Context, name string) (int64, bool) {
	raw := c.Param(name)
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id < 1 {
		response.ValidateFail(c, "非法参数: "+name)
		return 0, false
	}
	return id, true
}
