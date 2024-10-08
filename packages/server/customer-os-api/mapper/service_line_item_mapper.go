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
		Metadata: &model.Metadata{
			ID:            entity.ID,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        MapDataSourceToModel(entity.Source),
			SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
			AppSource:     entity.AppSource,
		},
		BillingCycle:   MapBilledTypeToModel(entity.Billed),
		Comments:       entity.Comments,
		Description:    entity.Name,
		ParentID:       entity.ParentID,
		Price:          entity.Price,
		Quantity:       entity.Quantity,
		ServiceEnded:   entity.EndedAt,
		ServiceStarted: entity.StartedAt,
		Tax: &model.Tax{
			TaxRate: entity.VatRate,
		},
		Closed: entity.Canceled,
		Paused: entity.Paused,
	}
}

func MapEntitiesToServiceLineItems(entities *neo4jentity.ServiceLineItemEntities) []*model.ServiceLineItem {
	var ServiceLineItems []*model.ServiceLineItem
	for _, ServiceLineItemEntity := range *entities {
		ServiceLineItems = append(ServiceLineItems, MapEntityToServiceLineItem(&ServiceLineItemEntity))
	}
	return ServiceLineItems
}
