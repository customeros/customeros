package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type IndustryService interface {
	GetAllForOrganizationIds(ctx context.Context, organizationIds []string) (*entity.IndustryEntities, error)
	GetInUseIndustries(ctx context.Context) (*entity.IndustryEntities, error)
	GetByCode(ctx context.Context, code string) (*entity.IndustryEntity, error)
	GetClosestByCode(ctx context.Context, code string) (*entity.IndustryEntity, error)
}
