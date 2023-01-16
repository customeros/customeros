package service

import (
	"github.com/gin-gonic/gin"
	repository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
)

type App string

const (
	CUSTOMER_OS_API  App = "customer-os-api"
	FILE_STORAGE_API App = "file-storage-api"
	SETTINGS_API     App = "settings-api"
)

func ApiKeyChecker(appKeyRepo repository.AppKeyRepository, app App) func(c *gin.Context) {
	return func(c *gin.Context) {
		kh := c.GetHeader("X-Openline-API-KEY")
		if kh != "" {

			keyResult := appKeyRepo.FindByKey(c, string(app), kh)

			if keyResult.Error != nil {
				c.AbortWithStatus(401)
				return
			}

			appKey := keyResult.Result.(*entity.AppKey)

			if appKey == nil {
				c.AbortWithStatus(401)
				return
			}

			c.Next()
			// illegal request, terminate the current process
		} else {
			c.AbortWithStatus(401)
			return
		}

	}
}
