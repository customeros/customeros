package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var opportunityRenewalLikelihoodByModel = map[model.OpportunityRenewalLikelihood]neo4jenum.RenewalLikelihood{
	model.OpportunityRenewalLikelihoodHighRenewal:   neo4jenum.RenewalLikelihoodHigh,
	model.OpportunityRenewalLikelihoodMediumRenewal: neo4jenum.RenewalLikelihoodMedium,
	model.OpportunityRenewalLikelihoodLowRenewal:    neo4jenum.RenewalLikelihoodLow,
	model.OpportunityRenewalLikelihoodZeroRenewal:   neo4jenum.RenewalLikelihoodZero,
}

var opportunityRenewalLikelihoodByValue = utils.ReverseMap(opportunityRenewalLikelihoodByModel)

func MapOpportunityRenewalLikelihoodFromModel(input *model.OpportunityRenewalLikelihood) neo4jenum.RenewalLikelihood {
	if input == nil {
		return ""
	}
	return opportunityRenewalLikelihoodByModel[*input]
}

func MapOpportunityRenewalLikelihoodToModel(input neo4jenum.RenewalLikelihood) model.OpportunityRenewalLikelihood {
	return opportunityRenewalLikelihoodByValue[input]
}

func MapOpportunityRenewalLikelihoodFromString(input *string) string {
	if input == nil {
		return ""
	}
	if v, exists := opportunityRenewalLikelihoodByModel[model.OpportunityRenewalLikelihood(*input)]; exists {
		return string(v)
	} else {
		return ""
	}
}

func MapOpportunityRenewalLikelihoodToModelPtr(input string) *model.OpportunityRenewalLikelihood {
	switch input {
	case string(neo4jenum.RenewalLikelihoodHigh):
		return utils.Ptr(model.OpportunityRenewalLikelihoodHighRenewal)
	case string(neo4jenum.RenewalLikelihoodMedium):
		return utils.Ptr(model.OpportunityRenewalLikelihoodMediumRenewal)
	case string(neo4jenum.RenewalLikelihoodLow):
		return utils.Ptr(model.OpportunityRenewalLikelihoodLowRenewal)
	case string(neo4jenum.RenewalLikelihoodZero):
		return utils.Ptr(model.OpportunityRenewalLikelihoodZeroRenewal)
	default:
		return nil
	}
}
