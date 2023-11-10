package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var probabilityByModel = map[model.RenewalLikelihoodProbability]string{
	model.RenewalLikelihoodProbabilityHigh:   string(entity.RenewalLikelihoodProbabilityHigh),
	model.RenewalLikelihoodProbabilityMedium: string(entity.RenewalLikelihoodProbabilityMedium),
	model.RenewalLikelihoodProbabilityLow:    string(entity.RenewalLikelihoodProbabilityLow),
	model.RenewalLikelihoodProbabilityZero:   string(entity.RenewalLikelihoodProbabilityZero),
}

var probabilityByValue = utils.ReverseMap(probabilityByModel)

func MapRenewalLikelihoodFromModel(input *model.RenewalLikelihoodProbability) string {
	if input == nil {
		return ""
	}
	if v, exists := probabilityByModel[*input]; exists {
		return v
	} else {
		return ""
	}
}

func MapRenewalLikelihoodFromString(input *string) string {
	if input == nil {
		return ""
	}
	if v, exists := probabilityByModel[model.RenewalLikelihoodProbability(*input)]; exists {
		return v
	} else {
		return ""
	}
}

func MapRenewalLikelihoodToModel(input string) *model.RenewalLikelihoodProbability {

	if v, exists := probabilityByValue[input]; exists {
		return &v
	} else {
		return nil
	}
}
