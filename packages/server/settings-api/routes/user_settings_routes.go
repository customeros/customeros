package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	commonUtils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/service"
)

func InitUserSettingsRoutes(r *gin.Engine, services *service.Services) {
	r.GET("/user/settings/oauth/:tenant",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API, security.WithCache(services.CommonServices.Cache)),

		func(c *gin.Context) {
			contextWithTimeout, cancel := commonUtils.GetLongLivedContext(context.Background())
			defer cancel()

			tenant := c.Param("tenant")

			userSettings, err := services.OAuthUserSettingsService.GetTenantOAuthUserSettings(contextWithTimeout, tenant)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, userSettings)
		})

	r.GET("/user/settings/slack",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API, security.WithCache(services.CommonServices.Cache)),

		func(c *gin.Context) {
			tenant, _ := c.Get(security.KEY_TENANT_NAME)
			userSettings, err := services.SlackSettingsService.GetSlackSettings(c, tenant.(string))
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, userSettings)
		})
}
