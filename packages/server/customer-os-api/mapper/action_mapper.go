package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

func MapEntityToAction(entity *neo4jentity.ActionEntity) *model.Action {
	if entity == nil {
		return nil
	}
	return &model.Action{
		ID:         entity.Id,
		CreatedAt:  entity.CreatedAt,
		ActionType: MapActionTypeToModel(entity.Type),
		AppSource:  entity.AppSource,
		Source:     MapDataSourceToModel(entity.Source),
		Content:    utils.StringPtrNillable(entity.Content),
		Metadata:   utils.StringPtrNillable(entity.Metadata),
	}
}

func MapEntitiesToAction(entities *neo4jentity.ActionEntities) []*model.Action {
	var mappedEntities []*model.Action
	for _, entity := range *entities {
		mappedEntities = append(mappedEntities, MapEntityToAction(&entity))
	}
	return mappedEntities
}

var actionTypeByValue = map[neo4jenum.ActionType]model.ActionType{
	neo4jenum.ActionCreated:                                   model.ActionTypeCreated,
	neo4jenum.ActionRenewalForecastUpdated:                    model.ActionTypeRenewalForecastUpdated,
	neo4jenum.ActionRenewalLikelihoodUpdated:                  model.ActionTypeRenewalLikelihoodUpdated,
	neo4jenum.ActionContractStatusUpdated:                     model.ActionTypeContractStatusUpdated,
	neo4jenum.ActionServiceLineItemPriceUpdated:               model.ActionTypeServiceLineItemPriceUpdated,
	neo4jenum.ActionServiceLineItemQuantityUpdated:            model.ActionTypeServiceLineItemQuantityUpdated,
	neo4jenum.ActionServiceLineItemBilledTypeUpdated:          model.ActionTypeServiceLineItemBilledTypeUpdated,
	neo4jenum.ActionServiceLineItemBilledTypeRecurringCreated: model.ActionTypeServiceLineItemBilledTypeRecurringCreated,
	neo4jenum.ActionServiceLineItemBilledTypeOnceCreated:      model.ActionTypeServiceLineItemBilledTypeOnceCreated,
	neo4jenum.ActionServiceLineItemBilledTypeUsageCreated:     model.ActionTypeServiceLineItemBilledTypeUsageCreated,
	neo4jenum.ActionContractRenewed:                           model.ActionTypeContractRenewed,
	neo4jenum.ActionServiceLineItemRemoved:                    model.ActionTypeServiceLineItemRemoved,
	neo4jenum.ActionOnboardingStatusChanged:                   model.ActionTypeOnboardingStatusChanged,
	neo4jenum.ActionInvoiceIssued:                             model.ActionTypeInvoiceIssued,
	neo4jenum.ActionInvoicePaid:                               model.ActionTypeInvoicePaid,
	neo4jenum.ActionInvoiceVoided:                             model.ActionTypeInvoiceVoided,
	neo4jenum.ActionInvoiceSent:                               model.ActionTypeInvoiceSent,
	neo4jenum.ActionInvoiceOverdue:                            model.ActionTypeInvoiceOverdue,
	neo4jenum.ActionInteractionEventRead:                      model.ActionTypeInteractionEventRead,
}

func MapActionTypeToModel(input neo4jenum.ActionType) model.ActionType {
	return actionTypeByValue[input]
}
