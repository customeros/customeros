package resolver

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
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

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusDraft,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
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
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contract1EndeddAt := neo4jt.FirstTimeOfMonth(2023, 9)
	insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1EndeddAt,
	}, entity.OpportunityEntity{
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
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

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Hide:       true,
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
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

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: false,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
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
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{IsCustomer: true})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     entity.InternalStageOpen,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})

	assertFor1Organization(ctx, t, driver, float64(10), float64(0))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Organization_With_1_Live_Renewal_Contract_Medium_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     entity.InternalStageOpen,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
	})

	assertFor1Organization(ctx, t, driver, float64(0), float64(10))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Organization_With_1_Live_Renewal_Contract_Low_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     entity.InternalStageOpen,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodLow,
	})

	assertFor1Organization(ctx, t, driver, float64(0), float64(10))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Organization_With_1_Live_Contract_Zero_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     entity.InternalStageOpen,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodZero,
	})

	assertFor1Organization(ctx, t, driver, float64(0), float64(10))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_Organization_With_1_High_1_At_Risk(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     entity.InternalStageOpen,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})

	insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ServiceStartedAt: &contractServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	}, entity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     entity.InternalStageOpen,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
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
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	})

	neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		InternalStage:     entity.InternalStageClosedWon,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
		MaxAmount:         5,
	})
	opId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		InternalStage:     entity.InternalStageOpen,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
		MaxAmount:         12,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opId)

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
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	})

	neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		InternalStage:     entity.InternalStageClosedWon,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
		MaxAmount:         5,
	})
	opId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		InternalStage:     entity.InternalStageOpen,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
		MaxAmount:         12,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opId)

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
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	org1Id := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	insertContractWithActiveRenewalOpportunity(ctx, driver, org1Id, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     entity.InternalStageOpen,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})

	org2Id := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	insertContractWithActiveRenewalOpportunity(ctx, driver, org2Id, entity.ContractEntity{
		ServiceStartedAt: &contractServiceStartedAt,
		ContractStatus:   entity.ContractStatusLive,
	}, entity.OpportunityEntity{
		MaxAmount:         10,
		InternalStage:     entity.InternalStageOpen,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodMedium,
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
