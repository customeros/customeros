package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

var contractStatusByModel = map[model.ContractStatus]entity.ContractStatus{
	model.ContractStatusUndefined: entity.ContractStatusUndefined,
	model.ContractStatusDraft:     entity.ContractStatusDraft,
	model.ContractStatusLive:      entity.ContractStatusLive,
	model.ContractStatusEnded:     entity.ContractStatusEnded,
}

var contractStatusByValue = utils.ReverseMap(contractStatusByModel)

func MapContractStatusFromModel(input model.ContractStatus) entity.ContractStatus {
	return contractStatusByModel[input]
}

func MapContractStatusToModel(input entity.ContractStatus) model.ContractStatus {
	return contractStatusByValue[input]
}
