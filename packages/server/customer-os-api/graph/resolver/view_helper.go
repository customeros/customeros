package resolver

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

const (
	TableViewDefinitionOrganizationsName = "Organization"
	TableViewDefinitionCustomersName     = "Customers"
	TableViewDefinitionMyPortfolioName   = "My portfolio"
)

// ColumnView represents a column in a table view with type and width.
func DefaultTableViewDefinitions(userId string) []postgresEntity.TableViewDefinition {
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

	upcomingInvoicesTableViewDefinition, err := DefaultTableViewDefinitionUpcomingInvoices()
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	organizationsTableViewDefinition, err := DefaultTableViewDefinitionOrganization()
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	customersTableViewDefinition, err := DefaultTableViewDefinitionCustomers()
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	myPortfolioTableViewDefinition, err := DefaultTableViewDefinitionMyPortfolio(userId)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	pastInvoicesTableViewDefinition, err := DefaultTableViewDefinitionPastInvoices()
	if err != nil {
		fmt.Println("Error: ", err)
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
		upcomingInvoicesTableViewDefinition,
		pastInvoicesTableViewDefinition,
		organizationsTableViewDefinition,
		customersTableViewDefinition,
		myPortfolioTableViewDefinition,
	}
}

func DefaultTableViewDefinitionPastInvoices() (postgresEntity.TableViewDefinition, error) {
	columns := postgresEntity.Columns{
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
	jsonData, err := json.Marshal(columns)
	if err != nil {
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeInvoices.String(),
		Name:        "Past",
		ColumnsJson: string(jsonData),
		Order:       5,
		Icon:        "InvoiceCheck",
		Filters:     `{"AND":[{"filter":{"property":"INVOICE_DRY_RUN","value":false}}]}`,
		Sorting:     "",
	}, nil
}

func DefaultTableViewDefinitionUpcomingInvoices() (postgresEntity.TableViewDefinition, error) {
	columns := postgresEntity.Columns{
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
	jsonData, err := json.Marshal(columns)
	if err != nil {
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeInvoices.String(),
		Name:        "Upcoming",
		ColumnsJson: string(jsonData),
		Order:       4,
		Icon:        "InvoiceUpcoming",
		Filters:     `{"AND":[{"filter":{"property":"INVOICE_PREVIEW","value":true}}]}`,
		Sorting:     "",
	}, nil
}

func DefaultTableViewDefinitionOrganization() (postgresEntity.TableViewDefinition, error) {
	columns := postgresEntity.Columns{
		Columns: []postgresEntity.ColumnView{
			{ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsRelationship.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsRenewalLikelihood.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsRenewalDate.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsOnboardingStatus.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsForecastArr.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsOwner.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsLastTouchpoint.String(), Width: 100, Visible: true},
		},
	}
	jsonData, err := json.Marshal(columns)
	if err != nil {
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeOrganizations.String(),
		Name:        TableViewDefinitionOrganizationsName,
		ColumnsJson: string(jsonData),
		Order:       1,
		Icon:        "Building07",
		Filters:     ``,
		Sorting:     "",
	}, nil
}

func DefaultTableViewDefinitionCustomers() (postgresEntity.TableViewDefinition, error) {
	columns := postgresEntity.Columns{
		Columns: []postgresEntity.ColumnView{
			{ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsRelationship.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsOnboardingStatus.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsRenewalLikelihood.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsRenewalDate.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsForecastArr.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsOwner.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsLastTouchpoint.String(), Width: 100, Visible: true},
		},
	}
	jsonData, err := json.Marshal(columns)
	if err != nil {
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeOrganizations.String(),
		Name:        TableViewDefinitionCustomersName,
		ColumnsJson: string(jsonData),
		Order:       2,
		Icon:        "CheckHeart",
		Filters:     `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"IS_CUSTOMER","value":[true]}}]}`,
		Sorting:     "",
	}, nil
}

func DefaultTableViewDefinitionMyPortfolio(userId string) (postgresEntity.TableViewDefinition, error) {
	columns := postgresEntity.Columns{
		Columns: []postgresEntity.ColumnView{
			{ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsRelationship.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsOnboardingStatus.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsRenewalLikelihood.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsRenewalDate.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsForecastArr.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsOwner.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsLastTouchpoint.String(), Width: 100, Visible: true},
		},
	}
	jsonData, err := json.Marshal(columns)
	if err != nil {
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeOrganizations.String(),
		Name:        TableViewDefinitionMyPortfolioName,
		ColumnsJson: string(jsonData),
		Order:       3,
		Icon:        "Briefcase01",
		Filters:     fmt.Sprintf(`{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"OWNER_ID","value":["%s"]}}]}`, userId),
		Sorting:     "",
	}, nil
}
