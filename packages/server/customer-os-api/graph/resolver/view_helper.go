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

	contactsTableViewDefinition, err := DefaultTableViewDefinitionContacts(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	targetOrganizationContactsTableViewDefinition, err := DefaultTableViewDefinitionTargetOrganizationsContacts(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	opportunitiesTableViewDeifnition, err := DefaultTableViewDefinitionOpportunities(span)
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
		contactsTableViewDefinition,
		targetOrganizationContactsTableViewDefinition,
		opportunitiesTableViewDeifnition,
	}
}

func DefaultTableViewDefinitionMonthlyRenewals(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeMonthlyRenewals.String())
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
		Sorting:     ``,
		IsPreset:    false,
	}, nil
}

func DefaultTableViewDefinitionQuarterlyRenewals(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeQuarterlyRenewals.String())
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
		Sorting:     ``,
		IsPreset:    false,
	}, nil
}

func DefaultTableViewDefinitionAnnualRenewals(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeAnnualRenewals.String())
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeRenewals.String(),
		TableId:     model.TableIDTypeAnnualRenewals.String(),
		Name:        "Annual renewals",
		ColumnsJson: string(jsonData),
		Order:       3,
		Icon:        "ClockFastForward",
		Filters:     `{"AND":[{"filter":{"property":"RENEWAL_CYCLE","value":"ANNUALLY","operation":"EQ","includeEmpty":false}}]}`,
		Sorting:     "",
		IsPreset:    false,
	}, nil
}

func DefaultTableViewDefinitionPastInvoices(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypePastInvoices.String())
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
		Sorting:     ``,
		IsPreset:    false,
	}, nil
}

func DefaultTableViewDefinitionUpcomingInvoices(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeUpcomingInvoices.String())
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
		Sorting:     ``,
		IsPreset:    false,
	}, nil
}

func DefaultTableViewDefinitionOrganization(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeOrganizations.String())
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
		Sorting:     `{"id": "ORGANIZATIONS_LAST_TOUCHPOINT", "desc": true}`,
		IsPreset:    false,
	}, nil
}

func DefaultTableViewDefinitionCustomers(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeCustomers.String())
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
		Sorting:     `{"id": "ORGANIZATIONS_LAST_TOUCHPOINT", "desc": true}`,
		IsPreset:    false,
	}, nil
}

func DefaultTableViewDefinitionMyPortfolio(userId string, span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeMyPortfolio.String())
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
		Sorting:     `{"id": "ORGANIZATIONS_LAST_TOUCHPOINT", "desc": true}`,
		IsPreset:    false,
	}, nil
}

func DefaultTableViewDefinitionLeads(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeLeads.String())
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
		Order:       3,
		Icon:        "SwitchHorizontal01",
		Filters:     `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"STAGE","value":["LEAD"]}}]}`,
		Sorting:     `{"id": "ORGANIZATIONS_LAST_TOUCHPOINT", "desc": true}`,
		IsPreset:    false,
	}, nil
}

func DefaultTableViewDefinitionNurture(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeNurture.String())
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeOrganizations.String(),
		TableId:     model.TableIDTypeNurture.String(),
		Name:        "Targets",
		ColumnsJson: string(jsonData),
		Order:       4,
		Icon:        "HeartHand",
		Filters:     fmt.Sprintf(`{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"STAGE","value":["%s"]}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"RELATIONSHIP","value":["%s"]}}]}`, neo4jenum.Target.String(), neo4jenum.Prospect.String()),
		Sorting:     `{"id": "ORGANIZATIONS_LAST_TOUCHPOINT", "desc": true}`,
		IsPreset:    false,
	}, nil
}

func DefaultTableViewDefinitionChurn(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeChurn.String())
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
		Icon:        "BrokenHeart",
		Filters:     fmt.Sprintf(`{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"RELATIONSHIP","value":["%s"]}}]}`, neo4jenum.FormerCustomer.String()),
		Sorting:     `{"id": "ORGANIZATIONS_CHURN_DATE", "desc": true}`,
		IsPreset:    false,
	}, nil
}

func DefaultTableViewDefinitionContacts(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeContacts.String())
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeContacts.String(),
		TableId:     model.TableIDTypeContacts.String(),
		Name:        "All Contacts",
		ColumnsJson: string(jsonData),
		Order:       0,
		Icon:        "HeartHand",
		Filters:     ``,
		Sorting:     ``,
		IsPreset:    false,
	}, nil
}

func DefaultTableViewDefinitionTargetOrganizationsContacts(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeContactsForTargetOrganizations.String())
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeContacts.String(),
		TableId:     model.TableIDTypeContactsForTargetOrganizations.String(),
		Name:        "Contacts",
		ColumnsJson: string(jsonData),
		Order:       0,
		Icon:        "HeartHand",
		Filters:     fmt.Sprintf(`{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"STAGE","value":["%s"]}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"RELATIONSHIP","value":["%s"]}}]}`, neo4jenum.Target.String(), neo4jenum.Prospect.String()),
		Sorting:     ``,
		IsPreset:    false,
	}, nil
}

func DefaultTableViewDefinitionOpportunities(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeOpportunities.String())
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeOpportunities.String(),
		TableId:     model.TableIDTypeOpportunities.String(),
		Name:        "Opportunities",
		ColumnsJson: string(jsonData),
		Order:       6,
		Icon:        "CoinsStacked01",
		Filters:     ``,
		Sorting:     ``,
		IsPreset:    true,
	}, nil
}

func DefaultColumns(tableId string) postgresEntity.Columns {
	switch tableId {
	case model.TableIDTypeChurn.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsChurnDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsLtv.String(), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeNurture.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsSocials.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsCreatedDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsLastTouchpoint.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsLeadSource.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsEmployeeCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsYearFounded.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsIndustry.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsCity.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsIsPublic.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsStage.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsLinkedinFollowerCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsTags.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsContactCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeLeads.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsSocials.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsCreatedDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsLastTouchpointDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsLeadSource.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsEmployeeCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsYearFounded.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsIndustry.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsCity.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsIsPublic.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsStage.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsLinkedinFollowerCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsTags.String(), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeMyPortfolio.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsRelationship.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsRenewalLikelihood.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsRenewalDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsOnboardingStatus.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsForecastArr.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsOwner.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsLastTouchpoint.String(), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeCustomers.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsRelationship.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsRenewalLikelihood.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsRenewalDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsOnboardingStatus.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsForecastArr.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsOwner.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsLastTouchpoint.String(), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeOrganizations.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsRelationship.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsRenewalLikelihood.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsRenewalDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsOnboardingStatus.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsForecastArr.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsOwner.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsContactCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsStage.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeOrganizationsLastTouchpoint.String(), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeUpcomingInvoices.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnType: model.ColumnViewTypeInvoicesInvoicePreview.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesContract.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesBillingCycle.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesIssueDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesDueDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesAmount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesInvoiceStatus.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesIssueDatePast.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesPaymentStatus.String(), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypePastInvoices.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnType: model.ColumnViewTypeInvoicesInvoiceNumber.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesContract.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesBillingCycle.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesIssueDatePast.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesDueDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesAmount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesPaymentStatus.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesIssueDate.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeInvoicesInvoiceStatus.String(), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeAnnualRenewals.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnType: model.ColumnViewTypeRenewalsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsRenewalDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsForecastArr.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsRenewalLikelihood.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsOwner.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsLastTouchpoint.String(), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeQuarterlyRenewals.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnType: model.ColumnViewTypeRenewalsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsRenewalDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsForecastArr.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsRenewalLikelihood.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsOwner.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsLastTouchpoint.String(), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeMonthlyRenewals.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnType: model.ColumnViewTypeRenewalsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsRenewalDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsForecastArr.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsRenewalLikelihood.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsOwner.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeRenewalsLastTouchpoint.String(), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeContacts.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnType: model.ColumnViewTypeContactsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsOrganization.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsEmails.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsPhoneNumbers.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsLinkedin.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsCountry.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsRegion.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsCity.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsPersona.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsLastInteraction.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsSkills.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsSchools.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsLanguages.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsTimeInCurrentRole.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsExperience.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsLinkedinFollowerCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsJobTitle.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsConnections.String(), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeContactsForTargetOrganizations.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnType: model.ColumnViewTypeContactsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsOrganization.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsEmails.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsPhoneNumbers.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsLinkedin.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsCountry.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsRegion.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsCity.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsPersona.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsLastInteraction.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsSkills.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsSchools.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsLanguages.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsTimeInCurrentRole.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsExperience.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsLinkedinFollowerCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsJobTitle.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnType: model.ColumnViewTypeContactsConnections.String(), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeOpportunities.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnType: model.ColumnViewTypeOpportunitiesCommonColumn.String(), Width: 100, Visible: true, Name: "Identified", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"externalStage","value":"STAGE1"}}]}`},
				{ColumnType: model.ColumnViewTypeOpportunitiesCommonColumn.String(), Width: 100, Visible: true, Name: "Qualified", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"externalStage","value":"STAGE2"}}]}`},
				{ColumnType: model.ColumnViewTypeOpportunitiesCommonColumn.String(), Width: 100, Visible: true, Name: "Committed", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"externalStage","value":"STAGE3"}}]}`},
				{ColumnType: model.ColumnViewTypeOpportunitiesCommonColumn.String(), Width: 100, Visible: true, Name: "Won", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"CLOSED_WON"}}]}`},
				{ColumnType: model.ColumnViewTypeOpportunitiesCommonColumn.String(), Width: 100, Visible: true, Name: "Lost", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"CLOSED_LOST"}}]}`},
			},
		}
	}
	return postgresEntity.Columns{}
}
