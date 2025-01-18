package interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type TenantService interface {
	GetAllTenants(ctx context.Context) ([]*neo4j_entity.TenantEntity, error)
	GetTenantForUserEmail(ctx context.Context, email string) (*neo4j_entity.TenantEntity, error)
	Merge(ctx context.Context, tenantEntity neo4j_entity.TenantEntity) (*neo4j_entity.TenantEntity, error)
	HardDelete(ctx context.Context, tenant string) error
}
