package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToAction(entity *entity.ActionEntity) *model.Action {
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

var actionTypeByValue = map[entity.ActionType]model.ActionType{
	entity.ActionCreated:                                   model.ActionTypeCreated,
	entity.ActionRenewalForecastUpdated:                    model.ActionTypeRenewalForecastUpdated,
	entity.ActionRenewalLikelihoodUpdated:                  model.ActionTypeRenewalLikelihoodUpdated,
	entity.ActionContractStatusUpdated:                     model.ActionTypeContractStatusUpdated,
	entity.ActionServiceLineItemPriceUpdated:               model.ActionTypeServiceLineItemPriceUpdated,
	entity.ActionServiceLineItemQuantityUpdated:            model.ActionTypeServiceLineItemQuantityUpdated,
	entity.ActionServiceLineItemBilledTypeUpdated:          model.ActionTypeServiceLineItemBilledTypeUpdated,
	entity.ActionServiceLineItemBilledTypeRecurringCreated: model.ActionTypeServiceLineItemBilledTypeRecurringCreated,
	entity.ActionServiceLineItemBilledTypeOnceCreated:      model.ActionTypeServiceLineItemBilledTypeOnceCreated,
	entity.ActionServiceLineItemBilledTypeUsageCreated:     model.ActionTypeServiceLineItemBilledTypeUsageCreated,
	entity.ActionContractRenewed:                           model.ActionTypeContractRenewed,
	entity.ActionServiceLineItemRemoved:                    model.ActionTypeServiceLineItemRemoved,
	entity.ActionOnboardingStatusChanged:                   model.ActionTypeOnboardingStatusChanged,
}

func MapActionTypeToModel(input entity.ActionType) model.ActionType {
	return actionTypeByValue[input]
}
