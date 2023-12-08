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

func Test_Dashboard_ARR_Breakdown_No_Period_No_Data_In_DB(t *testing.T) {
	ctx := context.Background()
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

func Test_Dashboard_ARR_Breakdown_InvalidPeriod(t *testing.T) {
	ctx := context.Background()
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

func Test_Dashboard_ARR_Breakdown_PeriodIntervals(t *testing.T) {
	ctx := context.Background()
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

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_Not_Canceled(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
		require.Equal(t, float64(0), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_Canceled_Before_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	sli1EndedAt := time.Date(2023, 6, 30, 23, 59, 59, 999999999, time.UTC)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

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
		require.Equal(t, float64(0), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_Started_Before_Canceled_End_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

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
		require.Equal(t, float64(2), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_Started_In_Month_Canceled_End_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

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
		require.Equal(t, float64(2), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_Canceled_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

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
		require.Equal(t, float64(2), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_Canceled_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeQuarterly, 4, 2, sli1StartedAt, sli1EndedAt)

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
		require.Equal(t, float64(2), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_Canceled_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 2, sli1StartedAt, sli1EndedAt)

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
		require.Equal(t, float64(2), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_Not_Customer(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: false,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 2, sli1StartedAt, sli1EndedAt)

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
		require.Equal(t, float64(0), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_No_Recurring_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)

	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:      1,
		Quantity:   2,
		IsCanceled: true,
		StartedAt:  sli1StartedAt,
		EndedAt:    &sli1EndedAt,
	})

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
		require.Equal(t, float64(0), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_Started_In_Month_Canceled_Next_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.MiddleTimeOfMonth(2023, 8)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

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
		require.Equal(t, float64(0), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_2_Versions_Started_Before_Not_Canceled(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 6)
	contractId := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 4, entity.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1Id)

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
		require.Equal(t, float64(0), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_2_Versions_Started_Before_Canceled_Before_Month(t *testing.T) {
	ctx := context.Background()
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

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemCanceledWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 4, entity.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1EndedAt, sli1Id)

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
		require.Equal(t, float64(0), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_2_Versions_Started_Before_Canceled_End_Month(t *testing.T) {
	ctx := context.Background()
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

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemCanceledWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 4, entity.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1EndedAt, sli1Id)

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
		require.Equal(t, float64(4), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_2_Versions_Started_In_Canceled_End_Month(t *testing.T) {
	ctx := context.Background()
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

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemCanceledWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 4, entity.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1EndedAt, sli1Id)

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
		require.Equal(t, float64(4), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_2_Contracts_1_Active_SLI_1_Canceled_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)

	contract1Id := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemCanceledWithParent(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 4, entity.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1EndedAt, sli1Id)

	contract2Id := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
		require.Equal(t, float64(4), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_2_Contracts_With_1_Canceled_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)

	contract1Id := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)

	contract2Id := insertARRBreakdownContractWithOpportunity(ctx, driver, orgId)
	insertARRBreakdownServiceLineItemCanceled(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)

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
		require.Equal(t, float64(4), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_2_Versions_Started_In_Canceled_Next_Month(t *testing.T) {
	ctx := context.Background()
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

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemCanceledWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 4, entity.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1EndedAt, sli1Id)

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
		require.Equal(t, float64(0), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Draft_Contract_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus: entity.ContractStatusDraft,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Live_Contract_In_Month_No_Recurring_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})

	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:        "1",
		ParentID:  "1",
		Price:     12,
		Quantity:  2,
		StartedAt: sli1StartedAt,
	})

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
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Ended_Contract_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, entity.OpportunityEntity{})

	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:        "1",
		ParentID:  "1",
		Billed:    entity.BilledTypeAnnually,
		Price:     12,
		Quantity:  2,
		StartedAt: sli1StartedAt,
	})

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
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_Before_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_Beginning_Of_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
		require.Equal(t, float64(2), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_End_Of_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
		require.Equal(t, float64(2), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_Not_Customer(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: false,
	})

	sli1StartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_No_Recurring_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: false,
	})

	sli1StartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     12,
		Quantity:  2,
		StartedAt: sli1StartedAt,
	})

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
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_Next_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_1_Contract_1_SLI_2_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 4, entity.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1Id)

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
		require.Equal(t, float64(4), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
		require.Equal(t, float64(2), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeQuarterly, 4, 2, sli1StartedAt)

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
		require.Equal(t, float64(2), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 2, sli1StartedAt)

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
		require.Equal(t, float64(2), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_1_Contract_1_SLI_3_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	sli1Middle2At := sli1MiddleAt.AddDate(0, 0, 1)

	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemEndedWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 24, 4, entity.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1Middle2At, sli1Id)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 36, 8, entity.BilledTypeAnnually, 12, 4, sli1MiddleAt, sli1Id)

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
		require.Equal(t, float64(24), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_2_Contracts_1_SLI_1_Version(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)

	contract1Id := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	contract2Id := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 24, 2, sli1StartedAt)

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
		require.Equal(t, float64(6), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_2_Contracts_1_SLI_2_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)

	contract1Id := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contract1Id, entity.BilledTypeAnnually, 24, 4, entity.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1Id)

	contract2Id := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})

	sli2Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contract2Id, entity.BilledTypeAnnually, 120, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contract2Id, entity.BilledTypeAnnually, 240, 4, entity.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli2Id)

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
		require.Equal(t, float64(88), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_1_Contract_2_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)

	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})

	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 120, 2, sli1StartedAt)

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
		require.Equal(t, float64(22), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_1_Contract_2_SLI_2_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)

	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 24, 4, entity.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1Id)

	sli2Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 120, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 240, 4, entity.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli2Id)

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
		require.Equal(t, float64(88), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_1_Contract_1_Active_SLI_1_Contract_1_Canceled_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	sli1EndedAt := neo4jt.LastTimeOfMonth(2023, 7)

	contract1Id := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})

	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contract1Id, entity.BilledTypeAnnually, 24, 4, entity.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1Id)

	contract2Id := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})

	sli2Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contract2Id, entity.BilledTypeAnnually, 120, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemCanceledWithParent(ctx, driver, contract2Id, entity.BilledTypeAnnually, 240, 4, entity.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1EndedAt, sli2Id)

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
		require.Equal(t, float64(8), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Draft_Contract_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus: entity.ContractStatusDraft,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
		require.Equal(t, float64(0), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Live_Contract_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
		require.Equal(t, float64(0), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
		require.Equal(t, float64(2), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_In_Month_No_Recurring_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, entity.OpportunityEntity{})
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Price:     12,
		Quantity:  2,
		StartedAt: sli1StartedAt,
	})

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
		require.Equal(t, float64(0), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_In_Month_Not_Customer(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: false,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
		require.Equal(t, float64(0), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_Before_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
		require.Equal(t, float64(0), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_Beginning_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
		require.Equal(t, float64(2), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_End_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.LastTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
		require.Equal(t, float64(2), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_Next_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 8)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
		require.Equal(t, float64(0), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
		require.Equal(t, float64(2), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeQuarterly, 4, 2, sli1StartedAt)

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
		require.Equal(t, float64(2), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 2, sli1StartedAt)

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
		require.Equal(t, float64(2), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_1_Contract_1_SLI_2_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, entity.OpportunityEntity{})
	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 24, 4, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1Id)

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
		require.Equal(t, float64(8), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_2_Contracts_1_SLI_1_Version(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract1Id := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	sli2StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract2Id := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli2StartedAt,
		EndedAt:          &sli2StartedAt,
	}, entity.OpportunityEntity{})
	insertARRBreakdownServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 24, 2, sli2StartedAt)

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
		require.Equal(t, float64(6), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_2_Contracts_1_SLI_2_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	contract1Id := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, entity.OpportunityEntity{})
	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contract1Id, entity.BilledTypeAnnually, 24, 4, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1Id)

	sli2StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli2MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	contract2Id := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli2StartedAt,
		EndedAt:          &sli2StartedAt,
	}, entity.OpportunityEntity{})
	sli2Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contract2Id, entity.BilledTypeAnnually, 120, 2, sli2StartedAt, sli2MiddleAt)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contract2Id, entity.BilledTypeAnnually, 240, 4, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli2Id)

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
		require.Equal(t, float64(88), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_1_Contract_2_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)

	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contractStartedAt,
		EndedAt:          &contractStartedAt,
	}, entity.OpportunityEntity{})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	sli2StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	insertARRBreakdownServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 24, 2, sli2StartedAt)

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
		require.Equal(t, float64(6), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_1_Contract_2_SLI_2_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)

	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contractStartedAt,
		EndedAt:          &contractStartedAt,
	}, entity.OpportunityEntity{})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 24, 4, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1Id)

	sli2StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli2MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	sli2Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 120, 2, sli2StartedAt, sli2MiddleAt)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 240, 4, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli2Id)

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
		require.Equal(t, float64(88), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_2_Organizations_1_Contract_Each(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	org1Id := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)

	contract1Id := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, org1Id, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contractStartedAt,
		EndedAt:          &contractStartedAt,
	}, entity.OpportunityEntity{})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	insertARRBreakdownServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	org2Id := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract2Id := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, org2Id, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contractStartedAt,
		EndedAt:          &contractStartedAt,
	}, entity.OpportunityEntity{})

	sli2StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	insertARRBreakdownServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 24, 2, sli2StartedAt)

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
		require.Equal(t, float64(6), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_1_Contract_2_SLI_1_Canceled(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := neo4jt.FirstTimeOfMonth(2023, 7)

	contractId := insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &contractStartedAt,
		EndedAt:          &contractStartedAt,
	}, entity.OpportunityEntity{})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	sli1Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	insertARRBreakdownServiceLineItemWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 24, 4, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1Id)

	sli2StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	sli2MiddleAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	sli2EndedAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	sli2Id := insertARRBreakdownServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeAnnually, 120, 2, sli2StartedAt, sli2MiddleAt)
	insertARRBreakdownServiceLineItemCanceledWithParent(ctx, driver, contractId, entity.BilledTypeAnnually, 240, 4, entity.BilledTypeAnnually, 12, 2, sli1StartedAt, sli2EndedAt, sli2Id)

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
		require.Equal(t, float64(8), month.Churned)
	}
}

func insertARRBreakdownContractWithOpportunity(ctx context.Context, driver *neo4j.DriverWithContext, orgId string) string {
	return insertARRBreakdownContractWithOpportunityExtended(ctx, driver, orgId, entity.ContractEntity{}, entity.OpportunityEntity{})
}

func insertARRBreakdownContractWithOpportunityExtended(ctx context.Context, driver *neo4j.DriverWithContext, orgId string, contract entity.ContractEntity, opportunity entity.OpportunityEntity) string {
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, contract)
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, opportunity)
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opportunityId)
	return contractId
}

func insertARRBreakdownServiceLineItem(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, quantity int64, startedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:        id,
		ParentID:  id,
		Billed:    billedType,
		Price:     price,
		Quantity:  quantity,
		StartedAt: startedAt,
	})
	return id
}

func insertARRBreakdownServiceLineItemEnded(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, quantity int64, startedAt, endedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:        id,
		ParentID:  id,
		Billed:    billedType,
		Price:     price,
		Quantity:  quantity,
		StartedAt: startedAt,
		EndedAt:   &endedAt,
	})
	return id
}

func insertARRBreakdownServiceLineItemCanceled(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, quantity int64, startedAt, endedAt time.Time) string {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:         id,
		ParentID:   id,
		Billed:     billedType,
		Price:      price,
		Quantity:   quantity,
		IsCanceled: true,
		StartedAt:  startedAt,
		EndedAt:    &endedAt,
	})
	return id
}

func insertARRBreakdownServiceLineItemWithParent(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, quantity int64, previousBilledType entity.BilledType, previousPrice float64, previousQuantity int64, startedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:               id,
		ParentID:         parentId,
		Billed:           billedType,
		Price:            price,
		Quantity:         quantity,
		PreviousBilled:   previousBilledType,
		PreviousPrice:    previousPrice,
		PreviousQuantity: previousQuantity,
		StartedAt:        startedAt,
	})
}

func insertARRBreakdownServiceLineItemEndedWithParent(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, quantity int64, previousBilledType entity.BilledType, previousPrice float64, previousQuantity int64, startedAt, endedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:               id,
		ParentID:         parentId,
		Billed:           billedType,
		Price:            price,
		Quantity:         quantity,
		PreviousBilled:   previousBilledType,
		PreviousPrice:    previousPrice,
		PreviousQuantity: previousQuantity,
		StartedAt:        startedAt,
		EndedAt:          &endedAt,
	})
}

func insertARRBreakdownServiceLineItemCanceledWithParent(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, quantity int64, previousBilledType entity.BilledType, previousPrice float64, previousQuantity int64, startedAt, endedAt time.Time, parentId string) {
	rand, _ := uuid.NewRandom()
	id := rand.String()
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		ID:               id,
		ParentID:         parentId,
		Billed:           billedType,
		Price:            price,
		Quantity:         quantity,
		PreviousBilled:   previousBilledType,
		PreviousPrice:    previousPrice,
		PreviousQuantity: previousQuantity,
		IsCanceled:       true,
		StartedAt:        startedAt,
		EndedAt:          &endedAt,
	})
}
