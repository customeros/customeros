package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapEntityToLocation(entity *entity.LocationEntity) *model.Location {
	location := model.Location{
		ID:        entity.Id,
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		Source:    utils.ToPtr(MapDataSourceToModel(entity.Source)),
		AppSource: utils.StringPtr(entity.AppSource),
	}
	return &location
}

func MapEntitiesToLocations(entities *entity.LocationEntities) []*model.Location {
	var locations []*model.Location
	for _, locationEntity := range *entities {
		locations = append(locations, MapEntityToLocation(&locationEntity))
	}
	return locations
}
