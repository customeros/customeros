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

func TestQueryResolver_Dashboard_MRR_Per_Customer_No_Period_No_Data_In_DB(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer_no_period",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "0%", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 12, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, float64(0), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_InvalidPeriod(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	response := callGraphQLExpectError(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2020-02-01T00:00:00.000Z",
			"end":   "2020-01-01T00:00:00.000Z",
		})

	require.Contains(t, "Failed to get the data for period", response.Message)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_PeriodIntervals(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	assert_Dashboard_MRR_Per_Customer_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-01-31T00:00:00.000Z", 1)
	assert_Dashboard_MRR_Per_Customer_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-01-01T00:00:00.000Z", 1)
	assert_Dashboard_MRR_Per_Customer_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-02-01T00:00:00.000Z", 2)
	assert_Dashboard_MRR_Per_Customer_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-02-28T00:00:00.000Z", 2)
	assert_Dashboard_MRR_Per_Customer_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2029-12-01T00:00:00.000Z", 120)
}

func assert_Dashboard_MRR_Per_Customer_PeriodIntervals(t *testing.T, start, end string, months int) {
	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": start,
			"end":   end,
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "0%", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, months, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_InMonth_Prospect(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: false,
	})

	contractServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "0%", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 0)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_InMonth_Hidden_Organization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
		Hide:       true,
	})

	contractServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "0%", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 0)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_Closed_Contract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractEndedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contractStartedAt,
		EndedAt:          &contractEndedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "0%", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_Canceled(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.FirstTimeOfMonth(2023, 9)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-09-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "-100%", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 3, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 2)
	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 8, 2)
	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 9, 0)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_BeforeMonth(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(2), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "+2", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 2)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_AfterMonth(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "0%", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 0)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_AtBeginningOfMonth(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(2), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "+2", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 2)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_AtEndOfMonth(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(2), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "+2", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 2)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_InMonth_EndedImmediately(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "0%", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 0)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_InMonth_EndedAtEndOfMonth(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "0%", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 0)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_InMonth_EndedNextMonth(t *testing.T) {

	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 8)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(2), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "+2", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 2)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_Yearly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(2), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "+2", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 2)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeQuarterly, 3, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(1), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "+100%", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 1)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 2, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(2), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "+2", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 2)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_2_SLI_SameMonth_SameOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract1Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contractServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	contract2ServiceStartedAt := neo4jt.MiddleTimeOfMonth(2023, 7).Add(10 * 24 * time.Hour)
	sli2StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract2Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, 2, sli2StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(4), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "+4", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 4)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_2_SLI_SameMonth_DifferentOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	contract1ServiceStartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract1Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId1, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	orgId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	contract2ServiceStartedAt := neo4jt.MiddleTimeOfMonth(2023, 7).Add(10 * 24 * time.Hour)
	sli2StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract2Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId2, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, 2, sli2StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(2), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "+2", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 2)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_2_SLI_DifferentMonths_SameOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract1Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	contract2ServiceStartedAt := neo4jt.MiddleTimeOfMonth(2023, 9)
	sli2StartedAt := neo4jt.FirstTimeOfMonth(2023, 9)
	contract2Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, 2, sli2StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(2), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "+2", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 2)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_2_SLI_DifferentMonths_DifferentOrganization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	contract1ServiceStartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract1Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId1, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	orgId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	contract2ServiceStartedAt := neo4jt.LastTimeOfMonth(2023, 8)
	sli2StartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contract2Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId2, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, 2, sli2StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(1), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "+100%", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 1)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_2_SLI_SameOrganization_Overlaps_2_Months(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.LastTimeOfMonth(2023, 6)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := neo4jt.FirstTimeOfMonth(2023, 9)
	contract1Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItemEnded(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

	contract2ServiceStartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	sli2StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli2EndedAt := neo4jt.FirstTimeOfMonth(2023, 10)
	contract2Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItemEnded(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, 2, sli2StartedAt, sli2EndedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-05-01T00:00:00.000Z",
			"end":   "2023-10-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "-100%", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 6, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 5, 0)
	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 6, 2)
	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 4)
	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 8, 4)
	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 9, 2)
	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 10, 0)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_2_SLI_DifferentOrganization_Overlaps_2_Months(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := neo4jt.FirstTimeOfMonth(2023, 9)
	contract1Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId1, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract1ServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItemEnded(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

	orgId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	contract2ServiceStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli2StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli2EndedAt := neo4jt.FirstTimeOfMonth(2023, 10)
	contract2Id := insertContractWithActiveRenewalOpportunity(ctx, driver, orgId2, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &contract2ServiceStartedAt,
	}, entity.OpportunityEntity{})
	insertServiceLineItemEnded(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, 2, sli2StartedAt, sli2EndedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_mrr_per_customer",
		map[string]interface{}{
			"start": "2023-05-01T00:00:00.000Z",
			"end":   "2023-10-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_MRRPerCustomer model.DashboardMRRPerCustomer
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, "-100%", dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 6, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 5, 0)
	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 6, 1)
	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 7, 2)
	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 8, 2)
	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 9, 1)
	assertMRRMonthData(t, &dashboardReport.Dashboard_MRRPerCustomer, 2023, 10, 0)
}

func assertMRRMonthData(t *testing.T, dashboardReport *model.DashboardMRRPerCustomer, year, month int, expected float64) {
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
	require.Equal(t, expected, dashboardReport.PerMonth[index].Value)
}
