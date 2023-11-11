package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var contractRenewalCycleByModel = map[model.ContractRenewalCycle]entity.ContractRenewalCycle{
	model.ContractRenewalCycleNone:           entity.ContractRenewalCycleNone,
	model.ContractRenewalCycleMonthlyRenewal: entity.ContractRenewalCycleMonthlyRenewal,
	model.ContractRenewalCycleAnnualRenewal:  entity.ContractRenewalCycleAnnualRenewal,
}

var contractRenewalCycleByValue = utils.ReverseMap(contractRenewalCycleByModel)

func MapContractRenewalCycleFromModel(input model.ContractRenewalCycle) entity.ContractRenewalCycle {
	return contractRenewalCycleByModel[input]
}

func MapContractRenewalCycleToModel(input entity.ContractRenewalCycle) model.ContractRenewalCycle {
	return contractRenewalCycleByValue[input]
}
