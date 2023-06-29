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
	}
}

var actionTypeByModel = map[model.ActionType]entity.ActionType{
	model.ActionTypeCreated: entity.ActionCreated,
}

var actionTypeByValue = utils.ReverseMap(actionTypeByModel)

func MapActionTypeToModel(input entity.ActionType) model.ActionType {
	return actionTypeByValue[input]
}
