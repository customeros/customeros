package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type IndustryService interface {
	GetAllForOrganizationIds(ctx context.Context, organizationIds []string) (*entity.IndustryEntities, error)
	GetByCode(ctx context.Context, code string) (*entity.IndustryEntity, error)
}
