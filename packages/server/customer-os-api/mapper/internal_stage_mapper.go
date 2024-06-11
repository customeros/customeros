package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var internalStageByModel = map[model.InternalStage]neo4jenum.OpportunityInternalStage{
	model.InternalStageOpen:       neo4jenum.OpportunityInternalStageOpen,
	model.InternalStageClosedLost: neo4jenum.OpportunityInternalStageClosedLost,
	model.InternalStageClosedWon:  neo4jenum.OpportunityInternalStageClosedWon,
}

var internalStageByValue = utils.ReverseMap(internalStageByModel)

func MapInternalStageFromModel(input model.InternalStage) neo4jenum.OpportunityInternalStage {
	return internalStageByModel[input]
}

func MapInternalStageToModel(input neo4jenum.OpportunityInternalStage) model.InternalStage {
	return internalStageByValue[input]
}
