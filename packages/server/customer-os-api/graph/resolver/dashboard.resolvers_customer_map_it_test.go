package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_Dashboard_Customer_Map_No_Data_In_DB(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, len(dashboardReport.Dashboard_CustomerMap))
}

func TestQueryResolver_Dashboard_Customer_Map_Empty_Contract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: neo4jenum.Customer})

	neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, len(dashboardReport.Dashboard_CustomerMap))
}

func TestQueryResolver_Dashboard_Customer_Map_Draft_Contract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: neo4jenum.Customer})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusDraft,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, len(dashboardReport.Dashboard_CustomerMap))
}

func TestQueryResolver_Dashboard_Customer_Map_Scheduled_Contract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusScheduled,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, model.DashboardCustomerMapStateOk, org1.State)
}

func TestQueryResolver_Dashboard_Customer_Map_OutOfContract_Contract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: neo4jenum.Customer})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusOutOfContract,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	org1 := dashboardReport.Dashboard_CustomerMap[0]
	require.Equal(t, 1, len(dashboardReport.Dashboard_CustomerMap))
	require.Equal(t, orgId, org1.Organization.ID)
	require.Equal(t, model.DashboardCustomerMapStateMediumRisk, org1.State)
}

func TestQueryResolver_Dashboard_Customer_Map_Hidden_Organization_With_Contract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Hide:         true,
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalType: neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Hide:         false,
		Relationship: neo4jenum.Prospect,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalType: neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_1_Live_Contract_Medium_Should_Be_MEDIUM_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, model.DashboardCustomerMapStateMediumRisk, org1.State)
	require.Equal(t, float64(12), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_1_Live_Contract_Low_Should_Be_HIGH_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodLow,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, model.DashboardCustomerMapStateHighRisk, org1.State)
	require.Equal(t, float64(12), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_1_Live_Contract_Zero_Should_Be_HIGH_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodZero,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, model.DashboardCustomerMapStateHighRisk, org1.State)
	require.Equal(t, float64(12), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_1_Live_Contract_No_Recurring_Revenue(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodZero,
		},
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, len(dashboardReport.Dashboard_CustomerMap))
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_1_Live_Contract_Closed_SLI_AT_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodZero,
		},
	})

	sliStartedAt := utils.FirstTimeOfMonth(2023, 7)
	sliEndedAt := utils.FirstTimeOfMonth(2023, 8)
	neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, sliStartedAt, sliEndedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_customer_map",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_CustomerMap []model.DashboardCustomerMap
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, len(dashboardReport.Dashboard_CustomerMap))
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_1_Closed_Contract_Should_Be_CHURNED_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contract1ServiceEndedAt := utils.FirstTimeOfMonth(2023, 8)
	contractId := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceEndedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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

func TestQueryResolver_Dashboard_Customer_Map_Organization_1_Contract_2_Opportunities_Should_Be_OK(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	})
	neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageClosedWon,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		},
	})
	opId := neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opId)

	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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

func TestQueryResolver_Dashboard_Customer_Map_Organization_1_Contract_3_Opportunities_Should_Be_OK(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	})
	neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageClosedLost,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodZero,
		},
	})
	neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageClosedWon,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		},
	})
	opId := neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opId)

	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 3})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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

func TestQueryResolver_Dashboard_Customer_Map_Organization_1_Contract_2_Opportunities_Should_Be_Medium_Risk(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	})
	neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageClosedWon,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	opId := neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		},
	})
	neo4jtest.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opId)

	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, model.DashboardCustomerMapStateMediumRisk, org1.State)
	require.Equal(t, float64(12), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_1_Contract_2_Opportunities_Should_Be_Churned(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contract1ServiceEndedAt := utils.FirstTimeOfMonth(2023, 8)
	contractId := neo4jtest.CreateContractForOrganization(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceEndedAt,
	})
	neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageClosedWon,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	opId := neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		},
	})
	neo4jtest.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opId)

	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	//contract 1
	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := utils.FirstTimeOfMonth(2023, 8)
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalType: neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_2_Live_Contracts_One_HIGH_One_MEDIUM_Should_Be_MEDIUM_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	//contract 1
	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := utils.FirstTimeOfMonth(2023, 8)
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 12, 1, contract2ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	require.Equal(t, model.DashboardCustomerMapStateMediumRisk, org1.State)
	require.Equal(t, float64(24), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_2_Live_Contracts_One_HIGH_One_LOW_Should_Be_HIGH_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: neo4jenum.Customer})

	//contract 1
	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := utils.FirstTimeOfMonth(2023, 8)
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodLow,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	require.Equal(t, model.DashboardCustomerMapStateHighRisk, org1.State)
	require.Equal(t, float64(24), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_2_Live_Contracts_One_HIGH_One_ZERO_Should_Be_HIGH_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: neo4jenum.Customer})

	//contract 1
	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := utils.FirstTimeOfMonth(2023, 8)
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodZero,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	require.Equal(t, model.DashboardCustomerMapStateHighRisk, org1.State)
	require.Equal(t, float64(24), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_2_Live_Contracts_One_MEDIUM_One_LOW_Should_Be_HIGH_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: neo4jenum.Customer})

	//contract 1
	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := utils.FirstTimeOfMonth(2023, 8)
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodLow,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	require.Equal(t, model.DashboardCustomerMapStateHighRisk, org1.State)
	require.Equal(t, float64(24), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_First_Live_Second_Closed_Contract_Live_HIGH_Should_Be_OK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: neo4jenum.Customer})

	//contract 1
	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 100, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contract2ServiceEndedAt := utils.FirstTimeOfMonth(2023, 8)
	contractId := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2ServiceEndedAt,
	}, neo4jentity.OpportunityEntity{
		InternalType: neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract2ServiceEndedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: neo4jenum.Customer})

	//contract 1
	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contract1ServiceEndedAt := utils.FirstTimeOfMonth(2023, 8)
	contractId := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceEndedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := utils.FirstTimeOfMonth(2023, 8)
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 100, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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

func TestQueryResolver_Dashboard_Customer_Map_Organization_With_First_Live_Second_Closed_Contract_Live_MEDIUM_Should_Be_MEDIUM_RISK_State(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: neo4jenum.Customer})

	//contract 1
	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contract1ServiceEndedAt := utils.FirstTimeOfMonth(2023, 8)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceEndedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := utils.FirstTimeOfMonth(2023, 8)
	contract2Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 100, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	require.Equal(t, model.DashboardCustomerMapStateMediumRisk, org1.State)
	require.Equal(t, float64(100), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_One_Organization_With_2_Closed_Contracts_Should_Be_CHURNED_1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	//contract 1
	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contract1ServiceEndedAt := utils.FirstTimeOfMonth(2023, 8)
	contractId := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceEndedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := utils.FirstTimeOfMonth(2023, 8)
	contract2ServiceEndedAt := utils.FirstTimeOfMonth(2023, 9)
	contract2Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2ServiceEndedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 100, 1, contract2ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{Relationship: neo4jenum.Customer})

	//contract 1
	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 8)
	contract1ServiceEndedAt := utils.FirstTimeOfMonth(2023, 9)
	contractId := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1ServiceEndedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//contract 2
	contract2ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contract2ServiceEndedAt := utils.FirstTimeOfMonth(2023, 8)
	contract2Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2ServiceEndedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 100, 1, contract2ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	//org 1
	org1Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})
	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, org1Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	//org 2
	org2Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})
	contract2ServiceStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract2ServiceEndedAt := utils.FirstTimeOfMonth(2023, 7)
	contract2Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, org2Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2ServiceEndedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 10, 1, contract2ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	//org 1
	org1Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})
	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 8)
	contract1EndedAt := utils.FirstTimeOfMonth(2023, 9)
	contract1Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, org1Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1EndedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, org1Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 24, 1, contract1ServiceStartedAt)

	//org 2
	org2Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})
	contract2ServiceStartedAt := utils.FirstTimeOfMonth(2023, 5)
	contract2ServiceEndedAt := utils.FirstTimeOfMonth(2023, 6)
	contract3Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, org2Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract3Id, neo4jenum.BilledTypeAnnually, 15, 1, contract2ServiceStartedAt)

	contract4Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, org2Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2ServiceEndedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract4Id, neo4jenum.BilledTypeAnnually, 10, 1, contract2ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 4})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 4})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 4})

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
	require.Equal(t, model.DashboardCustomerMapStateMediumRisk, org2.State)
	require.Equal(t, float64(15), org2.Arr)

	require.Equal(t, org1Id, org1.Organization.ID)
	require.Equal(t, contract1ServiceStartedAt, org1.ContractSignedDate)
	require.Equal(t, model.DashboardCustomerMapStateOk, org1.State)
	require.Equal(t, float64(24), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Annually_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodLow,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, model.DashboardCustomerMapStateHighRisk, org1.State)
	require.Equal(t, float64(12), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Quarterly_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodLow,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, model.DashboardCustomerMapStateHighRisk, org1.State)
	require.Equal(t, float64(48), org1.Arr)
}

func TestQueryResolver_Dashboard_Customer_Map_Monthly_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1ServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		InternalStage: neo4jenum.OpportunityInternalStageOpen,
		InternalType:  neo4jenum.OpportunityInternalTypeRenewal,
		RenewalDetails: neo4jentity.RenewalDetails{
			RenewalLikelihood: neo4jenum.RenewalLikelihoodLow,
		},
	})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 12, 1, contract1ServiceStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, model.DashboardCustomerMapStateHighRisk, org1.State)
	require.Equal(t, float64(144), org1.Arr)
}
