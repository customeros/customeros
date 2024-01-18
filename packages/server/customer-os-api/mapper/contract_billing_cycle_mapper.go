package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var contractBillingCycleByModel = map[model.ContractBillingCycle]neo4jenum.BillingCycle{
	model.ContractBillingCycleNone:             neo4jenum.BillingCycleNone,
	model.ContractBillingCycleMonthlyBilling:   neo4jenum.BillingCycleMonthlyBilling,
	model.ContractBillingCycleQuarterlyBilling: neo4jenum.BillingCycleQuarterlyBilling,
	model.ContractBillingCycleAnnualBilling:    neo4jenum.BillingCycleAnnualBilling,
}

var contractBillingCycleByValue = utils.ReverseMap(contractBillingCycleByModel)

func MapContractBillingCycleFromModel(input model.ContractBillingCycle) neo4jenum.BillingCycle {
	return contractBillingCycleByModel[input]
}

func MapContractBillingCycleToModel(input neo4jenum.BillingCycle) model.ContractBillingCycle {
	return contractBillingCycleByValue[input]
}
