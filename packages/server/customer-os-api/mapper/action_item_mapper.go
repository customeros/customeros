package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapEntityToActionItem(entity *entity.ActionItemEntity) *model.ActionItem {
	return &model.ActionItem{
		ID:        entity.Id,
		CreatedAt: *entity.CreatedAt,
		Content:   entity.Content,

		Source:    MapDataSourceToModel(entity.Source),
		AppSource: entity.AppSource,
	}
}

func MapEntitiesToActionItem(entities *entity.ActionItemEntities) []*model.ActionItem {
	var mappedEntities []*model.ActionItem
	for _, entity := range *entities {
		mappedEntities = append(mappedEntities, MapEntityToActionItem(&entity))
	}
	return mappedEntities
}
