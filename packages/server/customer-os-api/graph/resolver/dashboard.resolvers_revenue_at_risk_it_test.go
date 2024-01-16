package resolver

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_Dashboard_Revenue_At_Risk_No_Period_No_Data_In_Db(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_revenue_at_risk_no_period",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_RevenueAtRisk model.DashboardRevenueAtRisk
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_RevenueAtRisk.HighConfidence)
	require.Equal(t, float64(0), dashboardReport.Dashboard_RevenueAtRisk.AtRisk)
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_No_Data_In_DB(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_revenue_at_risk",
		map[string]interface{}{
			"start": "2000-02-01T00:00:00.000Z",
			"end":   "2500-01-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RevenueAtRisk model.DashboardRevenueAtRisk
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_RevenueAtRisk.HighConfidence)
	require.Equal(t, float64(0), dashboardReport.Dashboard_RevenueAtRisk.AtRisk)
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Draft_Contract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	contract1ServiceStartedAt := neo4jtest.FirstTimeOfMonth(2023, 8)
	neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusDraft,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		MaxAmount:         10,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_revenue_at_risk",
		map[string]interface{}{
			"start": "2000-02-01T00:00:00.000Z",
			"end":   "2500-01-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RevenueAtRisk model.DashboardRevenueAtRisk
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_RevenueAtRisk.HighConfidence)
	require.Equal(t, float64(0), dashboardReport.Dashboard_RevenueAtRisk.AtRisk)
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Closed_Contract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	contract1ServiceStartedAt := neo4jtest.FirstTimeOfMonth(2023, 8)
	contract1EndeddAt := neo4jtest.FirstTimeOfMonth(2023, 9)
	neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1EndeddAt,
	}, neo4jentity.OpportunityEntity{
		MaxAmount:         10,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_revenue_at_risk",
		map[string]interface{}{
			"start": "2000-02-01T00:00:00.000Z",
			"end":   "2500-01-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RevenueAtRisk model.DashboardRevenueAtRisk
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_RevenueAtRisk.HighConfidence)
	require.Equal(t, float64(0), dashboardReport.Dashboard_RevenueAtRisk.AtRisk)
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Hidden_Organization_With_Contract_Is_Not_Returned(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Hide:       true,
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jtest.FirstTimeOfMonth(2023, 8)
	neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		MaxAmount:         10,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_revenue_at_risk",
		map[string]interface{}{
			"start": "2000-02-01T00:00:00.000Z",
			"end":   "2500-01-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RevenueAtRisk model.DashboardRevenueAtRisk
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_RevenueAtRisk.HighConfidence)
	require.Equal(t, float64(0), dashboardReport.Dashboard_RevenueAtRisk.AtRisk)
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Prospect_Organization_With_Contract_Is_Not_Returned(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: false,
	})

	contract1ServiceStartedAt := neo4jtest.FirstTimeOfMonth(2023, 8)
	neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		MaxAmount:         10,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_revenue_at_risk",
		map[string]interface{}{
			"start": "2000-02-01T00:00:00.000Z",
			"end":   "2500-01-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RevenueAtRisk model.DashboardRevenueAtRisk
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_RevenueAtRisk.HighConfidence)
	require.Equal(t, float64(0), dashboardReport.Dashboard_RevenueAtRisk.AtRisk)
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Organization_With_1_Live_Renewal_Contract_High_Should_Be_HIGH(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	contract1ServiceStartedAt := neo4jtest.FirstTimeOfMonth(2023, 8)
	neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     neo4jenum.InternalStageOpen,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
	})

	assertFor1Organization(ctx, t, driver, float64(10), float64(0))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Organization_With_1_Live_Renewal_Contract_Medium_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jtest.FirstTimeOfMonth(2023, 8)
	neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     neo4jenum.InternalStageOpen,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
	})

	assertFor1Organization(ctx, t, driver, float64(0), float64(10))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Organization_With_1_Live_Renewal_Contract_Low_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jtest.FirstTimeOfMonth(2023, 8)
	neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     neo4jenum.InternalStageOpen,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodLow,
	})

	assertFor1Organization(ctx, t, driver, float64(0), float64(10))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Organization_With_1_Live_Contract_Zero_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jtest.FirstTimeOfMonth(2023, 8)
	neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     neo4jenum.InternalStageOpen,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodZero,
	})

	assertFor1Organization(ctx, t, driver, float64(0), float64(10))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Organization_With_1_High_1_At_Risk(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jtest.FirstTimeOfMonth(2023, 8)
	neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     neo4jenum.InternalStageOpen,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
	})

	neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ServiceStartedAt: &contractServiceStartedAt,
		ContractStatus:   neo4jenum.ContractStatusLive,
	}, neo4jentity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     neo4jenum.InternalStageOpen,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_revenue_at_risk",
		map[string]interface{}{
			"start": "2000-02-01T00:00:00.000Z",
			"end":   "2500-01-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RevenueAtRisk model.DashboardRevenueAtRisk
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(10), dashboardReport.Dashboard_RevenueAtRisk.HighConfidence)
	require.Equal(t, float64(10), dashboardReport.Dashboard_RevenueAtRisk.AtRisk)
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Organization_With_2_Opportunities_Ok(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jtest.FirstTimeOfMonth(2023, 8)
	contractId := neo4jtest.CreateContract(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	})

	neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage:     neo4jenum.InternalStageClosedWon,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		MaxAmount:         5,
	})
	opId := neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage:     neo4jenum.InternalStageOpen,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		MaxAmount:         12,
	})
	neo4jtest.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opId)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_revenue_at_risk",
		map[string]interface{}{
			"start": "2000-02-01T00:00:00.000Z",
			"end":   "2500-01-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RevenueAtRisk model.DashboardRevenueAtRisk
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(12), dashboardReport.Dashboard_RevenueAtRisk.HighConfidence)
	require.Equal(t, float64(0), dashboardReport.Dashboard_RevenueAtRisk.AtRisk)
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Organization_With_2_Opportunities_At_Risk(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jtest.FirstTimeOfMonth(2023, 8)
	contractId := neo4jtest.CreateContract(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	})

	neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage:     neo4jenum.InternalStageClosedWon,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
		MaxAmount:         5,
	})
	opId := neo4jtest.CreateOpportunityForContract(ctx, driver, tenantName, contractId, neo4jentity.OpportunityEntity{
		InternalStage:     neo4jenum.InternalStageOpen,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
		MaxAmount:         12,
	})
	neo4jtest.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opId)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_revenue_at_risk",
		map[string]interface{}{
			"start": "2000-02-01T00:00:00.000Z",
			"end":   "2500-01-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RevenueAtRisk model.DashboardRevenueAtRisk
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_RevenueAtRisk.HighConfidence)
	require.Equal(t, float64(12), dashboardReport.Dashboard_RevenueAtRisk.AtRisk)
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_2_Organizations_With_1_High_1_At_Risk(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	org1Id := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jtest.FirstTimeOfMonth(2023, 8)
	neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, org1Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, neo4jentity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     neo4jenum.InternalStageOpen,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodHigh,
	})

	org2Id := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, org2Id, neo4jentity.ContractEntity{
		ServiceStartedAt: &contractServiceStartedAt,
		ContractStatus:   neo4jenum.ContractStatusLive,
	}, neo4jentity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     neo4jenum.InternalStageOpen,
		InternalType:      neo4jenum.InternalTypeRenewal,
		RenewalLikelihood: neo4jenum.RenewalLikelihoodMedium,
	})

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_revenue_at_risk",
		map[string]interface{}{
			"start": "2000-02-01T00:00:00.000Z",
			"end":   "2500-01-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RevenueAtRisk model.DashboardRevenueAtRisk
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(10), dashboardReport.Dashboard_RevenueAtRisk.HighConfidence)
	require.Equal(t, float64(10), dashboardReport.Dashboard_RevenueAtRisk.AtRisk)
}

func assertFor1Organization(ctx context.Context, t *testing.T, driver *neo4j.DriverWithContext, expectedHighConfidence float64, expectedAtRisk float64) {
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_revenue_at_risk",
		map[string]interface{}{
			"start": "2000-02-01T00:00:00.000Z",
			"end":   "2500-01-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RevenueAtRisk model.DashboardRevenueAtRisk
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, expectedHighConfidence, dashboardReport.Dashboard_RevenueAtRisk.HighConfidence)
	require.Equal(t, expectedAtRisk, dashboardReport.Dashboard_RevenueAtRisk.AtRisk)
}
