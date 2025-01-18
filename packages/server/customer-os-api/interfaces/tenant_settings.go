package cosapi_interfaces

import "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

type TenantSettingsService interface {
	GetForTenant(tenantName string) (*entity.TenantSettings, map[string]bool, error)
	SaveIntegrationData(tenantName string, request map[string]interface{}) (*entity.TenantSettings, map[string]bool, error)
	ClearIntegrationData(tenantName, identifier string) (*entity.TenantSettings, map[string]bool, error)
}
