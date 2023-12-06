package resolver

import (
	"context"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueryResolver_Dashboard_ARR_Breakdown_No_Period_No_Data_In_DB(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown_no_period",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 12, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(0), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func TestQueryResolver_Dashboard_ARR_Breakdown_InvalidPeriod(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	response := callGraphQLExpectError(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2020-02-01T00:00:00.000Z",
			"end":   "2020-01-01T00:00:00.000Z",
		})

	require.Contains(t, "Failed to get the data for period", response.Message)
}

func TestQueryResolver_Dashboard_ARR_Breakdown_PeriodIntervals(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	assert_Dashboard_ARR_Breakdown_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-01-31T00:00:00.000Z", 1)
	assert_Dashboard_ARR_Breakdown_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-01-01T00:00:00.000Z", 1)
	assert_Dashboard_ARR_Breakdown_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-02-01T00:00:00.000Z", 2)
	assert_Dashboard_ARR_Breakdown_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-02-28T00:00:00.000Z", 2)
	assert_Dashboard_ARR_Breakdown_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2029-12-01T00:00:00.000Z", 120)
}

func assert_Dashboard_ARR_Breakdown_PeriodIntervals(t *testing.T, start, end string, months int) {
	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": start,
			"end":   end,
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, months, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(0), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func TestQueryResolver_Dashboard_ARR_Breakdown_Cancellations_SLI_Not_Canceled(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1StartedAt)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(0), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func TestQueryResolver_Dashboard_ARR_Breakdown_Cancellations_SLI_Canceled_Before_Month(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	sli1EndedAt := time.Date(2023, 6, 30, 23, 59, 59, 999999999, time.UTC)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1StartedAt, sli1EndedAt)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(0), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func TestQueryResolver_Dashboard_ARR_Breakdown_Cancellations_SLI_Started_Before_Canceled_End_Month(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1StartedAt, sli1EndedAt)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(2), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func TestQueryResolver_Dashboard_ARR_Breakdown_Cancellations_SLI_Started_In_Month_Canceled_End_Month(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1StartedAt, sli1EndedAt)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(2), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func TestQueryResolver_Dashboard_ARR_Breakdown_Cancellations_SLI_Canceled_Annually(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1StartedAt, sli1EndedAt)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(2), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func TestQueryResolver_Dashboard_ARR_Breakdown_Cancellations_SLI_Canceled_Quarterly(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeQuarterly, 4, sli1StartedAt, sli1EndedAt)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(2), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func TestQueryResolver_Dashboard_ARR_Breakdown_Cancellations_SLI_Canceled_Monthly(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeMonthly, 1, sli1StartedAt, sli1EndedAt)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(2), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func TestQueryResolver_Dashboard_ARR_Breakdown_Cancellations_SLI_Started_In_Month_Canceled_Next_Month(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.MiddleTimeOfMonth(2023, 8)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1StartedAt, sli1EndedAt)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(0), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func TestQueryResolver_Dashboard_ARR_Breakdown_Cancellations_SLI_2_Versions_Started_Before_Not_Canceled(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 6)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1MiddleAt, sli1Id)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(0), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func TestQueryResolver_Dashboard_ARR_Breakdown_Cancellations_SLI_2_Versions_Started_Before_Canceled_Before_Month(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 6)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemCanceledWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1MiddleAt, sli1EndedAt, sli1Id)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(0), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func TestQueryResolver_Dashboard_ARR_Breakdown_Cancellations_SLI_2_Versions_Started_Before_Canceled_End_Month(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 6)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemCanceledWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1MiddleAt, sli1EndedAt, sli1Id)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(2), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func TestQueryResolver_Dashboard_ARR_Breakdown_Cancellations_SLI_2_Versions_Started_In_Canceled_End_Month(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemCanceledWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1MiddleAt, sli1EndedAt, sli1Id)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(2), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func TestQueryResolver_Dashboard_ARR_Breakdown_Cancellations_SLI_2_Versions_Started_In_Canceled_Next_Month(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.MiddleTimeOfMonth(2023, 8)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemCanceledWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 12, sli1MiddleAt, sli1EndedAt, sli1Id)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(0), month.Renewals)
		require.Equal(t, float64(0), month.Upsells)
		require.Equal(t, float64(0), month.Downgrades)
		require.Equal(t, float64(0), month.Cancellations)
		require.Equal(t, float64(0), month.Churned)
	}
}

func insertARRBreakdownContractWithOpportunity(ctx context.Context, driver *neo4j.DriverWithContext, orgId string) string {
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{})
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opportunityId)
	return contractId
}

func insertARRBreakdownServiceLineItem(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, startedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:        id,
		ParentID:  id,
		Billed:    billedType,
		Price:     price,
		Quantity:  2,
		StartedAt: startedAt,
	})
	return id
}

func insertARRBreakdownServiceLineItemEnded(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, startedAt, endedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:        id,
		ParentID:  id,
		Billed:    billedType,
		Price:     price,
		Quantity:  2,
		StartedAt: startedAt,
		EndedAt:   &endedAt,
	})
	return id
}

func insertARRBreakdownServiceLineItemCanceled(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, startedAt, endedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:         id,
		ParentID:   id,
		Billed:     billedType,
		Price:      price,
		Quantity:   2,
		IsCanceled: true,
		StartedAt:  startedAt,
		EndedAt:    &endedAt,
	})
	return id
}

func insertARRBreakdownServiceLineItemWithParent(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, startedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:        id,
		ParentID:  parentId,
		Billed:    billedType,
		Price:     price,
		Quantity:  2,
		StartedAt: startedAt,
	})
}

func insertARRBreakdownServiceLineItemEndedWithParent(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, startedAt, endedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:        id,
		ParentID:  parentId,
		Billed:    billedType,
		Price:     price,
		Quantity:  2,
		StartedAt: startedAt,
		EndedAt:   &endedAt,
	})
}

func insertARRBreakdownServiceLineItemCanceledWithParent(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, startedAt, endedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:         id,
		ParentID:   parentId,
		Billed:     billedType,
		Price:      price,
		Quantity:   2,
		IsCanceled: true,
		StartedAt:  startedAt,
		EndedAt:    &endedAt,
	})
}
