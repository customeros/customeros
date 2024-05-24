package enummapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var stageByModel = map[model.OrganizationStage]neo4jenum.OrganizationStage{
	model.OrganizationStageLead:           neo4jenum.Lead,
	model.OrganizationStageTarget:         neo4jenum.Target,
	model.OrganizationStageEngaged:        neo4jenum.Engaged,
	model.OrganizationStageUnqualified:    neo4jenum.Unqualified,
	model.OrganizationStageReadyToBuy:     neo4jenum.ReadyToBuy,
	model.OrganizationStageOnboarding:     neo4jenum.Onboarding,
	model.OrganizationStageInitialValue:   neo4jenum.InitialValue,
	model.OrganizationStageRecurringValue: neo4jenum.RecurringValue,
	model.OrganizationStageMaxValue:       neo4jenum.MaxValue,
	model.OrganizationStagePendingChurn:   neo4jenum.PendingChurn,
}

var stageByValue = utils.ReverseMap(stageByModel)

func MapStageFromModel(input model.OrganizationStage) neo4jenum.OrganizationStage {
	return stageByModel[input]
}

func MapStageToModel(input neo4jenum.OrganizationStage) model.OrganizationStage {
	return stageByValue[input]
}
