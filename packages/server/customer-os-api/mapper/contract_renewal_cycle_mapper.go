package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var contractRenewalCycleByModel = map[model.ContractRenewalCycle]neo4jenum.RenewalCycle{
	model.ContractRenewalCycleNone:             neo4jenum.RenewalCycleNone,
	model.ContractRenewalCycleMonthlyRenewal:   neo4jenum.RenewalCycleMonthlyRenewal,
	model.ContractRenewalCycleQuarterlyRenewal: neo4jenum.RenewalCycleQuarterlyRenewal,
	model.ContractRenewalCycleAnnualRenewal:    neo4jenum.RenewalCycleAnnualRenewal,
}

var contractRenewalCycleByValue = utils.ReverseMap(contractRenewalCycleByModel)

func MapContractRenewalCycleFromModel(input model.ContractRenewalCycle) neo4jenum.RenewalCycle {
	return contractRenewalCycleByModel[input]
}

func MapContractRenewalCycleToModel(input neo4jenum.RenewalCycle) model.ContractRenewalCycle {
	return contractRenewalCycleByValue[input]
}
