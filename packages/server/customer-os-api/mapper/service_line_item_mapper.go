package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToServiceLineItem(entity *entity.ServiceLineItemEntity) *model.ServiceLineItem {
	if entity == nil {
		return nil
	}
	return &model.ServiceLineItem{
		ID:            entity.ID,
		Name:          entity.Name,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		Source:        MapDataSourceToModel(entity.Source),
		SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:     entity.AppSource,
		Billed:        MapBilledTypeToModel(entity.Billed),
		Price:         entity.Price,
		Quantity:      entity.Quantity,
	}
}

func MapServiceLineItemInputToEntity(input model.ServiceLineItemInput) *entity.ServiceLineItemEntity {
	serviceLineItemEntity := entity.ServiceLineItemEntity{
		Name:          utils.IfNotNilString(input.Name),
		Price:         utils.IfNotNilFloat64(input.Price),
		Quantity:      utils.IfNotNilInt64(input.Quantity),
		Source:        entity.DataSourceOpenline,
		SourceOfTruth: entity.DataSourceOpenline,
		AppSource:     utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	if input.Billed != nil {
		billedType := MapBilledTypeFromModel(*input.Billed)
		serviceLineItemEntity.Billed = billedType
	} else {
		billedType := entity.BilledTypeNone
		serviceLineItemEntity.Billed = billedType
	}
	return &serviceLineItemEntity
}

func MapEntitiesToServiceLineItems(entities *entity.ServiceLineItemEntities) []*model.ServiceLineItem {
	var ServiceLineItems []*model.ServiceLineItem
	for _, ServiceLineItemEntity := range *entities {
		ServiceLineItems = append(ServiceLineItems, MapEntityToServiceLineItem(&ServiceLineItemEntity))
	}
	return ServiceLineItems
}

func MapServiceLineItemUpdateInputToEntity(input model.ServiceLineItemUpdateInput) *entity.ServiceLineItemEntity {
	serviceLineItemEntity := entity.ServiceLineItemEntity{
		ID:            input.ServiceLineItemID,
		Name:          utils.IfNotNilString(input.Name),
		Price:         utils.IfNotNilFloat64(input.Price),
		Quantity:      utils.IfNotNilInt64(input.Quantity),
		Source:        entity.DataSourceOpenline,
		SourceOfTruth: entity.DataSourceOpenline,
		AppSource:     utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	if input.Billed != nil {
		serviceLineItemEntity.Billed = MapBilledTypeFromModel(*input.Billed)
	} else {
		serviceLineItemEntity.Billed = entity.BilledTypeMonthly
	}
	return &serviceLineItemEntity
}
