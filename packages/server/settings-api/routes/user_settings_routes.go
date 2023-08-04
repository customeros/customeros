package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/service"
)

func InitUserSettingsRoutes(r *gin.Engine, ctx context.Context, commonRepositoryContainer *commonRepository.Repositories, services *service.Services) {
	r.GET("/user/settings",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.SETTINGS_API),

		func(c *gin.Context) {
			username := c.Keys["Username"].(string)

			userSettings, err := services.UserSettingsService.GetByUserName(username)

			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, userSettings)
		})

	r.POST("/user/settings",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.SETTINGS_API),
		func(c *gin.Context) {
			var request model.UserSettings

			if err := c.BindJSON(&request); err != nil {
				println(err.Error())
				c.AbortWithStatus(500) //todo
				return
			}
			request.TenantName = c.Keys["TenantName"].(string)
			services.UserSettingsService.Save(&request)
			c.JSON(200, request)
		})

	r.DELETE("/user/settings/:identifier",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.SETTINGS_API),
		func(c *gin.Context) {
			//identifier := c.Param("identifier")
			//if identifier == "" {
			//	c.JSON(500, gin.H{"error": "integration identifier is empty"})
			//	return
			//}
			//tenantName := c.Keys["TenantName"].(string)
			//
			//data, activeServices, err := services.TenantSettingsService.ClearIntegrationData(tenantName, identifier)
			//if err != nil {
			//	c.JSON(500, gin.H{"error": err.Error()})
			//	return
			//}
			//
			//c.JSON(200, mapper.MapTenantSettingsEntityToDTO(data, activeServices))
		})
}
