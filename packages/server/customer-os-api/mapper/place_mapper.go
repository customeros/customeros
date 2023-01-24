package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapEntityToAddress(entity *entity.PlaceEntity) *model.Place {
	address := model.Place{
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
		CreatedAt: entity.CreatedAt,
	}
	return &address
}

func MapEntitiesToPlaces(entities *entity.PlaceEntities) []*model.Place {
	var addresses []*model.Place
	for _, addressEntity := range *entities {
		addresses = append(addresses, MapEntityToAddress(&addressEntity))
	}
	return addresses
}
