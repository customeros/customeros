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
}

var stageByValue = utils.ReverseMap(stageByModel)

func MapStageFromModel(input model.OrganizationStage) commonenum.OrganizationStage {
	return stageByModel[input]
}

func MapStageToModel(input commonenum.OrganizationStage) model.OrganizationStage {
	return stageByValue[input]
}
