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
		return nil
	}
	props := utils.GetPropsFromNode(*dbNode)
	invoiceEntity := entity.InvoiceEntity{
		Id:               utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:        utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:        utils.GetTimePropOrEpochStart(props, "updatedAt"),
		DryRun:           utils.GetBoolPropOrFalse(props, "dryRun"),
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
		InvoiceInternalFields: entity.InvoiceInternalFields{
			PaymentRequestedAt: utils.GetTimePropOrNil(props, "techPaymentRequestedAt"),
		},
	}
	return &invoiceEntity
}

func MapDbNodeToInvoiceLineEntity(dbNode *dbtype.Node) *entity.InvoiceLineEntity {
	if dbNode == nil {
		return nil
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
		return nil
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
		return nil
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
		return nil
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
		return nil
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
		return nil
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
		return nil
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

func MapDbNodeToTenantSettingsEntity(dbNode *dbtype.Node) *entity.TenantSettingsEntity {
	if dbNode == nil {
		return &entity.TenantSettingsEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	tenantSettingsEntity := entity.TenantSettingsEntity{
		CreatedAt:        utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:        utils.GetTimePropOrEpochStart(props, "updatedAt"),
		LogoUrl:          utils.GetStringPropOrEmpty(props, "logoUrl"),
		InvoicingEnabled: utils.GetBoolPropOrFalse(props, "invoicingEnabled"),
		DefaultCurrency:  enum.DecodeCurrency(utils.GetStringPropOrEmpty(props, "defaultCurrency")),
	}
	return &tenantSettingsEntity
}

func MapDbNodeToTenantBillingProfileEntity(dbNode *dbtype.Node) *entity.TenantBillingProfile {
	if dbNode == nil {
		return &entity.TenantBillingProfile{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	tenantBillingProfile := entity.TenantBillingProfile{
		Id:                                utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:                         utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:                         utils.GetTimePropOrEpochStart(props, "updatedAt"),
		LegalName:                         utils.GetStringPropOrEmpty(props, "legalName"),
		Email:                             utils.GetStringPropOrEmpty(props, "email"),
		Phone:                             utils.GetStringPropOrEmpty(props, "phone"),
		AddressLine1:                      utils.GetStringPropOrEmpty(props, "addressLine1"),
		AddressLine2:                      utils.GetStringPropOrEmpty(props, "addressLine2"),
		AddressLine3:                      utils.GetStringPropOrEmpty(props, "addressLine3"),
		DomesticPaymentsBankName:          utils.GetStringPropOrEmpty(props, "domesticPaymentsBankName"),
		DomesticPaymentsAccountNumber:     utils.GetStringPropOrEmpty(props, "domesticPaymentsAccountNumber"),
		DomesticPaymentsSortCode:          utils.GetStringPropOrEmpty(props, "domesticPaymentsSortCode"),
		InternationalPaymentsSwiftBic:     utils.GetStringPropOrEmpty(props, "internationalPaymentsSwiftBic"),
		InternationalPaymentsBankName:     utils.GetStringPropOrEmpty(props, "internationalPaymentsBankName"),
		InternationalPaymentsBankAddress:  utils.GetStringPropOrEmpty(props, "internationalPaymentsBankAddress"),
		InternationalPaymentsInstructions: utils.GetStringPropOrEmpty(props, "internationalPaymentsInstructions"),
		Source:                            entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:                     entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:                         utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &tenantBillingProfile
}

func MapDbNodeToCountryEntity(dbNode *dbtype.Node) *entity.CountryEntity {
	if dbNode == nil {
		return nil
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
		return nil
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
		OrganizationLegalName:           utils.GetStringPropOrEmpty(props, "organizationLegalName"),
		InvoiceEmail:                    utils.GetStringPropOrEmpty(props, "invoiceEmail"),
		InvoiceNote:                     utils.GetStringPropOrEmpty(props, "invoiceNote"),
	}
	return &contract
}

func MapDbNodeToOrganizationPlanEntity(dbNode *dbtype.Node) *entity.OrganizationPlanEntity {
	if dbNode == nil {
		return nil
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
	}
	return &orgPlanEntity
}

func MapDbNodeToOrganizationPlanMilestoneEntity(dbNode *dbtype.Node) *entity.OrganizationPlanMilestoneEntity {
	if dbNode == nil {
		return nil
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
		return nil
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
	}
	return &serviceLineItem
}
