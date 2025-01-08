package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
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

var actionTypeByValue = map[enum.ActionType]model.ActionType{
	enum.ActionCreated:                                   model.ActionTypeCreated,
	enum.ActionGeneric:                                   model.ActionTypeGeneric,
	enum.ActionRenewalForecastUpdated:                    model.ActionTypeRenewalForecastUpdated,
	enum.ActionRenewalLikelihoodUpdated:                  model.ActionTypeRenewalLikelihoodUpdated,
	enum.ActionContractStatusUpdated:                     model.ActionTypeContractStatusUpdated,
	enum.ActionServiceLineItemPriceUpdated:               model.ActionTypeServiceLineItemPriceUpdated,
	enum.ActionServiceLineItemQuantityUpdated:            model.ActionTypeServiceLineItemQuantityUpdated,
	enum.ActionServiceLineItemBilledTypeUpdated:          model.ActionTypeServiceLineItemBilledTypeUpdated,
	enum.ActionServiceLineItemBilledTypeRecurringCreated: model.ActionTypeServiceLineItemBilledTypeRecurringCreated,
	enum.ActionServiceLineItemBilledTypeOnceCreated:      model.ActionTypeServiceLineItemBilledTypeOnceCreated,
	enum.ActionServiceLineItemBilledTypeUsageCreated:     model.ActionTypeServiceLineItemBilledTypeUsageCreated,
	enum.ActionContractRenewed:                           model.ActionTypeContractRenewed,
	enum.ActionServiceLineItemRemoved:                    model.ActionTypeServiceLineItemRemoved,
	enum.ActionOnboardingStatusChanged:                   model.ActionTypeOnboardingStatusChanged,
	enum.ActionInvoiceIssued:                             model.ActionTypeInvoiceIssued,
	enum.ActionInvoicePaid:                               model.ActionTypeInvoicePaid,
	enum.ActionInvoiceVoided:                             model.ActionTypeInvoiceVoided,
	enum.ActionInvoiceSent:                               model.ActionTypeInvoiceSent,
	enum.ActionInvoiceOverdue:                            model.ActionTypeInvoiceOverdue,
	enum.ActionInteractionEventRead:                      model.ActionTypeInteractionEventRead,
}

func MapActionTypeToModel(input enum.ActionType) model.ActionType {
	return actionTypeByValue[input]
}
