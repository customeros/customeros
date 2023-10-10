package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/service"
)

func InitUserSettingsRoutes(r *gin.Engine, ctx context.Context, services *service.Services) {
	r.GET("/user/settings/google/:playerIdentityId",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, services.Repositories.CommonRepositories),
		commonService.ApiKeyCheckerHTTP(services.Repositories.CommonRepositories.AppKeyRepository, commonService.SETTINGS_API),

		func(c *gin.Context) {
			playerIdentityId := c.Param("playerIdentityId")
			userSettings, err := services.OAuthUserSettingsService.GetOAuthUserSettings(playerIdentityId)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, userSettings)
		})

	r.GET("/user/settings/slack",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, services.Repositories.CommonRepositories),
		commonService.ApiKeyCheckerHTTP(services.Repositories.CommonRepositories.AppKeyRepository, commonService.SETTINGS_API),

		func(c *gin.Context) {
			tenant, _ := c.Get(commonService.KEY_TENANT_NAME)
			userSettings, err := services.SlackSettingsService.GetSlackSettings(tenant.(string))
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, userSettings)
		})
}
