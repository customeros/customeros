package mapper

import (
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

func MapPersonalIntegrationToDTO(personalIntegration *postgresEntity.PersonalIntegration) *map[string]interface{} {
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
