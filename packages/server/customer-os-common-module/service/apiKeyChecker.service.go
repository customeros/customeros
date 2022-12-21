package service

import (
	"github.com/gin-gonic/gin"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
)

func ApiKeyChecker(appKeyRepo repository.AppKeyRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		kh := c.GetHeader("X-Openline-API-KEY")
		if kh != "" {

			keyResult := appKeyRepo.FindByKey(c, kh)

			if keyResult.Error != nil {
				c.AbortWithStatus(401)
				return
			}

			appKey := keyResult.Result.(*entity.AppKeyEntity)

			if appKey == nil {
				c.AbortWithStatus(401)
				return
			} else {
				// todo set tenant in context
			}

			c.Next()
			// illegal request, terminate the current process
		} else {
			c.AbortWithStatus(401)
			return
		}

	}
}
