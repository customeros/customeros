package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	commonRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/service"
)

func InitIntegrationRoutes(r *gin.Engine, ctx context.Context, commonRepositoryContainer *commonRepository.Repositories, services *service.Services) {
	r.GET("/integrations",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.SETTINGS_API),
		func(c *gin.Context) {
			tenantName := c.Keys["TenantName"].(string)

			tenantIntegrationSettings, activeServices, err := services.TenantSettingsService.GetForTenant(tenantName)

			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(tenantIntegrationSettings, activeServices))
		})

	r.POST("/integration",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.SETTINGS_API),
		func(c *gin.Context) {
			var request map[string]interface{}

			if err := c.BindJSON(&request); err != nil {
				println(err.Error())
				c.AbortWithStatus(500) //todo
				return
			}

			tenantName := c.Keys["TenantName"].(string)

			tenantIntegrationSettings, activeServices, err := services.TenantSettingsService.SaveIntegrationData(tenantName, request)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(tenantIntegrationSettings, activeServices))
		})

	r.DELETE("/integration/:identifier",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, commonRepositoryContainer),
		commonService.ApiKeyCheckerHTTP(commonRepositoryContainer.AppKeyRepository, commonService.SETTINGS_API),
		func(c *gin.Context) {
			identifier := c.Param("identifier")
			if identifier == "" {
				c.JSON(500, gin.H{"error": "integration identifier is empty"})
				return
			}
			tenantName := c.Keys["TenantName"].(string)

			data, activeServices, err := services.TenantSettingsService.ClearIntegrationData(tenantName, identifier)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapTenantSettingsEntityToDTO(data, activeServices))
		})
}
