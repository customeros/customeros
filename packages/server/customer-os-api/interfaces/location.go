package cosapi_interfaces

import (
	"context"

	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	cosapiEntity "github.com/customeros/customeros/packages/server/customer-os-api/entity"
)

type LocationService interface {
	CreateLocationForEntity(ctx context.Context, entityType commonModel.EntityType, entityId string, source cosapiEntity.SourceFields) (*neo4j_entity.LocationEntity, error)
	Update(ctx context.Context, entity neo4j_entity.LocationEntity) (*neo4j_entity.LocationEntity, error)
	DetachFromEntity(ctx context.Context, entityType commonModel.EntityType, entityId, locationId string) error
}
