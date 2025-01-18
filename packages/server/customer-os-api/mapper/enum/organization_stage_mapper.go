package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
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
	model.OrganizationStageTrial:          neo4jenum.Trial,
}

var stageByValue = utils.ReverseMap(stageByModel)

func MapStageFromModel(input model.OrganizationStage) neo4jenum.OrganizationStage {
	return stageByModel[input]
}

func MapStageToModel(input neo4jenum.OrganizationStage) model.OrganizationStage {
	return stageByValue[input]
}
