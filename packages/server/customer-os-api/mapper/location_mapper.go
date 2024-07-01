package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapLocationUpdateInputToEntity(input *model.LocationUpdateInput) *neo4jentity.LocationEntity {
	if input == nil {
		return nil
	}
	return &neo4jentity.LocationEntity{
		Id:            input.ID,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
		Name:          utils.IfNotNilString(input.Name),
		RawAddress:    utils.IfNotNilString(input.RawAddress),
		Country:       utils.IfNotNilString(input.Country),
		Region:        utils.IfNotNilString(input.Region),
		Locality:      utils.IfNotNilString(input.Locality),
		Address:       utils.IfNotNilString(input.Address),
		Address2:      utils.IfNotNilString(input.Address2),
		Zip:           utils.IfNotNilString(input.Zip),
		AddressType:   utils.IfNotNilString(input.AddressType),
		HouseNumber:   utils.IfNotNilString(input.HouseNumber),
		PostalCode:    utils.IfNotNilString(input.PostalCode),
		PlusFour:      utils.IfNotNilString(input.PlusFour),
		Commercial:    utils.IfNotNilBool(input.Commercial),
		Predirection:  utils.IfNotNilString(input.Predirection),
		District:      utils.IfNotNilString(input.District),
		Street:        utils.IfNotNilString(input.Street),
		TimeZone:      utils.IfNotNilString(input.TimeZone),
		UtcOffset:     input.UtcOffset,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
	}
}

func MapEntityToLocation(entity *neo4jentity.LocationEntity) *model.Location {
	return &model.Location{
		ID:            entity.Id,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Name:          utils.StringPtr(entity.Name),
		RawAddress:    utils.StringPtr(entity.RawAddress),
		Country:       utils.StringPtr(entity.Country),
		CountryCodeA2: utils.StringPtr(entity.CountryCodeA2),
		CountryCodeA3: utils.StringPtr(entity.CountryCodeA3),
		Region:        utils.StringPtr(entity.Region),
		Locality:      utils.StringPtr(entity.Locality),
		Address:       utils.StringPtr(entity.Address),
		Address2:      utils.StringPtr(entity.Address2),
		Zip:           utils.StringPtr(entity.Zip),
		AddressType:   utils.StringPtr(entity.AddressType),
		HouseNumber:   utils.StringPtr(entity.HouseNumber),
		PostalCode:    utils.StringPtr(entity.PostalCode),
		PlusFour:      utils.StringPtr(entity.PlusFour),
		Commercial:    utils.BoolPtr(entity.Commercial),
		Predirection:  utils.StringPtr(entity.Predirection),
		District:      utils.StringPtr(entity.District),
		Street:        utils.StringPtr(entity.Street),
		Latitude:      entity.Latitude,
		Longitude:     entity.Longitude,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
		TimeZone:      utils.StringPtr(entity.TimeZone),
		UtcOffset:     entity.UtcOffset,
	}
}

func MapEntitiesToLocations(entities *neo4jentity.LocationEntities) []*model.Location {
	var locations []*model.Location
	for _, locationEntity := range *entities {
		locations = append(locations, MapEntityToLocation(&locationEntity))
	}
	return locations
}
