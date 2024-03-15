package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapEntityToOrder(entity *entity.OrderEntity) *model.Order {
	return &model.Order{
		ID:            entity.Id,
		CreatedAt:     entity.CreatedAt,
		ConfirmedAt:   entity.ConfirmedAt,
		PaidAt:        entity.PaidAt,
		FulfilledAt:   entity.FulfilledAt,
		CancelledAt:   entity.CancelledAt,
		Source:        MapDataSourceToModel(entity.SourceFields.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceFields.SourceOfTruth),
		AppSource:     entity.SourceFields.AppSource,
	}
}

func MapEntitiesToOrders(entities *entity.OrderEntities) []*model.Order {
	var socials []*model.Order
	for _, entity := range *entities {
		socials = append(socials, MapEntityToOrder(&entity))
	}
	return socials
}
