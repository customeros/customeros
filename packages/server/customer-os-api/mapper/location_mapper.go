package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapEntityToLocation(entity *entity.LocationEntity) *model.Location {
	location := model.Location{
		ID:           entity.Id,
		Name:         entity.Name,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
		Country:      utils.StringPtr(entity.Country),
		Region:       utils.StringPtr(entity.Region),
		Locality:     utils.StringPtr(entity.Locality),
		Address:      utils.StringPtr(entity.Address),
		Address2:     utils.StringPtr(entity.Address2),
		Zip:          utils.StringPtr(entity.Zip),
		AddressType:  utils.StringPtr(entity.AddressType),
		HouseNumber:  utils.StringPtr(entity.HouseNumber),
		PostalCode:   utils.StringPtr(entity.PostalCode),
		PlusFour:     utils.StringPtr(entity.PlusFour),
		Commercial:   utils.BoolPtr(entity.Commercial),
		Predirection: utils.StringPtr(entity.Predirection),
		District:     utils.StringPtr(entity.District),
		Street:       utils.StringPtr(entity.Street),
		RawAddress:   utils.StringPtr(entity.RawAddress),
		Latitude:     entity.Latitude,
		Longitude:    entity.Longitude,
		Source:       utils.ToPtr(MapDataSourceToModel(entity.Source)),
		AppSource:    utils.StringPtr(entity.AppSource),
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
