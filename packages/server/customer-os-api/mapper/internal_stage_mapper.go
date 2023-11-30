package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var internalStageByModel = map[model.InternalStage]entity.InternalStage{
	model.InternalStageOpen:       entity.InternalStageOpen,
	model.InternalStageClosedLost: entity.InternalStageClosedLost,
	model.InternalStageClosedWon:  entity.InternalStageClosedWon,
	model.InternalStageEvaluating: entity.InternalStageEvaluating,
}

var internalStageByValue = utils.ReverseMap(internalStageByModel)

func MapInternalStageFromModel(input model.InternalStage) entity.InternalStage {
	return internalStageByModel[input]
}

func MapInternalStageToModel(input entity.InternalStage) model.InternalStage {
	return internalStageByValue[input]
}
