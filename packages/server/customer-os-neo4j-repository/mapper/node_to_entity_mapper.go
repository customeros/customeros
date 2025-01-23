package mapper

import (
	commonenum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"golang.org/x/exp/slices"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
)

func MapDbNodeToJobRoleEntity(dbNode *dbtype.Node) *neo4j_entity.JobRoleEntity {
	if dbNode == nil {
		return &neo4j_entity.JobRoleEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	result := neo4j_entity.JobRoleEntity{
		Id:          utils.GetStringPropOrEmpty(props, "id"),
		JobTitle:    utils.GetStringPropOrEmpty(props, "jobTitle"),
		Description: utils.GetStringPropOrNil(props, "description"),
		Company:     utils.GetStringPropOrNil(props, "company"),
		Primary:     utils.GetBoolPropOrFalse(props, "primary"),
		Source:      neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		AppSource:   utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:   utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:   utils.GetTimePropOrEpochStart(props, "updatedAt"),
		StartedAt:   utils.GetTimePropOrNil(props, "startedAt"),
		EndedAt:     utils.GetTimePropOrNil(props, "endedAt"),
	}
	return &result
}

func MapDbNodeToAttachmentEntity(dbNode *dbtype.Node) *neo4j_entity.AttachmentEntity {
	if dbNode == nil {
		return &neo4j_entity.AttachmentEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	createdAt := utils.GetTimePropOrEpochStart(props, "createdAt")
	attachmentEntity := neo4j_entity.AttachmentEntity{
		Id:                 utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:          &createdAt,
		FileName:           utils.GetStringPropOrEmpty(props, "fileName"),
		MimeType:           utils.GetStringPropOrEmpty(props, "mimeType"),
		CdnUrl:             utils.GetStringPropOrEmpty(props, "cdnUrl"),
		PublicUrl:          utils.GetStringPropOrEmpty(props, string(neo4j_entity.AttachmentPropertyPublicUrl)),
		PublicUrlExpiresAt: utils.GetTimePropOrNil(props, string(neo4j_entity.AttachmentPropertyPublicUrlExpiresAt)),
		BasePath:           utils.GetStringPropOrEmpty(props, "basePath"),
		Size:               utils.GetInt64PropOrZero(props, "size"),
		AppSource:          utils.GetStringPropOrEmpty(props, "appSource"),
		Source:             neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:      neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &attachmentEntity
}

func MapDbNodeToWorkspaceEntity(dbNode *dbtype.Node) *neo4j_entity.WorkspaceEntity {
	if dbNode == nil {
		return &neo4j_entity.WorkspaceEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	workspace := neo4j_entity.WorkspaceEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "domain"),
		Provider:      utils.GetStringPropOrEmpty(props, "provider"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &workspace
}

func MapDbNodeToPlayerEntity(node *neo4j.Node) *neo4j_entity.PlayerEntity {
	if node == nil {
		return &neo4j_entity.PlayerEntity{}
	}
	props := utils.GetPropsFromNode(*node)

	return &neo4j_entity.PlayerEntity{
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

func MapDbNodeToInvoiceEntity(dbNode *dbtype.Node) *neo4j_entity.InvoiceEntity {
	if dbNode == nil {
		return &neo4j_entity.InvoiceEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	invoiceEntity := neo4j_entity.InvoiceEntity{
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
		Source:               neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:        neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:            utils.GetStringPropOrEmpty(props, "appSource"),
		Status:               enum.DecodeInvoiceStatus(utils.GetStringPropOrEmpty(props, "status")),
		Note:                 utils.GetStringPropOrEmpty(props, "note"),
		Customer: neo4j_entity.InvoiceCustomer{
			Name:         utils.GetStringPropOrEmpty(props, "customerName"),
			Email:        utils.GetStringPropOrEmpty(props, "customerEmail"),
			AddressLine1: utils.GetStringPropOrEmpty(props, "customerAddressLine1"),
			AddressLine2: utils.GetStringPropOrEmpty(props, "customerAddressLine2"),
			Zip:          utils.GetStringPropOrEmpty(props, "customerAddressZip"),
			Locality:     utils.GetStringPropOrEmpty(props, "customerAddressLocality"),
			Country:      utils.GetStringPropOrEmpty(props, "customerAddressCountry"),
			Region:       utils.GetStringPropOrEmpty(props, "customerAddressRegion"),
		},
		Provider: neo4j_entity.InvoiceProvider{
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
		PaymentDetails: neo4j_entity.PaymentDetails{
			PaymentLink:           utils.GetStringPropOrEmpty(props, string(neo4j_entity.InvoicePropertyPaymentLink)),
			PaymentLinkValidUntil: utils.GetTimePropOrNil(props, string(neo4j_entity.InvoicePropertyPaymentLinkValidUntil)),
		},
		InvoiceInternalFields: neo4j_entity.InvoiceInternalFields{
			InvoiceFinalizedSentAt:               utils.GetTimePropOrNil(props, string(neo4j_entity.InvoicePropertyInvoiceFinalizedEventSentAt)),
			InvoiceFinalizedWebhookProcessedAt:   utils.GetTimePropOrNil(props, string(neo4j_entity.InvoicePropertyFinalizedWebhookProcessedAt)),
			InvoicePaidWebhookProcessedAt:        utils.GetTimePropOrNil(props, string(neo4j_entity.InvoicePropertyPaidWebhookProcessedAt)),
			PaymentLinkRequestedAt:               utils.GetTimePropOrNil(props, "techPaymentLinkRequestedAt"),
			PayInvoiceNotificationRequestedAt:    utils.GetTimePropOrNil(props, "techPayNotificationRequestedAt"),
			PayInvoiceNotificationSentAt:         utils.GetTimePropOrNil(props, "techPayInvoiceNotificationSentAt"),
			RemindInvoiceNotificationRequestedAt: utils.GetTimePropOrNil(props, string(neo4j_entity.InvoicePropertyRemindInvoiceNotificationRequestedAt)),
			LastRemindInvoiceNotificationSentAt:  utils.GetTimePropOrNil(props, string(neo4j_entity.InvoicePropertyLastRemindInvoiceNotificationSentAt)),
			PaidInvoiceNotificationSentAt:        utils.GetTimePropOrNil(props, "techPaidInvoiceNotificationSentAt"),
			VoidInvoiceNotificationSentAt:        utils.GetTimePropOrNil(props, "techVoidInvoiceNotificationSentAt"),
		},
		EventStoreAggregate: neo4j_entity.EventStoreAggregate{
			AggregateVersion: utils.GetInt64PropOrNil(props, "aggregateVersion"),
		},
	}
	return &invoiceEntity
}

func MapDbNodeToInvoiceLineEntity(dbNode *dbtype.Node) *neo4j_entity.InvoiceLineEntity {
	if dbNode == nil {
		return &neo4j_entity.InvoiceLineEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	invoiceLineEntity := neo4j_entity.InvoiceLineEntity{
		Id:                      utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:               utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:               utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Name:                    utils.GetStringPropOrEmpty(props, "name"),
		Price:                   utils.GetFloatPropOrZero(props, "price"),
		Quantity:                utils.GetInt64PropOrZero(props, "quantity"),
		Amount:                  utils.GetFloatPropOrZero(props, "amount"),
		Vat:                     utils.GetFloatPropOrZero(props, "vat"),
		TotalAmount:             utils.GetFloatPropOrZero(props, "totalAmount"),
		Source:                  neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:           neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:               utils.GetStringPropOrEmpty(props, "appSource"),
		ServiceLineItemId:       utils.GetStringPropOrEmpty(props, "serviceLineItemId"),
		ServiceLineItemParentId: utils.GetStringPropOrEmpty(props, "serviceLineItemParentId"),
		BilledType:              enum.DecodeBilledType(utils.GetStringPropOrEmpty(props, "billedType")),
	}
	return &invoiceLineEntity
}

func MapDbNodeToUserEntity(dbNode *dbtype.Node) *neo4j_entity.UserEntity {
	if dbNode == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*dbNode)
	userEntity := neo4j_entity.UserEntity{
		Id:              utils.GetStringPropOrEmpty(props, "id"),
		FirstName:       utils.GetStringPropOrEmpty(props, "firstName"),
		LastName:        utils.GetStringPropOrEmpty(props, "lastName"),
		Name:            utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:       utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:       utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:          neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:   neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:       utils.GetStringPropOrEmpty(props, "appSource"),
		Roles:           utils.GetListStringPropOrEmpty(props, "roles"),
		Internal:        utils.GetBoolPropOrFalse(props, "internal"),
		Test:            utils.GetBoolPropOrFalse(props, "test"),
		Bot:             utils.GetBoolPropOrFalse(props, "bot"),
		ProfilePhotoUrl: utils.GetStringPropOrEmpty(props, "profilePhotoUrl"),
		Timezone:        utils.GetStringPropOrEmpty(props, "timezone"),
		LastLogin:       utils.GetTimePropOrNil(props, string(neo4j_entity.UserPropertyLastLogin)),
		FirstLogin:      utils.GetTimePropOrNil(props, string(neo4j_entity.UserPropertyFirstLogin)),
		OnboardingDetails: neo4j_entity.UserOnboardingDetails{
			ShowOnboardingPage:               utils.GetBoolPropOrTrue(props, string(neo4j_entity.UserPropertyShowOnboardingPage)),
			OnboardingInboundStepCompleted:   utils.GetBoolPropOrFalse(props, string(neo4j_entity.UserPropertyOnboardingInboundStepCompleted)),
			OnboardingOutboundStepCompleted:  utils.GetBoolPropOrFalse(props, string(neo4j_entity.UserPropertyOnboardingOutboundStepCompleted)),
			OnboardingCrmStepCompleted:       utils.GetBoolPropOrFalse(props, string(neo4j_entity.UserPropertyOnboardingCrmStepCompleted)),
			OnboardingMailstackStepCompleted: utils.GetBoolPropOrFalse(props, string(neo4j_entity.UserPropertyOnboardingMailstackStepCompleted)),
		},
	}
	return &userEntity
}

func MapDbNodeToOrganizationEntity(dbNode *dbtype.Node) *neo4j_entity.OrganizationEntity {
	if dbNode == nil {
		return &neo4j_entity.OrganizationEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	organizationEntity := neo4j_entity.OrganizationEntity{
		ID:                   utils.GetStringPropOrEmpty(props, "id"),
		CustomerOsId:         utils.GetStringPropOrEmpty(props, "customerOsId"),
		ReferenceId:          utils.GetStringPropOrEmpty(props, "referenceId"),
		Name:                 utils.GetStringPropOrEmpty(props, "name"),
		Description:          utils.GetStringPropOrEmpty(props, "description"),
		Website:              utils.GetStringPropOrEmpty(props, "website"),
		Industry:             utils.GetStringPropOrEmpty(props, string(neo4j_entity.OrganizationPropertyIndustry)),
		LastFundingRound:     utils.GetStringPropOrEmpty(props, "lastFundingRound"),
		LastFundingAmount:    utils.GetStringPropOrEmpty(props, "lastFundingAmount"),
		Note:                 utils.GetStringPropOrEmpty(props, "note"),
		IsPublic:             utils.GetBoolPropOrFalse(props, string(neo4j_entity.OrganizationPropertyIsPublic)),
		Hide:                 utils.GetBoolPropOrFalse(props, string(neo4j_entity.OrganizationPropertyHide)),
		Employees:            utils.GetInt64PropOrZero(props, string(neo4j_entity.OrganizationPropertyEmployees)),
		Market:               utils.GetStringPropOrEmpty(props, "market"),
		CreatedAt:            utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:            utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:               neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		LastTouchpointAt:     utils.GetTimePropOrNil(props, "lastTouchpointAt"),
		LastTouchpointId:     utils.GetStringPropOrNil(props, "lastTouchpointId"),
		LastTouchpointType:   utils.GetStringPropOrNil(props, "lastTouchpointType"),
		YearFounded:          utils.GetInt64PropOrNil(props, string(neo4j_entity.OrganizationPropertyYearFounded)),
		Headquarters:         utils.GetStringPropOrEmpty(props, "headquarters"),
		EmployeeGrowthRate:   utils.GetStringPropOrEmpty(props, "employeeGrowthRate"),
		SlackChannelId:       utils.GetStringPropOrEmpty(props, "slackChannelId"),
		LogoUrl:              utils.GetStringPropOrEmpty(props, "logoUrl"),
		IconUrl:              utils.GetStringPropOrEmpty(props, "iconUrl"),
		Relationship:         enum.DecodeOrganizationRelationship(utils.GetStringPropOrEmpty(props, "relationship")),
		Stage:                enum.DecodeOrganizationStage(utils.GetStringPropOrEmpty(props, "stage")),
		StageUpdatedAt:       utils.GetTimePropOrNil(props, "stageUpdatedAt"),
		LeadSource:           utils.GetStringPropOrEmpty(props, "leadSource"),
		IcpFit:               utils.GetBoolPropOrFalse(props, string(neo4j_entity.OrganizationPropertyIcpFit)),
		WrongIndustry:        utils.GetBoolPropOrFalse(props, string(neo4j_entity.OrganizationPropertyWrongIndustry)),
		QuickbooksCustomerId: utils.GetStringPropOrEmpty(props, string(neo4j_entity.OrganizationPropertyQuickbooksCustomerId)),
		RenewalSummary: neo4j_entity.RenewalSummary{
			ArrForecast:            utils.GetFloatPropOrNil(props, "renewalForecastArr"),
			MaxArrForecast:         utils.GetFloatPropOrNil(props, "renewalForecastMaxArr"),
			RenewalLikelihood:      utils.GetStringPropOrEmpty(props, string(neo4j_entity.OrganizationPropertyRenewalLikelihood)),
			RenewalLikelihoodOrder: utils.GetInt64PropOrNil(props, "derivedRenewalLikelihoodOrder"),
			NextRenewalAt:          utils.GetTimePropOrNil(props, "derivedNextRenewalAt"),
		},
		DerivedData: neo4j_entity.DerivedData{
			ChurnedAt:    utils.GetTimePropOrNil(props, "derivedChurnedAt"),
			Ltv:          utils.GetFloatPropOrZero(props, "derivedLtv"),
			LtvCurrency:  enum.DecodeCurrency(utils.GetStringPropOrEmpty(props, "derivedLtvCurrency")),
			ContactCount: utils.GetInt64PropOrZero(props, string(neo4j_entity.OrganizationPropertyContactCount)),
		},
		OnboardingDetails: neo4j_entity.OnboardingDetails{
			Status:       utils.GetStringPropOrEmpty(props, "onboardingStatus"),
			SortingOrder: utils.GetInt64PropOrNil(props, "onboardingStatusOrder"),
			UpdatedAt:    utils.GetTimePropOrNil(props, "onboardingUpdatedAt"),
			Comments:     utils.GetStringPropOrEmpty(props, "onboardingComments"),
		},
		EnrichDetails: neo4j_entity.OrganizationEnrichDetails{
			EnrichRequestedAt: utils.GetTimePropOrNil(props, string(neo4j_entity.OrganizationPropertyEnrichRequestedAt)),
			EnrichedAt:        utils.GetTimePropOrNil(props, string(neo4j_entity.OrganizationPropertyEnrichedAt)),
			EnrichFailedAt:    utils.GetTimePropOrNil(props, string(neo4j_entity.OrganizationPropertyEnrichFailedAt)),
			EnrichAttempts:    utils.GetInt64PropOrZero(props, string(neo4j_entity.OrganizationPropertyEnrichAttempts)),
			EnrichSource:      enum.DecodeDomainEnrichSource(utils.GetStringPropOrEmpty(props, "enrichSource")),
			EnrichDomain:      utils.GetStringPropOrEmpty(props, "enrichDomain"),
		},
		OrganizationInternalFields: neo4j_entity.OrganizationInternalFields{
			DomainCheckedAt: utils.GetTimePropOrNil(props, string(neo4j_entity.OrganizationPropertyDomainCheckedAt)),
			HiddenAt:        utils.GetTimePropOrNil(props, string(neo4j_entity.OrganizationPropertyHiddenAt)),
		},
		EventStoreAggregate: neo4j_entity.EventStoreAggregate{
			AggregateVersion: utils.GetInt64PropOrNil(props, "aggregateVersion"),
		},
	}
	return &organizationEntity
}

func MapDbNodeToBillingProfileEntity(dbNode *dbtype.Node) *neo4j_entity.BillingProfileEntity {
	if dbNode == nil {
		return &neo4j_entity.BillingProfileEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	billingProfileEntity := neo4j_entity.BillingProfileEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		LegalName:     utils.GetStringPropOrEmpty(props, "legalName"),
		TaxId:         utils.GetStringPropOrEmpty(props, "taxId"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &billingProfileEntity
}

func MapDbNodeToTenantEntity(dbNode *dbtype.Node) *neo4j_entity.TenantEntity {
	if dbNode == nil {
		return &neo4j_entity.TenantEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	tenant := neo4j_entity.TenantEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource: utils.GetStringPropOrEmpty(props, "appSource"),
		Source:    neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		Active:    utils.GetBoolPropOrTrue(props, "active"),
	}
	return &tenant
}

func MapDbNodeToTenantSettingsEntity(dbNode *dbtype.Node) *neo4j_entity.TenantSettingsEntity {
	if dbNode == nil {
		return &neo4j_entity.TenantSettingsEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	tenantSettingsEntity := neo4j_entity.TenantSettingsEntity{
		Id:                       utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:                utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:                utils.GetTimePropOrEpochStart(props, "updatedAt"),
		LogoRepositoryFileId:     utils.GetStringPropOrEmpty(props, string(neo4j_entity.TenantSettingsPropertyLogoRepositoryFileId)),
		InvoicingEnabled:         utils.GetBoolPropOrFalse(props, string(neo4j_entity.TenantSettingsPropertyInvoicingEnabled)),
		InvoicingPostpaid:        utils.GetBoolPropOrFalse(props, string(neo4j_entity.TenantSettingsPropertyInvoicingPostpaid)),
		WorkspaceLogo:            utils.GetStringPropOrEmpty(props, string(neo4j_entity.TenantSettingsPropertyWorkspaceLogo)),
		WorkspaceName:            utils.GetStringPropOrEmpty(props, string(neo4j_entity.TenantSettingsPropertyWorkspaceName)),
		BaseCurrency:             enum.DecodeCurrency(utils.GetStringPropOrEmpty(props, string(neo4j_entity.TenantSettingsPropertyBaseCurrency))),
		EnrichContacts:           utils.GetBoolPropOrFalse(props, string(neo4j_entity.TenantSettingsPropertyEnrichContacts)),
		StripeCustomerPortalLink: utils.GetStringPropOrEmpty(props, string(neo4j_entity.TenantSettingsPropertyStripeCustomerPortalLink)),
		SharedSlackChannelUrl:    utils.GetStringPropOrEmpty(props, string(neo4j_entity.TenantSettingsPropertySlackChannelUrl)),
	}
	return &tenantSettingsEntity
}

func MapDbNodeToTenantBillingProfileEntity(dbNode *dbtype.Node) *neo4j_entity.TenantBillingProfileEntity {
	if dbNode == nil {
		return &neo4j_entity.TenantBillingProfileEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	tenantBillingProfile := neo4j_entity.TenantBillingProfileEntity{
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
		Source:                 neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:          neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:              utils.GetStringPropOrEmpty(props, "appSource"),
		Check:                  utils.GetBoolPropOrFalse(props, "check"),
	}
	return &tenantBillingProfile
}

func MapDbNodeToCountryEntity(dbNode *dbtype.Node) *neo4j_entity.CountryEntity {
	if dbNode == nil {
		return &neo4j_entity.CountryEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	result := neo4j_entity.CountryEntity{
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

func MapDbNodeToContractEntity(dbNode *dbtype.Node) *neo4j_entity.ContractEntity {
	if dbNode == nil {
		return &neo4j_entity.ContractEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	contract := neo4j_entity.ContractEntity{
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
		Source:                          neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:                   neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
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
		ContractInternalFields: neo4j_entity.ContractInternalFields{
			StatusRenewalRequestedAt:      utils.GetTimePropOrNil(props, "techStatusRenewalRequestedAt"),
			RolloutRenewalRequestedAt:     utils.GetTimePropOrNil(props, "techRolloutRenewalRequestedAt"),
			NextPreviewInvoiceRequestedAt: utils.GetTimePropOrNil(props, "techNextPreviewInvoiceRequestedAt"),
		},
		EventStoreAggregate: neo4j_entity.EventStoreAggregate{
			AggregateVersion: utils.GetInt64PropOrNil(props, "aggregateVersion"),
		},
	}

	return &contract
}

func MapDbNodeToServiceLineItemEntity(dbNode *dbtype.Node) *neo4j_entity.ServiceLineItemEntity {
	if dbNode == nil {
		return &neo4j_entity.ServiceLineItemEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	serviceLineItem := neo4j_entity.ServiceLineItemEntity{
		ID:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		StartedAt:     utils.GetTimePropOrEpochStart(props, "startedAt"),
		EndedAt:       utils.GetTimePropOrNil(props, "endedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		Billed:        enum.DecodeBilledType(utils.GetStringPropOrEmpty(props, "billed")),
		Price:         utils.GetFloatPropOrZero(props, "price"),
		Quantity:      utils.GetInt64PropOrZero(props, "quantity"),
		Comments:      utils.GetStringPropOrEmpty(props, "comments"),
		ParentID:      utils.GetStringPropOrEmpty(props, "parentId"),
		Canceled:      utils.GetBoolPropOrFalse(props, "isCanceled"),
		VatRate:       utils.GetFloatPropOrZero(props, "vatRate"),
		Paused:        utils.GetBoolPropOrFalse(props, string(neo4j_entity.SLIPropertyPaused)),
	}
	return &serviceLineItem
}

func MapDbNodeToTagEntity(dbNode *dbtype.Node) *neo4j_entity.TagEntity {
	if dbNode == nil {
		return &neo4j_entity.TagEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	tag := neo4j_entity.TagEntity{
		Id:         utils.GetStringPropOrEmpty(props, string(neo4j_entity.TagPropertyId)),
		Name:       utils.GetStringPropOrEmpty(props, string(neo4j_entity.TagPropertyName)),
		CreatedAt:  utils.GetTimePropOrEpochStart(props, string(neo4j_entity.TagPropertyCreatedAt)),
		UpdatedAt:  utils.GetTimePropOrEpochStart(props, string(neo4j_entity.TagPropertyUpdatedAt)),
		Source:     neo4j_entity.DataSource(utils.GetStringPropOrEmpty(props, string(neo4j_entity.TagPropertySource))),
		AppSource:  utils.GetStringPropOrEmpty(props, string(neo4j_entity.TagPropertyAppSource)),
		EntityType: model.DecodeEntityType(utils.GetStringPropOrEmpty(props, string(neo4j_entity.TagPropertyEntityType))),
		ColorCode:  utils.GetStringPropOrEmpty(props, string(neo4j_entity.TagPropertyColorCode)),
	}
	return &tag
}

func MapDbNodeToIssueEntity(dbNode *dbtype.Node) *neo4j_entity.IssueEntity {
	if dbNode == nil {
		return &neo4j_entity.IssueEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	issue := neo4j_entity.IssueEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:     utils.GetTimePropOrNow(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrNow(props, "updatedAt"),
		Subject:       utils.GetStringPropOrEmpty(props, "subject"),
		Status:        utils.GetStringPropOrEmpty(props, "status"),
		Priority:      utils.GetStringPropOrEmpty(props, "priority"),
		Description:   utils.GetStringPropOrEmpty(props, "description"),
		Source:        neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &issue
}

func MapDbNodeToCommentEntity(dbNode *dbtype.Node) *neo4j_entity.CommentEntity {
	if dbNode == nil {
		return &neo4j_entity.CommentEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	comment := neo4j_entity.CommentEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &comment
}

func MapDbNodeToSocialEntity(dbNode *dbtype.Node) *neo4j_entity.SocialEntity {
	if dbNode == nil {
		return &neo4j_entity.SocialEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	social := neo4j_entity.SocialEntity{
		Id:             utils.GetStringPropOrEmpty(props, string(neo4j_entity.SocialPropertyId)),
		Url:            utils.GetStringPropOrEmpty(props, string(neo4j_entity.SocialPropertyUrl)),
		Alias:          utils.GetStringPropOrEmpty(props, string(neo4j_entity.SocialPropertyAlias)),
		FollowersCount: utils.GetInt64PropOrZero(props, string(neo4j_entity.SocialPropertyFollowersCount)),
		ExternalId:     utils.GetStringPropOrEmpty(props, string(neo4j_entity.SocialPropertyExternalId)),
		CreatedAt:      utils.GetTimePropOrEpochStart(props, string(neo4j_entity.SocialPropertyCreatedAt)),
		UpdatedAt:      utils.GetTimePropOrEpochStart(props, string(neo4j_entity.SocialPropertyUpdatedAt)),
		AppSource:      utils.GetStringPropOrEmpty(props, string(neo4j_entity.SocialPropertyAppSource)),
		Source:         neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, string(neo4j_entity.SocialPropertySource))),
	}
	return &social
}

func MapDbNodeToReminderEntity(dbNode *dbtype.Node) *neo4j_entity.ReminderEntity {
	if dbNode == nil {
		return &neo4j_entity.ReminderEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	reminder := neo4j_entity.ReminderEntity{
		Id:             utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:      utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:      utils.GetTimePropOrEpochStart(props, "updatedAt"),
		UserId:         utils.GetStringPropOrEmpty(props, "userId"),
		OrganizationId: utils.GetStringPropOrEmpty(props, "organizationId"),
		Content:        utils.GetStringPropOrEmpty(props, "content"),
		DueDate:        utils.GetTimePropOrEpochStart(props, "dueDate"),
		Dismissed:      utils.GetBoolPropOrFalse(props, "dismissed"),
		Sent:           utils.GetBoolPropOrFalse(props, "sent"),
	}
	return &reminder
}

func MapDbNodeToBankAccountEntity(dbNode *dbtype.Node) *neo4j_entity.BankAccountEntity {
	if dbNode == nil {
		return &neo4j_entity.BankAccountEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	bankAccount := neo4j_entity.BankAccountEntity{
		Id:                  utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:           utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:           utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:              neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:       neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
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
func MapDbNodeToEmailEntity(node *dbtype.Node) *neo4j_entity.EmailEntity {
	if node == nil {
		return &neo4j_entity.EmailEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	return &neo4j_entity.EmailEntity{
		Id:                utils.GetStringPropOrEmpty(props, "id"),
		Email:             utils.GetStringPropOrEmpty(props, string(neo4j_entity.EmailPropertyEmail)),
		RawEmail:          utils.GetStringPropOrEmpty(props, string(neo4j_entity.EmailPropertyRawEmail)),
		Work:              utils.GetBoolPropOrNil(props, string(neo4j_entity.EmailPropertyWork)),
		CreatedAt:         utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:         utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Primary:           utils.GetBoolPropOrFalse(props, "primary"),
		Source:            neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		IsValidSyntax:     utils.GetBoolPropOrNil(props, string(neo4j_entity.EmailPropertyIsValidSyntax)),
		IsCatchAll:        utils.GetBoolPropOrNil(props, string(neo4j_entity.EmailPropertyIsCatchAll)),
		Deliverable:       utils.GetStringPropOrNil(props, string(neo4j_entity.EmailPropertyDeliverable)),
		IsRoleAccount:     utils.GetBoolPropOrNil(props, string(neo4j_entity.EmailPropertyIsRoleAccount)),
		IsSystemGenerated: utils.GetBoolPropOrNil(props, string(neo4j_entity.EmailPropertyIsSystemGenerated)),
		EmailInternalFields: neo4j_entity.EmailInternalFields{
			ValidatedAt:           utils.GetTimePropOrNil(props, string(neo4j_entity.EmailPropertyValidatedAt)),
			ValidationRequestedAt: utils.GetTimePropOrNil(props, string(neo4j_entity.EmailPropertyValidationRequestedAt)),
		},
		Username:        utils.GetStringPropOrNil(props, string(neo4j_entity.EmailPropertyUsername)),
		IsRisky:         utils.GetBoolPropOrNil(props, string(neo4j_entity.EmailPropertyIsRisky)),
		IsFirewalled:    utils.GetBoolPropOrNil(props, string(neo4j_entity.EmailPropertyIsFirewalled)),
		Provider:        utils.GetStringPropOrNil(props, string(neo4j_entity.EmailPropertyProvider)),
		Firewall:        utils.GetStringPropOrNil(props, string(neo4j_entity.EmailPropertyFirewall)),
		IsMailboxFull:   utils.GetBoolPropOrNil(props, string(neo4j_entity.EmailPropertyIsMailboxFull)),
		IsFreeAccount:   utils.GetBoolPropOrNil(props, string(neo4j_entity.EmailPropertyIsFreeAccount)),
		SmtpSuccess:     utils.GetBoolPropOrNil(props, string(neo4j_entity.EmailPropertySmtpSuccess)),
		ResponseCode:    utils.GetStringPropOrNil(props, string(neo4j_entity.EmailPropertyResponseCode)),
		ErrorCode:       utils.GetStringPropOrNil(props, string(neo4j_entity.EmailPropertyErrorCode)),
		Description:     utils.GetStringPropOrNil(props, string(neo4j_entity.EmailPropertyDescription)),
		IsPrimaryDomain: utils.GetBoolPropOrNil(props, string(neo4j_entity.EmailPropertyIsPrimaryDomain)),
		PrimaryDomain:   utils.GetStringPropOrNil(props, string(neo4j_entity.EmailPropertyPrimaryDomain)),
		AlternateEmail:  utils.GetStringPropOrNil(props, string(neo4j_entity.EmailPropertyAlternateEmail)),
		RetryValidation: utils.GetBoolPropOrNil(props, string(neo4j_entity.EmailPropertyRetryValidation)),
	}
}

func MapDbNodeToPhoneNumberEntity(node *dbtype.Node) *neo4j_entity.PhoneNumberEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)
	return &neo4j_entity.PhoneNumberEntity{
		Id:             utils.GetStringPropOrEmpty(props, "id"),
		E164:           utils.GetStringPropOrEmpty(props, "e164"),
		RawPhoneNumber: utils.GetStringPropOrEmpty(props, "rawPhoneNumber"),
		Validated:      utils.GetBoolPropOrNil(props, "validated"),
		Source:         neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:  neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:      utils.GetStringPropOrEmpty(props, "appSource"),
		CreatedAt:      utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:      utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
}

func MapDbNodeToExternalSystem(node *dbtype.Node) *neo4j_entity.ExternalSystemEntity {
	if node == nil {
		return &neo4j_entity.ExternalSystemEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	externalSystemEntity := neo4j_entity.ExternalSystemEntity{
		ExternalSystemId: commonenum.DecodeSource(utils.GetStringPropOrEmpty(props, "id")),
		Name:             utils.GetStringPropOrEmpty(props, "name"),
	}
	if externalSystemEntity.ExternalSystemId == commonenum.SourceStripe {
		externalSystemEntity.Stripe.PaymentMethodTypes = utils.GetListStringPropOrEmpty(props, neo4j_entity.PropertyExternalSystemStripePaymentMethodTypes)
	}
	return &externalSystemEntity
}

func AddDbRelationshipToExternalSystemEntity(relationship dbtype.Relationship, neo4j_entity *neo4j_entity.ExternalSystemEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	neo4j_entity.Relationship.SyncDate = utils.GetTimePropOrNil(props, "syncDate")
	neo4j_entity.Relationship.ExternalId = utils.GetStringPropOrEmpty(props, "externalId")
	neo4j_entity.Relationship.ExternalUrl = utils.GetStringPropOrNil(props, "externalUrl")
	neo4j_entity.Relationship.ExternalSource = utils.GetStringPropOrNil(props, "externalSource")
	neo4j_entity.Relationship.Primary = utils.GetBoolPropOrFalse(props, "primary")
}

// TODO use nil
func MapDbNodeToOpportunityEntity(node *dbtype.Node) *neo4j_entity.OpportunityEntity {
	if node == nil {
		return &neo4j_entity.OpportunityEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	opportunity := neo4j_entity.OpportunityEntity{
		Id:                utils.GetStringPropOrEmpty(props, "id"),
		Name:              utils.GetStringPropOrEmpty(props, "name"),
		Amount:            utils.GetFloatPropOrZero(props, string(neo4j_entity.OpportunityPropertyAmount)),
		MaxAmount:         utils.GetFloatPropOrZero(props, string(neo4j_entity.OpportunityPropertyMaxAmount)),
		InternalType:      enum.DecodeOpportunityInternalType(utils.GetStringPropOrEmpty(props, "internalType")),
		ExternalType:      utils.GetStringPropOrEmpty(props, "externalType"),
		InternalStage:     enum.DecodeOpportunityInternalStage(utils.GetStringPropOrEmpty(props, "internalStage")),
		ExternalStage:     utils.GetStringPropOrEmpty(props, "externalStage"),
		EstimatedClosedAt: utils.GetTimePropOrNil(props, "estimatedClosedAt"),
		ClosedAt:          utils.GetTimePropOrNil(props, "closedAt"),
		GeneralNotes:      utils.GetStringPropOrEmpty(props, "generalNotes"),
		NextSteps:         utils.GetStringPropOrEmpty(props, string(neo4j_entity.OpportunityPropertyNextSteps)),
		Comments:          utils.GetStringPropOrEmpty(props, "comments"),
		CreatedAt:         utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:         utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:         utils.GetStringPropOrEmpty(props, "appSource"),
		Source:            neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:     neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		OwnerUserId:       utils.GetStringPropOrEmpty(props, "ownerUserId"),
		Currency:          enum.DecodeCurrency(utils.GetStringPropOrEmpty(props, string(neo4j_entity.OpportunityPropertyCurrency))),
		LikelihoodRate:    utils.GetInt64PropOrDefault(props, string(neo4j_entity.OpportunityPropertyLikelihoodRate), 0),
		StageUpdatedAt:    utils.GetTimePropOrNil(props, string(neo4j_entity.OpportunityPropertyStageUpdatedAt)),
		RenewalDetails: neo4j_entity.RenewalDetails{
			RenewedAt:              utils.GetTimePropOrNil(props, "renewedAt"),
			RenewalLikelihood:      enum.DecodeRenewalLikelihood(utils.GetStringPropOrEmpty(props, "renewalLikelihood")),
			RenewalUpdatedByUserId: utils.GetStringPropOrEmpty(props, "renewalUpdatedByUserId"),
			RenewalUpdatedByUserAt: utils.GetTimePropOrNil(props, "renewalUpdatedByUserAt"),
			RenewalApproved:        utils.GetBoolPropOrFalse(props, "renewalApproved"),
			RenewalAdjustedRate:    utils.GetInt64PropOrDefault(props, "renewalAdjustedRate", 100),
		},
		InternalFields: neo4j_entity.OpportunityInternalFields{
			RolloutRenewalRequestedAt: utils.GetTimePropOrNil(props, "techRolloutRenewalRequestedAt"),
		},
	}
	return &opportunity
}

func MapDbNodeToStateEntity(node dbtype.Node) *neo4j_entity.StateEntity {
	props := utils.GetPropsFromNode(node)
	result := neo4j_entity.StateEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		Code:      utils.GetStringPropOrEmpty(props, "code"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
	return &result
}

func MapDbNodeToPageView(node *dbtype.Node) *neo4j_entity.PageViewEntity {
	if node == nil {
		return &neo4j_entity.PageViewEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	pageViewAction := neo4j_entity.PageViewEntity{
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
		Source:         neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:  neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:      utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &pageViewAction
}

func MapDbNodeToLogEntryEntity(node *dbtype.Node) *neo4j_entity.LogEntryEntity {
	if node == nil {
		return &neo4j_entity.LogEntryEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	logEntry := neo4j_entity.LogEntryEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		StartedAt:     utils.GetTimePropOrEpochStart(props, "startedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &logEntry
}

func MapDbNodeToMarkdownEventEntity(node *dbtype.Node) *neo4j_entity.MarkdownEventEntity {
	if node == nil {
		return &neo4j_entity.MarkdownEventEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	event := neo4j_entity.MarkdownEventEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Content:   utils.GetStringPropOrEmpty(props, "content"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource: utils.GetStringPropOrEmpty(props, "appSource"),
		Source:    neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
	}
	return &event
}

func MapDbNodeToMeetingEntity(node *dbtype.Node) *neo4j_entity.MeetingEntity {
	if node == nil {
		return &neo4j_entity.MeetingEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	status := enum.DecodeMeetingStatus(utils.GetStringPropOrEmpty(props, "status"))
	meetingEntity := neo4j_entity.MeetingEntity{
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
		Source:             neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:      neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
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

func MapDbNodeToActionEntity(node *dbtype.Node) *neo4j_entity.ActionEntity {
	if node == nil {
		return &neo4j_entity.ActionEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	action := neo4j_entity.ActionEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Type:      commonenum.GetActionType(utils.GetStringPropOrEmpty(props, "type")),
		Content:   utils.GetStringPropOrEmpty(props, "content"),
		Metadata:  utils.GetStringPropOrEmpty(props, "metadata"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		AppSource: utils.GetStringPropOrEmpty(props, "appSource"),
		Source:    neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
	}
	return &action
}

func MapDbNodeToNoteEntity(node *dbtype.Node) *neo4j_entity.NoteEntity {
	if node == nil {
		return &neo4j_entity.NoteEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	note := neo4j_entity.NoteEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &note
}

func MapDbNodeToInteractionEventEntity(node *neo4j.Node) *neo4j_entity.InteractionEventEntity {
	if node == nil {
		return &neo4j_entity.InteractionEventEntity{}
	}
	props := utils.GetPropsFromNode(*node)

	return MapDbPropsToInteractionEventEntity(props)
}

func MapDbPropsToInteractionEventEntity(props map[string]interface{}) *neo4j_entity.InteractionEventEntity {
	interactionEventEntity := neo4j_entity.InteractionEventEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Identifier:    utils.GetStringPropOrEmpty(props, "identifier"),
		Channel:       commonenum.DecodeInteractionEventChannel(utils.GetStringPropOrEmpty(props, "channel")),
		ChannelData:   utils.GetStringPropOrEmpty(props, "channelData"),
		EventType:     utils.GetStringPropOrEmpty(props, "eventType"),
		Hide:          utils.GetBoolPropOrFalse(props, "hide"),
		Content:       utils.GetStringPropOrEmpty(props, "content"),
		ContentType:   utils.GetStringPropOrEmpty(props, "contentType"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &interactionEventEntity
}

func MapDbNodeToInteractionSessionEntity(node *dbtype.Node) *neo4j_entity.InteractionSessionEntity {
	if node == nil {
		return &neo4j_entity.InteractionSessionEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	interactionSession := neo4j_entity.InteractionSessionEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Channel:       commonenum.DecodeInteractionSessionChannel(utils.GetStringPropOrEmpty(props, "channel")),
		ChannelData:   utils.GetStringPropOrEmpty(props, "channelData"),
		Identifier:    utils.GetStringPropOrEmpty(props, "identifier"),
		Type:          commonenum.DecodeInteractionSessionType(utils.GetStringPropOrEmpty(props, "type")),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		Status:        commonenum.DecodeInteractionSessionStatus(utils.GetStringPropOrEmpty(props, "status")),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
		Source:        neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &interactionSession
}

func MapDbNodeToContactEntity(dbNode *dbtype.Node) *neo4j_entity.ContactEntity {
	props := utils.GetPropsFromNode(*dbNode)
	contact := neo4j_entity.ContactEntity{
		Id:              utils.GetStringPropOrEmpty(props, "id"),
		FirstName:       utils.GetStringPropOrEmpty(props, string(neo4j_entity.ContactPropertyFirstName)),
		LastName:        utils.GetStringPropOrEmpty(props, string(neo4j_entity.ContactPropertyLastName)),
		Name:            utils.GetStringPropOrEmpty(props, string(neo4j_entity.ContactPropertyName)),
		Description:     utils.GetStringPropOrEmpty(props, string(neo4j_entity.ContactPropertyDescription)),
		Timezone:        utils.GetStringPropOrEmpty(props, string(neo4j_entity.ContactPropertyTimezone)),
		ProfilePhotoUrl: utils.GetStringPropOrEmpty(props, string(neo4j_entity.ContactPropertyProfilePhotoUrl)),
		Username:        utils.GetStringPropOrEmpty(props, string(neo4j_entity.ContactPropertyUsername)),
		Prefix:          utils.GetStringPropOrEmpty(props, string(neo4j_entity.ContactPropertyPrefix)),
		Hide:            utils.GetBoolPropOrFalse(props, string(neo4j_entity.ContactPropertyHide)),
		HiddenAt:        utils.GetTimePropOrNil(props, string(neo4j_entity.ContactPropertyHiddenAt)),
		CreatedAt:       utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:       utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:          neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		AppSource:       utils.GetStringPropOrEmpty(props, "appSource"),
		EventStoreAggregate: neo4j_entity.EventStoreAggregate{
			AggregateVersion: utils.GetInt64PropOrNil(props, "aggregateVersion"),
		},
		ContactInternalFields: neo4j_entity.ContactInternalFields{
			CheckedAt: utils.GetTimePropOrNil(props, string(neo4j_entity.ContactPropertyCheckedAt)),
		},
		EnrichDetails: neo4j_entity.ContactEnrichDetails{
			EnrichRequestedAt:         utils.GetTimePropOrNil(props, string(neo4j_entity.ContactPropertyEnrichRequestedAt)),
			EnrichedAt:                utils.GetTimePropOrNil(props, string(neo4j_entity.ContactPropertyEnrichedAt)),
			EnrichFailedAt:            utils.GetTimePropOrNil(props, string(neo4j_entity.ContactPropertyEnrichFailedAt)),
			EnrichAttempts:            utils.GetInt64PropOrZero(props, string(neo4j_entity.ContactPropertyEnrichAttempts)),
			BettercontactFoundEmailAt: utils.GetTimePropOrNil(props, string(neo4j_entity.ContactPropertyBettercontactFoundEmailAt)),
			EnrichedScrapinRecordId:   utils.GetStringPropOrEmpty(props, string(neo4j_entity.ContactPropertyEnrichedScrapinRecordId)),
			FindWorkEmailWithBetterContactRequestedId:   utils.GetStringPropOrNil(props, string(neo4j_entity.ContactPropertyFindWorkEmailWithBetterContactRequestedId)),
			FindWorkEmailWithBetterContactRequestedAt:   utils.GetTimePropOrNil(props, string(neo4j_entity.ContactPropertyFindWorkEmailWithBetterContactRequestedAt)),
			FindWorkEmailWithBetterContactCompletedAt:   utils.GetTimePropOrNil(props, string(neo4j_entity.ContactPropertyFindWorkEmailWithBetterContactCompletedAt)),
			FindWorkEmailWithBetterContactFound:         utils.GetBoolPropOrNil(props, string(neo4j_entity.ContactPropertyFindWorkEmailWithBetterContactFound)),
			FindMobilePhoneWithBetterContactRequestedId: utils.GetStringPropOrNil(props, string(neo4j_entity.ContactPropertyFindMobilePhoneWithBetterContactRequestedId)),
			FindMobilePhoneWithBetterContactRequestedAt: utils.GetTimePropOrNil(props, string(neo4j_entity.ContactPropertyFindMobilePhoneWithBetterContactRequestedAt)),
			FindMobilePhoneWithBetterContactCompletedAt: utils.GetTimePropOrNil(props, string(neo4j_entity.ContactPropertyFindMobilePhoneWithBetterContactCompletedAt)),
			FindMobilePhoneWithBetterContactFound:       utils.GetBoolPropOrNil(props, string(neo4j_entity.ContactPropertyFindMobilePhoneWithBetterContactFound)),
		},
	}
	return &contact
}

func MapDbNodeToTimelineEvent(dbNode *dbtype.Node) neo4j_entity.TimelineEvent {
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

func MapDbNodeToDomainEntity(node *dbtype.Node) *neo4j_entity.DomainEntity {
	if node == nil {
		return &neo4j_entity.DomainEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	domain := neo4j_entity.DomainEntity{
		CreatedAt:     utils.GetTimePropOrEpochStart(props, string(neo4j_entity.DomainPropertyCreatedAt)),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, string(neo4j_entity.DomainPropertyUpdatedAt)),
		AppSource:     utils.GetStringPropOrEmpty(props, string(neo4j_entity.DomainPropertyAppSource)),
		Source:        neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, string(neo4j_entity.DomainPropertySource))),
		Domain:        utils.GetStringPropOrEmpty(props, string(neo4j_entity.DomainPropertyDomain)),
		IsPrimary:     utils.GetBoolPropOrNil(props, string(neo4j_entity.DomainPropertyIsPrimary)),
		Accessible:    utils.GetBoolPropOrNil(props, string(neo4j_entity.DomainPropertyAccessible)),
		PrimaryDomain: utils.GetStringPropOrEmpty(props, string(neo4j_entity.DomainPropertyPrimaryDomain)),
		InternalFields: neo4j_entity.DomainInternalFields{
			PrimaryDomainCheckRequestedAt: utils.GetTimePropOrNil(props, string(neo4j_entity.DomainPropertyPrimaryDomainCheckRequestedAt)),
		},
	}
	return &domain
}

func MapDbNodeToLocationEntity(node *dbtype.Node) *neo4j_entity.LocationEntity {
	if node == nil {
		return &neo4j_entity.LocationEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	location := neo4j_entity.LocationEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyName)),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Country:       utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyCountry)),
		CountryCodeA2: utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyCountryCodeA2)),
		CountryCodeA3: utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyCountryCodeA3)),
		Region:        utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyRegion)),
		Locality:      utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyLocality)),
		Address:       utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyAddress)),
		Address2:      utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyAddress2)),
		Zip:           utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyZip)),
		AddressType:   utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyAddressType)),
		HouseNumber:   utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyHouseNumber)),
		PostalCode:    utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyPostalCode)),
		PlusFour:      utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyPlusFour)),
		Commercial:    utils.GetBoolPropOrFalse(props, string(neo4j_entity.LocationPropertyCommercial)),
		Predirection:  utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyPredirection)),
		District:      utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyDistrict)),
		Street:        utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyStreet)),
		RawAddress:    utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyRawAddress)),
		Latitude:      utils.GetFloatPropOrNil(props, string(neo4j_entity.LocationPropertyLatitude)),
		Longitude:     utils.GetFloatPropOrNil(props, string(neo4j_entity.LocationPropertyLongitude)),
		TimeZone:      utils.GetStringPropOrEmpty(props, string(neo4j_entity.LocationPropertyTimeZone)),
		UtcOffset:     utils.GetFloatPropOrNil(props, string(neo4j_entity.LocationPropertyUtcOffset)),
		Source:        neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4j_entity.DecodeDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &location
}

func MapDbNodeToFlowEntity(node *dbtype.Node) *neo4j_entity.FlowEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)
	domain := neo4j_entity.FlowEntity{
		Id:             utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:      utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:      utils.GetTimePropOrEpochStart(props, "updatedAt"),
		TableViewDefId: utils.GetStringPropOrEmpty(props, "tableViewDefId"),
		DefaultName:    utils.GetStringPropOrEmpty(props, "defaultName"),
		Name:           utils.GetStringPropOrEmpty(props, "name"),
		Nodes:          utils.GetStringPropOrEmpty(props, "nodes"),
		Edges:          utils.GetStringPropOrEmpty(props, "edges"),
		FirstStartedAt: utils.GetTimePropOrNil(props, "firstStartedAt"),
		Status:         neo4j_entity.GetFlowStatus(utils.GetStringPropOrEmpty(props, "status")),
		Total:          utils.GetInt64PropOrZero(props, "total"),
		OnHold:         utils.GetInt64PropOrZero(props, "onHold"),
		Ready:          utils.GetInt64PropOrZero(props, "ready"),
		Scheduled:      utils.GetInt64PropOrZero(props, "scheduled"),
		InProgress:     utils.GetInt64PropOrZero(props, "inProgress"),
		Completed:      utils.GetInt64PropOrZero(props, "completed"),
		GoalAchieved:   utils.GetInt64PropOrZero(props, "goalAchieved"),
	}
	return &domain
}

func MapDbNodeToFlowParticipantEntity(node *dbtype.Node) *neo4j_entity.FlowParticipantEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)
	e := neo4j_entity.FlowParticipantEntity{
		Id:                 utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:          utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:          utils.GetTimePropOrEpochStart(props, "updatedAt"),
		EntityId:           utils.GetStringPropOrEmpty(props, "entityId"),
		EntityType:         model.DecodeEntityType(utils.GetStringPropOrEmpty(props, "entityType")),
		Status:             neo4j_entity.GetFlowContactStatus(utils.GetStringPropOrEmpty(props, "status")),
		RequirementsUnmeet: neo4j_entity.GetFlowParticipantRequirementsUnmeet(utils.GetListStringPropOrEmpty(props, "requirementsUnmeet")),
	}
	return &e
}

func MapDbNodeToFlowSenderEntity(node *dbtype.Node) *neo4j_entity.FlowSenderEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)
	e := neo4j_entity.FlowSenderEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
		UserId:    utils.GetStringPropOrNil(props, "userId"),
	}
	return &e
}

func MapDbNodeToFlowActionEntity(node *dbtype.Node) *neo4j_entity.FlowActionEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)

	e := neo4j_entity.FlowActionEntity{
		Id:         utils.GetStringPropOrEmpty(props, "id"),
		ExternalId: utils.GetStringPropOrEmpty(props, "externalId"),
		CreatedAt:  utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:  utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Json:       utils.GetStringPropOrEmpty(props, "json"),
		Type:       utils.GetStringPropOrEmpty(props, "type"),
	}

	e.Data.Action = neo4j_entity.GetFlowActionType(utils.GetStringPropOrEmpty(props, "action"))

	e.Data.WaitBefore = utils.GetInt64PropOrZero(props, "waitBefore")

	e.Data.Entity = utils.GetStringPropOrNil(props, "data_neo4j_entity")
	e.Data.TriggerType = utils.GetStringPropOrNil(props, "data_triggerType")
	e.Data.Subject = utils.GetStringPropOrNil(props, "data_subject")
	e.Data.BodyTemplate = utils.GetStringPropOrNil(props, "data_bodyTemplate")
	e.Data.MessageTemplate = utils.GetStringPropOrNil(props, "data_messageTemplate")

	return &e
}

func MapDbNodeToFlowExecutionSettingsEntity(node *dbtype.Node) *neo4j_entity.FlowExecutionSettingsEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)
	e := neo4j_entity.FlowExecutionSettingsEntity{
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

func MapDbNodeToFlowActionExecutionEntity(node *dbtype.Node) *neo4j_entity.FlowActionExecutionEntity {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)
	e := neo4j_entity.FlowActionExecutionEntity{
		Id:              utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:       utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:       utils.GetTimePropOrEpochStart(props, "updatedAt"),
		FlowId:          utils.GetStringPropOrEmpty(props, "flowId"),
		ParticipantId:   utils.GetStringPropOrEmpty(props, "participantId"),
		EntityId:        utils.GetStringPropOrEmpty(props, "entityId"),
		EntityType:      model.DecodeEntityType(utils.GetStringPropOrEmpty(props, "entityType")),
		ActionId:        utils.GetStringPropOrEmpty(props, "actionId"),
		ScheduledAt:     utils.GetTimePropOrNow(props, "scheduledAt"),
		ExecutedAt:      utils.GetTimePropOrNil(props, "executedAt"),
		StatusUpdatedAt: utils.GetTimePropOrNow(props, "statusUpdatedAt"),
		Status:          neo4j_entity.GetFlowActionExecutionStatus(utils.GetStringPropOrEmpty(props, "status")),
		Error:           utils.GetStringPropOrNil(props, "error"),

		Mailbox:   utils.GetStringPropOrNil(props, "mailbox"),
		UserId:    utils.GetStringPropOrNil(props, "userId"),
		SocialUrl: utils.GetStringPropOrNil(props, "socialUrl"),
	}
	return &e
}

func MapDbNodeToLinkedinConnectionRequestEntity(node *dbtype.Node) *neo4j_entity.LinkedinConnectionRequest {
	if node == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*node)
	e := neo4j_entity.LinkedinConnectionRequest{
		Id:           utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:    utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:    utils.GetTimePropOrEpochStart(props, "updatedAt"),
		ProducerId:   utils.GetStringPropOrEmpty(props, "producerId"),
		ProducerType: utils.GetStringPropOrEmpty(props, "producerType"),
		ScheduledAt:  utils.GetTimePropOrNow(props, "scheduledAt"),
		SocialUrl:    utils.GetStringPropOrEmpty(props, "socialUrl"),
		UserId:       utils.GetStringPropOrEmpty(props, "userId"),
		Status:       neo4j_entity.GetLinkedinConnectionRequestStatus(utils.GetStringPropOrEmpty(props, "status")),
	}
	return &e
}

func MapDbNodeToCustomFieldTemplateEntity(node *dbtype.Node) *neo4j_entity.CustomFieldTemplateEntity {
	if node == nil {
		return &neo4j_entity.CustomFieldTemplateEntity{}
	}
	props := utils.GetPropsFromNode(*node)
	customFieldTemplateEntity := neo4j_entity.CustomFieldTemplateEntity{
		Id:          utils.GetStringPropOrEmpty(props, string(neo4j_entity.CustomFieldTemplatePropertyId)),
		Name:        utils.GetStringPropOrEmpty(props, string(neo4j_entity.CustomFieldTemplatePropertyName)),
		EntityType:  model.DecodeEntityType(utils.GetStringPropOrEmpty(props, string(neo4j_entity.CustomFieldTemplatePropertyEntityType))),
		Type:        utils.GetStringPropOrEmpty(props, string(neo4j_entity.CustomFieldTemplatePropertyType)),
		ValidValues: utils.GetListStringPropOrEmpty(props, string(neo4j_entity.CustomFieldTemplatePropertyValidValues)),
		Order:       utils.GetInt64PropOrNil(props, string(neo4j_entity.CustomFieldTemplatePropertyOrder)),
		Required:    utils.GetBoolPropOrNil(props, string(neo4j_entity.CustomFieldTemplatePropertyRequired)),
		Length:      utils.GetInt64PropOrNil(props, string(neo4j_entity.CustomFieldTemplatePropertyLength)),
		Min:         utils.GetInt64PropOrNil(props, string(neo4j_entity.CustomFieldTemplatePropertyMin)),
		Max:         utils.GetInt64PropOrNil(props, string(neo4j_entity.CustomFieldTemplatePropertyMax)),
		CreatedAt:   utils.GetTimePropOrEpochStart(props, string(neo4j_entity.CustomFieldTemplatePropertyCreatedAt)),
		UpdatedAt:   utils.GetTimePropOrEpochStart(props, string(neo4j_entity.CustomFieldTemplatePropertyUpdatedAt)),
	}
	return &customFieldTemplateEntity
}

func MapDbNodeToIndustryEntity(dbNode *dbtype.Node) *neo4j_entity.IndustryEntity {
	if dbNode == nil {
		return &neo4j_entity.IndustryEntity{}
	}
	props := utils.GetPropsFromNode(*dbNode)
	industry := neo4j_entity.IndustryEntity{
		CreatedAt: utils.GetTimePropOrEpochStart(props, string(neo4j_entity.IndustryPropertyCreatedAt)),
		Code:      utils.GetStringPropOrEmpty(props, string(neo4j_entity.IndustryPropertyCode)),
		Name:      utils.GetStringPropOrEmpty(props, string(neo4j_entity.IndustryPropertyName)),
	}
	return &industry
}
