package mapper

import (
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
	masterPlanEntity := entity.InvoiceEntity{
		Id:               utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:        utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:        utils.GetTimePropOrEpochStart(props, "updatedAt"),
		DryRun:           utils.GetBoolPropOrFalse(props, "dryRun"),
		Date:             utils.GetTimePropOrEpochStart(props, "date"),
		DueDate:          utils.GetTimePropOrEpochStart(props, "dueDate"),
		Amount:           utils.GetFloatPropOrZero(props, "amount"),
		Vat:              utils.GetFloatPropOrZero(props, "vat"),
		Total:            utils.GetFloatPropOrZero(props, "total"),
		RepositoryFileId: utils.GetStringPropOrEmpty(props, "repositoryFileId"),
		PdfGenerated:     utils.GetBoolPropOrFalse(props, "pdfGenerated"),
		Source:           entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:    entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:        utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &masterPlanEntity
}
func MapDbNodeToInvoiceLineEntity(dbNode *dbtype.Node) *entity.InvoiceLineEntity {
	if dbNode == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*dbNode)
	masterPlanEntity := entity.InvoiceLineEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Index:         utils.GetInt64PropOrZero(props, "index"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		Price:         utils.GetFloatPropOrZero(props, "price"),
		Quantity:      utils.GetInt64PropOrZero(props, "quantity"),
		Amount:        utils.GetFloatPropOrZero(props, "amount"),
		Vat:           utils.GetFloatPropOrZero(props, "vat"),
		Total:         utils.GetFloatPropOrZero(props, "total"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &masterPlanEntity
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
		return nil
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

func MapDbNodeToCurrencyEntity(dbNode *dbtype.Node) *entity.CurrencyEntity {
	if dbNode == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*dbNode)
	e := entity.CurrencyEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),

		Name:   utils.GetStringPropOrEmpty(props, "name"),
		Symbol: utils.GetStringPropOrEmpty(props, "symbol"),

		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &e
}
