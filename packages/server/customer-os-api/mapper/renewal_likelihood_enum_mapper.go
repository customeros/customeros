package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapRenewalLikelihoodFromModel(input *model.RenewalLikelihoodProbability) string {
	if input == nil {
		return ""
	}
	return input.String()
}

func MapRenewalLikelihoodToModel(input string) *model.RenewalLikelihoodProbability {
	if input == "" {
		return nil
	}
	v := model.RenewalLikelihoodProbability(input)
	if v.IsValid() {
		return &v
	} else {
		return nil
	}
}
