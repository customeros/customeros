package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/pkg/errors"
	"strings"
)

func GetTagId(ctx context.Context, services *service.Services, tagId, tagName *string) string {
	outputTagId := ""
	if tagId != nil && *tagId != "" {
		tagEntity, _ := services.TagService.GetById(ctx, *tagId)
		if tagEntity != nil {
			outputTagId = tagEntity.Id
		}
	}
	if outputTagId == "" && tagName != nil && strings.TrimSpace(*tagName) != "" {
		tagEntity, _ := services.TagService.GetByNameOptional(ctx, strings.TrimSpace(*tagName))
		if tagEntity != nil {
			outputTagId = tagEntity.Id
		}
	}
	return outputTagId
}

func CreateTag(ctx context.Context, services *service.Services, tagName *string) (*neo4jentity.TagEntity, error) {
	if tagName == nil || strings.TrimSpace(*tagName) == "" {
		return nil, errors.New("tag name is empty")
	}
	return services.TagService.Merge(ctx, &neo4jentity.TagEntity{
		Name:          strings.TrimSpace(*tagName),
		Source:        neo4jentity.DataSourceOpenline,
		SourceOfTruth: neo4jentity.DataSourceOpenline,
		AppSource:     constants.AppSourceCustomerOsApi,
	})
}

func DefaultTableViewDefinitions() []postgresEntity.TableViewDefinition {
	return []postgresEntity.TableViewDefinition{
		{
			TableType: model.TableViewTypeOrganizations.String(),
			Name:      "Organizations",
			Columns: strings.Join([]string{
				model.ColumnViewTypeOrganizationsAvatar.String(),
				model.ColumnViewTypeOrganizationsName.String(),
				model.ColumnViewTypeOrganizationsWebsite.String(),
				model.ColumnViewTypeOrganizationsRelationship.String(),
				model.ColumnViewTypeOrganizationsOnboardingStatus.String(),
				model.ColumnViewTypeOrganizationsRenewalLikelihood.String(),
				model.ColumnViewTypeOrganizationsRenewlDate.String(),
				model.ColumnViewTypeOrganizationsForecastArr.String(),
				model.ColumnViewTypeOrganizationsOwner.String(),
				model.ColumnViewTypeOrganizationsLastTouchpoint.String(),
			}, ","),
			Order:   1,
			Icon:    "",
			Filters: "",
			Sorting: "",
		},
		{
			TableType: model.TableViewTypeInvoices.String(),
			Name:      "Invoices",
			Columns: strings.Join([]string{
				model.ColumnViewTypeInvoicesIssueDate.String(),
				model.ColumnViewTypeInvoicesIssueDatePast.String(),
				model.ColumnViewTypeInvoicesDueDate.String(),
				model.ColumnViewTypeInvoicesContract.String(),
				model.ColumnViewTypeInvoicesBillingCycle.String(),
				model.ColumnViewTypeInvoicesPaymentStatus.String(),
				model.ColumnViewTypeInvoicesInvoiceNumber.String(),
				model.ColumnViewTypeInvoicesAmount.String(),
				model.ColumnViewTypeInvoicesInvoiceStatus.String(),
				model.ColumnViewTypeInvoicesInvoicePreview.String(),
			}, ","),
			Order:   2,
			Icon:    "",
			Filters: "",
			Sorting: "",
		},
		{
			TableType: model.TableViewTypeRenewals.String(),
			Name:      "Renewals",
			Columns: strings.Join([]string{
				model.ColumnViewTypeRenewalsAvatar.String(),
				model.ColumnViewTypeRenewalsName.String(),
				model.ColumnViewTypeRenewalsRenewalLikelihood.String(),
				model.ColumnViewTypeRenewalsRenewalDate.String(),
				model.ColumnViewTypeRenewalsForecastArr.String(),
				model.ColumnViewTypeRenewalsOwner.String(),
				model.ColumnViewTypeRenewalsLastTouchpoint.String(),
			}, ","),
			Order:   3,
			Icon:    "",
			Filters: "",
			Sorting: "",
		},
	}

}
