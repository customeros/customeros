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
		ID:                    entity.Id,
		Name:                  entity.Name,
		CreatedAt:             entity.CreatedAt,
		UpdatedAt:             entity.UpdatedAt,
		Source:                MapDataSourceToModel(entity.Source),
		SourceOfTruth:         MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:             entity.AppSource,
		Status:                MapContractStatusToModel(entity.ContractStatus),
		RenewalCycle:          MapContractRenewalCycleToModel(entity.RenewalCycle),
		RenewalPeriods:        entity.RenewalPeriods,
		ServiceStartedAt:      entity.ServiceStartedAt,
		SignedAt:              entity.SignedAt,
		EndedAt:               entity.EndedAt,
		ContractURL:           utils.StringPtrNillable(entity.ContractUrl),
		InvoicingStartDate:    entity.InvoicingStartDate,
		Currency:              utils.ToPtr(MapCurrencyToModel(entity.Currency)),
		BillingCycle:          utils.ToPtr(MapContractBillingCycleToModel(entity.BillingCycle)),
		AddressLine1:          &entity.AddressLine1,
		AddressLine2:          &entity.AddressLine2,
		Zip:                   &entity.Zip,
		Country:               &entity.Country,
		Locality:              &entity.Locality,
		OrganizationLegalName: &entity.OrganizationLegalName,
		InvoiceEmail:          &entity.InvoiceEmail,
	}
}

func MapContractInputToEntity(input model.ContractInput) *neo4jentity.ContractEntity {
	contractEntity := neo4jentity.ContractEntity{
		Name:               utils.IfNotNilString(input.Name),
		ContractUrl:        utils.IfNotNilString(input.ContractURL),
		SignedAt:           input.SignedAt,
		ServiceStartedAt:   input.ServiceStartedAt,
		InvoicingStartDate: input.InvoicingStartDate,
		Source:             neo4jentity.DataSourceOpenline,
		SourceOfTruth:      neo4jentity.DataSourceOpenline,
		AppSource:          utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
		RenewalPeriods:     input.RenewalPeriods,
	}

	if input.RenewalCycle != nil {
		contractEntity.RenewalCycle = MapContractRenewalCycleFromModel(*input.RenewalCycle)
	} else {
		contractEntity.RenewalCycle = neo4jenum.RenewalCycleNone
	}

	if input.Currency != nil {
		contractEntity.Currency = MapCurrencyFromModel(*input.Currency)
	}

	if input.BillingCycle != nil {
		contractEntity.BillingCycle = MapContractBillingCycleFromModel(*input.BillingCycle)
	} else {
		contractEntity.BillingCycle = neo4jenum.BillingCycleNone
	}

	return &contractEntity
}

func MapContractUpdateInputToEntity(input model.ContractUpdateInput) *neo4jentity.ContractEntity {
	contractEntity := neo4jentity.ContractEntity{
		Id:                 input.ContractID,
		Name:               utils.IfNotNilString(input.Name),
		ContractUrl:        utils.IfNotNilString(input.ContractURL),
		ServiceStartedAt:   input.ServiceStartedAt,
		SignedAt:           input.SignedAt,
		EndedAt:            input.EndedAt,
		InvoicingStartDate: input.InvoicingStartDate,
		Source:             neo4jentity.DataSourceOpenline,
		SourceOfTruth:      neo4jentity.DataSourceOpenline,
		AppSource:          utils.IfNotNilStringWithDefault(input.AppSource, constants.AppSourceCustomerOsApi),
		RenewalPeriods:     input.RenewalPeriods,
	}
	if input.RenewalCycle != nil {
		contractEntity.RenewalCycle = MapContractRenewalCycleFromModel(*input.RenewalCycle)
	} else {
		contractEntity.RenewalCycle = neo4jenum.RenewalCycleNone
	}

	if input.Currency != nil {
		contractEntity.Currency = MapCurrencyFromModel(*input.Currency)
	}

	if input.BillingCycle != nil {
		contractEntity.BillingCycle = MapContractBillingCycleFromModel(*input.BillingCycle)
	} else {
		contractEntity.BillingCycle = neo4jenum.BillingCycleNone
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
