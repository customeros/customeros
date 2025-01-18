package cosapi_interfaces

import postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

type TenantSettingsService interface {
	GetForTenant(tenantName string) (*postgres_entity.TenantSettings, map[string]bool, error)
	SaveIntegrationData(tenantName string, request map[string]interface{}) (*postgres_entity.TenantSettings, map[string]bool, error)
	ClearIntegrationData(tenantName, identifier string) (*postgres_entity.TenantSettings, map[string]bool, error)
}
