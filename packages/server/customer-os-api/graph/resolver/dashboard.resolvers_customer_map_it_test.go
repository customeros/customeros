package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueryResolver_Dashboard_Customer_Map_No_Data_In_DB(t *testing.T) {
	ctx := context.TODO()
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
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

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
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{})

	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		Name:              "oppo 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         1,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})

	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opportunityId)
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_NewCustomers []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, len(dashboardReport.Dashboard_NewCustomers))
}

func TestQueryResolver_Dashboard_Customer_Map_One_HIDDEN_Organization_With_Contract_Is_Not_Returned(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Hide: true,
	})

	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		Name:              "oppo 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         1,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opportunityId)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, len(dashboardReport.Dashboard_CustomerMap))
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_1_Live_Renewal_Contract_High_Renewal_Should_Be_OK_State(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		Name:              "oppo 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         1,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opportunityId)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

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
	require.Equal(t, float64(1), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_1_Live_Renewal_Contract_Medium_Renewal_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		Name:              "oppo 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         1,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opportunityId)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

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
	require.Equal(t, float64(1), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_1_Live_Renewal_Contract_Low_Renewal_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		Name:              "oppo 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         1,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opportunityId)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

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
	require.Equal(t, float64(1), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_1_Live_Renewal_Contract_Zero_Renewal_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		Name:              "oppo 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         1,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodZero,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opportunityId)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

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
	require.Equal(t, float64(1), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_1_Ended_Renewal_Contract_Should_Be_CHURNED_State(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusEnded,
	})
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		Name:              "oppo 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         100,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opportunityId)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

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

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_2_Live_Renewal_Contracts_High_Should_Be_OK_State(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	//contract 1
	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contract1Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity1Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract1Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         float64(100),
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract1Id, opportunity1Id)

	//contract 2
	contract2ServiceStartedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	contract2Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract2ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity2Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract2Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 2",
		CreatedAt:         contract2ServiceStartedAt,
		UpdatedAt:         contract2ServiceStartedAt,
		MaxAmount:         float64(150),
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract2Id, opportunity2Id)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

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
	require.Equal(t, float64(250), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_2_Live_Renewal_Contracts_One_HIGH_One_MEDIUM_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	//contract 1
	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contract1Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity1Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract1Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract1Id, opportunity1Id)

	//contract 2
	contract2ServiceStartedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	contract2Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract2ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity2Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract2Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 2",
		CreatedAt:         contract2ServiceStartedAt,
		UpdatedAt:         contract2ServiceStartedAt,
		MaxAmount:         20,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract2Id, opportunity2Id)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

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
	require.Equal(t, float64(30), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_2_Live_Renewal_Contracts_One_HIGH_One_LOW_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	//contract 1
	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contract1Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity1Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract1Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract1Id, opportunity1Id)

	//contract 2
	contract2ServiceStartedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	contract2Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract2ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity2Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract2Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 2",
		CreatedAt:         contract2ServiceStartedAt,
		UpdatedAt:         contract2ServiceStartedAt,
		MaxAmount:         15,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract2Id, opportunity2Id)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

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
	require.Equal(t, float64(25), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_2_Live_Renewal_Contracts_One_HIGH_One_ZERO_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	//contract 1
	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contract1Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity1Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract1Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract1Id, opportunity1Id)

	//contract 2
	contract2ServiceStartedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	contract2Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract2ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity2Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract2Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 2",
		CreatedAt:         contract2ServiceStartedAt,
		UpdatedAt:         contract2ServiceStartedAt,
		MaxAmount:         15,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodZero,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract2Id, opportunity2Id)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

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
	require.Equal(t, float64(25), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_2_Live_Renewal_Contracts_One_MEDIUM_One_LOW_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	//contract 1
	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contract1Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity1Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract1Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract1Id, opportunity1Id)

	//contract 2
	contract2ServiceStartedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	contract2Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract2ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity2Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract2Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 2",
		CreatedAt:         contract2ServiceStartedAt,
		UpdatedAt:         contract2ServiceStartedAt,
		MaxAmount:         15,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract2Id, opportunity2Id)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

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
	require.Equal(t, float64(25), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_1_Live_1_Closed_Renewal_Contracts_One_HIGH_Should_Be_OK_State(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	//contract 1
	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contract1Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity1Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract1Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract1Id, opportunity1Id)

	//contract 2
	contract2ServiceStartedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	contract2Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2ServiceStartedAt,
		ContractStatus:   entity.ContractStatusEnded,
	})
	opportunity2Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract2Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 2",
		CreatedAt:         contract2ServiceStartedAt,
		UpdatedAt:         contract2ServiceStartedAt,
		MaxAmount:         15,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract2Id, opportunity2Id)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

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
	require.Equal(t, float64(10), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_1_Live_1_Closed_Renewal_Contracts_One_MEDIUM_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	//contract 1
	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contract1Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity1Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract1Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract1Id, opportunity1Id)

	//contract 2
	contract2ServiceStartedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	contract2Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2ServiceStartedAt,
		ContractStatus:   entity.ContractStatusEnded,
	})
	opportunity2Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract2Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 2",
		CreatedAt:         contract2ServiceStartedAt,
		UpdatedAt:         contract2ServiceStartedAt,
		MaxAmount:         15,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract2Id, opportunity2Id)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

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
	require.Equal(t, float64(10), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_2_Closed_Renewal_Contracts_Should_Be_CHURNED_State(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	//contract 1
	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contract1Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusEnded,
	})
	opportunity1Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract1Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract1Id, opportunity1Id)

	//contract 2
	contract2ServiceStartedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	contract2Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2ServiceStartedAt,
		ContractStatus:   entity.ContractStatusEnded,
	})
	opportunity2Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract2Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 2",
		CreatedAt:         contract2ServiceStartedAt,
		UpdatedAt:         contract2ServiceStartedAt,
		MaxAmount:         15,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract2Id, opportunity2Id)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

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
	require.Equal(t, float64(15), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_1_Closed_1_Live_Renewal_Contracts_First_Close_Second_Live_Should_Be_AT_RISK_State(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	//contract 1
	contract1ServiceStartedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	contract1Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusEnded,
	})
	opportunity1Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract1Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract1Id, opportunity1Id)

	//contract 2
	contract2ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contract2Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract2ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity2Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract2Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 2",
		CreatedAt:         contract2ServiceStartedAt,
		UpdatedAt:         contract2ServiceStartedAt,
		MaxAmount:         15,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract2Id, opportunity2Id)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

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
	require.Equal(t, model.DashboardCustomerMapStateAtRisk, org1.State)
	require.Equal(t, float64(15), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Two_Organizations_With_1_Contract_Each(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	//org 1
	org1Id := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	contract1Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, org1Id, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusEnded,
	})
	opportunity1Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract1Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract1Id, opportunity1Id)

	//org 2
	org2Id := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	contract2ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contract2Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, org2Id, entity.ContractEntity{
		ServiceStartedAt: &contract2ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity2Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract2Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 2",
		CreatedAt:         contract2ServiceStartedAt,
		UpdatedAt:         contract2ServiceStartedAt,
		MaxAmount:         15,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract2Id, opportunity2Id)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

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
	require.Equal(t, model.DashboardCustomerMapStateChurned, org1.State)
	require.Equal(t, float64(10), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_MMR_1_Live_Contract(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	//org 1
	org1Id := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	contract1Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, org1Id, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	})
	opportunity1Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract1Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract1Id, opportunity1Id)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]

	require.Equal(t, org1Id, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateOk, org1.State)
	require.Equal(t, float64(10), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_MMR_1_Closed_Contract(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	//org 1
	org1Id := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	contract1Id := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, org1Id, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceStartedAt,
		ContractStatus:   entity.ContractStatusEnded,
	})
	opportunity1Id := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contract1Id, entity.OpportunityEntity{
		Name:              "opp 1 - ctr 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contract1Id, opportunity1Id)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))

	org1 := dashboardReport.Dashboard_CustomerMap[0]

	require.Equal(t, org1Id, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateChurned, org1.State)
	require.Equal(t, float64(10), org1.Arr)
}
