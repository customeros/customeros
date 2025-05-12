package enummapper

import (
	commonenum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

var stageByModel = map[model.OrganizationStage]commonenum.OrganizationStage{
	model.OrganizationStageTarget:      commonenum.Target,
	model.OrganizationStageEducation:   commonenum.Education,
	model.OrganizationStageSolution:    commonenum.Solution,
	model.OrganizationStageEvaluation:  commonenum.Evaluation,
	model.OrganizationStageReadyToBuy:  commonenum.ReadyToBuy,
	model.OrganizationStageOpportunity: commonenum.Opportunity,
	model.OrganizationStageCustomer:    commonenum.Customer,
	model.OrganizationStageNotAFit:     commonenum.NotAFit,

	model.OrganizationStageLead:           commonenum.Lead,
	model.OrganizationStageEngaged:        commonenum.Engaged,
	model.OrganizationStageUnqualified:    commonenum.Unqualified,
	model.OrganizationStageOnboarding:     commonenum.Onboarding,
	model.OrganizationStageInitialValue:   commonenum.InitialValue,
	model.OrganizationStageRecurringValue: commonenum.RecurringValue,
	model.OrganizationStageMaxValue:       commonenum.MaxValue,
	model.OrganizationStagePendingChurn:   commonenum.PendingChurn,
	model.OrganizationStageTrial:          commonenum.Trial,
}

var stageByValue = utils.ReverseMap(stageByModel)

func MapStageFromModel(input model.OrganizationStage) commonenum.OrganizationStage {
	return stageByModel[input]
}

func MapStageToModel(input commonenum.OrganizationStage) model.OrganizationStage {
	return stageByValue[input]
}
