package table_view

import (
	"encoding/json"
	"fmt"

	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

// ColumnView represents a column in a table view with type and width.
func DefaultTableViewDefinitions(hasSharedPresets bool, span opentracing.Span) []postgres_entity.TableViewDefinition {
	upcomingInvoicesTableViewDefinition, err := DefaultTableViewDefinitionUpcomingInvoices(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgres_entity.TableViewDefinition{}
	}

	organizationsTableViewDefinition, err := DefaultTableViewDefinitionOrganization(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgres_entity.TableViewDefinition{}
	}

	customersTableViewDefinition, err := DefaultTableViewDefinitionCustomers(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgres_entity.TableViewDefinition{}
	}

	targetsTableViewDefinition, err := DefaultTableViewDefinitionTargets(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgres_entity.TableViewDefinition{}
	}

	pastInvoicesTableViewDefinition, err := DefaultTableViewDefinitionPastInvoices(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgres_entity.TableViewDefinition{}
	}

	contactsTableViewDefinition, err := DefaultTableViewDefinitionContacts(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgres_entity.TableViewDefinition{}
	}

	contactsForTargetOrganizations, err := DefaultTableViewDefinitionTargetOrganizationsContacts(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgres_entity.TableViewDefinition{}
	}

	opportunitiesTableViewDefinition, err := DefaultTableViewDefinitionOpportunities(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgres_entity.TableViewDefinition{}
	}

	opportunitiesRecordsTableViewDefinition, err := DefaultTableViewDefinitionOpportunitiesRecords(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgres_entity.TableViewDefinition{}
	}

	contractsTableViewDefinition, err := DefaultTableViewDefinitionContracts(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgres_entity.TableViewDefinition{}
	}

	flowsTableViewDefinition, err := DefaultTableViewDefinitionFlows(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgres_entity.TableViewDefinition{}
	}

	flowContactsTableViewDefinition, err := DefaultTableViewDefinitionFlowContacts(span)
	if err != nil {
		fmt.Println("Error: ", err)
		return []postgres_entity.TableViewDefinition{}
	}

	defaultViewDefinitions := []postgres_entity.TableViewDefinition{
		upcomingInvoicesTableViewDefinition,
		pastInvoicesTableViewDefinition,
		organizationsTableViewDefinition,
		customersTableViewDefinition,
		targetsTableViewDefinition,
		contactsTableViewDefinition,
		contactsForTargetOrganizations,
		contractsTableViewDefinition,
		opportunitiesRecordsTableViewDefinition,
		flowsTableViewDefinition,
		flowContactsTableViewDefinition,
	}

	if !hasSharedPresets {
		defaultViewDefinitions = append(defaultViewDefinitions, opportunitiesTableViewDefinition)
	}

	return defaultViewDefinitions
}

func DefaultTableViewDefinitionPastInvoices(span opentracing.Span) (postgres_entity.TableViewDefinition, error) {
	columns := DefaultColumns(postgres_entity.TableIDTypePastInvoices)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgres_entity.TableViewDefinition{}, err
	}

	return postgres_entity.TableViewDefinition{
		TableType:      string(postgres_entity.TableViewTypeInvoices),
		TableId:        string(postgres_entity.TableIDTypePastInvoices),
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

func DefaultTableViewDefinitionUpcomingInvoices(span opentracing.Span) (postgres_entity.TableViewDefinition, error) {
	columns := DefaultColumns(postgres_entity.TableIDTypeUpcomingInvoices)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgres_entity.TableViewDefinition{}, err
	}

	return postgres_entity.TableViewDefinition{
		TableType:      string(postgres_entity.TableViewTypeInvoices),
		TableId:        string(postgres_entity.TableIDTypeUpcomingInvoices),
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

func DefaultTableViewDefinitionOrganization(span opentracing.Span) (postgres_entity.TableViewDefinition, error) {
	columns := DefaultColumns(postgres_entity.TableIDTypeOrganizations)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgres_entity.TableViewDefinition{}, err
	}

	return postgres_entity.TableViewDefinition{
		TableType:      string(postgres_entity.TableViewTypeOrganizations),
		TableId:        string(postgres_entity.TableIDTypeOrganizations),
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

func DefaultTableViewDefinitionCustomers(span opentracing.Span) (postgres_entity.TableViewDefinition, error) {
	columns := DefaultColumns(postgres_entity.TableIDTypeCustomers)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgres_entity.TableViewDefinition{}, err
	}

	return postgres_entity.TableViewDefinition{
		TableType:      string(postgres_entity.TableViewTypeOrganizations),
		TableId:        string(postgres_entity.TableIDTypeCustomers),
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

func DefaultTableViewDefinitionTargets(span opentracing.Span) (postgres_entity.TableViewDefinition, error) {
	columns := DefaultColumns(postgres_entity.TableIDTypeTargets)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgres_entity.TableViewDefinition{}, err
	}

	return postgres_entity.TableViewDefinition{
		TableType:      string(postgres_entity.TableViewTypeOrganizations),
		TableId:        string(postgres_entity.TableIDTypeTargets),
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

func DefaultTableViewDefinitionContacts(span opentracing.Span) (postgres_entity.TableViewDefinition, error) {
	columns := DefaultColumns(postgres_entity.TableIDTypeContacts)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgres_entity.TableViewDefinition{}, err
	}

	return postgres_entity.TableViewDefinition{
		TableType:      string(postgres_entity.TableViewTypeContacts),
		TableId:        string(postgres_entity.TableIDTypeContacts),
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

func DefaultTableViewDefinitionTargetOrganizationsContacts(span opentracing.Span) (postgres_entity.TableViewDefinition, error) {
	columns := DefaultColumns(postgres_entity.TableIDTypeContactsForTargetOrganizations)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgres_entity.TableViewDefinition{}, err
	}

	return postgres_entity.TableViewDefinition{
		TableType:      string(postgres_entity.TableViewTypeContacts),
		TableId:        string(postgres_entity.TableIDTypeContactsForTargetOrganizations),
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

func DefaultTableViewDefinitionOpportunities(span opentracing.Span) (postgres_entity.TableViewDefinition, error) {
	columns := DefaultColumns(postgres_entity.TableIDTypeOpportunities)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgres_entity.TableViewDefinition{}, err
	}

	return postgres_entity.TableViewDefinition{
		TableType:      string(postgres_entity.TableViewTypeOpportunities),
		TableId:        string(postgres_entity.TableIDTypeOpportunities),
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

func DefaultTableViewDefinitionOpportunitiesRecords(span opentracing.Span) (postgres_entity.TableViewDefinition, error) {
	columns := DefaultColumns(postgres_entity.TableIDTypeOpportunitiesRecords)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgres_entity.TableViewDefinition{}, err
	}

	return postgres_entity.TableViewDefinition{
		TableType:      string(postgres_entity.TableViewTypeOpportunities),
		TableId:        string(postgres_entity.TableIDTypeOpportunitiesRecords),
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

func DefaultTableViewDefinitionContracts(span opentracing.Span) (postgres_entity.TableViewDefinition, error) {
	columns := DefaultColumns(postgres_entity.TableIDTypeContracts)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgres_entity.TableViewDefinition{}, err
	}

	return postgres_entity.TableViewDefinition{
		TableType:      string(postgres_entity.TableViewTypeContracts),
		TableId:        string(postgres_entity.TableIDTypeContracts),
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

func DefaultTableViewDefinitionFlows(span opentracing.Span) (postgres_entity.TableViewDefinition, error) {
	columns := DefaultColumns(postgres_entity.TableIDTypeFlowActions)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgres_entity.TableViewDefinition{}, err
	}

	return postgres_entity.TableViewDefinition{
		TableType:      string(postgres_entity.TableViewTypeFlow),
		TableId:        string(postgres_entity.TableIDTypeFlowActions),
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

func DefaultTableViewDefinitionFlowContacts(span opentracing.Span) (postgres_entity.TableViewDefinition, error) {
	columns := DefaultColumns(postgres_entity.TableIDTypeFlowContacts)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgres_entity.TableViewDefinition{}, err
	}

	return postgres_entity.TableViewDefinition{
		TableType:      string(postgres_entity.TableViewTypeContacts),
		TableId:        string(postgres_entity.TableIDTypeFlowContacts),
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

func DefaultTableViewDefinitionFlowContactsV2(span opentracing.Span, flowId string) (postgres_entity.TableViewDefinition, error) {
	columns := DefaultColumns(postgres_entity.TableIDTypeFlowContacts)
	jsonData, err := json.Marshal(columns)
	if err != nil {
		tracing.TraceErr(span, err)
		fmt.Println("Error serializing data:", err)
		return postgres_entity.TableViewDefinition{}, err
	}

	return postgres_entity.TableViewDefinition{
		TableType:      string(postgres_entity.TableViewTypeContacts),
		TableId:        string(postgres_entity.TableIDTypeFlowContacts),
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

func DefaultColumns(tableId postgres_entity.TableIdType) postgres_entity.Columns {
	switch tableId {
	case postgres_entity.TableIDTypeCustomers:
		return postgres_entity.Columns{
			Columns: []postgres_entity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsAvatar), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsName), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsWebsite), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsPrimaryDomains), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsRelationship), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsRenewalLikelihood), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsRenewalDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsOnboardingStatus), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsForecastArr), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsOwner), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsLastTouchpoint), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsParentOrganization), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 12, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsCountry), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case postgres_entity.TableIDTypeOrganizations:
		return postgres_entity.Columns{
			Columns: []postgres_entity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsAvatar), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsName), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsWebsite), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 28, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsPrimaryDomains), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 22, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsIndustry), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 16, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsSocials), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 23, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsIsPublic), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 20, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsEmployeeCount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 21, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsYearFounded), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsRelationship), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsRenewalLikelihood), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsRenewalDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsOnboardingStatus), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsForecastArr), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsOwner), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsContactCount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsStage), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsChurnDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 14, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsLtv), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 15, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsTags), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 17, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsCreatedDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 18, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsLastTouchpointDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 19, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsLeadSource), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 25, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsParentOrganization), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 26, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsUpdatedDate), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 27, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsCountry), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case postgres_entity.TableIDTypeTargets:
		return postgres_entity.Columns{
			Columns: []postgres_entity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsAvatar), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsName), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsWebsite), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 20, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsPrimaryDomains), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsSocials), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsCreatedDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsLastTouchpoint), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsLeadSource), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsEmployeeCount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsYearFounded), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsIndustry), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsCity), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 12, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsIsPublic), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsStage), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 14, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsLinkedinFollowerCount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 15, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsTags), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 16, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsContactCount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 17, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsParentOrganization), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 18, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsUpdatedDate), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 19, ColumnType: string(postgres_entity.ColumnViewTypeOrganizationsCountry), Width: 100, Visible: true, Name: "", Filter: ""},
			},
		}
	case postgres_entity.TableIDTypeUpcomingInvoices:
		return postgres_entity.Columns{
			Columns: []postgres_entity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesInvoicePreview), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesContract), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesBillingCycle), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesIssueDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesDueDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesAmount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesInvoiceStatus), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesIssueDatePast), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesOrganization), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	case postgres_entity.TableIDTypePastInvoices:
		return postgres_entity.Columns{
			Columns: []postgres_entity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesInvoiceNumber), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesContract), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesBillingCycle), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesIssueDatePast), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesDueDate), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesAmount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesIssueDate), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesInvoiceStatus), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgres_entity.ColumnViewTypeInvoicesOrganization), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	case postgres_entity.TableIDTypeContacts:
		return postgres_entity.Columns{
			Columns: []postgres_entity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgres_entity.ColumnViewTypeContactsAvatar), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgres_entity.ColumnViewTypeContactsName), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgres_entity.ColumnViewTypeContactsOrganization), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgres_entity.ColumnViewTypeContactsLinkedin), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgres_entity.ColumnViewTypeContactsEmails), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 22, ColumnType: string(postgres_entity.ColumnViewTypeContactsPersonalEmails), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 23, ColumnType: string(postgres_entity.ColumnViewTypeContactsPrimaryEmail), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgres_entity.ColumnViewTypeContactsPhoneNumbers), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgres_entity.ColumnViewTypeContactsCountry), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgres_entity.ColumnViewTypeContactsRegion), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgres_entity.ColumnViewTypeContactsCity), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 18, ColumnType: string(postgres_entity.ColumnViewTypeContactsJobTitle), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 15, ColumnType: string(postgres_entity.ColumnViewTypeContactsTimeInCurrentRole), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 17, ColumnType: string(postgres_entity.ColumnViewTypeContactsLinkedinFollowerCount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 19, ColumnType: string(postgres_entity.ColumnViewTypeContactsConnections), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 21, ColumnType: string(postgres_entity.ColumnViewTypeContactsFlows), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgres_entity.ColumnViewTypeContactsPersona), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgres_entity.ColumnViewTypeContactsLastInteraction), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 12, ColumnType: string(postgres_entity.ColumnViewTypeContactsSkills), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: string(postgres_entity.ColumnViewTypeContactsSchools), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 14, ColumnType: string(postgres_entity.ColumnViewTypeContactsLanguages), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 16, ColumnType: string(postgres_entity.ColumnViewTypeContactsExperience), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 24, ColumnType: string(postgres_entity.ColumnViewTypeContactsUpdatedAt), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 26, ColumnType: string(postgres_entity.ColumnViewTypeContactsCreatedAt), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	case postgres_entity.TableIDTypeContactsForTargetOrganizations:
		return postgres_entity.Columns{
			Columns: []postgres_entity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgres_entity.ColumnViewTypeContactsAvatar), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgres_entity.ColumnViewTypeContactsName), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgres_entity.ColumnViewTypeContactsOrganization), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgres_entity.ColumnViewTypeContactsLinkedin), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgres_entity.ColumnViewTypeContactsEmails), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 22, ColumnType: string(postgres_entity.ColumnViewTypeContactsPersonalEmails), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 23, ColumnType: string(postgres_entity.ColumnViewTypeContactsPrimaryEmail), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgres_entity.ColumnViewTypeContactsPhoneNumbers), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgres_entity.ColumnViewTypeContactsCountry), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgres_entity.ColumnViewTypeContactsRegion), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgres_entity.ColumnViewTypeContactsCity), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 18, ColumnType: string(postgres_entity.ColumnViewTypeContactsJobTitle), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 15, ColumnType: string(postgres_entity.ColumnViewTypeContactsTimeInCurrentRole), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 17, ColumnType: string(postgres_entity.ColumnViewTypeContactsLinkedinFollowerCount), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 19, ColumnType: string(postgres_entity.ColumnViewTypeContactsConnections), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 21, ColumnType: string(postgres_entity.ColumnViewTypeContactsFlows), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgres_entity.ColumnViewTypeContactsPersona), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgres_entity.ColumnViewTypeContactsLastInteraction), Width: 100, Visible: true, Name: "", Filter: ""},
				{ColumnId: 12, ColumnType: string(postgres_entity.ColumnViewTypeContactsSkills), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 13, ColumnType: string(postgres_entity.ColumnViewTypeContactsSchools), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 14, ColumnType: string(postgres_entity.ColumnViewTypeContactsLanguages), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 16, ColumnType: string(postgres_entity.ColumnViewTypeContactsExperience), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 24, ColumnType: string(postgres_entity.ColumnViewTypeContactsUpdatedAt), Width: 100, Visible: false, Name: "", Filter: ""},
				{ColumnId: 26, ColumnType: string(postgres_entity.ColumnViewTypeContactsCreatedAt), Width: 100, Visible: false, Name: "", Filter: ""},
			},
		}
	case postgres_entity.TableIDTypeOpportunities:
		return postgres_entity.Columns{
			Columns: []postgres_entity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgres_entity.ColumnViewTypeOpportunitiesCommonColumn), Width: 100, Visible: true, Name: "Identified", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"externalStage","value":"STAGE1"}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"OPEN"}}]}`},
				{ColumnId: 2, ColumnType: string(postgres_entity.ColumnViewTypeOpportunitiesCommonColumn), Width: 100, Visible: true, Name: "Qualified", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"externalStage","value":"STAGE2"}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"OPEN"}}]}`},
				{ColumnId: 3, ColumnType: string(postgres_entity.ColumnViewTypeOpportunitiesCommonColumn), Width: 100, Visible: true, Name: "Committed", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"externalStage","value":"STAGE3"}},{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"OPEN"}}]}`},
				{ColumnId: 4, ColumnType: string(postgres_entity.ColumnViewTypeOpportunitiesCommonColumn), Width: 100, Visible: true, Name: "Won", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"CLOSED_WON"}}]}`},
				{ColumnId: 5, ColumnType: string(postgres_entity.ColumnViewTypeOpportunitiesCommonColumn), Width: 100, Visible: true, Name: "Lost", Filter: `{"AND":[{"filter":{"includeEmpty":false,"operation":"EQ","property":"internalStage","value":"CLOSED_LOST"}}]}`},
			},
		}
	case postgres_entity.TableIDTypeOpportunitiesRecords:
		return postgres_entity.Columns{
			Columns: []postgres_entity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgres_entity.ColumnViewTypeOpportunitiesName), Width: 100, Visible: true, Name: "Name", Filter: ``},
				{ColumnId: 2, ColumnType: string(postgres_entity.ColumnViewTypeOpportunitiesOrganization), Width: 100, Visible: true, Name: "Organization", Filter: ``},
				{ColumnId: 3, ColumnType: string(postgres_entity.ColumnViewTypeOpportunitiesStage), Width: 100, Visible: true, Name: "Stage", Filter: ``},
				{ColumnId: 4, ColumnType: string(postgres_entity.ColumnViewTypeOpportunitiesEstimatedArr), Width: 100, Visible: true, Name: "Estimated ARR", Filter: ``},
				{ColumnId: 5, ColumnType: string(postgres_entity.ColumnViewTypeOpportunitiesOwner), Width: 100, Visible: true, Name: "Owner", Filter: ``},
				{ColumnId: 6, ColumnType: string(postgres_entity.ColumnViewTypeOpportunitiesTimeInStage), Width: 100, Visible: true, Name: "Time in Stage", Filter: ``},
				{ColumnId: 7, ColumnType: string(postgres_entity.ColumnViewTypeOpportunitiesCreatedDate), Width: 100, Visible: true, Name: "Created", Filter: ``},
				{ColumnId: 8, ColumnType: string(postgres_entity.ColumnViewTypeOpportunitiesNextStep), Width: 100, Visible: true, Name: "Next Step", Filter: ``},
			},
		}
	case postgres_entity.TableIDTypeContracts:
		return postgres_entity.Columns{
			Columns: []postgres_entity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgres_entity.ColumnViewTypeContractsName), Width: 100, Visible: true, Name: "Name", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgres_entity.ColumnViewTypeContractsEnded), Width: 100, Visible: true, Name: "Ended", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgres_entity.ColumnViewTypeContractsPeriod), Width: 100, Visible: true, Name: "Period", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgres_entity.ColumnViewTypeContractsCurrency), Width: 100, Visible: true, Name: "Currency", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgres_entity.ColumnViewTypeContractsStatus), Width: 100, Visible: true, Name: "Status", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgres_entity.ColumnViewTypeContractsRenewal), Width: 100, Visible: true, Name: "Renewal", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgres_entity.ColumnViewTypeContractsLtv), Width: 100, Visible: true, Name: "LTV", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgres_entity.ColumnViewTypeContractsRenewalDate), Width: 100, Visible: true, Name: "Renewal Date", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgres_entity.ColumnViewTypeContractsForecastArr), Width: 100, Visible: true, Name: "ARR Forecast", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgres_entity.ColumnViewTypeContractsHealth), Width: 100, Visible: true, Name: "Health", Filter: ""},
				{ColumnId: 12, ColumnType: string(postgres_entity.ColumnViewTypeContractsOwner), Width: 100, Visible: true, Name: "Owner", Filter: ""},
			},
		}
	case postgres_entity.TableIDTypeFlowActions:
		return postgres_entity.Columns{
			Columns: []postgres_entity.ColumnView{
				{ColumnId: 2, ColumnType: string(postgres_entity.ColumnViewTypeFlowName), Width: 100, Visible: true, Name: "Flow", Filter: ""},
				{ColumnId: 1, ColumnType: string(postgres_entity.ColumnViewTypeFlowActionName), Width: 100, Visible: true, Name: "Status", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgres_entity.ColumnViewTypeFlowTotalCount), Width: 100, Visible: true, Name: "Contacts", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgres_entity.ColumnViewTypeFlowOnHoldCount), Width: 100, Visible: true, Name: "Blocked", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgres_entity.ColumnViewTypeFlowReadyCount), Width: 100, Visible: true, Name: "Ready", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgres_entity.ColumnViewTypeFlowScheduledCount), Width: 100, Visible: true, Name: "Scheduled", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgres_entity.ColumnViewTypeFlowInProgressCount), Width: 100, Visible: true, Name: "In Progress", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgres_entity.ColumnViewTypeFlowCompletedCount), Width: 100, Visible: true, Name: "Completed", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgres_entity.ColumnViewTypeFlowGoalAchievedCount), Width: 100, Visible: true, Name: "Goal achieved", Filter: ""},
			},
		}
	case postgres_entity.TableIDTypeFlowContacts:
		return postgres_entity.Columns{
			Columns: []postgres_entity.ColumnView{
				{ColumnId: 1, ColumnType: string(postgres_entity.ColumnViewTypeContactsAvatar), Width: 100, Visible: true, Name: "Avatar", Filter: ""},
				{ColumnId: 13, ColumnType: string(postgres_entity.ColumnViewTypeContactsFlowStatus), Width: 100, Visible: true, Name: "Status in Flow", Filter: ""},
				{ColumnId: 14, ColumnType: string(postgres_entity.ColumnViewTypeContactsFlowNextAction), Width: 100, Visible: true, Name: "Next action", Filter: ""},
				{ColumnId: 2, ColumnType: string(postgres_entity.ColumnViewTypeContactsName), Width: 100, Visible: true, Name: "Name", Filter: ""},
				{ColumnId: 3, ColumnType: string(postgres_entity.ColumnViewTypeContactsOrganization), Width: 100, Visible: true, Name: "Organization", Filter: ""},
				{ColumnId: 4, ColumnType: string(postgres_entity.ColumnViewTypeContactsPrimaryEmail), Width: 100, Visible: true, Name: "Primary email", Filter: ""},
				{ColumnId: 5, ColumnType: string(postgres_entity.ColumnViewTypeContactsEmails), Width: 100, Visible: true, Name: "Emails", Filter: ""},
				{ColumnId: 6, ColumnType: string(postgres_entity.ColumnViewTypeContactsPhoneNumbers), Width: 100, Visible: true, Name: "Phone numbers", Filter: ""},
				{ColumnId: 7, ColumnType: string(postgres_entity.ColumnViewTypeContactsLinkedin), Width: 100, Visible: true, Name: "LinkedIn", Filter: ""},
				{ColumnId: 8, ColumnType: string(postgres_entity.ColumnViewTypeContactsJobTitle), Width: 100, Visible: true, Name: "Job title", Filter: ""},
				{ColumnId: 9, ColumnType: string(postgres_entity.ColumnViewTypeContactsTimeInCurrentRole), Width: 100, Visible: true, Name: "Time In Current Role", Filter: ""},
				{ColumnId: 10, ColumnType: string(postgres_entity.ColumnViewTypeContactsCountry), Width: 100, Visible: true, Name: "Country", Filter: ""},
				{ColumnId: 11, ColumnType: string(postgres_entity.ColumnViewTypeContactsRegion), Width: 100, Visible: true, Name: "Region", Filter: ""},
				{ColumnId: 12, ColumnType: string(postgres_entity.ColumnViewTypeContactsCity), Width: 100, Visible: true, Name: "City", Filter: ""},
				{ColumnId: 15, ColumnType: string(postgres_entity.ColumnViewTypeContactsUpdatedAt), Width: 100, Visible: false, Name: "Updated at", Filter: ""},
				{ColumnId: 16, ColumnType: string(postgres_entity.ColumnViewTypeContactsCreatedAt), Width: 100, Visible: false, Name: "Created at", Filter: ""},
			},
		}
	}
	return postgres_entity.Columns{}
}

func CheckSharedPresetsExist(viewDefs []postgres_entity.TableViewDefinition) bool {
	for _, def := range viewDefs {
		if def.IsShared && def.TableType == string(postgres_entity.TableViewTypeOpportunities) {
			return true
		}
	}
	return false
}
