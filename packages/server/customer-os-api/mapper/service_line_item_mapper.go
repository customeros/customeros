package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToServiceLineItem(entity *neo4jentity.ServiceLineItemEntity) *model.ServiceLineItem {
	if entity == nil {
		return nil
	}
	return &model.ServiceLineItem{
		ID:            entity.ID,
		Name:          entity.Name,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		StartedAt:     entity.StartedAt,
		EndedAt:       entity.EndedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
		Billed:        MapBilledTypeToModel(entity.Billed),
		Price:         entity.Price,
		Quantity:      entity.Quantity,
		Comments:      entity.Comments,
		ParentID:      entity.ParentID,
		VatRate:       entity.VatRate,
	}
}

func MapEntitiesToServiceLineItems(entities *neo4jentity.ServiceLineItemEntities) []*model.ServiceLineItem {
	var ServiceLineItems []*model.ServiceLineItem
	for _, ServiceLineItemEntity := range *entities {
		ServiceLineItems = append(ServiceLineItems, MapEntityToServiceLineItem(&ServiceLineItemEntity))
	}
	return ServiceLineItems
}
