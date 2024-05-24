package resolver

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
)

// ColumnView represents a column in a table view with type and width.
func DefaultTableViewDefinitions(userId string, span opentracing.Span) []postgresEntity.TableViewDefinition {
	monthlyRenewalsTableViewDefinition, err := DefaultTableViewDefinitionMonthlyRenewals(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	quarterlyRenewalsTableViewDefinition, err := DefaultTableViewDefinitionQuarterlyRenewals(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	annualRenewalsTableViewDefinition, err := DefaultTableViewDefinitionAnnualRenewals(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	upcomingInvoicesTableViewDefinition, err := DefaultTableViewDefinitionUpcomingInvoices(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	organizationsTableViewDefinition, err := DefaultTableViewDefinitionOrganization(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	customersTableViewDefinition, err := DefaultTableViewDefinitionCustomers(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	myPortfolioTableViewDefinition, err := DefaultTableViewDefinitionMyPortfolio(userId, span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	pastInvoicesTableViewDefinition, err := DefaultTableViewDefinitionPastInvoices(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	leadsTableViewDefinition, err := DefaultTableViewDefinitionLeads(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	nurtureTableViewDefinition, err := DefaultTableViewDefinitionNurture(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	churnTableViewDefinition, err := DefaultTableViewDefinitionChurn(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	return []postgresEntity.TableViewDefinition{
		monthlyRenewalsTableViewDefinition,
		quarterlyRenewalsTableViewDefinition,
		annualRenewalsTableViewDefinition,
		upcomingInvoicesTableViewDefinition,
		pastInvoicesTableViewDefinition,
		organizationsTableViewDefinition,
		customersTableViewDefinition,
		myPortfolioTableViewDefinition,
		leadsTableViewDefinition,
		nurtureTableViewDefinition,
		churnTableViewDefinition,
	}
}

func DefaultTableViewDefinitionMonthlyRenewals(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := postgresEntity.Columns{
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
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeRenewals.String(),
		TableId:     model.TableIDTypeMonthlyRenewals.String(),
		Name:        "Monthly renewals",
		ColumnsJson: string(jsonData),
		Order:       1,
		Icon:        "ClockFastForward",
		Filters:     `{"AND":[{"filter":{"property":"RENEWAL_CYCLE","value":"MONTHLY","operation":"EQ","includeEmpty":false}}]}`,
		Sorting:     "",
	}, nil
}

func DefaultTableViewDefinitionQuarterlyRenewals(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := postgresEntity.Columns{
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
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeRenewals.String(),
		TableId:     model.TableIDTypeQuarterlyRenewals.String(),
		Name:        "Quarterly renewals",
		ColumnsJson: string(jsonData),
		Order:       2,
		Icon:        "ClockFastForward",
		Filters:     `{"AND":[{"filter":{"property":"RENEWAL_CYCLE","value":"QUARTERLY","operation":"EQ","includeEmpty":false}}]}`,
		Sorting:     "",
	}, nil
}

func DefaultTableViewDefinitionAnnualRenewals(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := postgresEntity.Columns{
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
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeRenewals.String(),
		TableId:     model.TableIDTypeQuarterlyRenewals.String(),
		Name:        "Annual renewals",
		ColumnsJson: string(jsonData),
		Order:       3,
		Icon:        "ClockFastForward",
		Filters:     `{"AND":[{"filter":{"property":"RENEWAL_CYCLE","value":"ANNUALLY","operation":"EQ","includeEmpty":false}}]}`,
		Sorting:     "",
	}, nil
}

func DefaultTableViewDefinitionPastInvoices(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
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
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeInvoices.String(),
		TableId:     model.TableIDTypePastInvoices.String(),
		Name:        "Past",
		ColumnsJson: string(jsonData),
		Order:       5,
		Icon:        "InvoiceCheck",
		Filters:     `{"AND":[{"filter":{"property":"INVOICE_DRY_RUN","value":false}}]}`,
		Sorting:     "",
	}, nil
}

func DefaultTableViewDefinitionUpcomingInvoices(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
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
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeInvoices.String(),
		TableId:     model.TableIDTypeUpcomingInvoices.String(),
		Name:        "Upcoming",
		ColumnsJson: string(jsonData),
		Order:       4,
		Icon:        "InvoiceUpcoming",
		Filters:     `{"AND":[{"filter":{"property":"INVOICE_PREVIEW","value":true}}]}`,
		Sorting:     "",
	}, nil
}

func DefaultTableViewDefinitionOrganization(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
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
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeOrganizations.String(),
		TableId:     model.TableIDTypeOrganizations.String(),
		Name:        "All orgs",
		ColumnsJson: string(jsonData),
		Order:       5,
		Icon:        "Building07",
		Filters:     ``,
		Sorting:     "",
	}, nil
}

func DefaultTableViewDefinitionCustomers(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
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
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeOrganizations.String(),
		TableId:     model.TableIDTypeCustomers.String(),
		Name:        "Customers",
		ColumnsJson: string(jsonData),
		Order:       1,
		Icon:        "CheckHeart",
		Filters:     fmt.Sprintf(`{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"RELATIONSHIP","value":["%s"]}}]}`, neo4jenum.Customer.String()),
		Sorting:     "",
	}, nil
}

func DefaultTableViewDefinitionMyPortfolio(userId string, span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
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
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeOrganizations.String(),
		TableId:     model.TableIDTypeMyPortfolio.String(),
		Name:        "My portfolio",
		ColumnsJson: string(jsonData),
		Order:       6,
		Icon:        "Briefcase01",
		Filters:     fmt.Sprintf(`{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"OWNER_ID","value":["%s"]}}]}`, userId),
		Sorting:     "",
	}, nil
}

func DefaultTableViewDefinitionLeads(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := postgresEntity.Columns{
		Columns: []postgresEntity.ColumnView{
			{ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsSocials.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsCreatedDate.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsLastTouchpointDate.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsLeadSource.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsEmployeeCount.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsYearFounded.String(), Width: 100, Visible: true},
		},
	}
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeOrganizations.String(),
		TableId:     model.TableIDTypeLeads.String(),
		Name:        "Leads",
		ColumnsJson: string(jsonData),
		Order:       4,
		Icon:        "SwitchHorizontal01",
		Filters:     `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"STAGE","value":["LEAD"]}}]}`,
		Sorting:     "",
	}, nil
}

func DefaultTableViewDefinitionNurture(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := postgresEntity.Columns{
		Columns: []postgresEntity.ColumnView{
			{ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsSocials.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsLastTouchpoint.String(), Width: 100, Visible: true},
		},
	}
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeOrganizations.String(),
		TableId:     model.TableIDTypeNurture.String(),
		Name:        "Nurture",
		ColumnsJson: string(jsonData),
		Order:       3,
		Icon:        "HeartHand",
		Filters:     fmt.Sprintf(`{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"STAGE","value":["%s"]}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"RELATIONSHIP","value":["%s"]}}]}`, neo4jenum.Target.String(), neo4jenum.Prospect.String()),
		Sorting:     "",
	}, nil
}

func DefaultTableViewDefinitionChurn(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := postgresEntity.Columns{
		Columns: []postgresEntity.ColumnView{
			{ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true},
			{ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 100, Visible: true},
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
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeOrganizations.String(),
		TableId:     model.TableIDTypeChurn.String(),
		Name:        "Churn",
		ColumnsJson: string(jsonData),
		Order:       5,
		Icon:        "HeartHand",
		Filters:     fmt.Sprintf(`{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"STAGE","value":["%s"]}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"RELATIONSHIP","value":["%s"]}}]}`, neo4jenum.PendingChurn.String(), neo4jenum.FormerCustomer.String()),
		Sorting:     "",
	}, nil
}
