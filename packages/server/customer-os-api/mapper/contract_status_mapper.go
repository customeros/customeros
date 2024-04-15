package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var contractStatusByModel = map[model.ContractStatus]neo4jenum.ContractStatus{
	model.ContractStatusUndefined:     neo4jenum.ContractStatusUndefined,
	model.ContractStatusDraft:         neo4jenum.ContractStatusDraft,
	model.ContractStatusLive:          neo4jenum.ContractStatusLive,
	model.ContractStatusEnded:         neo4jenum.ContractStatusEnded,
	model.ContractStatusOutOfContract: neo4jenum.ContractStatusOutOfContract,
	model.ContractStatusScheduled:     neo4jenum.ContractStatusScheduled,
}

var contractStatusByValue = utils.ReverseMap(contractStatusByModel)

func MapContractStatusFromModel(input model.ContractStatus) neo4jenum.ContractStatus {
	return contractStatusByModel[input]
}

func MapContractStatusToModel(input neo4jenum.ContractStatus) model.ContractStatus {
	return contractStatusByValue[input]
}
