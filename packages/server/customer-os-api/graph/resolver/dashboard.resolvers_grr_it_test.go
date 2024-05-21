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
	"time"
)

func TestQueryResolver_Dashboard_GRR_1_Contract_1_SLI_V1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})
	contractStartedAt := utils.FirstTimeOfMonth(2022, 1)
	sli1StartedAt := utils.FirstTimeOfMonth(2022, 6)

	contractId := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 150, 2, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_grr_rate",
		map[string]interface{}{
			"start": "2022-01-01T00:00:00.000Z",
			"end":   "2022-12-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_GrossRevenueRetention model.DashboardGrossRevenueRetention
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 12, len(dashboardReport.Dashboard_GrossRevenueRetention.PerMonth))
	require.Equal(t, float64(0), dashboardReport.Dashboard_GrossRevenueRetention.IncreasePercentageValue)
	require.Equal(t, float64(100), dashboardReport.Dashboard_GrossRevenueRetention.GrossRevenueRetention)

	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 5, 0)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 6, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 7, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 8, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 9, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 10, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 11, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 12, 100)
}

func TestQueryResolver_Dashboard_GRR_1_Contract_1_SLI_Contract_Ended_V1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})
	contractStartedAt := utils.FirstTimeOfMonth(2022, 1)
	contractEndedAt := utils.FirstTimeOfMonth(2022, 9)

	sli1StartedAt := utils.FirstTimeOfMonth(2022, 6)

	contractId := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
		EndedAt:          &contractEndedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 150, 2, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_grr_rate",
		map[string]interface{}{
			"start": "2022-01-01T00:00:00.000Z",
			"end":   "2022-12-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_GrossRevenueRetention model.DashboardGrossRevenueRetention
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 12, len(dashboardReport.Dashboard_GrossRevenueRetention.PerMonth))
	require.Equal(t, float64(0), dashboardReport.Dashboard_GrossRevenueRetention.IncreasePercentageValue)
	require.Equal(t, float64(0), dashboardReport.Dashboard_GrossRevenueRetention.GrossRevenueRetention)

	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 5, 0)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 6, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 7, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 8, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 9, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 10, 0)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 11, 0)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 12, 0)
}

func TestQueryResolver_Dashboard_GRR_2_Contracts_1_SLI_Each(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})
	contract1StartedAt := utils.FirstTimeOfMonth(2022, 1)
	contract2StartedAt := utils.FirstTimeOfMonth(2022, 1)

	sli1StartedAt := utils.FirstTimeOfMonth(2022, 6)
	sli2StartedAt := utils.FirstTimeOfMonth(2022, 6)

	contract1Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeMonthly, 150, 2, sli1StartedAt)

	contract2Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeMonthly, 150, 2, sli2StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_grr_rate",
		map[string]interface{}{
			"start": "2022-01-01T00:00:00.000Z",
			"end":   "2022-12-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_GrossRevenueRetention model.DashboardGrossRevenueRetention
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 12, len(dashboardReport.Dashboard_GrossRevenueRetention.PerMonth))
	require.Equal(t, float64(0), dashboardReport.Dashboard_GrossRevenueRetention.IncreasePercentageValue)
	require.Equal(t, float64(100), dashboardReport.Dashboard_GrossRevenueRetention.GrossRevenueRetention)

	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 5, 0)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 6, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 7, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 8, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 9, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 10, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 11, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 12, 100)
}

func TestQueryResolver_Dashboard_GRR_2_Contracts_1_SLI_Each_Both_Contracts_Ended(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1StartedAt := utils.FirstTimeOfMonth(2022, 1)
	contract1EndedAt := utils.FirstTimeOfMonth(2022, 9)
	sli1StartedAt := utils.FirstTimeOfMonth(2022, 6)
	contract1Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract1StartedAt,
		EndedAt:          &contract1EndedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	contract2StartedAt := utils.FirstTimeOfMonth(2022, 1)
	contract2EndedAt := utils.FirstTimeOfMonth(2022, 8)
	sli2StartedAt := utils.FirstTimeOfMonth(2022, 7)
	contract2Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract2StartedAt,
		EndedAt:          &contract2EndedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeMonthly, 2, 1, sli2StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_grr_rate",
		map[string]interface{}{
			"start": "2022-01-01T00:00:00.000Z",
			"end":   "2022-12-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_GrossRevenueRetention model.DashboardGrossRevenueRetention
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 12, len(dashboardReport.Dashboard_GrossRevenueRetention.PerMonth))
	require.Equal(t, float64(0), dashboardReport.Dashboard_GrossRevenueRetention.IncreasePercentageValue)
	require.Equal(t, float64(0), dashboardReport.Dashboard_GrossRevenueRetention.GrossRevenueRetention)

	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 5, 0)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 6, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 7, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 8, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 9, 33.33)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 10, 0)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 11, 0)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 12, 0)
}

func TestQueryResolver_Dashboard_GRR_2_Contracts_1_SLI_1_Contract_Ended(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1StartedAt := utils.FirstTimeOfMonth(2022, 1)
	contract1EndedAt := utils.FirstTimeOfMonth(2022, 9)
	sli1StartedAt := utils.FirstTimeOfMonth(2022, 6)
	contract1Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contract1StartedAt,
		EndedAt:          &contract1EndedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	contract2StartedAt := utils.FirstTimeOfMonth(2022, 1)
	sli2StartedAt := utils.FirstTimeOfMonth(2022, 7)
	contract2Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeMonthly, 2, 1, sli2StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_grr_rate",
		map[string]interface{}{
			"start": "2022-01-01T00:00:00.000Z",
			"end":   "2022-12-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_GrossRevenueRetention model.DashboardGrossRevenueRetention
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 12, len(dashboardReport.Dashboard_GrossRevenueRetention.PerMonth))
	require.Equal(t, float64(0), dashboardReport.Dashboard_GrossRevenueRetention.IncreasePercentageValue)
	require.Equal(t, float64(66.6), dashboardReport.Dashboard_GrossRevenueRetention.GrossRevenueRetention)

	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 5, 0)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 6, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 7, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 8, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 9, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 10, 66.67)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 11, 66.67)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 12, 66.67)
}

func TestQueryResolver_Dashboard_GRR_2_Contracts_1_SLI_1_SLI_Ended(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	//contract1StartedAt := utils.FirstTimeOfMonth(2022, 1)
	//sli1StartedAt := utils.FirstTimeOfMonth(2022, 6)
	//contract1Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
	//	ContractStatus:   neo4jenum.ContractStatusLive,
	//	ServiceStartedAt: &contract1StartedAt,
	//	RenewalCycle:     neo4jenum.RenewalCycleMonthlyRenewal,
	//}, entity.OpportunityEntity{})
	//neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	contract2StartedAt := utils.FirstTimeOfMonth(2022, 1)
	sli2StartedAt := utils.FirstTimeOfMonth(2022, 7)
	sli2EndedAt := utils.FirstTimeOfMonth(2022, 9)
	contract2Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeMonthly, 2, 1, sli2StartedAt, sli2EndedAt)
	neo4jtest.InsertServiceLineItemCanceledWithParent(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeMonthly, 5, 1, neo4jenum.BilledTypeMonthly, 2, 1, sli2EndedAt, sli2EndedAt.Add(time.Hour*24), sliId)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	//neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	//neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	//neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 3})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_grr_rate",
		map[string]interface{}{
			"start": "2022-01-01T00:00:00.000Z",
			"end":   "2022-12-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_GrossRevenueRetention model.DashboardGrossRevenueRetention
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)
	//
	//require.Equal(t, 12, len(dashboardReport.Dashboard_GrossRevenueRetention.PerMonth))
	//require.Equal(t, "0%", dashboardReport.Dashboard_GrossRevenueRetention.IncreasePercentageValue)
	//require.Equal(t, float64(33.33), dashboardReport.Dashboard_GrossRevenueRetention.GrossRevenueRetention)
	//
	//assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 5, 0)
	//assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 6, 100)
	//assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 7, 100)
	//assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 8, 100)
	//assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 9, 100)
	//assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 10, 33.33)
	//assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 11, 33.33)
	//assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 12, 33.33)
}

func TestQueryResolver_Dashboard_GRR_2_Contracts_1_SLI_1_SLI_Canceld(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	contract1StartedAt := utils.FirstTimeOfMonth(2022, 1)
	sli1StartedAt := utils.FirstTimeOfMonth(2022, 6)
	contract1Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	contract2StartedAt := utils.FirstTimeOfMonth(2022, 1)
	sli2StartedAt := utils.FirstTimeOfMonth(2022, 7)
	sli2EndedAt := utils.FirstTimeOfMonth(2022, 9)
	contract2Id := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeMonthly, 2, 1, sli2StartedAt, sli2EndedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_grr_rate",
		map[string]interface{}{
			"start": "2022-01-01T00:00:00.000Z",
			"end":   "2022-12-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_GrossRevenueRetention model.DashboardGrossRevenueRetention
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 12, len(dashboardReport.Dashboard_GrossRevenueRetention.PerMonth))
	require.Equal(t, float64(0), dashboardReport.Dashboard_GrossRevenueRetention.IncreasePercentageValue)
	require.Equal(t, 33.3, dashboardReport.Dashboard_GrossRevenueRetention.GrossRevenueRetention)

	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 5, 0)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 6, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 7, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 8, 100)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 9, 33.33)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 10, 33.33)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 11, 33.33)
	assertGRRMonthData(t, &dashboardReport.Dashboard_GrossRevenueRetention, 2022, 12, 33.33)
}

func assertGRRMonthData(t *testing.T, dashboardReport *model.DashboardGrossRevenueRetention, year, month int, expected float64) {
	// Find the index corresponding to the given year and month in the PerMonth slice
	var index int
	for i, data := range dashboardReport.PerMonth {
		if data.Year == year && data.Month == month {
			index = i
			break
		}
	}

	require.Equal(t, year, dashboardReport.PerMonth[index].Year)
	require.Equal(t, month, dashboardReport.PerMonth[index].Month)
	require.Equal(t, expected, dashboardReport.PerMonth[index].Percentage)
}
