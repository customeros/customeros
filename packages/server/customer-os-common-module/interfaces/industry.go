package interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type IndustryService interface {
	GetAllForOrganizationIds(ctx context.Context, organizationIds []string) (*neo4j_entity.IndustryEntities, error)
	GetInUseIndustries(ctx context.Context) (*neo4j_entity.IndustryEntities, error)
	GetByCode(ctx context.Context, code string) (*neo4j_entity.IndustryEntity, error)
	GetClosestByCode(ctx context.Context, code string) (*neo4j_entity.IndustryEntity, error)
}
