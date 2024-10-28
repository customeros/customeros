package mapper

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"golang.org/x/exp/slices"
)

func MapDbNodeToJobRoleEntity(dbNode *dbtype.Node) *entity.JobRoleEntity {
	if dbNode == nil {
		return &entity.JobRoleEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	result := entity.JobRoleEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		JobTitle:      utils.GetStringPropOrEmpty(props, "jobTitle"),
		Description:   utils.GetStringPropOrNil(props, "description"),
		Company:       utils.GetStringPropOrNil(props, "company"),
		Primary:       utils.GetBoolPropOrFalse(props, "primary"),
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		StartedAt:     utils.GetTimePropOrNil(props, "startedAt"),
		EndedAt:       utils.GetTimePropOrNil(props, "endedAt"),
	}
	return &result
}

func MapDbNodeToAttachmentEntity(dbNode *dbtype.Node) *entity.AttachmentEntity {
	if dbNode == nil {
		return &entity.AttachmentEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	createdAt := utils.GetTimePropOrEpochStart(props, "createdAt")
	attachmentEntity := entity.AttachmentEntity{
		Id:                 utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:          &createdAt,
		FileName:           utils.GetStringPropOrEmpty(props, "fileName"),
		MimeType:           utils.GetStringPropOrEmpty(props, "mimeType"),
		CdnUrl:             utils.GetStringPropOrEmpty(props, "cdnUrl"),
		PublicUrl:          utils.GetStringPropOrEmpty(props, string(entity.AttachmentPropertyPublicUrl)),
		PublicUrlExpiresAt: utils.GetTimePropOrNil(props, string(entity.AttachmentPropertyPublicUrlExpiresAt)),
		BasePath:           utils.GetStringPropOrEmpty(props, "basePath"),
		Size:               utils.GetInt64PropOrZero(props, "size"),
		AppSource:          utils.GetStringPropOrEmpty(props, "appSource"),
		Source:             entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:      entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &attachmentEntity
}

func MapDbNodeToWorkspaceEntity(dbNode *dbtype.Node) *entity.WorkspaceEntity {
	if dbNode == nil {
		return &entity.WorkspaceEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	workspace := entity.WorkspaceEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "domain"),
		Provider:      utils.GetStringPropOrEmpty(props, "provider"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &workspace
}

func MapDbNodeToPlayerEntity(node *neo4j.Node) *entity.PlayerEntity {
	if node == nil {
		return &entity.PlayerEntity{}
	}
	props := utils.GetPropsFromNode(*node)

	return &entity.PlayerEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		AuthId:        utils.GetStringPropOrEmpty(props, "authId"),
		Provider:      utils.GetStringPropOrEmpty(props, "provider"),
		IdentityId:    utils.GetStringPropOrEmpty(props, "identityId"),
		Source:        utils.GetStringPropOrEmpty(props, "source"),
		SourceOfTruth: utils.GetStringPropOrEmpty(props, "sourceOfTruth"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
}

func MapDbNodeToInvoiceEntity(dbNode *dbtype.Node) *entity.InvoiceEntity {
	if dbNode == nil {
		return &entity.InvoiceEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	invoiceEntity := entity.InvoiceEntity{
		Id:                   utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:            utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:            utils.GetTimePropOrEpochStart(props, "updatedAt"),
		DryRun:               utils.GetBoolPropOrFalse(props, "dryRun"),
		OffCycle:             utils.GetBoolPropOrFalse(props, "offCycle"),
		Postpaid:             utils.GetBoolPropOrFalse(props, "postpaid"),
		Preview:              utils.GetBoolPropOrFalse(props, "preview"),
		Number:               utils.GetStringPropOrEmpty(props, "number"),
		PeriodStartDate:      utils.GetTimePropOrEpochStart(props, "periodStartDate"),
		PeriodEndDate:        utils.GetTimePropOrEpochStart(props, "periodEndDate"),
		DueDate:              utils.GetTimePropOrEpochStart(props, "dueDate"),
		IssuedDate:           utils.GetTimePropOrEpochStart(props, "issuedDate"),
		Currency:             enum.DecodeCurrency(utils.GetStringPropOrEmpty(props, "currency")),
		BillingCycleInMonths: utils.GetInt64PropOrZero(props, "billingCycleInMonths"),
		Amount:               utils.GetFloatPropOrZero(props, "amount"),
		Vat:                  utils.GetFloatPropOrZero(props, "vat"),
		TotalAmount:          utils.GetFloatPropOrZero(props, "totalAmount"),
		RepositoryFileId:     utils.GetStringPropOrEmpty(props, "repositoryFileId"),
		Source:               entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:            utils.GetStringPropOrEmpty(props, "appSource"),
		Status:               enum.DecodeInvoiceStatus(utils.GetStringPropOrEmpty(props, "status")),
		Note:                 utils.GetStringPropOrEmpty(props, "note"),
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
			PaymentLink:           utils.GetStringPropOrEmpty(props, string(entity.InvoicePropertyPaymentLink)),
			PaymentLinkValidUntil: utils.GetTimePropOrNil(props, string(entity.InvoicePropertyPaymentLinkValidUntil)),
		},
		InvoiceInternalFields: entity.InvoiceInternalFields{
			InvoiceFinalizedSentAt:               utils.GetTimePropOrNil(props, string(entity.InvoicePropertyInvoiceFinalizedEventSentAt)),
			InvoiceFinalizedWebhookProcessedAt:   utils.GetTimePropOrNil(props, string(entity.InvoicePropertyFinalizedWebhookProcessedAt)),
			InvoicePaidWebhookProcessedAt:        utils.GetTimePropOrNil(props, string(entity.InvoicePropertyPaidWebhookProcessedAt)),
			PaymentLinkRequestedAt:               utils.GetTimePropOrNil(props, "techPaymentLinkRequestedAt"),
			PayInvoiceNotificationRequestedAt:    utils.GetTimePropOrNil(props, "techPayNotificationRequestedAt"),
			PayInvoiceNotificationSentAt:         utils.GetTimePropOrNil(props, "techPayInvoiceNotificationSentAt"),
			RemindInvoiceNotificationRequestedAt: utils.GetTimePropOrNil(props, string(entity.InvoicePropertyRemindInvoiceNotificationRequestedAt)),
			LastRemindInvoiceNotificationSentAt:  utils.GetTimePropOrNil(props, string(entity.InvoicePropertyLastRemindInvoiceNotificationSentAt)),
			PaidInvoiceNotificationSentAt:        utils.GetTimePropOrNil(props, "techPaidInvoiceNotificationSentAt"),
			VoidInvoiceNotificationSentAt:        utils.GetTimePropOrNil(props, "techVoidInvoiceNotificationSentAt"),
		},
		EventStoreAggregate: entity.EventStoreAggregate{
			AggregateVersion: utils.GetInt64PropOrNil(props, "aggregateVersion"),
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
		Source:                  entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:           entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:               utils.GetStringPropOrEmpty(props, "appSource"),
		ServiceLineItemId:       utils.GetStringPropOrEmpty(props, "serviceLineItemId"),
		ServiceLineItemParentId: utils.GetStringPropOrEmpty(props, "serviceLineItemParentId"),
		BilledType:              enum.DecodeBilledType(utils.GetStringPropOrEmpty(props, "billedType")),
	}
	return &invoiceLineEntity
}

// TODO return nil
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
		Source:          entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:   entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
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
		ReferenceId:        utils.GetStringPropOrEmpty(props, "referenceId"),
		Name:               utils.GetStringPropOrEmpty(props, "name"),
		Description:        utils.GetStringPropOrEmpty(props, "description"),
		Website:            utils.GetStringPropOrEmpty(props, "website"),
		Industry:           utils.GetStringPropOrEmpty(props, string(entity.OrganizationPropertyIndustry)),
		IndustryGroup:      utils.GetStringPropOrEmpty(props, "industryGroup"),
		SubIndustry:        utils.GetStringPropOrEmpty(props, "subIndustry"),
		TargetAudience:     utils.GetStringPropOrEmpty(props, "targetAudience"),
		ValueProposition:   utils.GetStringPropOrEmpty(props, "valueProposition"),
		LastFundingRound:   utils.GetStringPropOrEmpty(props, "lastFundingRound"),
		LastFundingAmount:  utils.GetStringPropOrEmpty(props, "lastFundingAmount"),
		Note:               utils.GetStringPropOrEmpty(props, "note"),
		IsPublic:           utils.GetBoolPropOrFalse(props, string(entity.OrganizationPropertyIsPublic)),
		Hide:               utils.GetBoolPropOrFalse(props, string(entity.OrganizationPropertyHide)),
		Employees:          utils.GetInt64PropOrZero(props, string(entity.OrganizationPropertyEmployees)),
		Market:             utils.GetStringPropOrEmpty(props, "market"),
		CreatedAt:          utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:          utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:             entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		LastTouchpointAt:   utils.GetTimePropOrNil(props, "lastTouchpointAt"),
		LastTouchpointId:   utils.GetStringPropOrNil(props, "lastTouchpointId"),
		LastTouchpointType: utils.GetStringPropOrNil(props, "lastTouchpointType"),
		YearFounded:        utils.GetInt64PropOrNil(props, string(entity.OrganizationPropertyYearFounded)),
		Headquarters:       utils.GetStringPropOrEmpty(props, "headquarters"),
		EmployeeGrowthRate: utils.GetStringPropOrEmpty(props, "employeeGrowthRate"),
		SlackChannelId:     utils.GetStringPropOrEmpty(props, "slackChannelId"),
		LogoUrl:            utils.GetStringPropOrEmpty(props, "logoUrl"),
		IconUrl:            utils.GetStringPropOrEmpty(props, "iconUrl"),
		Relationship:       enum.DecodeOrganizationRelationship(utils.GetStringPropOrEmpty(props, "relationship")),
		Stage:              enum.DecodeOrganizationStage(utils.GetStringPropOrEmpty(props, "stage")),
		StageUpdatedAt:     utils.GetTimePropOrNil(props, "stageUpdatedAt"),
		LeadSource:         utils.GetStringPropOrEmpty(props, "leadSource"),
		IcpFit:             utils.GetBoolPropOrFalse(props, string(entity.OrganizationPropertyIcpFit)),
		RenewalSummary: entity.RenewalSummary{
			ArrForecast:            utils.GetFloatPropOrNil(props, "renewalForecastArr"),
			MaxArrForecast:         utils.GetFloatPropOrNil(props, "renewalForecastMaxArr"),
			RenewalLikelihood:      utils.GetStringPropOrEmpty(props, "derivedRenewalLikelihood"),
			RenewalLikelihoodOrder: utils.GetInt64PropOrNil(props, "derivedRenewalLikelihoodOrder"),
			NextRenewalAt:          utils.GetTimePropOrNil(props, "derivedNextRenewalAt"),
		},
		DerivedData: entity.DerivedData{
			ChurnedAt:   utils.GetTimePropOrNil(props, "derivedChurnedAt"),
			Ltv:         utils.GetFloatPropOrZero(props, "derivedLtv"),
			LtvCurrency: enum.DecodeCurrency(utils.GetStringPropOrEmpty(props, "derivedLtvCurrency")),
		},
		OnboardingDetails: entity.OnboardingDetails{
			Status:       utils.GetStringPropOrEmpty(props, "onboardingStatus"),
			SortingOrder: utils.GetInt64PropOrNil(props, "onboardingStatusOrder"),
			UpdatedAt:    utils.GetTimePropOrNil(props, "onboardingUpdatedAt"),
			Comments:     utils.GetStringPropOrEmpty(props, "onboardingComments"),
		},
		EnrichDetails: entity.OrganizationEnrichDetails{
			EnrichRequestedAt: utils.GetTimePropOrNil(props, string(entity.OrganizationPropertyEnrichRequestedAt)),
			EnrichedAt:        utils.GetTimePropOrNil(props, string(entity.OrganizationPropertyEnrichedAt)),
			EnrichFailedAt:    utils.GetTimePropOrNil(props, string(entity.OrganizationPropertyEnrichFailedAt)),
			EnrichAttempts:    utils.GetInt64PropOrZero(props, string(entity.OrganizationPropertyEnrichAttempts)),
			EnrichSource:      enum.DecodeDomainEnrichSource(utils.GetStringPropOrEmpty(props, "enrichSource")),
			EnrichDomain:      utils.GetStringPropOrEmpty(props, "enrichDomain"),
		},
		OrganizationInternalFields: entity.OrganizationInternalFields{
			DomainCheckedAt:   utils.GetTimePropOrNil(props, string(entity.OrganizationPropertyDomainCheckedAt)),
			IndustryCheckedAt: utils.GetTimePropOrNil(props, string(entity.OrganizationPropertyIndustryCheckedAt)),
			HiddenAt:          utils.GetTimePropOrNil(props, string(entity.OrganizationPropertyHiddenAt)),
		},
		EventStoreAggregate: entity.EventStoreAggregate{
			AggregateVersion: utils.GetInt64PropOrNil(props, "aggregateVersion"),
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
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
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
		Source:    entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		Active:    utils.GetBoolPropOrTrue(props, "active"),
	}
	return &tenant
}

func MapDbNodeToTenantSettingsEntity(dbNode *dbtype.Node) *entity.TenantSettingsEntity {
	if dbNode == nil {
		return &entity.TenantSettingsEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	tenantSettingsEntity := entity.TenantSettingsEntity{
		Id:                       utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:                utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:                utils.GetTimePropOrEpochStart(props, "updatedAt"),
		LogoRepositoryFileId:     utils.GetStringPropOrEmpty(props, string(entity.TenantSettingsPropertyLogoRepositoryFileId)),
		InvoicingEnabled:         utils.GetBoolPropOrFalse(props, string(entity.TenantSettingsPropertyInvoicingEnabled)),
		InvoicingPostpaid:        utils.GetBoolPropOrFalse(props, string(entity.TenantSettingsPropertyInvoicingPostpaid)),
		WorkspaceLogo:            utils.GetStringPropOrEmpty(props, string(entity.TenantSettingsPropertyWorkspaceLogo)),
		WorkspaceName:            utils.GetStringPropOrEmpty(props, string(entity.TenantSettingsPropertyWorkspaceName)),
		BaseCurrency:             enum.DecodeCurrency(utils.GetStringPropOrEmpty(props, string(entity.TenantSettingsPropertyBaseCurrency))),
		EnrichContacts:           utils.GetBoolPropOrFalse(props, string(entity.TenantSettingsPropertyEnrichContacts)),
		StripeCustomerPortalLink: utils.GetStringPropOrEmpty(props, string(entity.TenantSettingsPropertyStripeCustomerPortalLink)),
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
		Source:                 entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:          entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
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
		TriggeredOnboardingStatusChange: utils.GetBoolPropOrFalse(props, "triggeredOnboardingStatusChange"),
		NextInvoiceDate:                 utils.GetTimePropOrNil(props, "nextInvoiceDate"),
		Source:                          entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:                   entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:                       utils.GetStringPropOrEmpty(props, "appSource"),
		InvoicingStartDate:              utils.GetTimePropOrNil(props, "invoicingStartDate"),
		Currency:                        enum.DecodeCurrency(utils.GetStringPropOrEmpty(props, "currency")),
		BillingCycleInMonths:            utils.GetInt64PropOrZero(props, "billingCycleInMonths"),
		AddressLine1:                    utils.GetStringPropOrEmpty(props, "addressLine1"),
		AddressLine2:                    utils.GetStringPropOrEmpty(props, "addressLine2"),
		Zip:                             utils.GetStringPropOrEmpty(props, "zip"),
		Locality:                        utils.GetStringPropOrEmpty(props, "locality"),
		Country:                         utils.GetStringPropOrEmpty(props, "country"),
		Region:                          utils.GetStringPropOrEmpty(props, "region"),
		OrganizationLegalName:           utils.GetStringPropOrEmpty(props, "organizationLegalName"),
		InvoiceEmail:                    utils.GetStringPropOrEmpty(props, "invoiceEmail"),
		InvoiceEmailCC:                  utils.GetListStringPropOrEmpty(props, "invoiceEmailCC"),
		InvoiceEmailBCC:                 utils.GetListStringPropOrEmpty(props, "invoiceEmailBCC"),
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
		LengthInMonths:                  utils.GetInt64PropOrZero(props, "lengthInMonths"),
		Approved:                        utils.GetBoolPropOrFalse(props, "approved"),
		Ltv:                             utils.GetFloatPropOrZero(props, "ltv"),
		ContractInternalFields: entity.ContractInternalFields{
			StatusRenewalRequestedAt:      utils.GetTimePropOrNil(props, "techStatusRenewalRequestedAt"),
			RolloutRenewalRequestedAt:     utils.GetTimePropOrNil(props, "techRolloutRenewalRequestedAt"),
			NextPreviewInvoiceRequestedAt: utils.GetTimePropOrNil(props, "techNextPreviewInvoiceRequestedAt"),
		},
		EventStoreAggregate: entity.EventStoreAggregate{
			AggregateVersion: utils.GetInt64PropOrNil(props, "aggregateVersion"),
		},
	}

	return &contract
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
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		Billed:        enum.DecodeBilledType(utils.GetStringPropOrEmpty(props, "billed")),
		Price:         utils.GetFloatPropOrZero(props, "price"),
		Quantity:      utils.GetInt64PropOrZero(props, "quantity"),
		Comments:      utils.GetStringPropOrEmpty(props, "comments"),
		ParentID:      utils.GetStringPropOrEmpty(props, "parentId"),
		Canceled:      utils.GetBoolPropOrFalse(props, "isCanceled"),
		VatRate:       utils.GetFloatPropOrZero(props, "vatRate"),
		Paused:        utils.GetBoolPropOrFalse(props, string(entity.SLIPropertyPaused)),
	}
	return &serviceLineItem
}

func MapDbNodeToTagEntity(dbNode *dbtype.Node) *entity.TagEntity {
	if dbNode == nil {
		return &entity.TagEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	tag := entity.TagEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:    entity.DataSource(utils.GetStringPropOrEmpty(props, "source")),
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
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
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
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &comment
}

func MapDbNodeToSocialEntity(dbNode *dbtype.Node) *entity.SocialEntity {
	if dbNode == nil {
		return &entity.SocialEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	social := entity.SocialEntity{
		Id:             utils.GetStringPropOrEmpty(props, string(entity.SocialPropertyId)),
		Url:            utils.GetStringPropOrEmpty(props, string(entity.SocialPropertyUrl)),
		Alias:          utils.GetStringPropOrEmpty(props, string(entity.SocialPropertyAlias)),
		FollowersCount: utils.GetInt64PropOrZero(props, string(entity.SocialPropertyFollowersCount)),
		ExternalId:     utils.GetStringPropOrEmpty(props, string(entity.SocialPropertyExternalId)),
		CreatedAt:      utils.GetTimePropOrEpochStart(props, string(entity.SocialPropertyCreatedAt)),
		UpdatedAt:      utils.GetTimePropOrEpochStart(props, string(entity.SocialPropertyUpdatedAt)),
		AppSource:      utils.GetStringPropOrEmpty(props, string(entity.SocialPropertyAppSource)),
		Source:         entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, string(entity.SocialPropertySource))),
		SourceOfTruth:  entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, string(entity.SocialPropertySourceOfTruth))),
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
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &reminder
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
		Source:              entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:       entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
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

// TODO RETURN NIL NOT EMPTY
func MapDbNodeToEmailEntity(node *dbtype.Node) *entity.EmailEntity {
	if node == nil {
		return &entity.EmailEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	return &entity.EmailEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Email:         utils.GetStringPropOrEmpty(props, string(entity.EmailPropertyEmail)),
		RawEmail:      utils.GetStringPropOrEmpty(props, string(entity.EmailPropertyRawEmail)),
		Work:          utils.GetBoolPropOrNil(props, string(entity.EmailPropertyWork)),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Primary:       utils.GetBoolPropOrFalse(props, "primary"),
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		IsValidSyntax: utils.GetBoolPropOrNil(props, string(entity.EmailPropertyIsValidSyntax)),
		IsCatchAll:    utils.GetBoolPropOrNil(props, string(entity.EmailPropertyIsCatchAll)),
		Deliverable:   utils.GetStringPropOrNil(props, string(entity.EmailPropertyDeliverable)),
		IsRoleAccount: utils.GetBoolPropOrNil(props, string(entity.EmailPropertyIsRoleAccount)),
		EmailInternalFields: entity.EmailInternalFields{
			ValidatedAt:           utils.GetTimePropOrNil(props, string(entity.EmailPropertyValidatedAt)),
			ValidationRequestedAt: utils.GetTimePropOrNil(props, string(entity.EmailPropertyValidationRequestedAt)),
		},
		Username:        utils.GetStringPropOrNil(props, string(entity.EmailPropertyUsername)),
		IsRisky:         utils.GetBoolPropOrNil(props, string(entity.EmailPropertyIsRisky)),
		IsFirewalled:    utils.GetBoolPropOrNil(props, string(entity.EmailPropertyIsFirewalled)),
		Provider:        utils.GetStringPropOrNil(props, string(entity.EmailPropertyProvider)),
		Firewall:        utils.GetStringPropOrNil(props, string(entity.EmailPropertyFirewall)),
		IsMailboxFull:   utils.GetBoolPropOrNil(props, string(entity.EmailPropertyIsMailboxFull)),
		IsFreeAccount:   utils.GetBoolPropOrNil(props, string(entity.EmailPropertyIsFreeAccount)),
		SmtpSuccess:     utils.GetBoolPropOrNil(props, string(entity.EmailPropertySmtpSuccess)),
		ResponseCode:    utils.GetStringPropOrNil(props, string(entity.EmailPropertyResponseCode)),
		ErrorCode:       utils.GetStringPropOrNil(props, string(entity.EmailPropertyErrorCode)),
		Description:     utils.GetStringPropOrNil(props, string(entity.EmailPropertyDescription)),
		IsPrimaryDomain: utils.GetBoolPropOrNil(props, string(entity.EmailPropertyIsPrimaryDomain)),
		PrimaryDomain:   utils.GetStringPropOrNil(props, string(entity.EmailPropertyPrimaryDomain)),
		AlternateEmail:  utils.GetStringPropOrNil(props, string(entity.EmailPropertyAlternateEmail)),
		RetryValidation: utils.GetBoolPropOrNil(props, string(entity.EmailPropertyRetryValidation)),
	}
}

func MapDbNodeToPhoneNumberEntity(node *dbtype.Node) *entity.PhoneNumberEntity {
	if node == nil {
		return &entity.PhoneNumberEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	return &entity.PhoneNumberEntity{
		Id:             utils.GetStringPropOrEmpty(props, "id"),
		E164:           utils.GetStringPropOrEmpty(props, "e164"),
		RawPhoneNumber: utils.GetStringPropOrEmpty(props, "rawPhoneNumber"),
		Validated:      utils.GetBoolPropOrNil(props, "validated"),
		Source:         entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:  entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:      utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:      utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:      utils.GetTimePropOrEpochStart(props, "updatedAt"),
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

// TODO use nil
func MapDbNodeToOpportunityEntity(node *dbtype.Node) *entity.OpportunityEntity {
	if node == nil {
		return &entity.OpportunityEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	opportunity := entity.OpportunityEntity{
		Id:                utils.GetStringPropOrEmpty(props, "id"),
		Name:              utils.GetStringPropOrEmpty(props, "name"),
		Amount:            utils.GetFloatPropOrZero(props, string(entity.OpportunityPropertyAmount)),
		MaxAmount:         utils.GetFloatPropOrZero(props, string(entity.OpportunityPropertyMaxAmount)),
		InternalType:      enum.DecodeOpportunityInternalType(utils.GetStringPropOrEmpty(props, "internalType")),
		ExternalType:      utils.GetStringPropOrEmpty(props, "externalType"),
		InternalStage:     enum.DecodeOpportunityInternalStage(utils.GetStringPropOrEmpty(props, "internalStage")),
		ExternalStage:     utils.GetStringPropOrEmpty(props, "externalStage"),
		EstimatedClosedAt: utils.GetTimePropOrNil(props, "estimatedClosedAt"),
		ClosedAt:          utils.GetTimePropOrNil(props, "closedAt"),
		GeneralNotes:      utils.GetStringPropOrEmpty(props, "generalNotes"),
		NextSteps:         utils.GetStringPropOrEmpty(props, string(entity.OpportunityPropertyNextSteps)),
		Comments:          utils.GetStringPropOrEmpty(props, "comments"),
		CreatedAt:         utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:         utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:         utils.GetStringPropOrEmpty(props, "appSource"),
		Source:            entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:     entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		OwnerUserId:       utils.GetStringPropOrEmpty(props, "ownerUserId"),
		Currency:          enum.DecodeCurrency(utils.GetStringPropOrEmpty(props, string(entity.OpportunityPropertyCurrency))),
		LikelihoodRate:    utils.GetInt64PropOrDefault(props, string(entity.OpportunityPropertyLikelihoodRate), 0),
		StageUpdatedAt:    utils.GetTimePropOrNil(props, string(entity.OpportunityPropertyStageUpdatedAt)),
		RenewalDetails: entity.RenewalDetails{
			RenewedAt:              utils.GetTimePropOrNil(props, "renewedAt"),
			RenewalLikelihood:      enum.DecodeRenewalLikelihood(utils.GetStringPropOrEmpty(props, "renewalLikelihood")),
			RenewalUpdatedByUserId: utils.GetStringPropOrEmpty(props, "renewalUpdatedByUserId"),
			RenewalUpdatedByUserAt: utils.GetTimePropOrNil(props, "renewalUpdatedByUserAt"),
			RenewalApproved:        utils.GetBoolPropOrFalse(props, "renewalApproved"),
			RenewalAdjustedRate:    utils.GetInt64PropOrDefault(props, "renewalAdjustedRate", 100),
		},
		InternalFields: entity.OpportunityInternalFields{
			RolloutRenewalRequestedAt: utils.GetTimePropOrNil(props, "techRolloutRenewalRequestedAt"),
		},
	}
	return &opportunity
}

func MapDbNodeToStateEntity(node dbtype.Node) *entity.StateEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.StateEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		Code:      utils.GetStringPropOrEmpty(props, "code"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
	return &result
}

func MapDbNodeToPageView(node *dbtype.Node) *entity.PageViewEntity {
	if node == nil {
		return &entity.PageViewEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	pageViewAction := entity.PageViewEntity{
		Id:             utils.GetStringPropOrEmpty(props, "id"),
		Application:    utils.GetStringPropOrEmpty(props, "application"),
		TrackerName:    utils.GetStringPropOrEmpty(props, "trackerName"),
		SessionId:      utils.GetStringPropOrEmpty(props, "sessionId"),
		PageUrl:        utils.GetStringPropOrEmpty(props, "pageUrl"),
		PageTitle:      utils.GetStringPropOrEmpty(props, "pageTitle"),
		OrderInSession: utils.GetInt64PropOrZero(props, "orderInSession"),
		EngagedTime:    utils.GetInt64PropOrZero(props, "engagedTime"),
		StartedAt:      utils.GetTimePropOrNow(props, "startedAt"),
		EndedAt:        utils.GetTimePropOrNow(props, "endedAt"),
		Source:         entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:  entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:      utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &pageViewAction
}

func MapDbNodeToLogEntryEntity(node *dbtype.Node) *entity.LogEntryEntity {
	if node == nil {
		return &entity.LogEntryEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	logEntry := entity.LogEntryEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		StartedAt:     utils.GetTimePropOrEpochStart(props, "startedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		EventStoreAggregate: entity.EventStoreAggregate{
			AggregateVersion: utils.GetInt64PropOrNil(props, "aggregateVersion"),
		},
	}
	return &logEntry
}

func MapDbNodeToMeetingEntity(node *dbtype.Node) *entity.MeetingEntity {
	if node == nil {
		return &entity.MeetingEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	status := enum.DecodeMeetingStatus(utils.GetStringPropOrEmpty(props, "status"))
	meetingEntity := entity.MeetingEntity{
		Id:                 utils.GetStringPropOrEmpty(props, "id"),
		Name:               utils.GetStringPropOrNil(props, "name"),
		ConferenceUrl:      utils.GetStringPropOrNil(props, "conferenceUrl"),
		MeetingExternalUrl: utils.GetStringPropOrNil(props, "meetingExternalUrl"),
		Agenda:             utils.GetStringPropOrNil(props, "agenda"),
		AgendaContentType:  utils.GetStringPropOrNil(props, "agendaContentType"),
		UpdatedAt:          utils.GetTimePropOrNow(props, "updatedAt"),
		StartedAt:          utils.GetTimePropOrNil(props, "startedAt"),
		EndedAt:            utils.GetTimePropOrNil(props, "endedAt"),
		Recording:          utils.GetStringPropOrNil(props, "recording"),
		AppSource:          utils.GetStringPropOrEmpty(props, "appSource"),
		Source:             entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:      entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		Status:             &status,
	}
	if props["createdAt"] != nil {
		meetingEntity.CreatedAt = utils.GetTimePropOrNow(props, "createdAt")
	} else if props["startedAt"] != nil {
		meetingEntity.CreatedAt = utils.GetTimePropOrNow(props, "startedAt")
	} else {
		meetingEntity.CreatedAt = utils.Now()
	}

	return &meetingEntity
}

func MapDbNodeToActionEntity(node *dbtype.Node) *entity.ActionEntity {
	if node == nil {
		return &entity.ActionEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	action := entity.ActionEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Type:          enum.GetActionType(utils.GetStringPropOrEmpty(props, "type")),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		Metadata:      utils.GetStringPropOrEmpty(props, "metadata"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &action
}

func MapDbNodeToNoteEntity(node *dbtype.Node) *entity.NoteEntity {
	if node == nil {
		return &entity.NoteEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	note := entity.NoteEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &note
}

func MapDbNodeToInteractionEventEntity(node *neo4j.Node) *entity.InteractionEventEntity {
	if node == nil {
		return &entity.InteractionEventEntity{}
	}
	props := utils.GetPropsFromNode(*node)

	return MapDbPropsToInteractionEventEntity(props)
}

func MapDbPropsToInteractionEventEntity(props map[string]interface{}) *entity.InteractionEventEntity {
	interactionEventEntity := entity.InteractionEventEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Identifier:    utils.GetStringPropOrEmpty(props, "identifier"),
		Channel:       utils.GetStringPropOrEmpty(props, "channel"),
		ChannelData:   utils.GetStringPropOrEmpty(props, "channelData"),
		EventType:     utils.GetStringPropOrEmpty(props, "eventType"),
		Hide:          utils.GetBoolPropOrFalse(props, "hide"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &interactionEventEntity
}

func MapDbNodeToInteractionSessionEntity(node *dbtype.Node) *entity.InteractionSessionEntity {
	if node == nil {
		return &entity.InteractionSessionEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	interactionSession := entity.InteractionSessionEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Channel:       utils.GetStringPropOrEmpty(props, "channel"),
		ChannelData:   utils.GetStringPropOrEmpty(props, "channelData"),
		Identifier:    utils.GetStringPropOrEmpty(props, "identifier"),
		Type:          utils.GetStringPropOrEmpty(props, "type"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		Status:        utils.GetStringPropOrEmpty(props, "status"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &interactionSession
}

func MapDbNodeToContactEntity(dbNode *dbtype.Node) *entity.ContactEntity {
	props := utils.GetPropsFromNode(*dbNode)
	contact := entity.ContactEntity{
		Id:              utils.GetStringPropOrEmpty(props, "id"),
		FirstName:       utils.GetStringPropOrEmpty(props, string(entity.ContactPropertyFirstName)),
		LastName:        utils.GetStringPropOrEmpty(props, string(entity.ContactPropertyLastName)),
		Name:            utils.GetStringPropOrEmpty(props, string(entity.ContactPropertyName)),
		Description:     utils.GetStringPropOrEmpty(props, string(entity.ContactPropertyDescription)),
		Timezone:        utils.GetStringPropOrEmpty(props, string(entity.ContactPropertyTimezone)),
		ProfilePhotoUrl: utils.GetStringPropOrEmpty(props, string(entity.ContactPropertyProfilePhotoUrl)),
		Username:        utils.GetStringPropOrEmpty(props, string(entity.ContactPropertyUsername)),
		Prefix:          utils.GetStringPropOrEmpty(props, string(entity.ContactPropertyPrefix)),
		Hide:            utils.GetBoolPropOrFalse(props, string(entity.ContactPropertyHide)),
		CreatedAt:       utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:       utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:          entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		AppSource:       utils.GetStringPropOrEmpty(props, "appSource"),
		EventStoreAggregate: entity.EventStoreAggregate{
			AggregateVersion: utils.GetInt64PropOrNil(props, "aggregateVersion"),
		},
		ContactInternalFields: entity.ContactInternalFields{
			FindWorkEmailWithBetterContactRequestedId: utils.GetStringPropOrNil(props, string(entity.ContactPropertyFindWorkEmailWithBetterContactRequestedId)),
			FindWorkEmailWithBetterContactRequestedAt: utils.GetTimePropOrNil(props, string(entity.ContactPropertyFindWorkEmailWithBetterContactRequestedAt)),
			FindWorkEmailWithBetterContactCompletedAt: utils.GetTimePropOrNil(props, string(entity.ContactPropertyFindWorkEmailWithBetterContactCompletedAt)),
			HiddenAt: utils.GetTimePropOrNil(props, string(entity.ContactPropertyHiddenAt)),
		},
		EnrichDetails: entity.ContactEnrichDetails{
			EnrichRequestedAt:         utils.GetTimePropOrNil(props, string(entity.ContactPropertyEnrichRequestedAt)),
			EnrichedAt:                utils.GetTimePropOrNil(props, string(entity.ContactPropertyEnrichedAt)),
			EnrichFailedAt:            utils.GetTimePropOrNil(props, string(entity.ContactPropertyEnrichFailedAt)),
			EnrichAttempts:            utils.GetInt64PropOrZero(props, string(entity.ContactPropertyEnrichAttempts)),
			BettercontactFoundEmailAt: utils.GetTimePropOrNil(props, string(entity.ContactPropertyBettercontactFoundEmailAt)),
			EnrichedScrapinRecordId:   utils.GetStringPropOrEmpty(props, string(entity.ContactPropertyEnrichedScrapinRecordId)),
		},
	}
	return &contact
}

func MapDbNodeToTimelineEvent(dbNode *dbtype.Node) entity.TimelineEvent {
	if slices.Contains(dbNode.Labels, model.NodeLabelPageView) {
		return MapDbNodeToPageView(dbNode)
	} else if slices.Contains(dbNode.Labels, model.NodeLabelInteractionSession) {
		return MapDbNodeToInteractionSessionEntity(dbNode)
	} else if slices.Contains(dbNode.Labels, model.NodeLabelIssue) {
		return MapDbNodeToIssueEntity(dbNode)
	} else if slices.Contains(dbNode.Labels, model.NodeLabelNote) {
		return MapDbNodeToNoteEntity(dbNode)
	} else if slices.Contains(dbNode.Labels, model.NodeLabelInteractionEvent) {
		return MapDbNodeToInteractionEventEntity(dbNode)
	} else if slices.Contains(dbNode.Labels, model.NodeLabelMeeting) {
		return MapDbNodeToMeetingEntity(dbNode)
	} else if slices.Contains(dbNode.Labels, model.NodeLabelAction) {
		return MapDbNodeToActionEntity(dbNode)
	} else if slices.Contains(dbNode.Labels, model.NodeLabelLogEntry) {
		return MapDbNodeToLogEntryEntity(dbNode)
	}
	return nil
}

func MapDbNodeToDomainEntity(node *dbtype.Node) *entity.DomainEntity {
	if node == nil {
		return &entity.DomainEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	domain := entity.DomainEntity{
		CreatedAt:     utils.GetTimePropOrEpochStart(props, string(entity.DomainPropertyCreatedAt)),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, string(entity.DomainPropertyUpdatedAt)),
		AppSource:     utils.GetStringPropOrEmpty(props, string(entity.DomainPropertyAppSource)),
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, string(entity.DomainPropertySource))),
		Domain:        utils.GetStringPropOrEmpty(props, string(entity.DomainPropertyDomain)),
		IsPrimary:     utils.GetBoolPropOrNil(props, string(entity.DomainPropertyIsPrimary)),
		PrimaryDomain: utils.GetStringPropOrEmpty(props, string(entity.DomainPropertyPrimaryDomain)),
		InternalFields: entity.DomainInternalFields{
			PrimaryDomainCheckRequestedAt: utils.GetTimePropOrNil(props, string(entity.DomainPropertyPrimaryDomainCheckRequestedAt)),
		},
	}
	return &domain
}

func MapDbNodeToLocationEntity(node *dbtype.Node) *entity.LocationEntity {
	if node == nil {
		return &entity.LocationEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	location := entity.LocationEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyName)),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Country:       utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyCountry)),
		CountryCodeA2: utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyCountryCodeA2)),
		CountryCodeA3: utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyCountryCodeA3)),
		Region:        utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyRegion)),
		Locality:      utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyLocality)),
		Address:       utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyAddress)),
		Address2:      utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyAddress2)),
		Zip:           utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyZip)),
		AddressType:   utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyAddressType)),
		HouseNumber:   utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyHouseNumber)),
		PostalCode:    utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyPostalCode)),
		PlusFour:      utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyPlusFour)),
		Commercial:    utils.GetBoolPropOrFalse(props, string(entity.LocationPropertyCommercial)),
		Predirection:  utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyPredirection)),
		District:      utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyDistrict)),
		Street:        utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyStreet)),
		RawAddress:    utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyRawAddress)),
		Latitude:      utils.GetFloatPropOrNil(props, string(entity.LocationPropertyLatitude)),
		Longitude:     utils.GetFloatPropOrNil(props, string(entity.LocationPropertyLongitude)),
		TimeZone:      utils.GetStringPropOrEmpty(props, string(entity.LocationPropertyTimeZone)),
		UtcOffset:     utils.GetFloatPropOrNil(props, string(entity.LocationPropertyUtcOffset)),
		Source:        entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &location
}

func MapDbNodeToFlowEntity(node *dbtype.Node) *entity.FlowEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)
	domain := entity.FlowEntity{
		Id:           utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:    utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:    utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Name:         utils.GetStringPropOrEmpty(props, "name"),
		Nodes:        utils.GetStringPropOrEmpty(props, "nodes"),
		Edges:        utils.GetStringPropOrEmpty(props, "edges"),
		Status:       entity.GetFlowStatus(utils.GetStringPropOrEmpty(props, "status")),
		Total:        utils.GetInt64PropOrZero(props, "total"),
		OnHold:       utils.GetInt64PropOrZero(props, "onHold"),
		Ready:        utils.GetInt64PropOrZero(props, "ready"),
		Scheduled:    utils.GetInt64PropOrZero(props, "scheduled"),
		InProgress:   utils.GetInt64PropOrZero(props, "inProgress"),
		Completed:    utils.GetInt64PropOrZero(props, "completed"),
		GoalAchieved: utils.GetInt64PropOrZero(props, "goalAchieved"),
	}
	return &domain
}

func MapDbNodeToFlowParticipantEntity(node *dbtype.Node) *entity.FlowParticipantEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)
	e := entity.FlowParticipantEntity{
		Id:         utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:  utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:  utils.GetTimePropOrEpochStart(props, "updatedAt"),
		EntityId:   utils.GetStringPropOrEmpty(props, "entityId"),
		EntityType: model.DecodeEntityType(utils.GetStringPropOrEmpty(props, "entityType")),
		Status:     entity.GetFlowContactStatus(utils.GetStringPropOrEmpty(props, "status")),
	}
	return &e
}

func MapDbNodeToFlowSenderEntity(node *dbtype.Node) *entity.FlowSenderEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)
	e := entity.FlowSenderEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
		UserId:    utils.GetStringPropOrNil(props, "userId"),
	}
	return &e
}

func MapDbNodeToFlowActionEntity(node *dbtype.Node) *entity.FlowActionEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)

	e := entity.FlowActionEntity{
		Id:         utils.GetStringPropOrEmpty(props, "id"),
		ExternalId: utils.GetStringPropOrEmpty(props, "externalId"),
		CreatedAt:  utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:  utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Json:       utils.GetStringPropOrEmpty(props, "json"),
		Type:       utils.GetStringPropOrEmpty(props, "type"),
	}

	e.Data.Action = entity.GetFlowActionType(utils.GetStringPropOrEmpty(props, "action"))

	e.Data.WaitBefore = utils.GetInt64PropOrZero(props, "waitBefore")

	e.Data.Entity = utils.GetStringPropOrNil(props, "data_entity")
	e.Data.TriggerType = utils.GetStringPropOrNil(props, "data_triggerType")
	e.Data.Subject = utils.GetStringPropOrNil(props, "data_subject")
	e.Data.BodyTemplate = utils.GetStringPropOrNil(props, "data_bodyTemplate")
	e.Data.MessageTemplate = utils.GetStringPropOrNil(props, "data_messageTemplate")

	return &e
}

func MapDbNodeToFlowExecutionSettingsEntity(node *dbtype.Node) *entity.FlowExecutionSettingsEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)
	e := entity.FlowExecutionSettingsEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
		FlowId:    utils.GetStringPropOrEmpty(props, "flowId"),
		EntityId:  utils.GetStringPropOrEmpty(props, "entityId"),
		Mailbox:   utils.GetStringPropOrNil(props, "mailbox"),
		UserId:    utils.GetStringPropOrNil(props, "userId"),
	}
	return &e
}

func MapDbNodeToFlowActionExecutionEntity(node *dbtype.Node) *entity.FlowActionExecutionEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)
	e := entity.FlowActionExecutionEntity{
		Id:              utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:       utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:       utils.GetTimePropOrEpochStart(props, "updatedAt"),
		FlowId:          utils.GetStringPropOrEmpty(props, "flowId"),
		EntityId:        utils.GetStringPropOrEmpty(props, "entityId"),
		EntityType:      model.DecodeEntityType(utils.GetStringPropOrEmpty(props, "entityType")),
		ActionId:        utils.GetStringPropOrEmpty(props, "actionId"),
		ScheduledAt:     utils.GetTimePropOrNow(props, "scheduledAt"),
		ExecutedAt:      utils.GetTimePropOrNil(props, "executedAt"),
		StatusUpdatedAt: utils.GetTimePropOrNow(props, "statusUpdatedAt"),
		Status:          entity.GetFlowActionExecutionStatus(utils.GetStringPropOrEmpty(props, "status")),
		Error:           utils.GetStringPropOrNil(props, "error"),

		Mailbox: utils.GetStringPropOrNil(props, "mailbox"),
		UserId:  utils.GetStringPropOrNil(props, "userId"),
	}
	return &e
}

func MapDbNodeToCustomFieldTemplateEntity(node *dbtype.Node) *entity.CustomFieldTemplateEntity {
	if node == nil {
		return &entity.CustomFieldTemplateEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	customFieldTemplateEntity := entity.CustomFieldTemplateEntity{
		Id:          utils.GetStringPropOrEmpty(props, string(entity.CustomFieldTemplatePropertyId)),
		Name:        utils.GetStringPropOrEmpty(props, string(entity.CustomFieldTemplatePropertyName)),
		EntityType:  model.DecodeEntityType(utils.GetStringPropOrEmpty(props, string(entity.CustomFieldTemplatePropertyEntityType))),
		Type:        utils.GetStringPropOrEmpty(props, string(entity.CustomFieldTemplatePropertyType)),
		ValidValues: utils.GetListStringPropOrEmpty(props, string(entity.CustomFieldTemplatePropertyValidValues)),
		Order:       utils.GetInt64PropOrNil(props, string(entity.CustomFieldTemplatePropertyOrder)),
		Required:    utils.GetBoolPropOrNil(props, string(entity.CustomFieldTemplatePropertyRequired)),
		Length:      utils.GetInt64PropOrNil(props, string(entity.CustomFieldTemplatePropertyLength)),
		Min:         utils.GetInt64PropOrNil(props, string(entity.CustomFieldTemplatePropertyMin)),
		Max:         utils.GetInt64PropOrNil(props, string(entity.CustomFieldTemplatePropertyMax)),
		CreatedAt:   utils.GetTimePropOrEpochStart(props, string(entity.CustomFieldTemplatePropertyCreatedAt)),
		UpdatedAt:   utils.GetTimePropOrEpochStart(props, string(entity.CustomFieldTemplatePropertyUpdatedAt)),
	}
	return &customFieldTemplateEntity
}
