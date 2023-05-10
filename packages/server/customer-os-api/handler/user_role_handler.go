package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func GetUserRoleHandlerEnhancer() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Set("Role", model.RoleUser)
		c.Next()
	}
}
