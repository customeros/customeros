package cosapi_interfaces

import (
	"context"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
)

type ExternalSystemService interface {
	GetAllExternalSystemInstances(ctx context.Context) (*neo4j_entity.ExternalSystemEntities, error)
}
