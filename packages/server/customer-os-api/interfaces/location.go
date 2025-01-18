package cosapi_interfaces

import (
	"context"

	commonModel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	cosapiEntity "github.com/customeros/customeros/packages/server/customer-os-api/entity"
)

type LocationService interface {
	CreateLocationForEntity(ctx context.Context, entityType commonModel.EntityType, entityId string, source cosapiEntity.SourceFields) (*entity.LocationEntity, error)
	Update(ctx context.Context, entity entity.LocationEntity) (*entity.LocationEntity, error)
	DetachFromEntity(ctx context.Context, entityType commonModel.EntityType, entityId, locationId string) error
}
