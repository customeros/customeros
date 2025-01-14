package service

import (
	"encoding/json"
	"fmt"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
)

// ColumnView represents a column in a table view with type and width.
func DefaultTableViewDefinitions(hasSharedPresets bool, span opentracing.Span) []postgresEntity.TableViewDefinition {
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
	columns := DefaultColumns(postgresEntity.TableIDTypePastInvoices)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:      string(postgresEntity.TableViewTypeInvoices),
		TableId:        string(postgresEntity.TableIDTypePastInvoices),
		Name:           "Past",
		ColumnsJson:    string(jsonData),
		Order:          5,
		Icon:           "InvoiceCheck",
		Filters:        ``,
		DefaultFilters: `{"AND":[{"filter":{"property":"INVOICE_DRY_RUN","value":false}}]}`,
		Sorting:        ``,
		IsPreset:       true,
		IsShared:       false,
	}, nil
}

func DefaultTableViewDefinitionUpcomingInvoices(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(postgresEntity.TableIDTypeUpcomingInvoices)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:      string(postgresEntity.TableViewTypeInvoices),
		TableId:        string(postgresEntity.TableIDTypeUpcomingInvoices),
		Name:           "Upcoming",
		ColumnsJson:    string(jsonData),
		Order:          4,
		Icon:           "InvoiceUpcoming",
		Filters:        ``,
		DefaultFilters: `{"AND":[{"filter":{"property":"INVOICE_PREVIEW","value":true}}]}`,
		Sorting:        ``,
		IsPreset:       true,
		IsShared:       false,
	}, nil
}

func DefaultTableViewDefinitionOrganization(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(postgresEntity.TableIDTypeOrganizations)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:      string(postgresEntity.TableViewTypeOrganizations),
		TableId:        string(postgresEntity.TableIDTypeOrganizations),
		Name:           "Organizations",
		ColumnsJson:    string(jsonData),
		Order:          5,
		Icon:           "Building07",
		Filters:        ``,
		DefaultFilters: ``,
		Sorting:        `{"id": "ORGANIZATIONS_UPDATED_DATE", "desc": true}`,
		IsPreset:       true,
		IsShared:       false,
	}, nil
}

func DefaultTableViewDefinitionCustomers(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(postgresEntity.TableIDTypeCustomers)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:      string(postgresEntity.TableViewTypeOrganizations),
		TableId:        string(postgresEntity.TableIDTypeCustomers),
		Name:           "Customers",
		ColumnsJson:    string(jsonData),
		Order:          1,
		Icon:           "CheckHeart",
		Filters:        ``,
		DefaultFilters: fmt.Sprintf(`{"AND":[{"filter":{"includeEmpty":false,"operation":"IN","property":"ORGANIZATIONS_RELATIONSHIP","value":["%s"],"active":true}}]}`, neo4jenum.OrganizationRelationshipCustomer),
		Sorting:        `{"id": "ORGANIZATIONS_UPDATED_DATE", "desc": true}`,
		IsPreset:       true,
		IsShared:       false,
	}, nil
}

func DefaultTableViewDefinitionTargets(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(postgresEntity.TableIDTypeTargets)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:      string(postgresEntity.TableViewTypeOrganizations),
		TableId:        string(postgresEntity.TableIDTypeTargets),
		Name:           "Targets",
		ColumnsJson:    string(jsonData),
		Order:          1,
		Icon:           "Target05",
		Filters:        ``,
		DefaultFilters: fmt.Sprintf(`{"AND":[{"filter":{"includeEmpty":false,"operation":"IN","property":"ORGANIZATIONS_STAGE","value":["%s"],"active":true}},{"filter":{"includeEmpty":false,"operation":"IN","property":"ORGANIZATIONS_RELATIONSHIP","value":["%s"],"active":true}}]}`, neo4jenum.Target, neo4jenum.OrganizationRelationshipProspect),
		Sorting:        `{"id": "ORGANIZATIONS_UPDATED_DATE", "desc": true}`,
		IsPreset:       true,
		IsShared:       false,
	}, nil
}

func DefaultTableViewDefinitionContacts(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(postgresEntity.TableIDTypeContacts)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:      string(postgresEntity.TableViewTypeContacts),
		TableId:        string(postgresEntity.TableIDTypeContacts),
		Name:           "Contacts",
		ColumnsJson:    string(jsonData),
		Order:          0,
		Icon:           "Users01",
		Filters:        ``,
		DefaultFilters: ``,
		Sorting:        `{"id": "CONTACTS_UPDATED_AT", "desc": true}`,
		IsPreset:       true,
		IsShared:       false,
	}, nil
}

func DefaultTableViewDefinitionTargetOrganizationsContacts(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(postgresEntity.TableIDTypeContactsForTargetOrganizations)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:      string(postgresEntity.TableViewTypeContacts),
		TableId:        string(postgresEntity.TableIDTypeContactsForTargetOrganizations),
		Name:           "Contacts",
		ColumnsJson:    string(jsonData),
		Order:          0,
		Icon:           "HeartHand",
		Filters:        ``,
		DefaultFilters: fmt.Sprintf(`{"AND":[{"filter":{"includeEmpty":false,"operation":"IN","property":"ORGANIZATIONS_STAGE","value":["%s"],"active":true}},{"filter":{"includeEmpty":false,"operation":"IN","property":"ORGANIZATIONS_RELATIONSHIP","value":["%s"],"active":true}}]}`, neo4jenum.Target, neo4jenum.OrganizationRelationshipProspect),
		Sorting:        `{"id": "CONTACTS_UPDATED_AT", "desc": true}`,
		IsPreset:       true,
		IsShared:       false,
	}, nil
}

func DefaultTableViewDefinitionOpportunities(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(postgresEntity.TableIDTypeOpportunities)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:      string(postgresEntity.TableViewTypeOpportunities),
		TableId:        string(postgresEntity.TableIDTypeOpportunities),
		Name:           "Opportunities",
		ColumnsJson:    string(jsonData),
		Order:          6,
		Icon:           "CoinsStacked01",
		Filters:        ``,
		DefaultFilters: ``,
		Sorting:        ``,
		IsPreset:       true,
		IsShared:       true,
	}, nil
}

func DefaultTableViewDefinitionOpportunitiesRecords(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(postgresEntity.TableIDTypeOpportunitiesRecords)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:      string(postgresEntity.TableViewTypeOpportunities),
		TableId:        string(postgresEntity.TableIDTypeOpportunitiesRecords),
		Name:           "Opportunities",
		ColumnsJson:    string(jsonData),
		Order:          7,
		Icon:           "CoinsStacked01",
		Filters:        ``,
		DefaultFilters: ``,
		Sorting:        ``,
		IsPreset:       true,
		IsShared:       false,
	}, nil
}

func DefaultTableViewDefinitionContracts(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(postgresEntity.TableIDTypeContracts)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:      string(postgresEntity.TableViewTypeContracts),
		TableId:        string(postgresEntity.TableIDTypeContracts),
		Name:           "Contracts",
		ColumnsJson:    string(jsonData),
		Order:          8,
		Icon:           "Signature",
		Filters:        ``,
		DefaultFilters: ``,
		Sorting:        ``,
		IsPreset:       true,
		IsShared:       false,
	}, nil
}

func DefaultTableViewDefinitionFlows(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(postgresEntity.TableIDTypeFlowActions)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:      string(postgresEntity.TableViewTypeFlow),
		TableId:        string(postgresEntity.TableIDTypeFlowActions),
		Name:           "Flows",
		ColumnsJson:    string(jsonData),
		Order:          9,
		Icon:           "Shuffle01",
		Filters:        ``,
		DefaultFilters: ``,
		Sorting:        ``,
		IsPreset:       true,
		IsShared:       false,
	}, nil
}

func DefaultTableViewDefinitionFlowContacts(span opentracing.Span) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(postgresEntity.TableIDTypeFlowContacts)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:      string(postgresEntity.TableViewTypeContacts),
		TableId:        string(postgresEntity.TableIDTypeFlowContacts),
		Name:           "Contacts",
		ColumnsJson:    string(jsonData),
		Order:          0,
		Icon:           "Users01",
		Filters:        ``,
		DefaultFilters: ``,
		Sorting:        `{"id": "CONTACTS_UPDATED_AT", "desc": true}`,
		IsPreset:       true,
		IsShared:       false,
	}, nil
}

func DefaultTableViewDefinitionFlowContactsV2(span opentracing.Span, flowId string) (postgresEntity.TableViewDefinition, error) {
	columns := DefaultColumns(postgresEntity.TableIDTypeFlowContacts)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgresEntity.TableViewDefinition{}, err
	}

	return postgresEntity.TableViewDefinition{
		TableType:      string(postgresEntity.TableViewTypeContacts),
		TableId:        string(postgresEntity.TableIDTypeFlowContacts),
		Name:           "Flow contacts",
		ColumnsJson:    string(jsonData),
		Order:          0,
		Icon:           "Users01",
		Filters:        ``,
		DefaultFilters: fmt.Sprintf(`{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"FLOW_ID","value":["%s"],"active":true}}]}`, flowId),
		Sorting:        `{"id": "CONTACTS_UPDATED_AT", "desc": true}`,
		IsPreset:       true,
		IsShared:       true,
	}, nil
}

func DefaultColumns(tableId postgresEntity.TableIdType) postgresEntity.Columns {
	switch tableId {
	case postgresEntity.TableIDTypeCustomers:
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsAvatar), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsName), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsWebsite), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsPrimaryDomains), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsRelationship), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsRenewalLikelihood), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsRenewalDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsOnboardingStatus), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsForecastArr), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsOwner), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsLastTouchpoint), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsParentOrganization), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 12, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsCountry), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case postgresEntity.TableIDTypeOrganizations:
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsAvatar), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsName), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsWebsite), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 28, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsPrimaryDomains), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 22, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsIndustry), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 16, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsSocials), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 23, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsIsPublic), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 20, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsEmployeeCount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 21, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsYearFounded), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsRelationship), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsRenewalLikelihood), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsRenewalDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsOnboardingStatus), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsForecastArr), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsOwner), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsContactCount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsStage), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsChurnDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 14, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsLtv), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 15, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsTags), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 17, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsCreatedDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 18, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsLastTouchpointDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 19, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsLeadSource), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 25, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsParentOrganization), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 26, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsUpdatedDate), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 27, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsCountry), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case postgresEntity.TableIDTypeTargets:
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsAvatar), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsName), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsWebsite), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 20, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsPrimaryDomains), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsSocials), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsCreatedDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsLastTouchpoint), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsLeadSource), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsEmployeeCount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsYearFounded), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsIndustry), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsCity), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 12, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsIsPublic), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsStage), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 14, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsLinkedinFollowerCount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 15, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsTags), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 16, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsContactCount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 17, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsParentOrganization), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 18, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsUpdatedDate), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 19, ColumnType: string(postgresEntity.ColumnViewTypeOrganizationsCountry), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case postgresEntity.TableIDTypeUpcomingInvoices:
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesInvoicePreview), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesContract), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesBillingCycle), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesIssueDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesDueDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesAmount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesInvoiceStatus), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesIssueDatePast), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesOrganization), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	case postgresEntity.TableIDTypePastInvoices:
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesInvoiceNumber), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesContract), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesBillingCycle), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesIssueDatePast), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesDueDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesAmount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesIssueDate), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesInvoiceStatus), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgresEntity.ColumnViewTypeInvoicesOrganization), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	case postgresEntity.TableIDTypeContacts:
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgresEntity.ColumnViewTypeContactsAvatar), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgresEntity.ColumnViewTypeContactsName), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgresEntity.ColumnViewTypeContactsOrganization), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgresEntity.ColumnViewTypeContactsLinkedin), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgresEntity.ColumnViewTypeContactsEmails), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 22, ColumnType: string(postgresEntity.ColumnViewTypeContactsPersonalEmails), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 23, ColumnType: string(postgresEntity.ColumnViewTypeContactsPrimaryEmail), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgresEntity.ColumnViewTypeContactsPhoneNumbers), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgresEntity.ColumnViewTypeContactsCountry), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgresEntity.ColumnViewTypeContactsRegion), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgresEntity.ColumnViewTypeContactsCity), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 18, ColumnType: string(postgresEntity.ColumnViewTypeContactsJobTitle), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 15, ColumnType: string(postgresEntity.ColumnViewTypeContactsTimeInCurrentRole), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 17, ColumnType: string(postgresEntity.ColumnViewTypeContactsLinkedinFollowerCount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 19, ColumnType: string(postgresEntity.ColumnViewTypeContactsConnections), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 21, ColumnType: string(postgresEntity.ColumnViewTypeContactsFlows), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgresEntity.ColumnViewTypeContactsPersona), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgresEntity.ColumnViewTypeContactsLastInteraction), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 12, ColumnType: string(postgresEntity.ColumnViewTypeContactsSkills), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: string(postgresEntity.ColumnViewTypeContactsSchools), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 14, ColumnType: string(postgresEntity.ColumnViewTypeContactsLanguages), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 16, ColumnType: string(postgresEntity.ColumnViewTypeContactsExperience), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 24, ColumnType: string(postgresEntity.ColumnViewTypeContactsUpdatedAt), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 26, ColumnType: string(postgresEntity.ColumnViewTypeContactsCreatedAt), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	case postgresEntity.TableIDTypeContactsForTargetOrganizations:
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgresEntity.ColumnViewTypeContactsAvatar), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgresEntity.ColumnViewTypeContactsName), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgresEntity.ColumnViewTypeContactsOrganization), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgresEntity.ColumnViewTypeContactsLinkedin), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgresEntity.ColumnViewTypeContactsEmails), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 22, ColumnType: string(postgresEntity.ColumnViewTypeContactsPersonalEmails), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 23, ColumnType: string(postgresEntity.ColumnViewTypeContactsPrimaryEmail), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgresEntity.ColumnViewTypeContactsPhoneNumbers), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgresEntity.ColumnViewTypeContactsCountry), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgresEntity.ColumnViewTypeContactsRegion), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgresEntity.ColumnViewTypeContactsCity), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 18, ColumnType: string(postgresEntity.ColumnViewTypeContactsJobTitle), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 15, ColumnType: string(postgresEntity.ColumnViewTypeContactsTimeInCurrentRole), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 17, ColumnType: string(postgresEntity.ColumnViewTypeContactsLinkedinFollowerCount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 19, ColumnType: string(postgresEntity.ColumnViewTypeContactsConnections), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 21, ColumnType: string(postgresEntity.ColumnViewTypeContactsFlows), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgresEntity.ColumnViewTypeContactsPersona), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgresEntity.ColumnViewTypeContactsLastInteraction), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 12, ColumnType: string(postgresEntity.ColumnViewTypeContactsSkills), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: string(postgresEntity.ColumnViewTypeContactsSchools), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 14, ColumnType: string(postgresEntity.ColumnViewTypeContactsLanguages), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 16, ColumnType: string(postgresEntity.ColumnViewTypeContactsExperience), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 24, ColumnType: string(postgresEntity.ColumnViewTypeContactsUpdatedAt), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 26, ColumnType: string(postgresEntity.ColumnViewTypeContactsCreatedAt), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	case postgresEntity.TableIDTypeOpportunities:
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgresEntity.ColumnViewTypeOpportunitiesCommonColumn), Width: 100, Visible: true, Name: "Identified", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"externalStage","value":"STAGE1"}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"OPEN"}}]}`},
				{ColumnId: 2, ColumnType: string(postgresEntity.ColumnViewTypeOpportunitiesCommonColumn), Width: 100, Visible: true, Name: "Qualified", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"externalStage","value":"STAGE2"}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"OPEN"}}]}`},
				{ColumnId: 3, ColumnType: string(postgresEntity.ColumnViewTypeOpportunitiesCommonColumn), Width: 100, Visible: true, Name: "Committed", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"externalStage","value":"STAGE3"}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"OPEN"}}]}`},
				{ColumnId: 4, ColumnType: string(postgresEntity.ColumnViewTypeOpportunitiesCommonColumn), Width: 100, Visible: true, Name: "Won", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"CLOSED_WON"}}]}`},
				{ColumnId: 5, ColumnType: string(postgresEntity.ColumnViewTypeOpportunitiesCommonColumn), Width: 100, Visible: true, Name: "Lost", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"CLOSED_LOST"}}]}`},
			},
		}
	case postgresEntity.TableIDTypeOpportunitiesRecords:
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgresEntity.ColumnViewTypeOpportunitiesName), Width: 100, Visible: true, Name: "Name", Filter: ``},
				{ColumnId: 2, ColumnType: string(postgresEntity.ColumnViewTypeOpportunitiesOrganization), Width: 100, Visible: true, Name: "Organization", Filter: ``},
				{ColumnId: 3, ColumnType: string(postgresEntity.ColumnViewTypeOpportunitiesStage), Width: 100, Visible: true, Name: "Stage", Filter: ``},
				{ColumnId: 4, ColumnType: string(postgresEntity.ColumnViewTypeOpportunitiesEstimatedArr), Width: 100, Visible: true, Name: "Estimated ARR", Filter: ``},
				{ColumnId: 5, ColumnType: string(postgresEntity.ColumnViewTypeOpportunitiesOwner), Width: 100, Visible: true, Name: "Owner", Filter: ``},
				{ColumnId: 6, ColumnType: string(postgresEntity.ColumnViewTypeOpportunitiesTimeInStage), Width: 100, Visible: true, Name: "Time in Stage", Filter: ``},
				{ColumnId: 7, ColumnType: string(postgresEntity.ColumnViewTypeOpportunitiesCreatedDate), Width: 100, Visible: true, Name: "Created", Filter: ``},
				{ColumnId: 8, ColumnType: string(postgresEntity.ColumnViewTypeOpportunitiesNextStep), Width: 100, Visible: true, Name: "Next Step", Filter: ``},
			},
		}
	case postgresEntity.TableIDTypeContracts:
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgresEntity.ColumnViewTypeContractsName), Width: 100, Visible: true, Name: "Name", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgresEntity.ColumnViewTypeContractsEnded), Width: 100, Visible: true, Name: "Ended", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgresEntity.ColumnViewTypeContractsPeriod), Width: 100, Visible: true, Name: "Period", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgresEntity.ColumnViewTypeContractsCurrency), Width: 100, Visible: true, Name: "Currency", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgresEntity.ColumnViewTypeContractsStatus), Width: 100, Visible: true, Name: "Status", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgresEntity.ColumnViewTypeContractsRenewal), Width: 100, Visible: true, Name: "Renewal", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgresEntity.ColumnViewTypeContractsLtv), Width: 100, Visible: true, Name: "LTV", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgresEntity.ColumnViewTypeContractsRenewalDate), Width: 100, Visible: true, Name: "Renewal Date", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgresEntity.ColumnViewTypeContractsForecastArr), Width: 100, Visible: true, Name: "ARR Forecast", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgresEntity.ColumnViewTypeContractsHealth), Width: 100, Visible: true, Name: "Health", Filter: ""},
				{ColumnId: 12, ColumnType: string(postgresEntity.ColumnViewTypeContractsOwner), Width: 100, Visible: true, Name: "Owner", Filter: ""},
			},
		}
	case postgresEntity.TableIDTypeFlowActions:
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 2, ColumnType: string(postgresEntity.ColumnViewTypeFlowName), Width: 100, Visible: true, Name: "Flow", Filter: ""},
				{ColumnId: 1, ColumnType: string(postgresEntity.ColumnViewTypeFlowActionName), Width: 100, Visible: true, Name: "Status", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgresEntity.ColumnViewTypeFlowTotalCount), Width: 100, Visible: true, Name: "Contacts", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgresEntity.ColumnViewTypeFlowOnHoldCount), Width: 100, Visible: true, Name: "Blocked", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgresEntity.ColumnViewTypeFlowReadyCount), Width: 100, Visible: true, Name: "Ready", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgresEntity.ColumnViewTypeFlowScheduledCount), Width: 100, Visible: true, Name: "Scheduled", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgresEntity.ColumnViewTypeFlowInProgressCount), Width: 100, Visible: true, Name: "In Progress", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgresEntity.ColumnViewTypeFlowCompletedCount), Width: 100, Visible: true, Name: "Completed", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgresEntity.ColumnViewTypeFlowGoalAchievedCount), Width: 100, Visible: true, Name: "Goal achieved", Filter: ""},
			},
		}
	case postgresEntity.TableIDTypeFlowContacts:
		return postgresEntity.Columns{
			Columns: []postgresEntity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgresEntity.ColumnViewTypeContactsAvatar), Width: 100, Visible: true, Name: "Avatar", Filter: ""},
				{ColumnId: 13, ColumnType: string(postgresEntity.ColumnViewTypeContactsFlowStatus), Width: 100, Visible: true, Name: "Status in Flow", Filter: ""},
				{ColumnId: 14, ColumnType: string(postgresEntity.ColumnViewTypeContactsFlowNextAction), Width: 100, Visible: true, Name: "Next action", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgresEntity.ColumnViewTypeContactsName), Width: 100, Visible: true, Name: "Name", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgresEntity.ColumnViewTypeContactsOrganization), Width: 100, Visible: true, Name: "Organization", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgresEntity.ColumnViewTypeContactsPrimaryEmail), Width: 100, Visible: true, Name: "Primary email", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgresEntity.ColumnViewTypeContactsEmails), Width: 100, Visible: true, Name: "Emails", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgresEntity.ColumnViewTypeContactsPhoneNumbers), Width: 100, Visible: true, Name: "Phone numbers", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgresEntity.ColumnViewTypeContactsLinkedin), Width: 100, Visible: true, Name: "LinkedIn", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgresEntity.ColumnViewTypeContactsJobTitle), Width: 100, Visible: true, Name: "Job title", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgresEntity.ColumnViewTypeContactsTimeInCurrentRole), Width: 100, Visible: true, Name: "Time In Current Role", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgresEntity.ColumnViewTypeContactsCountry), Width: 100, Visible: true, Name: "Country", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgresEntity.ColumnViewTypeContactsRegion), Width: 100, Visible: true, Name: "Region", Filter: ""},
				{ColumnId: 12, ColumnType: string(postgresEntity.ColumnViewTypeContactsCity), Width: 100, Visible: true, Name: "City", Filter: ""},
				{ColumnId: 15, ColumnType: string(postgresEntity.ColumnViewTypeContactsUpdatedAt), Width: 100, Visible: false, Name: "Updated at", Filter: ""},
				{ColumnId: 16, ColumnType: string(postgresEntity.ColumnViewTypeContactsCreatedAt), Width: 100, Visible: false, Name: "Created at", Filter: ""},
			},
		}
	}
	return postgresEntity.Columns{}
}

func CheckSharedPresetsExist(viewDefs []postgresEntity.TableViewDefinition) bool {
	for _, def := range viewDefs {
		if def.IsShared && def.TableType == string(postgresEntity.TableViewTypeOpportunities) {
			return true
		}
	}
	return false
}
