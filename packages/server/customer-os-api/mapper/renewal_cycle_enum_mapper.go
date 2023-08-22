package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	organization_grpc_service "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/organization"
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

func MapFrequencyFromModelToGrpc(input *model.RenewalCycle) *organization_grpc_service.Frequency {
	if input == nil {
		return nil
	}
	switch *input {
	case model.RenewalCycleWeekly:
		return utils.ToPtr(organization_grpc_service.Frequency_WEEKLY)
	case model.RenewalCycleBiweekly:
		return utils.ToPtr(organization_grpc_service.Frequency_BIWEEKLY)
	case model.RenewalCycleMonthly:
		return utils.ToPtr(organization_grpc_service.Frequency_MONTHLY)
	case model.RenewalCycleQuarterly:
		return utils.ToPtr(organization_grpc_service.Frequency_QUARTERLY)
	case model.RenewalCycleBiannually:
		return utils.ToPtr(organization_grpc_service.Frequency_BIANNUALLY)
	case model.RenewalCycleAnnually:
		return utils.ToPtr(organization_grpc_service.Frequency_ANNUALLY)
	}
	return nil
}
