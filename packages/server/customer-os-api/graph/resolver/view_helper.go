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
func DefaultTableViewDefinitions(userId string, hasSharedPresets bool, span opentracing.Span) []postgresEntity.TableViewDefinition {
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

	targetsTableViewDefinition, err := DefaultTableViewDefinitionTargets(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	pastInvoicesTableViewDefinition, err := DefaultTableViewDefinitionPastInvoices(span)
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

	opportunitiesTableViewDefinition, err := DefaultTableViewDefinitionOpportunities(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	opportunitiesRecordsTableViewDefinition, err := DefaultTableViewDefinitionOpportunitiesRecords(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	contractsTableViewDefinition, err := DefaultTableViewDefinitionContracts(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	flowsTableViewDefinition, err := DefaultTableViewDefinitionFlows(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	flowContactsTableViewDefinition, err := DefaultTableViewDefinitionFlowContacts(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgresEntity.TableViewDefinition{}
	}

	defaultViewDefinitions := []postgresEntity.TableViewDefinition{
		upcomingInvoicesTableViewDefinition,
		pastInvoicesTableViewDefinition,
		organizationsTableViewDefinition,
		customersTableViewDefinition,
		contactsTableViewDefinition,
		targetOrganizationContactsTableViewDefinition,
		contractsTableViewDefinition,
		targetsTableViewDefinition,
		opportunitiesRecordsTableViewDefinition,
		flowsTableViewDefinition,
		flowContactsTableViewDefinition,
	}

	if !hasSharedPresets {
		defaultViewDefinitions = append(defaultViewDefinitions, opportunitiesTableViewDefinition)
	}

	return defaultViewDefinitions
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
		IsPreset:    true,
		IsShared:    false,
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
		IsPreset:    true,
		IsShared:    false,
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
		Name:        "Organizations",
		ColumnsJson: string(jsonData),
		Order:       5,
		Icon:        "Building07",
		Filters:     ``,
		Sorting:     `{"id": "ORGANIZATIONS_LAST_TOUCHPOINT", "desc": true}`,
		IsPreset:    true,
		IsShared:    false,
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
		IsPreset:    true,
		IsShared:    false,
	}, nil
}

func DefaultTableViewDefinitionTargets(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeTargets.String())
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeOrganizations.String(),
		TableId:     model.TableIDTypeTargets.String(),
		Name:        "Targets",
		ColumnsJson: string(jsonData),
		Order:       1,
		Icon:        "Target05",
		Filters:     fmt.Sprintf(`{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"STAGE","value":["%s"]}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"RELATIONSHIP","value":["%s"]}}]}`, neo4jenum.Target.String(), neo4jenum.Prospect.String()),
		Sorting:     `{"id": "ORGANIZATIONS_LAST_TOUCHPOINT", "desc": true}`,
		IsPreset:    true,
		IsShared:    false,
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
		Name:        "Contacts",
		ColumnsJson: string(jsonData),
		Order:       0,
		Icon:        "Users01",
		Filters:     ``,
		Sorting:     `{"id": "CONTACTS_UPDATED_AT", "desc": true}`,
		IsPreset:    true,
		IsShared:    false,
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
		Sorting:     `{"id": "CONTACTS_UPDATED_AT", "desc": true}`,
		IsPreset:    true,
		IsShared:    false,
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
		IsShared:    true,
	}, nil
}

func DefaultTableViewDefinitionOpportunitiesRecords(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeOpportunitiesRecords.String())
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeOpportunities.String(),
		TableId:     model.TableIDTypeOpportunitiesRecords.String(),
		Name:        "Opportunities",
		ColumnsJson: string(jsonData),
		Order:       7,
		Icon:        "CoinsStacked01",
		Filters:     ``,
		Sorting:     ``,
		IsPreset:    true,
		IsShared:    false,
	}, nil
}

func DefaultTableViewDefinitionContracts(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeContracts.String())
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeContracts.String(),
		TableId:     model.TableIDTypeContracts.String(),
		Name:        "Contracts",
		ColumnsJson: string(jsonData),
		Order:       8,
		Icon:        "Signature",
		Filters:     ``,
		Sorting:     ``,
		IsPreset:    true,
		IsShared:    false,
	}, nil
}

func DefaultTableViewDefinitionFlows(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeFlowActions.String())
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeFlow.String(),
		TableId:     model.TableIDTypeFlowActions.String(),
		Name:        "Flows",
		ColumnsJson: string(jsonData),
		Order:       9,
		Icon:        "Shuffle01",
		Filters:     ``,
		Sorting:     ``,
		IsPreset:    true,
		IsShared:    false,
	}, nil
}

func DefaultTableViewDefinitionFlowContacts(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(model.TableIDTypeFlowContacts.String())
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:   model.TableViewTypeContacts.String(),
		TableId:     model.TableIDTypeFlowContacts.String(),
		Name:        "Contacts",
		ColumnsJson: string(jsonData),
		Order:       0,
		Icon:        "Users01",
		Filters:     ``,
		Sorting:     `{"id": "CONTACTS_UPDATED_AT", "desc": true}`,
		IsPreset:    true,
		IsShared:    false,
	}, nil
}

func DefaultColumns(tableId string) postgresEntity.Columns {
	switch tableId {
	case model.TableIDTypeCustomers.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: model.ColumnViewTypeOrganizationsRelationship.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: model.ColumnViewTypeOrganizationsRenewalLikelihood.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: model.ColumnViewTypeOrganizationsRenewalDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: model.ColumnViewTypeOrganizationsOnboardingStatus.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: model.ColumnViewTypeOrganizationsForecastArr.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: model.ColumnViewTypeOrganizationsOwner.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: model.ColumnViewTypeOrganizationsLastTouchpoint.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: model.ColumnViewTypeOrganizationsParentOrganization.String(), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeOrganizations.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 24, ColumnType: model.ColumnViewTypeOrganizationsHeadquarters.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 22, ColumnType: model.ColumnViewTypeOrganizationsIndustry.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 16, ColumnType: model.ColumnViewTypeOrganizationsSocials.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 23, ColumnType: model.ColumnViewTypeOrganizationsIsPublic.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 20, ColumnType: model.ColumnViewTypeOrganizationsEmployeeCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 21, ColumnType: model.ColumnViewTypeOrganizationsYearFounded.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: model.ColumnViewTypeOrganizationsRelationship.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: model.ColumnViewTypeOrganizationsRenewalLikelihood.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: model.ColumnViewTypeOrganizationsRenewalDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: model.ColumnViewTypeOrganizationsOnboardingStatus.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: model.ColumnViewTypeOrganizationsForecastArr.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: model.ColumnViewTypeOrganizationsOwner.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: model.ColumnViewTypeOrganizationsContactCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: model.ColumnViewTypeOrganizationsStage.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 12, ColumnType: model.ColumnViewTypeOrganizationsLastTouchpoint.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: model.ColumnViewTypeOrganizationsChurnDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 14, ColumnType: model.ColumnViewTypeOrganizationsLtv.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 15, ColumnType: model.ColumnViewTypeOrganizationsTags.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 17, ColumnType: model.ColumnViewTypeOrganizationsCreatedDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 18, ColumnType: model.ColumnViewTypeOrganizationsLastTouchpointDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 19, ColumnType: model.ColumnViewTypeOrganizationsLeadSource.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 25, ColumnType: model.ColumnViewTypeOrganizationsParentOrganization.String(), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeTargets.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: model.ColumnViewTypeOrganizationsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: model.ColumnViewTypeOrganizationsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: model.ColumnViewTypeOrganizationsWebsite.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: model.ColumnViewTypeOrganizationsSocials.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: model.ColumnViewTypeOrganizationsCreatedDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: model.ColumnViewTypeOrganizationsLastTouchpoint.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: model.ColumnViewTypeOrganizationsLeadSource.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: model.ColumnViewTypeOrganizationsEmployeeCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: model.ColumnViewTypeOrganizationsYearFounded.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: model.ColumnViewTypeOrganizationsIndustry.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: model.ColumnViewTypeOrganizationsCity.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 12, ColumnType: model.ColumnViewTypeOrganizationsIsPublic.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: model.ColumnViewTypeOrganizationsStage.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 14, ColumnType: model.ColumnViewTypeOrganizationsLinkedinFollowerCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 15, ColumnType: model.ColumnViewTypeOrganizationsTags.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 16, ColumnType: model.ColumnViewTypeOrganizationsContactCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 17, ColumnType: model.ColumnViewTypeOrganizationsParentOrganization.String(), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeUpcomingInvoices.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: model.ColumnViewTypeInvoicesInvoicePreview.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: model.ColumnViewTypeInvoicesContract.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: model.ColumnViewTypeInvoicesBillingCycle.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: model.ColumnViewTypeInvoicesIssueDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: model.ColumnViewTypeInvoicesDueDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: model.ColumnViewTypeInvoicesAmount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: model.ColumnViewTypeInvoicesInvoiceStatus.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: model.ColumnViewTypeInvoicesIssueDatePast.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: model.ColumnViewTypeInvoicesOrganization.String(), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypePastInvoices.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: model.ColumnViewTypeInvoicesInvoiceNumber.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: model.ColumnViewTypeInvoicesContract.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: model.ColumnViewTypeInvoicesBillingCycle.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: model.ColumnViewTypeInvoicesIssueDatePast.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: model.ColumnViewTypeInvoicesDueDate.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: model.ColumnViewTypeInvoicesAmount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: model.ColumnViewTypeInvoicesIssueDate.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: model.ColumnViewTypeInvoicesInvoiceStatus.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: model.ColumnViewTypeInvoicesOrganization.String(), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeContacts.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: model.ColumnViewTypeContactsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: model.ColumnViewTypeContactsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: model.ColumnViewTypeContactsOrganization.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: model.ColumnViewTypeContactsLinkedin.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: model.ColumnViewTypeContactsEmails.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 22, ColumnType: model.ColumnViewTypeContactsPersonalEmails.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 23, ColumnType: model.ColumnViewTypeContactsPrimaryEmail.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: model.ColumnViewTypeContactsPhoneNumbers.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: model.ColumnViewTypeContactsCountry.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: model.ColumnViewTypeContactsRegion.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: model.ColumnViewTypeContactsCity.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 18, ColumnType: model.ColumnViewTypeContactsJobTitle.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 15, ColumnType: model.ColumnViewTypeContactsTimeInCurrentRole.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 17, ColumnType: model.ColumnViewTypeContactsLinkedinFollowerCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 19, ColumnType: model.ColumnViewTypeContactsConnections.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 21, ColumnType: model.ColumnViewTypeContactsFlows.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: model.ColumnViewTypeContactsPersona.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: model.ColumnViewTypeContactsLastInteraction.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 12, ColumnType: model.ColumnViewTypeContactsSkills.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: model.ColumnViewTypeContactsSchools.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 14, ColumnType: model.ColumnViewTypeContactsLanguages.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 16, ColumnType: model.ColumnViewTypeContactsExperience.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 17, ColumnType: model.ColumnViewTypeContactsFlowStatus.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 24, ColumnType: model.ColumnViewTypeContactsUpdatedAt.String(), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeContactsForTargetOrganizations.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: model.ColumnViewTypeContactsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: model.ColumnViewTypeContactsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: model.ColumnViewTypeContactsOrganization.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: model.ColumnViewTypeContactsLinkedin.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: model.ColumnViewTypeContactsEmails.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 22, ColumnType: model.ColumnViewTypeContactsPersonalEmails.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 23, ColumnType: model.ColumnViewTypeContactsPrimaryEmail.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: model.ColumnViewTypeContactsPhoneNumbers.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: model.ColumnViewTypeContactsCountry.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: model.ColumnViewTypeContactsRegion.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: model.ColumnViewTypeContactsCity.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 18, ColumnType: model.ColumnViewTypeContactsJobTitle.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 15, ColumnType: model.ColumnViewTypeContactsTimeInCurrentRole.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 17, ColumnType: model.ColumnViewTypeContactsLinkedinFollowerCount.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 19, ColumnType: model.ColumnViewTypeContactsConnections.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 21, ColumnType: model.ColumnViewTypeContactsFlows.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: model.ColumnViewTypeContactsPersona.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: model.ColumnViewTypeContactsLastInteraction.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 12, ColumnType: model.ColumnViewTypeContactsSkills.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: model.ColumnViewTypeContactsSchools.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 14, ColumnType: model.ColumnViewTypeContactsLanguages.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 16, ColumnType: model.ColumnViewTypeContactsExperience.String(), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 17, ColumnType: model.ColumnViewTypeContactsFlowStatus.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 24, ColumnType: model.ColumnViewTypeContactsUpdatedAt.String(), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	case model.TableIDTypeOpportunities.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: model.ColumnViewTypeOpportunitiesCommonColumn.String(), Width: 100, Visible: true, Name: "Identified", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"externalStage","value":"STAGE1"}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"OPEN"}}]}`},
				{ColumnId: 2, ColumnType: model.ColumnViewTypeOpportunitiesCommonColumn.String(), Width: 100, Visible: true, Name: "Qualified", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"externalStage","value":"STAGE2"}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"OPEN"}}]}`},
				{ColumnId: 3, ColumnType: model.ColumnViewTypeOpportunitiesCommonColumn.String(), Width: 100, Visible: true, Name: "Committed", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"externalStage","value":"STAGE3"}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"OPEN"}}]}`},
				{ColumnId: 4, ColumnType: model.ColumnViewTypeOpportunitiesCommonColumn.String(), Width: 100, Visible: true, Name: "Won", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"CLOSED_WON"}}]}`},
				{ColumnId: 5, ColumnType: model.ColumnViewTypeOpportunitiesCommonColumn.String(), Width: 100, Visible: true, Name: "Lost", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"CLOSED_LOST"}}]}`},
			},
		}
	case model.TableIDTypeOpportunitiesRecords.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: model.ColumnViewTypeOpportunitiesName.String(), Width: 100, Visible: true, Name: "Name", Filter: ``},
				{ColumnId: 2, ColumnType: model.ColumnViewTypeOpportunitiesOrganization.String(), Width: 100, Visible: true, Name: "Organization", Filter: ``},
				{ColumnId: 3, ColumnType: model.ColumnViewTypeOpportunitiesStage.String(), Width: 100, Visible: true, Name: "Stage", Filter: ``},
				{ColumnId: 4, ColumnType: model.ColumnViewTypeOpportunitiesEstimatedArr.String(), Width: 100, Visible: true, Name: "Estimated ARR", Filter: ``},
				{ColumnId: 5, ColumnType: model.ColumnViewTypeOpportunitiesOwner.String(), Width: 100, Visible: true, Name: "Owner", Filter: ``},
				{ColumnId: 6, ColumnType: model.ColumnViewTypeOpportunitiesTimeInStage.String(), Width: 100, Visible: true, Name: "Time in Stage", Filter: ``},
				{ColumnId: 7, ColumnType: model.ColumnViewTypeOpportunitiesCreatedDate.String(), Width: 100, Visible: true, Name: "Created", Filter: ``},
				{ColumnId: 8, ColumnType: model.ColumnViewTypeOpportunitiesNextStep.String(), Width: 100, Visible: true, Name: "Next Step", Filter: ``},
			},
		}
	case model.TableIDTypeContracts.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: model.ColumnViewTypeContractsName.String(), Width: 100, Visible: true, Name: "Name", Filter: ""},
				{ColumnId: 2, ColumnType: model.ColumnViewTypeContractsEnded.String(), Width: 100, Visible: true, Name: "Ended", Filter: ""},
				{ColumnId: 3, ColumnType: model.ColumnViewTypeContractsPeriod.String(), Width: 100, Visible: true, Name: "Period", Filter: ""},
				{ColumnId: 5, ColumnType: model.ColumnViewTypeContractsCurrency.String(), Width: 100, Visible: true, Name: "Currency", Filter: ""},
				{ColumnId: 6, ColumnType: model.ColumnViewTypeContractsStatus.String(), Width: 100, Visible: true, Name: "Status", Filter: ""},
				{ColumnId: 7, ColumnType: model.ColumnViewTypeContractsRenewal.String(), Width: 100, Visible: true, Name: "Renewal", Filter: ""},
				{ColumnId: 8, ColumnType: model.ColumnViewTypeContractsLtv.String(), Width: 100, Visible: true, Name: "LTV", Filter: ""},
				{ColumnId: 9, ColumnType: model.ColumnViewTypeContractsRenewalDate.String(), Width: 100, Visible: true, Name: "Renewal Date", Filter: ""},
				{ColumnId: 10, ColumnType: model.ColumnViewTypeContractsForecastArr.String(), Width: 100, Visible: true, Name: "ARR Forecast", Filter: ""},
				{ColumnId: 11, ColumnType: model.ColumnViewTypeContractsHealth.String(), Width: 100, Visible: true, Name: "Health", Filter: ""},
				{ColumnId: 12, ColumnType: model.ColumnViewTypeContractsOwner.String(), Width: 100, Visible: true, Name: "Owner", Filter: ""},
			},
		}
	case model.TableIDTypeFlowActions.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 2, ColumnType: model.ColumnViewTypeFlowName.String(), Width: 100, Visible: true, Name: "Flow", Filter: ""},
				{ColumnId: 1, ColumnType: model.ColumnViewTypeFlowActionName.String(), Width: 100, Visible: true, Name: "Status", Filter: ""},
				{ColumnId: 5, ColumnType: model.ColumnViewTypeFlowTotalCount.String(), Width: 100, Visible: true, Name: "Contacts", Filter: ""},
				{ColumnId: 6, ColumnType: model.ColumnViewTypeFlowOnHoldCount.String(), Width: 100, Visible: true, Name: "On Hold", Filter: ""},
				{ColumnId: 7, ColumnType: model.ColumnViewTypeFlowReadyCount.String(), Width: 100, Visible: true, Name: "Ready", Filter: ""},
				{ColumnId: 8, ColumnType: model.ColumnViewTypeFlowScheduledCount.String(), Width: 100, Visible: true, Name: "Scheduled", Filter: ""},
				{ColumnId: 9, ColumnType: model.ColumnViewTypeFlowInProgressCount.String(), Width: 100, Visible: true, Name: "In Progress", Filter: ""},
				{ColumnId: 10, ColumnType: model.ColumnViewTypeFlowCompletedCount.String(), Width: 100, Visible: true, Name: "Completed", Filter: ""},
				{ColumnId: 11, ColumnType: model.ColumnViewTypeFlowGoalAchievedCount.String(), Width: 100, Visible: true, Name: "Goal achieved", Filter: ""},
			},
		}
	case model.TableIDTypeFlowContacts.String():
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: model.ColumnViewTypeContactsAvatar.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: model.ColumnViewTypeContactsName.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: model.ColumnViewTypeContactsOrganization.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 14, ColumnType: model.ColumnViewTypeContactsPrimaryEmail.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: model.ColumnViewTypeContactsEmails.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: model.ColumnViewTypeContactsPhoneNumbers.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: model.ColumnViewTypeContactsLinkedin.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: model.ColumnViewTypeContactsJobTitle.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: model.ColumnViewTypeContactsTimeInCurrentRole.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: model.ColumnViewTypeContactsCountry.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: model.ColumnViewTypeContactsRegion.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: model.ColumnViewTypeContactsCity.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 12, ColumnType: model.ColumnViewTypeContactsFlowStatus.String(), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: model.ColumnViewTypeContactsUpdatedAt.String(), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	}
	return postgresEntity.Columns{}
}

func CheckSharedPresetsExist(viewDefs []postgresEntity.TableViewDefinition) bool {
	for _, def := range viewDefs {
		if def.IsShared && def.TableType == model.TableViewTypeOpportunities.String() {
			return true
		}
	}
	return false
}
