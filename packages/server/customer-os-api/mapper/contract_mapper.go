package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

func MapEntityToContract(entity *neo4jentity.ContractEntity) *model.Contract {
	if entity == nil {
		return nil
	}
	return &model.Contract{
		ID:               entity.Id,
		Name:             entity.Name,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
		Source:           MapDataSourceToModel(entity.Source),
		SourceOfTruth:    MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:        entity.AppSource,
		Status:           MapContractStatusToModel(entity.ContractStatus),
		RenewalCycle:     MapContractRenewalCycleToModel(entity.RenewalCycle),
		RenewalPeriods:   entity.RenewalPeriods,
		ServiceStartedAt: entity.ServiceStartedAt,
		SignedAt:         entity.SignedAt,
		EndedAt:          entity.EndedAt,
		ContractURL:      utils.StringPtrNillable(entity.ContractUrl),
	}
}

func MapContractInputToEntity(input model.ContractInput) *neo4jentity.ContractEntity {
	contractEntity := neo4jentity.ContractEntity{
		Name:             utils.IfNotNilString(input.Name),
		ContractUrl:      utils.IfNotNilString(input.ContractURL),
		SignedAt:         input.SignedAt,
		ServiceStartedAt: input.ServiceStartedAt,
		Source:           neo4jentity.DataSourceOpenline,
		SourceOfTruth:    neo4jentity.DataSourceOpenline,
		AppSource:        utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
		RenewalPeriods:   input.RenewalPeriods,
	}
	if input.RenewalCycle != nil {
		contractRenewalCycle := MapContractRenewalCycleFromModel(*input.RenewalCycle)
		contractEntity.RenewalCycle = contractRenewalCycle
	} else {
		contractRenewalCycle := neo4jenum.RenewalCycleNone
		contractEntity.RenewalCycle = contractRenewalCycle
	}
	return &contractEntity
}

func MapContractUpdateInputToEntity(input model.ContractUpdateInput) *neo4jentity.ContractEntity {
	contractEntity := neo4jentity.ContractEntity{
		Id:               input.ContractID,
		Name:             utils.IfNotNilString(input.Name),
		ContractUrl:      utils.IfNotNilString(input.ContractURL),
		ServiceStartedAt: input.ServiceStartedAt,
		SignedAt:         input.SignedAt,
		EndedAt:          input.EndedAt,
		Source:           neo4jentity.DataSourceOpenline,
		SourceOfTruth:    neo4jentity.DataSourceOpenline,
		AppSource:        utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
		RenewalPeriods:   input.RenewalPeriods,
	}
	if input.RenewalCycle != nil {
		contractEntity.RenewalCycle = MapContractRenewalCycleFromModel(*input.RenewalCycle)
	} else {
		contractEntity.RenewalCycle = neo4jenum.RenewalCycleNone
	}
	return &contractEntity
}

func MapEntitiesToContracts(entities *neo4jentity.ContractEntities) []*model.Contract {
	var contracts []*model.Contract
	for _, contractEntity := range *entities {
		contracts = append(contracts, MapEntityToContract(&contractEntity))
	}
	return contracts
}
