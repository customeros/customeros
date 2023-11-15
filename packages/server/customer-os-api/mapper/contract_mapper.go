package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapEntityToContract(entity *entity.ContractEntity) *model.Contract {
	if entity == nil {
		return nil
	}
	return &model.Contract{
		ID:               entity.ID,
		Name:             entity.Name,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
		Source:           MapDataSourceToModel(entity.Source),
		SourceOfTruth:    MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:        entity.AppSource,
		Status:           MapContractStatusToModel(entity.ContractStatus),
		RenewalCycle:     MapContractRenewalCycleToModel(entity.ContractRenewalCycle),
		ServiceStartedAt: entity.ServiceStartedAt,
		SignedAt:         entity.SignedAt,
		EndedAt:          entity.EndedAt,
		ContractURL:      utils.StringPtrNillable(entity.ContractUrl),
	}
}

func MapContractInputToEntity(input model.ContractInput) *entity.ContractEntity {
	contractEntity := entity.ContractEntity{
		Name:             utils.IfNotNilString(input.Name),
		ContractUrl:      utils.IfNotNilString(input.ContractURL),
		SignedAt:         input.SignedAt,
		ServiceStartedAt: input.ServiceStartedAt,
		Source:           entity.DataSourceOpenline,
		SourceOfTruth:    entity.DataSourceOpenline,
		AppSource:        utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
	}
	if input.RenewalCycle != nil {
		contractRenewalCycle := MapContractRenewalCycleFromModel(*input.RenewalCycle)
		contractEntity.ContractRenewalCycle = contractRenewalCycle
	} else {
		contractRenewalCycle := entity.ContractRenewalCycleNone
		contractEntity.ContractRenewalCycle = contractRenewalCycle
	}
	return &contractEntity
}

func MapEntitiesToContracts(entities *entity.ContractEntities) []*model.Contract {
	var contracts []*model.Contract
	for _, contractEntity := range *entities {
		contracts = append(contracts, MapEntityToContract(&contractEntity))
	}
	return contracts
}
