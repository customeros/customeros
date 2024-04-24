package resolver

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

// ColumnView represents a column in a table view with type and width.
func DefaultTableViewDefinitions() []postgresEntity.TableViewDefinition {
	upcomingInvoiceColumns := postgresEntity.Columns{
		Columns: []postgresEntity.ColumnView{
			{ColumnType: model.ColumnViewTypeInvoicesInvoicePreview.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesContract.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesBillingCycle.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesIssueDate.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesDueDate.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesAmount.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesInvoiceStatus.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesIssueDatePast.String(), Width: 100, Visible: false},
			{ColumnType: model.ColumnViewTypeInvoicesPaymentStatus.String(), Width: 100, Visible: false},
		},
	}
	ucomingInvoiceColumnsJsonData, err := json.Marshal(upcomingInvoiceColumns)
	if err != nil {
		fmt.Println("Error serializing data:", err)
		return []postgresEntity.TableViewDefinition{}
	}

	pastInvoiceColumns := postgresEntity.Columns{
		Columns: []postgresEntity.ColumnView{
			{ColumnType: model.ColumnViewTypeInvoicesInvoiceNumber.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesContract.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesBillingCycle.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesIssueDatePast.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesDueDate.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesAmount.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesPaymentStatus.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeInvoicesIssueDate.String(), Width: 100, Visible: false},
			{ColumnType: model.ColumnViewTypeInvoicesInvoiceStatus.String(), Width: 100, Visible: false},
		},
	}
	pastInvoiceColumnsJsonData, err := json.Marshal(pastInvoiceColumns)
	if err != nil {
		fmt.Println("Error serializing data:", err)
		return []postgresEntity.TableViewDefinition{}
	}

	renewalColumns := postgresEntity.Columns{
		Columns: []postgresEntity.ColumnView{
			{ColumnType: model.ColumnViewTypeRenewalsAvatar.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeRenewalsName.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeRenewalsRenewalDate.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeRenewalsForecastArr.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeRenewalsRenewalLikelihood.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeRenewalsOwner.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeRenewalsLastTouchpoint.String(), Width: 100, Visible: true},
		},
	}
	renewalColumnsJsonData, err := json.Marshal(renewalColumns)
	if err != nil {
		fmt.Println("Error serializing data:", err)
		return []postgresEntity.TableViewDefinition{}
	}

	return []postgresEntity.TableViewDefinition{
		{
			TableType:   model.TableViewTypeRenewals.String(),
			Name:        "Monthly renewals",
			ColumnsJson: string(renewalColumnsJsonData),
			Order:       1,
			Icon:        "ClockFastForward",
			Filters:     `{"AND":[{"filter":{"property":"RENEWAL_CYCLE","value":"MONTHLY","operation":"EQ","includeEmpty":false}}]}`,
			Sorting:     "",
		},
		{
			TableType:   model.TableViewTypeRenewals.String(),
			Name:        "Quarterly renewals",
			ColumnsJson: string(renewalColumnsJsonData),
			Order:       2,
			Icon:        "ClockFastForward",
			Filters:     `{"AND":[{"filter":{"property":"RENEWAL_CYCLE","value":"QUARTERLY","operation":"EQ","includeEmpty":false}}]}`,
			Sorting:     "",
		},
		{
			TableType:   model.TableViewTypeRenewals.String(),
			Name:        "Annual renewals",
			ColumnsJson: string(renewalColumnsJsonData),
			Order:       3,
			Icon:        "ClockFastForward",
			Filters:     `{"AND":[{"filter":{"property":"RENEWAL_CYCLE","value":"ANNUALLY","operation":"EQ","includeEmpty":false}}]}`,
			Sorting:     "",
		},
		{
			TableType:   model.TableViewTypeInvoices.String(),
			Name:        "Upcoming",
			ColumnsJson: string(ucomingInvoiceColumnsJsonData),
			Order:       4,
			Icon:        "InvoiceUpcoming",
			Filters:     `{"AND":[{"filter":{"property":"INVOICE_PREVIEW","value":true}}]}`,
			Sorting:     "",
		},
		{
			TableType:   model.TableViewTypeInvoices.String(),
			Name:        "Past",
			ColumnsJson: string(pastInvoiceColumnsJsonData),
			Order:       5,
			Icon:        "InvoiceCheck",
			Filters:     `{"AND":[{"filter":{"property":"INVOICE_DRY_RUN","value":false}}]}`,
			Sorting:     "",
		},
	}

}
