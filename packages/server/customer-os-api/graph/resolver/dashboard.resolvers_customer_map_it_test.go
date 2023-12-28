package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_Dashboard_Customer_Map_No_Data_In_DB(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_NewCustomers []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, len(dashboardReport.Dashboard_NewCustomers))
}

func TestQueryResolver_Dashboard_Customer_Map_Empty_Contract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_NewCustomers []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, len(dashboardReport.Dashboard_NewCustomers))
}

func TestQueryResolver_Dashboard_Customer_Map_Draft_Contract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_NewCustomers []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, len(dashboardReport.Dashboard_NewCustomers))
}

func TestQueryResolver_Dashboard_Customer_Map_Hidden_Organization_With_Contract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Hide:       true,
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, len(dashboardReport.Dashboard_CustomerMap))
}

func TestQueryResolver_Dashboard_Customer_Map_Prospect_Organization_With_Contract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Hide:       false,
		IsCustomer: false,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, len(dashboardReport.Dashboard_CustomerMap))
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_1_Live_Contract_High_Should_Be_OK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateOk, org1.State)
	require.Equal(t, float64(12), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_1_Live_Contract_Medium_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateAtRisk, org1.State)
	require.Equal(t, float64(12), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_1_Live_Contract_Low_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateAtRisk, org1.State)
	require.Equal(t, float64(12), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_1_Live_Contract_Zero_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodZero,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateAtRisk, org1.State)
	require.Equal(t, float64(12), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_1_Closed_Contract_Should_Be_CHURNED_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract1ServiceEndedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceEndedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateChurned, org1.State)
	require.Equal(t, float64(12), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_2_Live_Contracts_High_Should_Be_OK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	//contract 1
	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contract2Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateOk, org1.State)
	require.Equal(t, float64(24), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_2_Live_Contracts_One_HIGH_One_MEDIUM_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	//contract 1
	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contract2Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, 1, contract2ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateAtRisk, org1.State)
	require.Equal(t, float64(24), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_2_Live_Contracts_One_HIGH_One_LOW_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	//contract 1
	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contract2Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateAtRisk, org1.State)
	require.Equal(t, float64(24), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_2_Live_Contracts_One_HIGH_One_ZERO_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	//contract 1
	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contract2Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodZero,
	})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateAtRisk, org1.State)
	require.Equal(t, float64(24), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_2_Live_Contracts_One_MEDIUM_One_LOW_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	//contract 1
	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contract2Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, entity.OpportunityEntity{
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateAtRisk, org1.State)
	require.Equal(t, float64(24), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_First_Live_Second_Closed_Contract_Live_HIGH_Should_Be_OK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	//contract 1
	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract1Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 100, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract2ServiceEndedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2ServiceEndedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract2ServiceEndedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateOk, org1.State)
	require.Equal(t, float64(100), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_First_Closed_Second_Live_Contract_Live_HIGH_Should_Be_OK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	//contract 1
	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract1ServiceEndedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceEndedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contract2Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 100, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateOk, org1.State)
	require.Equal(t, float64(100), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_First_Live_Second_Closed_Contract_Live_MEDIUM_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	//contract 1
	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract1ServiceEndedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceEndedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contract2Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 100, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateAtRisk, org1.State)
	require.Equal(t, float64(100), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_2_Closed_Contracts_Should_Be_CHURNED_1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	//contract 1
	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract1ServiceEndedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceEndedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contract2ServiceEndedAt := neo4jt.FirstTimeOfMonth(2023, 9)
	contract2Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2ServiceEndedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 100, 1, contract2ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateChurned, org1.State)
	require.Equal(t, float64(100), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_2_Closed_Contracts_Should_Be_CHURNED_2(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	//contract 1
	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contract1ServiceEndedAt := neo4jt.FirstTimeOfMonth(2023, 9)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceEndedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract2ServiceEndedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contract2Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2ServiceEndedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 100, 1, contract2ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract2ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateChurned, org1.State)
	require.Equal(t, float64(12), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Two_Organizations_With_1_Contract_Each(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	//org 1
	org1Id := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, org1Id, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//org 2
	org2Id := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	contract2ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract2ServiceEndedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract2Id := insertContractWithOpportunity(ctx, driver, org2Id, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2ServiceEndedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 10, 1, contract2ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 2, len(dashboardReport.Dashboard_CustomerMap))

	org2 := dashboardReport.Dashboard_CustomerMap[0]
	org1 := dashboardReport.Dashboard_CustomerMap[1]

	require.Equal(t, org2Id, org2.Organization.ID)
	require.Equal(t, contract2ServiceStartedAt, org2.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateChurned, org2.State)
	require.Equal(t, float64(10), org2.Arr)

	require.Equal(t, org1Id, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateOk, org1.State)
	require.Equal(t, float64(12), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Two_Organizations_With_2_Contracts_Each(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	//org 1
	org1Id := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contract1EndedAt := neo4jt.FirstTimeOfMonth(2023, 9)
	contract1Id := insertContractWithOpportunity(ctx, driver, org1Id, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1EndedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	contract2Id := insertContractWithActiveRenewalOpportunity(ctx, driver, org1Id, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 24, 1, contract1ServiceStartedAt)

	//org 2
	org2Id := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	contract2ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 5)
	contract2ServiceEndedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract3Id := insertContractWithActiveRenewalOpportunity(ctx, driver, org2Id, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	insertServiceLineItem(ctx, driver, contract3Id, entity.BilledTypeAnnually, 15, 1, contract2ServiceStartedAt)

	contract4Id := insertContractWithOpportunity(ctx, driver, org2Id, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2ServiceEndedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	insertServiceLineItem(ctx, driver, contract4Id, entity.BilledTypeAnnually, 10, 1, contract2ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 4})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 4})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 4})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 2, len(dashboardReport.Dashboard_CustomerMap))

	org2 := dashboardReport.Dashboard_CustomerMap[0]
	org1 := dashboardReport.Dashboard_CustomerMap[1]

	require.Equal(t, org2Id, org2.Organization.ID)
	require.Equal(t, contract2ServiceStartedAt, org2.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateAtRisk, org2.State)
	require.Equal(t, float64(15), org2.Arr)

	require.Equal(t, org1Id, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateOk, org1.State)
	require.Equal(t, float64(24), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Annually_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateAtRisk, org1.State)
	require.Equal(t, float64(12), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Quarterly_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeQuarterly, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateAtRisk, org1.State)
	require.Equal(t, float64(48), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Monthly_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeMonthly, 12, 1, contract1ServiceStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateAtRisk, org1.State)
	require.Equal(t, float64(144), org1.Arr)
}
