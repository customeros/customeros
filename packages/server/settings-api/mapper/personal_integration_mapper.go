package mapper

import (
	commonPgEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
)

func MapPersonalIntegrationToDTO(personalIntegration *commonPgEntity.PersonalIntegration) *map[string]interface{} {
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
