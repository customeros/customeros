package mapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
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
