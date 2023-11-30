package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var opportunityRenewalLikelihoodByModel = map[model.OpportunityRenewalLikelihood]entity.OpportunityRenewalLikelihood{
	model.OpportunityRenewalLikelihoodHighRenewal:   entity.OpportunityRenewalLikelihoodHigh,
	model.OpportunityRenewalLikelihoodMediumRenewal: entity.OpportunityRenewalLikelihoodMedium,
	model.OpportunityRenewalLikelihoodLowRenewal:    entity.OpportunityRenewalLikelihoodLow,
	model.OpportunityRenewalLikelihoodZeroRenewal:   entity.OpportunityRenewalLikelihoodZero,
}

var opportunityRenewalLikelihoodByValue = utils.ReverseMap(opportunityRenewalLikelihoodByModel)

func MapOpportunityRenewalLikelihoodFromModel(input model.OpportunityRenewalLikelihood) entity.OpportunityRenewalLikelihood {
	return opportunityRenewalLikelihoodByModel[input]
}

func MapOpportunityRenewalLikelihoodToModel(input entity.OpportunityRenewalLikelihood) model.OpportunityRenewalLikelihood {
	return opportunityRenewalLikelihoodByValue[input]
}

func MapOpportunityRenewalLikelihoodToModelPtr(input string) *model.OpportunityRenewalLikelihood {
	switch input {
	case string(entity.OpportunityRenewalLikelihoodHigh):
		return utils.Ptr(model.OpportunityRenewalLikelihoodHighRenewal)
	case string(entity.OpportunityRenewalLikelihoodMedium):
		return utils.Ptr(model.OpportunityRenewalLikelihoodMediumRenewal)
	case string(entity.OpportunityRenewalLikelihoodLow):
		return utils.Ptr(model.OpportunityRenewalLikelihoodLowRenewal)
	case string(entity.OpportunityRenewalLikelihoodZero):
		return utils.Ptr(model.OpportunityRenewalLikelihoodZeroRenewal)
	default:
		return nil
	}
}
