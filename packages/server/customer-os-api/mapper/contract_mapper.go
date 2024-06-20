package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	mapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityToContract(entity *neo4jentity.ContractEntity) *model.Contract {
	if entity == nil {
		return nil
	}
	contract := model.Contract{
		Metadata: &model.Metadata{
			ID:            entity.Id,
			Created:       entity.CreatedAt,
			LastUpdated:   entity.UpdatedAt,
			Source:        MapDataSourceToModel(entity.Source),
			SourceOfTruth: MapDataSourceToModel(entity.SourceOfTruth),
			AppSource:     entity.AppSource,
			Version:       entity.AggregateVersion,
		},
		BillingDetails: &model.BillingDetails{
			BillingCycleInMonths:   utils.ToPtr(entity.BillingCycleInMonths),
			InvoicingStarted:       entity.InvoicingStartDate,
			NextInvoicing:          entity.NextInvoiceDate,
			AddressLine1:           utils.ToPtr(entity.AddressLine1),
			AddressLine2:           utils.ToPtr(entity.AddressLine2),
			Locality:               utils.ToPtr(entity.Locality),
			Region:                 utils.ToPtr(entity.Region),
			Country:                utils.ToPtr(entity.Country),
			PostalCode:             utils.ToPtr(entity.Zip),
			OrganizationLegalName:  utils.ToPtr(entity.OrganizationLegalName),
			BillingEmail:           utils.ToPtr(entity.InvoiceEmail),
			BillingEmailCc:         entity.InvoiceEmailCC,
			BillingEmailBcc:        entity.InvoiceEmailBCC,
			InvoiceNote:            utils.ToPtr(entity.InvoiceNote),
			CanPayWithCard:         utils.ToPtr(entity.CanPayWithCard),
			CanPayWithDirectDebit:  utils.ToPtr(entity.CanPayWithDirectDebit),
			CanPayWithBankTransfer: utils.ToPtr(entity.CanPayWithBankTransfer),
			PayOnline:              utils.ToPtr(entity.PayOnline),
			PayAutomatically:       utils.ToPtr(entity.PayAutomatically),
			Check:                  utils.ToPtr(entity.Check),
			DueDays:                utils.ToPtr(entity.DueDays),
		},
		ContractEnded:           entity.EndedAt,
		ContractName:            entity.Name,
		ContractSigned:          entity.SignedAt,
		ContractURL:             utils.StringPtrNillable(entity.ContractUrl),
		Currency:                utils.ToPtr(mapper.MapCurrencyToModel(entity.Currency)),
		BillingEnabled:          entity.InvoicingEnabled,
		ServiceStarted:          entity.ServiceStartedAt,
		ContractStatus:          MapContractStatusToModel(entity.ContractStatus),
		AutoRenew:               entity.AutoRenew,
		CommittedPeriodInMonths: utils.ToPtr[int64](entity.LengthInMonths),
		Approved:                entity.Approved,
		Ltv:                     entity.Ltv,

		// All below are deprecated
		ID:                    entity.Id,
		Name:                  entity.Name,
		CreatedAt:             entity.CreatedAt,
		UpdatedAt:             entity.UpdatedAt,
		Source:                MapDataSourceToModel(entity.Source),
		SourceOfTruth:         MapDataSourceToModel(entity.SourceOfTruth),
		AppSource:             entity.AppSource,
		Status:                MapContractStatusToModel(entity.ContractStatus),
		ServiceStartedAt:      entity.ServiceStartedAt,
		SignedAt:              entity.SignedAt,
		EndedAt:               entity.EndedAt,
		InvoicingStartDate:    entity.InvoicingStartDate,
		AddressLine1:          utils.ToPtr(entity.AddressLine1),
		AddressLine2:          utils.ToPtr(entity.AddressLine2),
		Zip:                   utils.ToPtr(entity.Zip),
		Country:               utils.ToPtr(entity.Country),
		Locality:              utils.ToPtr(entity.Locality),
		OrganizationLegalName: utils.ToPtr(entity.OrganizationLegalName),
		InvoiceEmail:          utils.ToPtr(entity.InvoiceEmail),
		InvoiceNote:           utils.ToPtr(entity.InvoiceNote),
	}

	if entity.LengthInMonths == int64(0) {
		contract.RenewalPeriods = nil
		contract.CommittedPeriods = nil
		contract.ContractRenewalCycle = model.ContractRenewalCycleNone
		contract.RenewalCycle = model.ContractRenewalCycleNone
	} else if entity.LengthInMonths < int64(3) {
		contract.RenewalPeriods = utils.ToPtr(int64(1))
		contract.CommittedPeriods = utils.ToPtr(int64(1))
		contract.ContractRenewalCycle = model.ContractRenewalCycleMonthlyRenewal
		contract.RenewalCycle = model.ContractRenewalCycleMonthlyRenewal
	} else if entity.LengthInMonths < int64(12) {
		contract.RenewalPeriods = utils.ToPtr(int64(1))
		contract.CommittedPeriods = utils.ToPtr(int64(1))
		contract.ContractRenewalCycle = model.ContractRenewalCycleQuarterlyRenewal
		contract.RenewalCycle = model.ContractRenewalCycleQuarterlyRenewal
	} else {
		contract.RenewalPeriods = utils.ToPtr(entity.LengthInMonths / 12)
		contract.CommittedPeriods = utils.ToPtr(entity.LengthInMonths / 12)
		contract.ContractRenewalCycle = model.ContractRenewalCycleAnnualRenewal
		contract.RenewalCycle = model.ContractRenewalCycleAnnualRenewal
	}
	if entity.BillingCycleInMonths == int64(0) {
		contract.BillingCycle = utils.ToPtr(model.ContractBillingCycleNone)
		contract.BillingDetails.BillingCycle = utils.ToPtr(model.ContractBillingCycleNone)
	} else if entity.BillingCycleInMonths == int64(3) {
		contract.BillingCycle = utils.ToPtr(model.ContractBillingCycleQuarterlyBilling)
		contract.BillingDetails.BillingCycle = utils.ToPtr(model.ContractBillingCycleQuarterlyBilling)
	} else if entity.BillingCycleInMonths == int64(12) {
		contract.BillingCycle = utils.ToPtr(model.ContractBillingCycleAnnualBilling)
		contract.BillingDetails.BillingCycle = utils.ToPtr(model.ContractBillingCycleAnnualBilling)
	} else if entity.BillingCycleInMonths == int64(1) {
		contract.BillingCycle = utils.ToPtr(model.ContractBillingCycleAnnualBilling)
		contract.BillingDetails.BillingCycle = utils.ToPtr(model.ContractBillingCycleMonthlyBilling)
	} else {
		contract.BillingCycle = utils.ToPtr(model.ContractBillingCycleCustomBilling)
		contract.BillingDetails.BillingCycle = utils.ToPtr(model.ContractBillingCycleCustomBilling)
	}

	return &contract
}

func MapEntitiesToContracts(entities *neo4jentity.ContractEntities) []*model.Contract {
	var contracts []*model.Contract
	for _, contractEntity := range *entities {
		contracts = append(contracts, MapEntityToContract(&contractEntity))
	}
	return contracts
}
