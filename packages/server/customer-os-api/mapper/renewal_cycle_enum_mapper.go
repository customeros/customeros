package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapRenewalCycleFromModel(input *model.RenewalCycle) string {
	if input == nil {
		return ""
	}
	return input.String()
}

func MapRenewalCycleToModel(input string) *model.RenewalCycle {
	if input == "" {
		return nil
	}
	v := model.RenewalCycle(input)
	if v.IsValid() {
		return &v
	} else {
		return nil
	}
}
