package views

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "userId"

// GetUserNameFromJWT 从 JWT claims 中提取用户名
func GetUserNameFromJWT(c *gin.Context) (string, bool) {
	claims := jwt.ExtractClaims(c)
	raw, ok := claims["id"]
	if !ok {
		return "", false
	}
	username, ok := raw.(string)
	return username, ok
}
