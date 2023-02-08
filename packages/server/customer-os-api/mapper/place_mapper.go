package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapEntityToPlace(entity *entity.PlaceEntity) *model.Place {
	place := model.Place{
		ID:        entity.Id,
		Country:   utils.StringPtr(entity.Country),
		State:     utils.StringPtr(entity.State),
		City:      utils.StringPtr(entity.City),
		Address:   utils.StringPtr(entity.Address),
		Address2:  utils.StringPtr(entity.Address2),
		Zip:       utils.StringPtr(entity.Zip),
		Phone:     utils.StringPtr(entity.Phone),
		Fax:       utils.StringPtr(entity.Fax),
		Source:    utils.ToPtr(MapDataSourceToModel(entity.Source)),
		AppSource: utils.StringPtr(entity.AppSource),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
	return &place
}

func MapEntitiesToPlaces(entities *entity.PlaceEntities) []*model.Place {
	var places []*model.Place
	for _, placeEntity := range *entities {
		places = append(places, MapEntityToPlace(&placeEntity))
	}
	return places
}
