package resolver

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

// ColumnView represents a column in a table view with type and width.
func DefaultTableViewDefinitions() []postgresEntity.TableViewDefinition {
	organizationColumns := postgresEntity.Columns{
		Columns: []postgresEntity.ColumnView{
			{ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsRelationship.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsOnboardingStatus.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsRenewalLikelihood.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsRenewlDate.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsForecastArr.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsOwner.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsLastTouchpoint.String(), Width: 200, Visible: true},
		},
	}
	organizationColumnsJsonData, err := json.Marshal(organizationColumns)
	if err != nil {
		fmt.Println("Error serializing data:", err)
		return []postgresEntity.TableViewDefinition{}
	}

	invoiceColumns := postgresEntity.Columns{
		Columns: []postgresEntity.ColumnView{
			{ColumnType: model.ColumnViewTypeInvoicesIssueDate.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesIssueDatePast.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesDueDate.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesContract.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesBillingCycle.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesPaymentStatus.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesInvoiceNumber.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesAmount.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesInvoiceStatus.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesInvoicePreview.String(), Width: 200, Visible: true},
		},
	}
	invoiceColumnsJsonData, err := json.Marshal(invoiceColumns)
	if err != nil {
		fmt.Println("Error serializing data:", err)
		return []postgresEntity.TableViewDefinition{}
	}

	renewalColumns := postgresEntity.Columns{
		Columns: []postgresEntity.ColumnView{
			{ColumnType: model.ColumnViewTypeRenewalsAvatar.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeRenewalsName.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeRenewalsRenewalLikelihood.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeRenewalsRenewalDate.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeRenewalsForecastArr.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeRenewalsOwner.String(), Width: 200, Visible: true},
			{ColumnType: model.ColumnViewTypeRenewalsLastTouchpoint.String(), Width: 200, Visible: true},
		},
	}
	renewalColumnsJsonData, err := json.Marshal(renewalColumns)
	if err != nil {
		fmt.Println("Error serializing data:", err)
		return []postgresEntity.TableViewDefinition{}
	}

	return []postgresEntity.TableViewDefinition{
		{
			TableType:   model.TableViewTypeOrganizations.String(),
			Name:        "Organizations",
			ColumnsJson: string(organizationColumnsJsonData),
			Order:       1,
			Icon:        "",
			Filters:     "",
			Sorting:     "",
		},
		{
			TableType:   model.TableViewTypeInvoices.String(),
			Name:        "Invoices",
			ColumnsJson: string(invoiceColumnsJsonData),
			Order:       2,
			Icon:        "",
			Filters:     "",
			Sorting:     "",
		},
		{
			TableType:   model.TableViewTypeRenewals.String(),
			Name:        "Renewals",
			ColumnsJson: string(renewalColumnsJsonData),
			Order:       3,
			Icon:        "",
			Filters:     "",
			Sorting:     "",
		},
	}

}
