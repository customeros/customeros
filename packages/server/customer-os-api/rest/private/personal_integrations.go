package private

import (
	"github.com/gin-gonic/gin"
	cosapi_services "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

func CreatePersonalIntegrations(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request map[string]interface{}

		if err := c.BindJSON(&request); err != nil {
			println(err.Error())
			c.AbortWithStatus(500) // todo
			return
		}

		tenantName := c.Keys["TenantName"].(string)
		userEmail := c.Keys["UserEmail"].(string)
		integration := postgresEntity.PersonalIntegration{
			Name:       request["name"].(string),
			TenantName: tenantName,
			Email:      userEmail,
			Secret:     request["secret"].(string),
			Active:     true,
		}
		saved, err := s.PersonalIntegrationsService.SavePersonalIntegration(integration)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, mapPersonalIntegrationToDTO(saved))
	}
}

func GetPersonalIntegrations(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantName := c.Keys["TenantName"].(string)
		userMail := c.Keys["UserEmail"].(string)
		integrationName := c.Param("integrationName")
		if integrationName == "" {
			c.JSON(500, gin.H{"error": "integration name is empty"})
			return
		}
		integrations, err := s.PersonalIntegrationsService.GetPersonalIntegrations(tenantName, userMail)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var integrationsDTO []map[string]interface{}
		for _, integration := range integrations {
			integrationsDTO = append(integrationsDTO, *mapPersonalIntegrationToDTO(integration))
		}
		c.JSON(200, integrationsDTO)
	}
}

func GetPersonalIntegrationByName(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantName := c.Keys["TenantName"].(string)
		userMail := c.Keys["UserEmail"].(string)
		integrationName := c.Param("integrationName")
		if integrationName == "" {
			c.JSON(500, gin.H{"error": "integration name is empty"})
			return
		}
		integration, err := s.PersonalIntegrationsService.GetPersonalIntegration(tenantName, userMail, integrationName)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, mapPersonalIntegrationToDTO(integration))
	}
}

func mapPersonalIntegrationToDTO(personalIntegration *postgresEntity.PersonalIntegration) *map[string]interface{} {
	return &map[string]interface{}{
		"id":         personalIntegration.ID,
		"tenantName": personalIntegration.TenantName,
		"name":       personalIntegration.Name,
		"email":      personalIntegration.Email,
		"secret":     personalIntegration.Secret,
		"createdAt":  personalIntegration.CreatedAt,
		"updatedAt":  personalIntegration.UpdatedAt,
	}
}
