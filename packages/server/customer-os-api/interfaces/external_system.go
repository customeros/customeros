package cosapi_interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type ExternalSystemService interface {
	GetAllExternalSystemInstances(ctx context.Context) (*entity.ExternalSystemEntities, error)
}
