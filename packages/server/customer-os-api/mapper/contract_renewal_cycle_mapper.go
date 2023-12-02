package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var contractRenewalCycleByModel = map[model.ContractRenewalCycle]entity.RenewalCycle{
	model.ContractRenewalCycleNone:             entity.RenewalCycleNone,
	model.ContractRenewalCycleMonthlyRenewal:   entity.RenewalCycleMonthlyRenewal,
	model.ContractRenewalCycleQuarterlyRenewal: entity.RenewalCycleQuarterlyRenewal,
	model.ContractRenewalCycleAnnualRenewal:    entity.RenewalCycleAnnualRenewal,
}

var contractRenewalCycleByValue = utils.ReverseMap(contractRenewalCycleByModel)

func MapContractRenewalCycleFromModel(input model.ContractRenewalCycle) entity.RenewalCycle {
	return contractRenewalCycleByModel[input]
}

func MapContractRenewalCycleToModel(input entity.RenewalCycle) model.ContractRenewalCycle {
	return contractRenewalCycleByValue[input]
}
