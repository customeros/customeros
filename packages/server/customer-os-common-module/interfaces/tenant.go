package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type TenantService interface {
	GetAllTenants(ctx context.Context) ([]*entity.TenantEntity, error)
	GetTenantForUserEmail(ctx context.Context, email string) (*entity.TenantEntity, error)
	Merge(ctx context.Context, tenantEntity entity.TenantEntity) (*entity.TenantEntity, error)
	HardDelete(ctx context.Context, tenant string) error
}
