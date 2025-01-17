package private

import (
	"context"

	"github.com/gin-gonic/gin"
	cosapi_services "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/security"
	commonUtils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func GetOAuthSettings(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextWithTimeout, cancel := commonUtils.GetLongLivedContext(context.Background())
		defer cancel()

		tenant := c.Param("tenant")

		userSettings, err := s.OAuthUserSettingsService.GetTenantOAuthUserSettings(contextWithTimeout, tenant)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, userSettings)
	}
}

func GetSlackSettings(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant, _ := c.Get(security.KEY_TENANT_NAME)
		userSettings, err := s.SlackService.GetSlackSettings(c, tenant.(string))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, userSettings)
	}
}
