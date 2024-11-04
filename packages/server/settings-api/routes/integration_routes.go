package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/service"
)

func InitIntegrationRoutes(r *gin.Engine, services *service.Services) {
	r.GET("/integrations",
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API, security.WithCache(services.CommonServices.Cache)),
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
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API, security.WithCache(services.CommonServices.Cache)),
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
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API, security.WithCache(services.CommonServices.Cache)),
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
