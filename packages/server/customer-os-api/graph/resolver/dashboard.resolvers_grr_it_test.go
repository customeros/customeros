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

func TestQueryResolver_Dashboard_GRR_1_Contract_1_SLI_V1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	contractStartedAt := neo4jt.FirstTimeOfMonth(2022, 1)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2022, 6)

	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeMonthly, 150, 2, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, "0%", dashboardReport.Dashboard_GrossRevenueRetention.IncreasePercentage)
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
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	contractStartedAt := neo4jt.FirstTimeOfMonth(2022, 1)
	contractEndedAt := neo4jt.FirstTimeOfMonth(2022, 9)

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2022, 6)

	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
		EndedAt:          &contractEndedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeMonthly, 150, 2, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, "0%", dashboardReport.Dashboard_GrossRevenueRetention.IncreasePercentage)
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
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	contract1StartedAt := neo4jt.FirstTimeOfMonth(2022, 1)
	contract2StartedAt := neo4jt.FirstTimeOfMonth(2022, 1)

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2022, 6)
	sli2StartedAt := neo4jt.FirstTimeOfMonth(2022, 6)

	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeMonthly, 150, 2, sli1StartedAt)

	contract2Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeMonthly, 150, 2, sli2StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	require.Equal(t, "0%", dashboardReport.Dashboard_GrossRevenueRetention.IncreasePercentage)
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
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1StartedAt := neo4jt.FirstTimeOfMonth(2022, 1)
	contract1EndedAt := neo4jt.FirstTimeOfMonth(2022, 9)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2022, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract1StartedAt,
		EndedAt:          &contract1EndedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	contract2StartedAt := neo4jt.FirstTimeOfMonth(2022, 1)
	contract2EndedAt := neo4jt.FirstTimeOfMonth(2022, 8)
	sli2StartedAt := neo4jt.FirstTimeOfMonth(2022, 7)
	contract2Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract2StartedAt,
		EndedAt:          &contract2EndedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeMonthly, 2, 1, sli2StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	require.Equal(t, "0%", dashboardReport.Dashboard_GrossRevenueRetention.IncreasePercentage)
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
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1StartedAt := neo4jt.FirstTimeOfMonth(2022, 1)
	contract1EndedAt := neo4jt.FirstTimeOfMonth(2022, 9)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2022, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contract1StartedAt,
		EndedAt:          &contract1EndedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	contract2StartedAt := neo4jt.FirstTimeOfMonth(2022, 1)
	sli2StartedAt := neo4jt.FirstTimeOfMonth(2022, 7)
	contract2Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeMonthly, 2, 1, sli2StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	require.Equal(t, "0%", dashboardReport.Dashboard_GrossRevenueRetention.IncreasePercentage)
	require.Equal(t, float64(66.67), dashboardReport.Dashboard_GrossRevenueRetention.GrossRevenueRetention)

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
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	//contract1StartedAt := neo4jt.FirstTimeOfMonth(2022, 1)
	//sli1StartedAt := neo4jt.FirstTimeOfMonth(2022, 6)
	//contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
	//	ContractStatus:   entity.ContractStatusLive,
	//	ServiceStartedAt: &contract1StartedAt,
	//	RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	//}, entity.OpportunityEntity{})
	//insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	contract2StartedAt := neo4jt.FirstTimeOfMonth(2022, 1)
	sli2StartedAt := neo4jt.FirstTimeOfMonth(2022, 7)
	sli2EndedAt := neo4jt.FirstTimeOfMonth(2022, 9)
	contract2Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	sliId := insertServiceLineItemEnded(ctx, driver, contract2Id, entity.BilledTypeMonthly, 2, 1, sli2StartedAt, sli2EndedAt)
	insertServiceLineItemCanceledWithParent(ctx, driver, contract2Id, entity.BilledTypeMonthly, 5, 1, entity.BilledTypeMonthly, 2, 1, sli2EndedAt, sli2EndedAt.Add(time.Hour*24), sliId)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	//assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	//assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	//assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 3})

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
	//require.Equal(t, "0%", dashboardReport.Dashboard_GrossRevenueRetention.IncreasePercentage)
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
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1StartedAt := neo4jt.FirstTimeOfMonth(2022, 1)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2022, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	contract2StartedAt := neo4jt.FirstTimeOfMonth(2022, 1)
	sli2StartedAt := neo4jt.FirstTimeOfMonth(2022, 7)
	sli2EndedAt := neo4jt.FirstTimeOfMonth(2022, 9)
	contract2Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItemCanceled(ctx, driver, contract2Id, entity.BilledTypeMonthly, 2, 1, sli2StartedAt, sli2EndedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	require.Equal(t, "0%", dashboardReport.Dashboard_GrossRevenueRetention.IncreasePercentage)
	require.Equal(t, float64(33.33), dashboardReport.Dashboard_GrossRevenueRetention.GrossRevenueRetention)

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
