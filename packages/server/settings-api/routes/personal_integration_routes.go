package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	commonPgEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/service"
)

func InitPersonalIntegrationRoutes(r *gin.Engine, ctx context.Context, services *service.Services) {
	r.GET("/personal_integrations/:integrationName",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, services.Repositories.CommonRepositories),
		commonService.ApiKeyCheckerHTTP(services.Repositories.CommonRepositories.AppKeyRepository, commonService.SETTINGS_API),
		func(c *gin.Context) {
			tenantName := c.Keys["TenantName"].(string)
			userMail := c.Keys["UserEmail"].(string)
			integrationName := c.Param("integrationName")
			if integrationName == "" {
				c.JSON(500, gin.H{"error": "integration name is empty"})
				return
			}
			integration, err := services.PersonalIntegrationsService.GetPersonalIntegration(tenantName, userMail, integrationName)

			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapPersonalIntegrationToDTO(integration))
		})

	r.GET("/personal_integrations/",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, services.Repositories.CommonRepositories),
		commonService.ApiKeyCheckerHTTP(services.Repositories.CommonRepositories.AppKeyRepository, commonService.SETTINGS_API),
		func(c *gin.Context) {
			tenantName := c.Keys["TenantName"].(string)
			userMail := c.Keys["UserEmail"].(string)
			integrationName := c.Param("integrationName")
			if integrationName == "" {
				c.JSON(500, gin.H{"error": "integration name is empty"})
				return
			}
			integrations, err := services.PersonalIntegrationsService.GetPersonalIntegrations(tenantName, userMail)

			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			var integrationsDTO []map[string]interface{}
			for _, integration := range integrations {
				integrationsDTO = append(integrationsDTO, *mapper.MapPersonalIntegrationToDTO(integration))
			}
			c.JSON(200, integrationsDTO)
		})

	r.POST("/personal_integrations",
		commonService.TenantUserContextEnhancer(ctx, commonService.USERNAME, services.Repositories.CommonRepositories),
		commonService.ApiKeyCheckerHTTP(services.Repositories.CommonRepositories.AppKeyRepository, commonService.SETTINGS_API),
		func(c *gin.Context) {
			var request map[string]interface{}

			if err := c.BindJSON(&request); err != nil {
				println(err.Error())
				c.AbortWithStatus(500) //todo
				return
			}

			tenantName := c.Keys["TenantName"].(string)
			userEmail := c.Keys["UserEmail"].(string)
			integration := commonPgEntity.PersonalIntegration{
				Name:       request["name"].(string),
				TenantName: tenantName,
				Email:      userEmail,
				Secret:     request["secret"].(string),
			}
			saved, err := services.PersonalIntegrationsService.SavePersonalIntegration(integration)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, mapper.MapPersonalIntegrationToDTO(saved))
		})
}
