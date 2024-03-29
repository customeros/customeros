package mapper

import (
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

func MapDbNodeToInvoiceEntity(dbNode *dbtype.Node) *entity.InvoiceEntity {
	if dbNode == nil {
		return &entity.InvoiceEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	invoiceEntity := entity.InvoiceEntity{
		Id:               utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:        utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:        utils.GetTimePropOrEpochStart(props, "updatedAt"),
		DryRun:           utils.GetBoolPropOrFalse(props, "dryRun"),
		OffCycle:         utils.GetBoolPropOrFalse(props, "offCycle"),
		Postpaid:         utils.GetBoolPropOrFalse(props, "postpaid"),
		Number:           utils.GetStringPropOrEmpty(props, "number"),
		PeriodStartDate:  utils.GetTimePropOrEpochStart(props, "periodStartDate"),
		PeriodEndDate:    utils.GetTimePropOrEpochStart(props, "periodEndDate"),
		DueDate:          utils.GetTimePropOrEpochStart(props, "dueDate"),
		Currency:         enum.DecodeCurrency(utils.GetStringPropOrEmpty(props, "currency")),
		BillingCycle:     enum.DecodeBillingCycle(utils.GetStringPropOrEmpty(props, "billingCycle")),
		Amount:           utils.GetFloatPropOrZero(props, "amount"),
		Vat:              utils.GetFloatPropOrZero(props, "vat"),
		TotalAmount:      utils.GetFloatPropOrZero(props, "totalAmount"),
		RepositoryFileId: utils.GetStringPropOrEmpty(props, "repositoryFileId"),
		Source:           entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:    entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:        utils.GetStringPropOrEmpty(props, "appSource"),
		Status:           enum.DecodeInvoiceStatus(utils.GetStringPropOrEmpty(props, "status")),
		Note:             utils.GetStringPropOrEmpty(props, "note"),
		Customer: entity.InvoiceCustomer{
			Name:         utils.GetStringPropOrEmpty(props, "customerName"),
			Email:        utils.GetStringPropOrEmpty(props, "customerEmail"),
			AddressLine1: utils.GetStringPropOrEmpty(props, "customerAddressLine1"),
			AddressLine2: utils.GetStringPropOrEmpty(props, "customerAddressLine2"),
			Zip:          utils.GetStringPropOrEmpty(props, "customerAddressZip"),
			Locality:     utils.GetStringPropOrEmpty(props, "customerAddressLocality"),
			Country:      utils.GetStringPropOrEmpty(props, "customerAddressCountry"),
			Region:       utils.GetStringPropOrEmpty(props, "customerAddressRegion"),
		},
		Provider: entity.InvoiceProvider{
			LogoRepositoryFileId: utils.GetStringPropOrEmpty(props, "providerLogoRepositoryFileId"),
			Name:                 utils.GetStringPropOrEmpty(props, "providerName"),
			Email:                utils.GetStringPropOrEmpty(props, "providerEmail"),
			AddressLine1:         utils.GetStringPropOrEmpty(props, "providerAddressLine1"),
			AddressLine2:         utils.GetStringPropOrEmpty(props, "providerAddressLine2"),
			Zip:                  utils.GetStringPropOrEmpty(props, "providerAddressZip"),
			Locality:             utils.GetStringPropOrEmpty(props, "providerAddressLocality"),
			Country:              utils.GetStringPropOrEmpty(props, "providerAddressCountry"),
			Region:               utils.GetStringPropOrEmpty(props, "providerAddressRegion"),
		},
		PaymentDetails: entity.PaymentDetails{
			PaymentLink: utils.GetStringPropOrEmpty(props, "paymentLink"),
		},
		InvoiceInternalFields: entity.InvoiceInternalFields{
			InvoiceFinalizedSentAt:            utils.GetTimePropOrNil(props, "techInvoiceFinalizedSentAt"),
			PaymentLinkRequestedAt:            utils.GetTimePropOrNil(props, "techPaymentLinkRequestedAt"),
			PayInvoiceNotificationRequestedAt: utils.GetTimePropOrNil(props, "techPayNotificationRequestedAt"),
			PayInvoiceNotificationSentAt:      utils.GetTimePropOrNil(props, "techPayInvoiceNotificationSentAt"),
			PaidInvoiceNotificationSentAt:     utils.GetTimePropOrNil(props, "techPaidInvoiceNotificationSentAt"),
		},
	}
	return &invoiceEntity
}

func MapDbNodeToInvoiceLineEntity(dbNode *dbtype.Node) *entity.InvoiceLineEntity {
	if dbNode == nil {
		return &entity.InvoiceLineEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	invoiceLineEntity := entity.InvoiceLineEntity{
		Id:                      utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:               utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:               utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Name:                    utils.GetStringPropOrEmpty(props, "name"),
		Price:                   utils.GetFloatPropOrZero(props, "price"),
		Quantity:                utils.GetInt64PropOrZero(props, "quantity"),
		Amount:                  utils.GetFloatPropOrZero(props, "amount"),
		Vat:                     utils.GetFloatPropOrZero(props, "vat"),
		TotalAmount:             utils.GetFloatPropOrZero(props, "totalAmount"),
		Source:                  entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:           entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:               utils.GetStringPropOrEmpty(props, "appSource"),
		ServiceLineItemId:       utils.GetStringPropOrEmpty(props, "serviceLineItemId"),
		ServiceLineItemParentId: utils.GetStringPropOrEmpty(props, "serviceLineItemParentId"),
		BilledType:              enum.DecodeBilledType(utils.GetStringPropOrEmpty(props, "billedType")),
	}
	return &invoiceLineEntity
}

func MapDbNodeToInvoicingCycleEntity(dbNode *dbtype.Node) *entity.InvoicingCycleEntity {
	if dbNode == nil {
		return &entity.InvoicingCycleEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	masterPlanEntity := entity.InvoicingCycleEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Type:          entity.InvoicingCycleType(utils.GetStringPropOrEmpty(props, "type")),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &masterPlanEntity
}

func MapDbNodeToMasterPlanEntity(dbNode *dbtype.Node) *entity.MasterPlanEntity {
	if dbNode == nil {
		return &entity.MasterPlanEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	masterPlanEntity := entity.MasterPlanEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Retired:       utils.GetBoolPropOrFalse(props, "retired"),
	}
	return &masterPlanEntity
}

func MapDbNodeToMasterPlanMilestoneEntity(dbNode *dbtype.Node) *entity.MasterPlanMilestoneEntity {
	if dbNode == nil {
		return &entity.MasterPlanMilestoneEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	masterPlanMilestoneEntity := entity.MasterPlanMilestoneEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Order:         utils.GetInt64PropOrZero(props, "order"),
		DurationHours: utils.GetInt64PropOrZero(props, "durationHours"),
		Optional:      utils.GetBoolPropOrFalse(props, "optional"),
		Items:         utils.GetListStringPropOrEmpty(props, "items"),
		Retired:       utils.GetBoolPropOrFalse(props, "retired"),
	}
	return &masterPlanMilestoneEntity
}

func MapDbNodeToUserEntity(dbNode *dbtype.Node) *entity.UserEntity {
	if dbNode == nil {
		return &entity.UserEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	userEntity := entity.UserEntity{
		Id:              utils.GetStringPropOrEmpty(props, "id"),
		FirstName:       utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:        utils.GetStringPropOrEmpty(props, "lastName"),
		Name:            utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:       utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:       utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:          entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:   entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:       utils.GetStringPropOrEmpty(props, "appSource"),
		Roles:           utils.GetListStringPropOrEmpty(props, "roles"),
		Internal:        utils.GetBoolPropOrFalse(props, "internal"),
		Bot:             utils.GetBoolPropOrFalse(props, "bot"),
		ProfilePhotoUrl: utils.GetStringPropOrEmpty(props, "profilePhotoUrl"),
		Timezone:        utils.GetStringPropOrEmpty(props, "timezone"),
	}
	return &userEntity
}

func MapDbNodeToOrganizationEntity(dbNode *dbtype.Node) *entity.OrganizationEntity {
	if dbNode == nil {
		return &entity.OrganizationEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	organizationEntity := entity.OrganizationEntity{
		ID:                 utils.GetStringPropOrEmpty(props, "id"),
		CustomerOsId:       utils.GetStringPropOrEmpty(props, "customerOsId"),
		Name:               utils.GetStringPropOrEmpty(props, "name"),
		Description:        utils.GetStringPropOrEmpty(props, "description"),
		Website:            utils.GetStringPropOrEmpty(props, "website"),
		Industry:           utils.GetStringPropOrEmpty(props, "industry"),
		IndustryGroup:      utils.GetStringPropOrEmpty(props, "industryGroup"),
		SubIndustry:        utils.GetStringPropOrEmpty(props, "subIndustry"),
		TargetAudience:     utils.GetStringPropOrEmpty(props, "targetAudience"),
		ValueProposition:   utils.GetStringPropOrEmpty(props, "valueProposition"),
		LastFundingRound:   utils.GetStringPropOrEmpty(props, "lastFundingRound"),
		LastFundingAmount:  utils.GetStringPropOrEmpty(props, "lastFundingAmount"),
		ReferenceId:        utils.GetStringPropOrEmpty(props, "referenceId"),
		Note:               utils.GetStringPropOrEmpty(props, "note"),
		IsPublic:           utils.GetBoolPropOrFalse(props, "isPublic"),
		IsCustomer:         utils.GetBoolPropOrFalse(props, "isCustomer"),
		Hide:               utils.GetBoolPropOrFalse(props, "hide"),
		Employees:          utils.GetInt64PropOrZero(props, "employees"),
		Market:             utils.GetStringPropOrEmpty(props, "market"),
		CreatedAt:          utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:          utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:             entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:      entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:          utils.GetStringPropOrEmpty(props, "appSource"),
		LastTouchpointAt:   utils.GetTimePropOrNil(props, "lastTouchpointAt"),
		LastTouchpointId:   utils.GetStringPropOrNil(props, "lastTouchpointId"),
		YearFounded:        utils.GetInt64PropOrNil(props, "yearFounded"),
		Headquarters:       utils.GetStringPropOrEmpty(props, "headquarters"),
		EmployeeGrowthRate: utils.GetStringPropOrEmpty(props, "employeeGrowthRate"),
		SlackChannelId:     utils.GetStringPropOrEmpty(props, "slackChannelId"),
		LogoUrl:            utils.GetStringPropOrEmpty(props, "logoUrl"),
		RenewalSummary: entity.RenewalSummary{
			ArrForecast:            utils.GetFloatPropOrNil(props, "renewalForecastArr"),
			MaxArrForecast:         utils.GetFloatPropOrNil(props, "renewalForecastMaxArr"),
			RenewalLikelihood:      utils.GetStringPropOrEmpty(props, "derivedRenewalLikelihood"),
			RenewalLikelihoodOrder: utils.GetInt64PropOrNil(props, "derivedRenewalLikelihoodOrder"),
			NextRenewalAt:          utils.GetTimePropOrNil(props, "derivedNextRenewalAt"),
		},
		WebScrapeDetails: entity.WebScrapeDetails{
			WebScrapedUrl:             utils.GetStringPropOrEmpty(props, "webScrapedUrl"),
			WebScrapedAt:              utils.GetTimePropOrNil(props, "webScrapedAt"),
			WebScrapeLastRequestedAt:  utils.GetTimePropOrNil(props, "webScrapeLastRequestedAt"),
			WebScrapeLastRequestedUrl: utils.GetStringPropOrEmpty(props, "webScrapeLastRequestedUrl"),
			WebScrapeAttempts:         utils.GetInt64PropOrZero(props, "webScrapeAttempts"),
		},
		OnboardingDetails: entity.OnboardingDetails{
			Status:       utils.GetStringPropOrEmpty(props, "onboardingStatus"),
			SortingOrder: utils.GetInt64PropOrNil(props, "onboardingStatusOrder"),
			UpdatedAt:    utils.GetTimePropOrNil(props, "onboardingUpdatedAt"),
			Comments:     utils.GetStringPropOrEmpty(props, "onboardingComments"),
		},
	}
	return &organizationEntity
}

func MapDbNodeToBillingProfileEntity(dbNode *dbtype.Node) *entity.BillingProfileEntity {
	if dbNode == nil {
		return &entity.BillingProfileEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	billingProfileEntity := entity.BillingProfileEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		LegalName:     utils.GetStringPropOrEmpty(props, "legalName"),
		TaxId:         utils.GetStringPropOrEmpty(props, "taxId"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &billingProfileEntity
}

func MapDbNodeToTenantEntity(dbNode *dbtype.Node) *entity.TenantEntity {
	if dbNode == nil {
		return &entity.TenantEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	tenant := entity.TenantEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource: utils.GetStringPropOrEmpty(props, "appSource"),
		Source:    entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
	}
	return &tenant
}

func MapDbNodeToTenantSettingsEntity(dbNode *dbtype.Node) *entity.TenantSettingsEntity {
	if dbNode == nil {
		return &entity.TenantSettingsEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	tenantSettingsEntity := entity.TenantSettingsEntity{
		Id:                   utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:            utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:            utils.GetTimePropOrEpochStart(props, "updatedAt"),
		LogoRepositoryFileId: utils.GetStringPropOrEmpty(props, "logoRepositoryFileId"),
		InvoicingEnabled:     utils.GetBoolPropOrFalse(props, "invoicingEnabled"),
		InvoicingPostpaid:    utils.GetBoolPropOrFalse(props, "invoicingPostpaid"),
		BaseCurrency:         enum.DecodeCurrency(utils.GetStringPropOrEmpty(props, "baseCurrency")),
	}
	return &tenantSettingsEntity
}

func MapDbNodeToTenantBillingProfileEntity(dbNode *dbtype.Node) *entity.TenantBillingProfileEntity {
	if dbNode == nil {
		return &entity.TenantBillingProfileEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	tenantBillingProfile := entity.TenantBillingProfileEntity{
		Id:                     utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:              utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:              utils.GetTimePropOrEpochStart(props, "updatedAt"),
		LegalName:              utils.GetStringPropOrEmpty(props, "legalName"),
		Phone:                  utils.GetStringPropOrEmpty(props, "phone"),
		AddressLine1:           utils.GetStringPropOrEmpty(props, "addressLine1"),
		AddressLine2:           utils.GetStringPropOrEmpty(props, "addressLine2"),
		AddressLine3:           utils.GetStringPropOrEmpty(props, "addressLine3"),
		Locality:               utils.GetStringPropOrEmpty(props, "locality"),
		Country:                utils.GetStringPropOrEmpty(props, "country"),
		Region:                 utils.GetStringPropOrEmpty(props, "region"),
		Zip:                    utils.GetStringPropOrEmpty(props, "zip"),
		VatNumber:              utils.GetStringPropOrEmpty(props, "vatNumber"),
		SendInvoicesFrom:       utils.GetStringPropOrEmpty(props, "sendInvoicesFrom"),
		SendInvoicesBcc:        utils.GetStringPropOrEmpty(props, "sendInvoicesBcc"),
		CanPayWithPigeon:       utils.GetBoolPropOrFalse(props, "canPayWithPigeon"),
		CanPayWithBankTransfer: utils.GetBoolPropOrFalse(props, "canPayWithBankTransfer"),
		Source:                 entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:          entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:              utils.GetStringPropOrEmpty(props, "appSource"),
		Check:                  utils.GetBoolPropOrFalse(props, "check"),
	}
	return &tenantBillingProfile
}

func MapDbNodeToCountryEntity(dbNode *dbtype.Node) *entity.CountryEntity {
	if dbNode == nil {
		return &entity.CountryEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	result := entity.CountryEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		CodeA2:    utils.GetStringPropOrEmpty(props, "codeA2"),
		CodeA3:    utils.GetStringPropOrEmpty(props, "codeA3"),
		PhoneCode: utils.GetStringPropOrEmpty(props, "phoneCode"),
	}
	return &result
}

func MapDbNodeToContractEntity(dbNode *dbtype.Node) *entity.ContractEntity {
	if dbNode == nil {
		return &entity.ContractEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	contract := entity.ContractEntity{
		Id:                              utils.GetStringPropOrEmpty(props, "id"),
		Name:                            utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:                       utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:                       utils.GetTimePropOrEpochStart(props, "updatedAt"),
		ServiceStartedAt:                utils.GetTimePropOrNil(props, "serviceStartedAt"),
		SignedAt:                        utils.GetTimePropOrNil(props, "signedAt"),
		EndedAt:                         utils.GetTimePropOrNil(props, "endedAt"),
		ContractUrl:                     utils.GetStringPropOrEmpty(props, "contractUrl"),
		ContractStatus:                  enum.DecodeContractStatus(utils.GetStringPropOrEmpty(props, "status")),
		RenewalCycle:                    enum.DecodeRenewalCycle(utils.GetStringPropOrEmpty(props, "renewalCycle")),
		RenewalPeriods:                  utils.GetInt64PropOrNil(props, "renewalPeriods"),
		TriggeredOnboardingStatusChange: utils.GetBoolPropOrFalse(props, "triggeredOnboardingStatusChange"),
		NextInvoiceDate:                 utils.GetTimePropOrNil(props, "nextInvoiceDate"),
		Source:                          entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:                   entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:                       utils.GetStringPropOrEmpty(props, "appSource"),
		InvoicingStartDate:              utils.GetTimePropOrNil(props, "invoicingStartDate"),
		Currency:                        enum.DecodeCurrency(utils.GetStringPropOrEmpty(props, "currency")),
		BillingCycle:                    enum.DecodeBillingCycle(utils.GetStringPropOrEmpty(props, "billingCycle")),
		AddressLine1:                    utils.GetStringPropOrEmpty(props, "addressLine1"),
		AddressLine2:                    utils.GetStringPropOrEmpty(props, "addressLine2"),
		Zip:                             utils.GetStringPropOrEmpty(props, "zip"),
		Locality:                        utils.GetStringPropOrEmpty(props, "locality"),
		Country:                         utils.GetStringPropOrEmpty(props, "country"),
		Region:                          utils.GetStringPropOrEmpty(props, "region"),
		OrganizationLegalName:           utils.GetStringPropOrEmpty(props, "organizationLegalName"),
		InvoiceEmail:                    utils.GetStringPropOrEmpty(props, "invoiceEmail"),
		InvoiceNote:                     utils.GetStringPropOrEmpty(props, "invoiceNote"),
		CanPayWithCard:                  utils.GetBoolPropOrFalse(props, "canPayWithCard"),
		CanPayWithDirectDebit:           utils.GetBoolPropOrFalse(props, "canPayWithDirectDebit"),
		CanPayWithBankTransfer:          utils.GetBoolPropOrFalse(props, "canPayWithBankTransfer"),
		InvoicingEnabled:                utils.GetBoolPropOrFalse(props, "invoicingEnabled"),
		PayOnline:                       utils.GetBoolPropOrFalse(props, "payOnline"),
		PayAutomatically:                utils.GetBoolPropOrFalse(props, "payAutomatically"),
		AutoRenew:                       utils.GetBoolPropOrFalse(props, "autoRenew"),
		DueDays:                         utils.GetInt64PropOrZero(props, "dueDays"),
		Check:                           utils.GetBoolPropOrFalse(props, "check"),
		ContractInternalFields: entity.ContractInternalFields{
			StatusRenewalRequestedAt:  utils.GetTimePropOrNil(props, "techStatusRenewalRequestedAt"),
			RolloutRenewalRequestedAt: utils.GetTimePropOrNil(props, "techRolloutRenewalRequestedAt"),
		},
	}
	return &contract
}

func MapDbNodeToOrganizationPlanEntity(dbNode *dbtype.Node) *entity.OrganizationPlanEntity {
	if dbNode == nil {
		return &entity.OrganizationPlanEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	orgPlanEntity := entity.OrganizationPlanEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Retired:       utils.GetBoolPropOrFalse(props, "retired"),
		StatusDetails: entity.OrganizationPlanStatusDetails{
			Status:    utils.GetStringPropOrEmpty(props, "status"),
			UpdatedAt: utils.GetTimePropOrEpochStart(props, "statusUpdatedAt"),
			Comments:  utils.GetStringPropOrEmpty(props, "statusComments"),
		},
		MasterPlanId: utils.GetStringPropOrEmpty(props, "masterPlanId"),
	}
	return &orgPlanEntity
}

func MapDbNodeToOrganizationPlanMilestoneEntity(dbNode *dbtype.Node) *entity.OrganizationPlanMilestoneEntity {
	if dbNode == nil {
		return &entity.OrganizationPlanMilestoneEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	orgPlanMilestoneEntity := entity.OrganizationPlanMilestoneEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Order:         utils.GetInt64PropOrZero(props, "order"),
		DueDate:       utils.GetTimePropOrEpochStart(props, "dueDate"),
		Optional:      utils.GetBoolPropOrFalse(props, "optional"),
		Items:         MapOrganizationPlanMilestoneItemToEntity(props),
		Retired:       utils.GetBoolPropOrFalse(props, "retired"),
		StatusDetails: entity.OrganizationPlanMilestoneStatusDetails{
			Status:    utils.GetStringPropOrEmpty(props, "status"),
			UpdatedAt: utils.GetTimePropOrEpochStart(props, "statusUpdatedAt"),
			Comments:  utils.GetStringPropOrEmpty(props, "statusComments"),
		},
		Adhoc: utils.GetBoolPropOrFalse(props, "adhoc"),
	}
	return &orgPlanMilestoneEntity
}

func MapOrganizationPlanMilestoneItemToEntity(props map[string]any) []entity.OrganizationPlanMilestoneItem {
	items := props["items"].([]any)

	itemArray := make([]entity.OrganizationPlanMilestoneItem, 0)
	if items == nil {
		return itemArray
	}
	for _, anyitem := range items {
		item := anyitem.(string)
		itemEntity := entity.OrganizationPlanMilestoneItem{}
		json.Unmarshal([]byte(item), &itemEntity)
		itemArray = append(itemArray, itemEntity)
	}
	return itemArray
}

func MapDbNodeToServiceLineItemEntity(dbNode *dbtype.Node) *entity.ServiceLineItemEntity {
	if dbNode == nil {
		return &entity.ServiceLineItemEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	serviceLineItem := entity.ServiceLineItemEntity{
		ID:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		StartedAt:     utils.GetTimePropOrEpochStart(props, "startedAt"),
		EndedAt:       utils.GetTimePropOrNil(props, "endedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		Billed:        enum.DecodeBilledType(utils.GetStringPropOrEmpty(props, "billed")),
		Price:         utils.GetFloatPropOrZero(props, "price"),
		Quantity:      utils.GetInt64PropOrZero(props, "quantity"),
		Comments:      utils.GetStringPropOrEmpty(props, "comments"),
		ParentID:      utils.GetStringPropOrEmpty(props, "parentId"),
		IsCanceled:    utils.GetBoolPropOrFalse(props, "isCanceled"),
		VatRate:       utils.GetFloatPropOrZero(props, "vatRate"),
	}
	return &serviceLineItem
}

func MapDbNodeToTagEntity(dbNode *dbtype.Node) *entity.TagEntity {
	if dbNode == nil {
		return &entity.TagEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	tag := entity.TagEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &tag
}

func MapDbNodeToIssueEntity(dbNode *dbtype.Node) *entity.IssueEntity {
	if dbNode == nil {
		return &entity.IssueEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	issue := entity.IssueEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:     utils.GetTimePropOrNow(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrNow(props, "updatedAt"),
		Subject:       utils.GetStringPropOrEmpty(props, "subject"),
		Status:        utils.GetStringPropOrEmpty(props, "status"),
		Priority:      utils.GetStringPropOrEmpty(props, "priority"),
		Description:   utils.GetStringPropOrEmpty(props, "description"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &issue
}

func MapDbNodeToCommentEntity(dbNode *dbtype.Node) *entity.CommentEntity {
	if dbNode == nil {
		return &entity.CommentEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	comment := entity.CommentEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &comment
}

func MapDbNodeToSocialEntity(dbNode *dbtype.Node) *entity.SocialEntity {
	if dbNode == nil {
		return &entity.SocialEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	social := entity.SocialEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		PlatformName:  utils.GetStringPropOrEmpty(props, "platformName"),
		Url:           utils.GetStringPropOrEmpty(props, "url"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &social
}

func MapDbNodeToReminderEntity(dbNode *dbtype.Node) *entity.ReminderEntity {
	if dbNode == nil {
		return &entity.ReminderEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	reminder := entity.ReminderEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		DueDate:       utils.GetTimePropOrEpochStart(props, "dueDate"),
		Dismissed:     utils.GetBoolPropOrFalse(props, "dismissed"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &reminder
}

func MapDbNodeToOrderEntity(dbNode *dbtype.Node) *entity.OrderEntity {
	if dbNode == nil {
		return &entity.OrderEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	return &entity.OrderEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		ConfirmedAt:   utils.GetTimePropOrNil(props, "confirmedAt"),
		PaidAt:        utils.GetTimePropOrNil(props, "paidAt"),
		FulfilledAt:   utils.GetTimePropOrNil(props, "fulfilledAt"),
		CancelledAt:   utils.GetTimePropOrNil(props, "cancelledAt"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
}

func MapDbNodeToBankAccountEntity(dbNode *dbtype.Node) *entity.BankAccountEntity {
	if dbNode == nil {
		return &entity.BankAccountEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	bankAccount := entity.BankAccountEntity{
		Id:                  utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:           utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:           utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:              entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:       entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:           utils.GetStringPropOrEmpty(props, "appSource"),
		BankName:            utils.GetStringPropOrEmpty(props, "bankName"),
		Currency:            enum.DecodeCurrency(utils.GetStringPropOrEmpty(props, "currency")),
		BankTransferEnabled: utils.GetBoolPropOrFalse(props, "bankTransferEnabled"),
		AllowInternational:  utils.GetBoolPropOrFalse(props, "allowInternational"),
		AccountNumber:       utils.GetStringPropOrEmpty(props, "accountNumber"),
		SortCode:            utils.GetStringPropOrEmpty(props, "sortCode"),
		Iban:                utils.GetStringPropOrEmpty(props, "iban"),
		Bic:                 utils.GetStringPropOrEmpty(props, "bic"),
		RoutingNumber:       utils.GetStringPropOrEmpty(props, "routingNumber"),
		OtherDetails:        utils.GetStringPropOrEmpty(props, "otherDetails"),
	}
	return &bankAccount
}

func MapDbNodeToEmailEntity(node *dbtype.Node) *entity.EmailEntity {
	if node == nil {
		return &entity.EmailEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	return &entity.EmailEntity{
		Id:             utils.GetStringPropOrEmpty(props, "id"),
		Email:          utils.GetStringPropOrEmpty(props, "email"),
		RawEmail:       utils.GetStringPropOrEmpty(props, "rawEmail"),
		CreatedAt:      utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:      utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Primary:        utils.GetBoolPropOrFalse(props, "primary"),
		Source:         entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:  entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:      utils.GetStringPropOrEmpty(props, "appSource"),
		Label:          utils.GetStringPropOrEmpty(props, "label"),
		Validated:      utils.GetBoolPropOrNil(props, "validated"),
		IsReachable:    utils.GetStringPropOrNil(props, "isReachable"),
		IsValidSyntax:  utils.GetBoolPropOrNil(props, "isValidSyntax"),
		CanConnectSMTP: utils.GetBoolPropOrNil(props, "canConnectSMTP"),
		AcceptsMail:    utils.GetBoolPropOrNil(props, "acceptsMail"),
		HasFullInbox:   utils.GetBoolPropOrNil(props, "hasFullInbox"),
		IsCatchAll:     utils.GetBoolPropOrNil(props, "isCatchAll"),
		IsDeliverable:  utils.GetBoolPropOrNil(props, "isDeliverable"),
		IsDisabled:     utils.GetBoolPropOrNil(props, "isDisabled"),
		Error:          utils.GetStringPropOrNil(props, "error"),
	}
}

func MapDbNodeToExternalSystem(node *dbtype.Node) *entity.ExternalSystemEntity {
	if node == nil {
		return &entity.ExternalSystemEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	externalSystemEntity := entity.ExternalSystemEntity{
		ExternalSystemId: enum.DecodeExternalSystemId(utils.GetStringPropOrEmpty(props, "id")),
		Name:             utils.GetStringPropOrEmpty(props, "name"),
	}
	if externalSystemEntity.ExternalSystemId == enum.Stripe {
		externalSystemEntity.Stripe.PaymentMethodTypes = utils.GetListStringPropOrEmpty(props, entity.PropertyExternalSystemStripePaymentMethodTypes)
	}
	return &externalSystemEntity
}

func MapDbNodeToOpportunityEntity(node *dbtype.Node) *entity.OpportunityEntity {
	if node == nil {
		return &entity.OpportunityEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	opportunity := entity.OpportunityEntity{
		Id:                utils.GetStringPropOrEmpty(props, "id"),
		Name:              utils.GetStringPropOrEmpty(props, "name"),
		Amount:            utils.GetFloatPropOrZero(props, "amount"),
		MaxAmount:         utils.GetFloatPropOrZero(props, "maxAmount"),
		InternalType:      enum.DecodeOpportunityInternalType(utils.GetStringPropOrEmpty(props, "internalType")),
		ExternalType:      utils.GetStringPropOrEmpty(props, "externalType"),
		InternalStage:     enum.DecodeOpportunityInternalStage(utils.GetStringPropOrEmpty(props, "internalStage")),
		ExternalStage:     utils.GetStringPropOrEmpty(props, "externalStage"),
		EstimatedClosedAt: utils.GetTimePropOrNil(props, "estimatedClosedAt"),
		ClosedAt:          utils.GetTimePropOrNil(props, "closedAt"),
		GeneralNotes:      utils.GetStringPropOrEmpty(props, "generalNotes"),
		NextSteps:         utils.GetStringPropOrEmpty(props, "nextSteps"),
		Comments:          utils.GetStringPropOrEmpty(props, "comments"),
		CreatedAt:         utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:         utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:         utils.GetStringPropOrEmpty(props, "appSource"),
		Source:            entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:     entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		OwnerUserId:       utils.GetStringPropOrEmpty(props, "ownerUserId"),
		RenewalDetails: entity.RenewalDetails{
			RenewedAt:              utils.GetTimePropOrNil(props, "renewedAt"),
			RenewalLikelihood:      enum.DecodeRenewalLikelihood(utils.GetStringPropOrEmpty(props, "renewalLikelihood")),
			RenewalUpdatedByUserId: utils.GetStringPropOrEmpty(props, "renewalUpdatedByUserId"),
			RenewalUpdatedByUserAt: utils.GetTimePropOrNil(props, "renewalUpdatedByUserAt"),
			RenewalApproved:        utils.GetBoolPropOrFalse(props, "renewalApproved"),
		},
		InternalFields: entity.OpportunityInternalFields{
			RolloutRenewalRequestedAt: utils.GetTimePropOrNil(props, "techRolloutRenewalRequestedAt"),
		},
	}
	return &opportunity
}

func MapDbNodeToOfferingEntity(node *dbtype.Node) *entity.OfferingEntity {
	if node == nil {
		return &entity.OfferingEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	offering := entity.OfferingEntity{
		Id:                    utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:             utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:             utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:                entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:         entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:             utils.GetStringPropOrEmpty(props, "appSource"),
		Name:                  utils.GetStringPropOrEmpty(props, "name"),
		Active:                utils.GetBoolPropOrFalse(props, "active"),
		Type:                  enum.DecodeOfferingType(utils.GetStringPropOrEmpty(props, "type")),
		PricingModel:          enum.DecodePricingModel(utils.GetStringPropOrEmpty(props, "pricingModel")),
		PricingPeriodInMonths: utils.GetInt64PropOrZero(props, "pricingPeriodInMonths"),
		Currency:              enum.DecodeCurrency(utils.GetStringPropOrEmpty(props, "currency")),
		Price:                 utils.GetFloatPropOrZero(props, "price"),
		PriceCalculated:       utils.GetBoolPropOrFalse(props, "priceCalculated"),
		Conditional:           utils.GetBoolPropOrFalse(props, "conditional"),
		Taxable:               utils.GetBoolPropOrFalse(props, "taxable"),
		PriceCalculation: entity.PriceCalculation{
			Type:                   enum.DecodePriceCalculationType(utils.GetStringPropOrEmpty(props, "priceCalculationType")),
			RevenueSharePercentage: utils.GetFloatPropOrZero(props, "priceCalculationRevenueSharePercentage"),
		},
		Conditionals: entity.Conditionals{
			MinimumChargePeriod: enum.DecodeChargePeriod(utils.GetStringPropOrEmpty(props, "conditionalsMinimumChargePeriod")),
			MinimumChargeAmount: utils.GetFloatPropOrZero(props, "conditionalsMinimumChargeAmount"),
		},
	}
	return &offering
}
