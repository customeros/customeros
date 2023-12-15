package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

//func TestQueryResolver_Dashboard_Retention_Rate_No_Period_No_Data_In_DB(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//	neo4jt.CreateTenant(ctx, driver, tenantName)
//
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
//
//	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate_no_period",
//		map[string]interface{}{})
//
//	var dashboardReport struct {
//		Dashboard_RetentionRate model.DashboardRetentionRate
//	}
//
//	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
//	require.Nil(t, err)
//
//	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.RetentionRate)
//	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
//	require.Equal(t, 12, len(dashboardReport.Dashboard_RetentionRate.PerMonth))
//
//	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
//		require.Equal(t, float64(0), month.RenewCount)
//		require.Equal(t, float64(0), month.ChurnCount)
//	}
//}

func TestQueryResolver_Dashboard_Retention_Rate_InvalidPeriod(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	response := callGraphQLExpectError(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2020-02-01T00:00:00.000Z",
			"end":   "2020-01-01T00:00:00.000Z",
		})

	require.Contains(t, "Failed to get the data for period", response.Message)
}

//func TestQueryResolver_Dashboard_Retention_Rate_PeriodIntervals(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//	neo4jt.CreateTenant(ctx, driver, tenantName)
//
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
//
//	assert_Dashboard_Retention_Rate_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-01-31T00:00:00.000Z", 1)
//	assert_Dashboard_Retention_Rate_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-01-01T00:00:00.000Z", 1)
//	assert_Dashboard_Retention_Rate_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-02-01T00:00:00.000Z", 2)
//	assert_Dashboard_Retention_Rate_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-02-28T00:00:00.000Z", 2)
//	assert_Dashboard_Retention_Rate_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2029-12-01T00:00:00.000Z", 120)
//}

func assert_Dashboard_Retention_Rate_PeriodIntervals(t *testing.T, start, end string, months int) {
	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": start,
			"end":   end,
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, months, len(dashboardReport.Dashboard_RetentionRate.PerMonth))
}

//
//func TestQueryResolver_Dashboard_Retention_Rate_Hidden_Organization(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//	neo4jt.CreateTenant(ctx, driver, tenantName)
//	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
//
//	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
//		Hide: true,
//	})
//
//	sli1StartedAt := neo4jt.FirstTimeOfMonth(2024, 7)
//	contractId := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
//	insertMRRPerCustomerServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)
//
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})
//
//	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
//		map[string]interface{}{
//			"start": "2023-07-01T00:00:00.000Z",
//			"end":   "2023-07-01T00:00:00.000Z",
//		})
//
//	var dashboardReport struct {
//		Dashboard_RetentionRate model.DashboardRetentionRate
//	}
//
//	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
//	require.Nil(t, err)
//
//	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.RetentionRate)
//	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
//	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))
//
//	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
//		require.Equal(t, 2023, month.Year)
//		require.Equal(t, 7, month.Month)
//		require.Equal(t, float64(0), month.RenewCount)
//		require.Equal(t, float64(0), month.ChurnCount)
//	}
//}
//
//func TestQueryResolver_Dashboard_Retention_Rate_Prospect_Organization(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//	neo4jt.CreateTenant(ctx, driver, tenantName)
//	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
//
//	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
//		IsCustomer: false,
//	})
//
//	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
//	contractId := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
//	insertMRRPerCustomerServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)
//
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})
//
//	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
//		map[string]interface{}{
//			"start": "2023-07-01T00:00:00.000Z",
//			"end":   "2023-07-01T00:00:00.000Z",
//		})
//
//	var dashboardReport struct {
//		Dashboard_RetentionRate model.DashboardRetentionRate
//	}
//
//	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
//	require.Nil(t, err)
//
//	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.RetentionRate)
//	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
//	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))
//
//	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
//		require.Equal(t, 2023, month.Year)
//		require.Equal(t, 7, month.Month)
//		require.Equal(t, float64(0), month.RenewCount)
//		require.Equal(t, float64(0), month.ChurnCount)
//	}
//}
//
//func TestQueryResolver_Dashboard_Retention_Rate_Not_A_Customer(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//	neo4jt.CreateTenant(ctx, driver, tenantName)
//	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
//
//	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
//		IsCustomer: true,
//		Hide:       true,
//	})
//
//	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
//	contractId := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
//	insertMRRPerCustomerServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)
//
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})
//
//	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
//		map[string]interface{}{
//			"start": "2023-07-01T00:00:00.000Z",
//			"end":   "2023-07-01T00:00:00.000Z",
//		})
//
//	var dashboardReport struct {
//		Dashboard_RetentionRate model.DashboardRetentionRate
//	}
//
//	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
//	require.Nil(t, err)
//
//	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.RetentionRate)
//	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
//	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))
//
//	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
//		require.Equal(t, 2023, month.Year)
//		require.Equal(t, 7, month.Month)
//		require.Equal(t, float64(0), month.RenewCount)
//		require.Equal(t, float64(0), month.ChurnCount)
//	}
//}
//
//func TestQueryResolver_Dashboard_Retention_Renewals_1_Contract_1_SLI_V1(t *testing.T) {
//	ctx := context.Background()
//	defer tearDownTestCase(ctx)(t)
//	neo4jt.CreateTenant(ctx, driver, tenantName)
//	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)
//
//	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
//		IsCustomer: true,
//		Hide:       true,
//	})
//
//	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
//	contractId := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
//	insertMRRPerCustomerServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)
//
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
//	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})
//
//	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
//		map[string]interface{}{
//			"start": "2023-07-01T00:00:00.000Z",
//			"end":   "2023-07-01T00:00:00.000Z",
//		})
//
//	var dashboardReport struct {
//		Dashboard_RetentionRate model.DashboardRetentionRate
//	}
//
//	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
//	require.Nil(t, err)
//
//	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.RetentionRate)
//	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
//	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))
//
//	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 7, 2, 0)
//}

func assertRetentionRateMonthData(t *testing.T, dashboardReport *model.DashboardRetentionRate, year, month int, renewExpected, churnExpected int64) {
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
	require.Equal(t, renewExpected, dashboardReport.PerMonth[index].RenewCount)
	require.Equal(t, churnExpected, dashboardReport.PerMonth[index].ChurnCount)
}
