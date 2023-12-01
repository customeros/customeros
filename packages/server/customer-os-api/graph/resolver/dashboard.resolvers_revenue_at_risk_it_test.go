package resolver

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueryResolver_Dashboard_Revenue_At_Risk_No_Period_No_Data_In_Db(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

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
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

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
		MaxAmount:         1,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: entity.OpportunityRenewalLikelihoodHigh,
	})

	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opportunityId)
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

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

func TestQueryResolver_Dashboard_Revenue_At_Risk_One_HIDDEN_Organization_With_Contract_Is_Not_Returned(t *testing.T) {
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

func TestQueryResolver_Dashboard_Revenue_At_Risk_One_Organization_With_1_Live_Annual_Renewal_Contract_High_Renewal_Should_Be_HIGH(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	insert1OrganizationWith1ContractWithOpportunity(ctx, driver, entity.ContractRenewalCycleAnnualRenewal, entity.OpportunityRenewalLikelihoodHigh)

	assertFor1Organization(ctx, t, driver, float64(10), float64(0))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_One_Organization_With_1_Live_Monthly_Renewal_Contract_High_Renewal_Should_Be_HIGH(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	insert1OrganizationWith1ContractWithOpportunity(ctx, driver, entity.ContractRenewalCycleMonthlyRenewal, entity.OpportunityRenewalLikelihoodHigh)

	assertFor1Organization(ctx, t, driver, float64(120), float64(0))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_One_Organization_With_1_Live_Quarterly_Renewal_Contract_High_Renewal_Should_Be_HIGH(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	insert1OrganizationWith1ContractWithOpportunity(ctx, driver, entity.ContractRenewalCycleQuarterlyRenewal, entity.OpportunityRenewalLikelihoodHigh)

	assertFor1Organization(ctx, t, driver, float64(40), float64(0))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_One_Organization_With_1_Live_Annual_Renewal_Contract_Medium_Renewal_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	insert1OrganizationWith1ContractWithOpportunity(ctx, driver, entity.ContractRenewalCycleAnnualRenewal, entity.OpportunityRenewalLikelihoodMedium)

	assertFor1Organization(ctx, t, driver, float64(0), float64(10))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_One_Organization_With_1_Live_Monthly_Renewal_Contract_Medium_Renewal_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	insert1OrganizationWith1ContractWithOpportunity(ctx, driver, entity.ContractRenewalCycleMonthlyRenewal, entity.OpportunityRenewalLikelihoodMedium)

	assertFor1Organization(ctx, t, driver, float64(0), float64(120))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_One_Organization_With_1_Live_Quarterly_Renewal_Contract_Medium_Renewal_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	insert1OrganizationWith1ContractWithOpportunity(ctx, driver, entity.ContractRenewalCycleQuarterlyRenewal, entity.OpportunityRenewalLikelihoodMedium)

	assertFor1Organization(ctx, t, driver, float64(0), float64(40))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_One_Organization_With_1_Live_Annual_Renewal_Contract_Low_Renewal_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	insert1OrganizationWith1ContractWithOpportunity(ctx, driver, entity.ContractRenewalCycleAnnualRenewal, entity.OpportunityRenewalLikelihoodLow)

	assertFor1Organization(ctx, t, driver, float64(0), float64(10))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_One_Organization_With_1_Live_Monthly_Renewal_Contract_Low_Renewal_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	insert1OrganizationWith1ContractWithOpportunity(ctx, driver, entity.ContractRenewalCycleMonthlyRenewal, entity.OpportunityRenewalLikelihoodLow)

	assertFor1Organization(ctx, t, driver, float64(0), float64(120))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_One_Organization_With_1_Live_Quarterly_Renewal_Contract_Low_Renewal_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	insert1OrganizationWith1ContractWithOpportunity(ctx, driver, entity.ContractRenewalCycleQuarterlyRenewal, entity.OpportunityRenewalLikelihoodLow)

	assertFor1Organization(ctx, t, driver, float64(0), float64(40))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_One_Organization_With_1_Live_Annual_Renewal_Contract_Zero_Renewal_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	insert1OrganizationWith1ContractWithOpportunity(ctx, driver, entity.ContractRenewalCycleAnnualRenewal, entity.OpportunityRenewalLikelihoodZero)

	assertFor1Organization(ctx, t, driver, float64(0), float64(10))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_One_Organization_With_1_Live_Monthly_Renewal_Contract_Zero_Renewal_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	insert1OrganizationWith1ContractWithOpportunity(ctx, driver, entity.ContractRenewalCycleMonthlyRenewal, entity.OpportunityRenewalLikelihoodZero)

	assertFor1Organization(ctx, t, driver, float64(0), float64(120))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_One_Organization_With_1_Live_Quarterly_Renewal_Contract_Zero_Renewal_Should_Be_AT_RISK(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	insert1OrganizationWith1ContractWithOpportunity(ctx, driver, entity.ContractRenewalCycleQuarterlyRenewal, entity.OpportunityRenewalLikelihoodZero)

	assertFor1Organization(ctx, t, driver, float64(0), float64(40))
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_One_Organization_With_1_High_1_At_Risk(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	insertContractWithOpportunity(ctx, driver, orgId, entity.ContractRenewalCycleAnnualRenewal, entity.OpportunityRenewalLikelihoodHigh)
	insertContractWithOpportunity(ctx, driver, orgId, entity.ContractRenewalCycleQuarterlyRenewal, entity.OpportunityRenewalLikelihoodMedium)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

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
	require.Equal(t, float64(40), dashboardReport.Dashboard_RevenueAtRisk.AtRisk)
}

func TestQueryResolver_Dashboard_Revenue_At_Risk_2_Organizations_With_1_High_1_At_Risk(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	insert1OrganizationWith1ContractWithOpportunity(ctx, driver, entity.ContractRenewalCycleAnnualRenewal, entity.OpportunityRenewalLikelihoodHigh)
	insert1OrganizationWith1ContractWithOpportunity(ctx, driver, entity.ContractRenewalCycleQuarterlyRenewal, entity.OpportunityRenewalLikelihoodMedium)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})

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
	require.Equal(t, float64(40), dashboardReport.Dashboard_RevenueAtRisk.AtRisk)
}

func insert1OrganizationWith1ContractWithOpportunity(ctx context.Context, driver *neo4j.DriverWithContext, renewalCycle entity.ContractRenewalCycle, renewalLikelihood entity.OpportunityRenewalLikelihood) {
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	insertContractWithOpportunity(ctx, driver, orgId, renewalCycle, renewalLikelihood)
}

func insertContractWithOpportunity(ctx context.Context, driver *neo4j.DriverWithContext, orgId string, renewalCycle entity.ContractRenewalCycle, renewalLikelihood entity.OpportunityRenewalLikelihood) {
	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt:     &contract1ServiceStartedAt,
		ContractStatus:       entity.ContractStatusLive,
		ContractRenewalCycle: renewalCycle,
	})
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{
		Name:              "opportunity 1",
		CreatedAt:         contract1ServiceStartedAt,
		UpdatedAt:         contract1ServiceStartedAt,
		MaxAmount:         10,
		InternalType:      entity.InternalTypeRenewal,
		RenewalLikelihood: renewalLikelihood,
	})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opportunityId)
}

func assertFor1Organization(ctx context.Context, t *testing.T, driver *neo4j.DriverWithContext, expectedHighConfidence float64, expectedAtRisk float64) {
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})

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
