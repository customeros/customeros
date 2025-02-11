package cosapi_interfaces

import (
	"context"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type TenantSettingsService interface {
	GetForTenant(ctx context.Context) (*postgres_entity.TenantSettings, map[string]bool, error)
	SaveIntegrationData(ctx context.Context, request map[string]interface{}) (*postgres_entity.TenantSettings, map[string]bool, error)
	ClearIntegrationData(ctx context.Context, identifier string) (*postgres_entity.TenantSettings, map[string]bool, error)
}
