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

func Test_Dashboard_ARR_Breakdown_No_Period_No_Data_In_DB(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown_no_period",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	response := callGraphQLExpectError(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": "2020-02-01T00:00:00.000Z",
			"end":   "2020-01-01T00:00:00.000Z",
		})

	require.Contains(t, response.Message, "Failed to get the data for period")
}

func Test_Dashboard_ARR_Breakdown_PeriodIntervals(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 6)
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

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
	require.Equal(t, "-100%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_Started_In_Month_Canceled_End_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_Canceled_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_Canceled_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 3, 2, sli1StartedAt, sli1EndedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_Canceled_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 2, sli1StartedAt, sli1EndedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_Hidden_Organization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Hide:       true,
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 2, sli1StartedAt, sli1EndedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_Prospects(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: false,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 2, sli1StartedAt, sli1EndedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)

	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Price:     1,
		Quantity:  2,
		Canceled:  true,
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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := utils.MiddleTimeOfMonth(2023, 8)
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+24", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.MiddleTimeOfMonth(2023, 6)
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndedAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndedAt, sli1Id)

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

	require.Equal(t, float64(48), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 6)
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemCanceledWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1EndedAt, sli1Id)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 6)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemCanceledWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1EndedAt, sli1Id)

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
	require.Equal(t, "-100%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(48), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_2_Versions_Started_In_Canceled_End_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)

	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemCanceledWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1EndedAt, sli1Id)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(48), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_2_Contracts_1_Active_SLI_1_Canceled_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)

	contract1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemCanceledWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1EndedAt, sli1Id)

	contract2StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
	}, neo4jentity.OpportunityEntity{})

	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+24", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(48), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_2_Contracts_With_1_Canceled_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)

	contract1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)

	contract2StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contract2StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(48), month.Cancellations)
	}
}

func Test_Dashboard_ARR_Breakdown_Cancellations_SLI_2_Versions_Started_In_Canceled_Next_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	sli1EndedAt := utils.MiddleTimeOfMonth(2023, 8)

	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemCanceledWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1EndedAt, sli1Id)

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

	require.Equal(t, float64(48), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+48", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus: neo4jenum.ContractStatusDraft,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
		require.Equal(t, float64(24), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_Before_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+24", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_Churned_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractEndedAt := utils.MiddleTimeOfMonth(2023, 7)
	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contractServiceStartedAt,
		EndedAt:          &contractEndedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_End_Of_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.LastTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+24", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_Hidden_Organization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Hide:       true,
		IsCustomer: true,
	})

	sli1StartedAt := utils.LastTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_Prospect(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: false,
	})

	sli1StartedAt := utils.LastTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: false,
	})

	sli1StartedAt := utils.LastTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 8)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1Id)

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

	require.Equal(t, float64(48), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+48", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(48), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+24", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 4, 2, sli1StartedAt)

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

	require.Equal(t, float64(32), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+32", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(32), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_Contract_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 2, sli1StartedAt)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+24", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_1_Contract_1_SLI_3_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	sli1Middle2At := sli1MiddleAt.AddDate(0, 0, 1)

	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemEndedWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1Middle2At, sli1Id)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 36, 8, neo4jenum.BilledTypeAnnually, 12, 4, sli1MiddleAt, sli1Id)

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

	require.Equal(t, float64(288), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+288", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(288), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_2_Contracts_1_SLI_1_Version(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)

	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 24, 2, sli1StartedAt)

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

	require.Equal(t, float64(72), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+72", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(72), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_2_Contracts_1_SLI_2_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)

	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1Id)

	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 120, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 240, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli2Id)

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

	require.Equal(t, float64(1056), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+1056", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(1056), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_1_Contract_2_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)

	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 120, 2, sli1StartedAt)

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

	require.Equal(t, float64(264), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+264", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(264), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_1_Contract_2_SLI_2_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)

	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1Id)

	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 120, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 240, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli2Id)

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

	require.Equal(t, float64(1056), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+1056", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(1056), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Newly_Contracted_1_Contract_1_Active_SLI_1_Contract_1_Canceled_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)

	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1Id)

	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 120, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemCanceledWithParent(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 240, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1EndedAt, sli2Id)

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

	require.Equal(t, float64(96), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+96", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(96), month.NewlyContracted)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Draft_Contract_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus: neo4jenum.ContractStatusDraft,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+24", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.NewlyContracted)
		require.Equal(t, float64(0), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Hidden_Organization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Hide:       true,
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Prospect(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: false,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractServiceStartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractEndedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contractServiceStartedAt,
		EndedAt:          &contractEndedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_In_Month_No_Recurring_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: false,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_End_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.LastTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_Next_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 8)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 4, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(32), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_Contract_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_1_Contract_1_SLI_2_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1Id)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(96), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_2_Contracts_1_SLI_1_Version(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

	sli2StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli2StartedAt,
		EndedAt:          &sli2StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 24, 2, sli2StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(72), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_2_Contracts_1_SLI_2_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1MiddleAt,
	}, neo4jentity.OpportunityEntity{})
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemEndedWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1MiddleAt, sli1Id)

	sli2StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli2MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli2StartedAt,
		EndedAt:          &sli2MiddleAt,
	}, neo4jentity.OpportunityEntity{})
	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 120, 2, sli2StartedAt, sli2MiddleAt)
	neo4jtest.InsertServiceLineItemEndedWithParent(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 240, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt, sli2Id)

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
	require.Equal(t, "-100%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_1_Contract_2_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 7)

	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contractStartedAt,
		EndedAt:          &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

	sli2StartedAt := utils.FirstTimeOfMonth(2023, 7)
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 2, sli2StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(72), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_1_Contract_2_SLI_2_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 7)

	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contractStartedAt,
		EndedAt:          &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1Id)

	sli2StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli2MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 120, 2, sli2StartedAt, sli2MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 240, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli2Id)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(1056), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_2_Organizations_1_Contract_Each(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	org1Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 7)

	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, org1Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contractStartedAt,
		EndedAt:          &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

	org2Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, org2Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contractStartedAt,
		EndedAt:          &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli2StartedAt := utils.FirstTimeOfMonth(2023, 7)
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 24, 2, sli2StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(72), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Churned_1_Contract_2_SLI_1_Canceled(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	contractStartedAt := utils.FirstTimeOfMonth(2023, 7)

	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &contractStartedAt,
		EndedAt:          &contractStartedAt,
	}, neo4jentity.OpportunityEntity{})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1Id)

	sli2StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli2MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	sli2EndedAt := utils.MiddleTimeOfMonth(2023, 7)
	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 120, 2, sli2StartedAt, sli2MiddleAt)
	neo4jtest.InsertServiceLineItemCanceledWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 240, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli2EndedAt, sli2Id)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(96), month.Churned)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Draft_Contract_No_Upsell_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus: neo4jenum.ContractStatusDraft,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Draft_Contract_With_Upsell_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus: neo4jenum.ContractStatusDraft,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sliId)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Live_Contract_No_Upsell_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+24", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Live_Contract_With_Upsell_In_Month_Should_Be_0(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 24 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	// 96 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sliId)

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

	require.Equal(t, float64(96), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+96", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Live_Contract_With_Upsell_And_SLI_Canceled_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 24 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	// 96 / year
	neo4jtest.InsertServiceLineItemCanceledWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1EndedAt, sliId)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Live_Contract_With_Downgrade_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 24 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	// 12 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sliId)

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

	require.Equal(t, float64(12), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+12", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Ended_Contract_No_Upsell_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1EndedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Ended_Contract_With_Upsell_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1EndedAt,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sliId)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Contract_No_Upsell_No_Recurring_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Contract_With_Upsell_No_Recurring_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Price:     12,
		Quantity:  2,
		StartedAt: sli1StartedAt,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		ParentID:         sliId,
		Price:            24,
		Quantity:         4,
		PreviousPrice:    12,
		PreviousQuantity: 2,
		StartedAt:        sli1StartedAt,
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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Hidden_Organization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Hide:       true,
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sliId)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Prospect(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: false,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sliId)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Before_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sliId)

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

	require.Equal(t, float64(96), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Beginning_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 24 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	// 96 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sliId)

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

	require.Equal(t, float64(96), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "3×", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(72), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_End_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.LastTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 24 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 96 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(96), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "3×", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(72), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_Next_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 8)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+24", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_2_SLI_Versions_Both_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 24 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 96 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 24, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(96), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "3×", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(72), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_2_SLI_Versions_1_Annually_1_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 24 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 96 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 12, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(96), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "3×", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(72), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_2_SLI_Versions_1_Annually_1_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 24 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 576 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 12, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(576), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "23×", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(552), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_2_SLI_Versions_Both_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 96 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 12, 2, sli1StartedAt, sli1EndAt)
	//384 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 24, 4, neo4jenum.BilledTypeQuarterly, 12, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(384), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "3×", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(288), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_2_SLI_Versions_1_Quarterly_1_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 48 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 6, 2, sli1StartedAt, sli1EndAt)
	// 288 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 12, 2, neo4jenum.BilledTypeQuarterly, 6, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(288), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "5×", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(240), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Upsells_2_SLI_Versions_Both_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 288 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 12, 2, sli1StartedAt, sli1EndAt)
	// 1152 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 24, 4, neo4jenum.BilledTypeMonthly, 12, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(1152), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "3×", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(864), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_Draft_Contract_With_Downgrade_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus: neo4jenum.ContractStatusDraft,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sliId)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_Live_Contract_No_Downgrade_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+24", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_Live_Contract_With_Downgrade_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 24 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	// 12 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sliId)

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

	require.Equal(t, float64(12), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "-50%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(12), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_Live_Contract_With_Downgrade_And_SLI_Canceled_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 24 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	// 12 / year
	neo4jtest.InsertServiceLineItemCanceledWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1EndedAt, sliId)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_Live_Contract_With_Upsell_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 12 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	// 24 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sliId)

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

	require.Equal(t, float64(48), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+48", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_Ended_Contract_No_Downgrade_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+24", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_Ended_Contract_With_Downgrade_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1EndedAt,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sliId)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_Contract_No_Downgrade_No_Recurring_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_Contract_With_Downgrade_No_Recurring_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		Price:     12,
		Quantity:  2,
		StartedAt: sli1StartedAt,
	})
	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		ParentID:         sliId,
		Price:            6,
		Quantity:         2,
		PreviousPrice:    12,
		PreviousQuantity: 2,
		StartedAt:        sli1StartedAt,
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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_Hidden_Organization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: false,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sliId)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_Prospect(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Hide:       false,
		IsCustomer: false,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sliId)

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
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_Before_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sliId)

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

	require.Equal(t, float64(12), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_Beginning_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 24 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	// 12 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sliId)

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

	require.Equal(t, float64(12), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "-50%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(12), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_End_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.LastTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 24 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 12 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(12), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "-50%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(12), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_Next_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 8)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+24", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_2_SLI_Versions_Both_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 24 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 12 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(12), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "-50%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(12), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_2_SLI_Versions_1_Annually_1_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 24 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 16 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 2, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(16), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "-33%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(8), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_2_SLI_Versions_1_Annually_1_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 32 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 16, 2, sli1StartedAt, sli1EndAt)
	// 24 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 2, neo4jenum.BilledTypeAnnually, 16, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "-25%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(8), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_2_SLI_Versions_Both_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 96 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 12, 2, sli1StartedAt, sli1EndAt)
	// 48 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 6, 2, neo4jenum.BilledTypeQuarterly, 12, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(48), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "-50%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(48), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_2_SLI_Versions_1_Quarterly_1_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 40 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 5, 2, sli1StartedAt, sli1EndAt)
	// 24 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 2, neo4jenum.BilledTypeQuarterly, 5, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "-40%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(16), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Downgrades_2_SLI_Versions_Both_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	// 288 / year
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 12, 2, sli1StartedAt, sli1EndAt)
	// 144 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 6, 2, neo4jenum.BilledTypeMonthly, 12, 2, sli1EndAt, sliId)

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

	require.Equal(t, float64(144), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "-50%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(144), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_1_Contract_With_Upsell_1_Contract_Without_Upsell(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	//contract 1
	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 24 / year
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 48 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sli1Id)

	//contract 2
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 12 / year
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 2, sli1StartedAt)

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

	require.Equal(t, float64(60), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+67%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_1_Contract_With_Downgrade_1_Contract_Without_Downgrade(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	//contract 1
	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 24 / year
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 12 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sli1Id)

	//contract 2
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 12 / year
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 2, sli1StartedAt)

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

	require.Equal(t, float64(24), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "-33%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(12), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_2_Contracts_With_Upsells(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	//contract 1
	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 24 / year
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 48 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sli1Id)

	//contract 2
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 12 / year
	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 2, sli1StartedAt, sli1EndAt)
	// 24 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 4, neo4jenum.BilledTypeAnnually, 6, 2, sli1EndAt, sli2Id)

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

	require.Equal(t, float64(72), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+100%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(36), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_2_Contracts_With_2_Upsells(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	//contract 1
	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 24 / year
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 48 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sli1Id)

	// 12 / year
	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 6, 2, sli1StartedAt, sli1EndAt)
	// 24 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 6, 4, neo4jenum.BilledTypeAnnually, 6, 2, sli1EndAt, sli2Id)

	//contract 2
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 240 / year
	sli3Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 120, 2, sli1StartedAt, sli1EndAt)
	// 480 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 120, 4, neo4jenum.BilledTypeAnnually, 120, 2, sli1EndAt, sli3Id)

	// 120 / year
	sli4Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 60, 2, sli1StartedAt, sli1EndAt)
	// 240 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 60, 4, neo4jenum.BilledTypeAnnually, 60, 2, sli1EndAt, sli4Id)

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

	require.Equal(t, float64(792), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+100%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(396), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_2_Contracts_With_Downgrades(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	//contract 1
	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 24 / year
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 12 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sli1Id)

	//contract 2
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 12 / year
	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 2, sli1StartedAt, sli1EndAt)
	// 6 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 1, neo4jenum.BilledTypeAnnually, 6, 2, sli1EndAt, sli2Id)

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

	require.Equal(t, float64(18), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "-50%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(18), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_2_Contracts_With_2_Downgrades(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	//contract 1
	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 24 / year
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 12 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sli1Id)

	// 12 / year
	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 6, 2, sli1StartedAt, sli1EndAt)
	// 6 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 6, 1, neo4jenum.BilledTypeAnnually, 6, 2, sli1EndAt, sli2Id)

	//contract 2
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 240 / year
	sli3Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 120, 2, sli1StartedAt, sli1EndAt)
	// 120 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 60, 2, neo4jenum.BilledTypeAnnually, 120, 2, sli1EndAt, sli3Id)

	// 120 / year
	sli4Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 60, 2, sli1StartedAt, sli1EndAt)
	// 60 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 60, 1, neo4jenum.BilledTypeAnnually, 60, 2, sli1EndAt, sli4Id)

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

	require.Equal(t, float64(198), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "-50%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(198), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_2_Organizations_1_Contract_With_Upsell_Each(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	//contract 1
	org1Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})
	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, org1Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 24 / year
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 48 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 4, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sli1Id)

	//contract 2
	org2Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, org2Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 12 / year
	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 2, sli1StartedAt, sli1EndAt)
	// 24 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 4, neo4jenum.BilledTypeAnnually, 6, 2, sli1EndAt, sli2Id)

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

	require.Equal(t, float64(72), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+100%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(36), month.Upsells)
	}
}

func Test_Dashboard_ARR_Breakdown_2_Organizations_1_Contract_With_Downgrade_Each(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	//contract 1
	org1Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})
	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, org1Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 24 / year
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 12 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sli1Id)

	//contract 2
	org2Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, org2Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 12 / year
	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 2, sli1StartedAt, sli1EndAt)
	// 6 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 1, neo4jenum.BilledTypeAnnually, 6, 2, sli1EndAt, sli2Id)

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

	require.Equal(t, float64(18), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "-50%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(18), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_2_Contracts_With_Downgrade_Each(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	//contract 1
	org1Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})
	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, org1Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 24 / year
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1EndAt)
	// 12 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1EndAt, sli1Id)

	//contract 2
	org2Id := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, org2Id, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 12 / year
	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 2, sli1StartedAt, sli1EndAt)
	// 6 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 1, neo4jenum.BilledTypeAnnually, 6, 2, sli1EndAt, sli2Id)

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

	require.Equal(t, float64(18), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "-50%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(18), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_1_Organization_With_1_Contract_With_Downgrade_1_Contract_With_Upsell_SLI_2_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	//contract 1
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})
	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1MiddleAt := utils.MiddleTimeOfMonth(2023, 7)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 24 / year
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1MiddleAt)
	// 12 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1MiddleAt, sli1Id)

	//contract 2
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 12 / year
	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 2, sli1StartedAt, sli1MiddleAt)
	// 24 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 12, 2, neo4jenum.BilledTypeAnnually, 6, 2, sli1MiddleAt, sli2Id)

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

	require.Equal(t, float64(36), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "0%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(12), month.Upsells)
		require.Equal(t, float64(12), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_1_Organization_With_1_Contract_With_Downgrade_1_Contract_With_Upsell_SLI_3_Versions(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	//contract 1
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})
	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1End1At := utils.MiddleTimeOfMonth(2023, 7)
	sli1End2At := sli1End1At.Add(time.Hour * 24)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 24 / year
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1End1At)
	// 12 / year
	neo4jtest.InsertServiceLineItemEndedWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1End1At, sli1End2At, sli1Id)
	// 6 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 3, 2, neo4jenum.BilledTypeAnnually, 6, 2, sli1End2At, sli1Id)

	//contract 2
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 12 / year
	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 2, sli1StartedAt, sli1End1At)
	// 24 / year
	neo4jtest.InsertServiceLineItemEndedWithParent(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 4, neo4jenum.BilledTypeAnnually, 6, 2, sli1End1At, sli1End2At, sli2Id)
	// 36 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 6, neo4jenum.BilledTypeAnnually, 6, 4, sli1End2At, sli2Id)

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

	require.Equal(t, float64(42), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+17%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.Upsells)
		require.Equal(t, float64(18), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_1_Organization_With_1_Contract_With_Downgrade_1_Contract_With_Upsell_SLI_4_Versions_UP_DOWN_IN_MIDDLE(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	//contract 1
	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})
	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1End1At := utils.MiddleTimeOfMonth(2023, 7)
	sli1End2At := sli1End1At.Add(time.Hour * 24)
	sli1End3At := sli1End2At.Add(time.Hour * 24)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 24 / year
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt, sli1End1At)
	// 12 / year
	neo4jtest.InsertServiceLineItemEndedWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 6, 2, neo4jenum.BilledTypeAnnually, 12, 2, sli1End1At, sli1End2At, sli1Id)
	// 48 / year
	neo4jtest.InsertServiceLineItemEndedWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 24, 2, neo4jenum.BilledTypeAnnually, 6, 2, sli1End2At, sli1End3At, sli1Id)
	// 6 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 3, 2, neo4jenum.BilledTypeAnnually, 24, 2, sli1End3At, sli1Id)

	//contract 2
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})
	// 12 / year
	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 2, sli1StartedAt, sli1End1At)
	// 24 / year
	neo4jtest.InsertServiceLineItemEndedWithParent(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 4, neo4jenum.BilledTypeAnnually, 6, 2, sli1End1At, sli1End2At, sli2Id)
	// 6 / year
	neo4jtest.InsertServiceLineItemEndedWithParent(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 1, neo4jenum.BilledTypeAnnually, 6, 4, sli1End2At, sli1End3At, sli2Id)
	// 36 / year
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeAnnually, 6, 6, neo4jenum.BilledTypeAnnually, 6, 1, sli1End3At, sli2Id)

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

	require.Equal(t, float64(42), dashboardReport.Dashboard_ARRBreakdown.ArrBreakdown)
	require.Equal(t, "+17%", dashboardReport.Dashboard_ARRBreakdown.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	for _, month := range dashboardReport.Dashboard_ARRBreakdown.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(24), month.Upsells)
		require.Equal(t, float64(18), month.Downgrades)
	}
}

func Test_Dashboard_ARR_Breakdown_Renewals_Draft_Contract_Monthly_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus: neo4jenum.ContractStatusDraft,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 12, 2, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2024, 12)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 18, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Draft_Contract_Quarterly_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus: neo4jenum.ContractStatusDraft,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 12, 2, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2024, 12)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 18, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Draft_Contract_Annually_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus: neo4jenum.ContractStatusDraft,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2024, 12)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 18, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_In_Month_No_Recurring_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
	}, neo4jentity.OpportunityEntity{})

	neo4jtest.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, neo4jentity.ServiceLineItemEntity{
		ID:        "1",
		ParentID:  "1",
		Price:     1,
		Quantity:  1,
		StartedAt: sli1StartedAt,
	})

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2024, 12)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 18, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Prospect(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: false,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 2, sli1StartedAt, sli1EndedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2024, 12)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 18, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Hidden_Organization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
		Hide:       true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.LastTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 2, sli1StartedAt, sli1EndedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2024, 12)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 18, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Monthly_Renewal_1_SLI_Monthly_Canceled(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.FirstTimeOfMonth(2023, 10)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 1, sli1StartedAt, sli1EndedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2023, 11)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 5, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Ended_Contract_Monthly_Renewal_1_SLI_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.FirstTimeOfMonth(2023, 10)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1EndedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2023, 11)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 6, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Monthly_Renewal_1_SLI_V1_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2024, 3)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 10, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 5)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Monthly_Renewal_1_SLI_V2_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2023, 6)
	sli2date := utils.FirstTimeOfMonth(2023, 9)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 2, neo4jenum.BilledTypeMonthly, 5, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2024, 3)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 10, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 10)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 10)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 10)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 10)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 10)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 10)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Monthly_Renewal_1_SLI_Quarterly_Canceled(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.FirstTimeOfMonth(2024, 1)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 5, 1, sli1StartedAt, sli1EndedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2024, 3)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 10, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Ended_Contract_Monthly_Renewal_1_SLI_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.FirstTimeOfMonth(2024, 1)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1EndedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2024, 3)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 10, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Monthly_Renewal_1_SLI_V1_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2024, 3)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 10, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 5)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Monthly_Renewal_1_SLI_V2_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2023, 6)
	sli2date := utils.FirstTimeOfMonth(2024, 1)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 3, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 5, 1, neo4jenum.BilledTypeQuarterly, 3, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2024, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 13, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 3)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 3)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 5)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Monthly_Renewal_1_SLI_Annually_Canceled(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.FirstTimeOfMonth(2024, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItemCanceled(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 5, 1, sli1StartedAt, sli1EndedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 0)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Ended_Contract_Monthly_Renewal_1_SLI_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.FirstTimeOfMonth(2024, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1EndedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 0)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Monthly_Renewal_1_SLI_V1_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 5)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Monthly_Renewal_1_SLI_V2_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2023, 6)
	sli2date := utils.FirstTimeOfMonth(2024, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 5, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 10, 1, neo4jenum.BilledTypeAnnually, 5, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 10)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Quarterly_Renewal_1_SLI_V1_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   3,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2024, 3)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 10, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 15)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 15)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 15)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Quarterly_Renewal_1_SLI_V2_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2023, 6)
	sli2date := utils.FirstTimeOfMonth(2024, 1)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   3,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 2, neo4jenum.BilledTypeMonthly, 5, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2024, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 13, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 15)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 15)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 30)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 30)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Quarterly_Renewal_1_SLI_V1_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   3,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2024, 3)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 10, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 5)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Quarterly_Renewal_1_SLI_V2_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2023, 6)
	sli2date := utils.FirstTimeOfMonth(2024, 1)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   3,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 3, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 5, 1, neo4jenum.BilledTypeQuarterly, 3, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2024, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 13, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 3)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 3)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 5)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Quarterly_Renewal_1_SLI_V1_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   3,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 5)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Quarterly_Renewal_1_SLI_V2_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2023, 6)
	sli2date := utils.FirstTimeOfMonth(2024, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   3,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 5, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 10, 1, neo4jenum.BilledTypeAnnually, 5, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 10)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Annual_Renewal_1_SLI_V1_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2024, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 13, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 60)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Annual_Renewal_1_SLI_V2_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2023, 6)
	sli2date := utils.FirstTimeOfMonth(2024, 1)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 2, neo4jenum.BilledTypeMonthly, 5, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2024, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 13, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 120)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Annual_Renewal_1_SLI_V1_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2024, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 13, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 20)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Annual_Renewal_1_SLI_V2_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2023, 6)
	sli2date := utils.FirstTimeOfMonth(2024, 1)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 3, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 5, 1, neo4jenum.BilledTypeQuarterly, 3, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2024, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 13, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 20)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Annual_Renewal_1_SLI_V1_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 5)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_Annual_Renewal_1_SLI_V2_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2023, 6)
	sli2date := utils.FirstTimeOfMonth(2024, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 5, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 10, 1, neo4jenum.BilledTypeAnnually, 5, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 10)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_1_Multi_Year_Renewal_1_SLI_V1_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 60)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 60)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_2_Multi_Year_Renewal_1_SLI_V1_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   24,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 120)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_1_Multi_Year_Renewal_1_SLI_V2_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2023, 6)
	sli2date := utils.FirstTimeOfMonth(2024, 1)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 2, neo4jenum.BilledTypeMonthly, 5, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 120)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 120)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_2_Multi_Year_Renewal_1_SLI_V2_Monthly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2023, 6)
	sli2date := utils.FirstTimeOfMonth(2024, 1)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   24,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 5, 2, neo4jenum.BilledTypeMonthly, 5, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 240)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_1_Multi_Year_Renewal_1_SLI_V1_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 20)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 20)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_2_Multi_Year_Renewal_1_SLI_V1_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   24,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 40)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_1_Multi_Year_Renewal_1_SLI_V2_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2023, 6)
	sli2date := utils.FirstTimeOfMonth(2024, 1)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 3, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 5, 1, neo4jenum.BilledTypeQuarterly, 3, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 20)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 20)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_2_Multi_Year_Renewal_1_SLI_V2_Quarterly(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2023, 6)
	sli2date := utils.FirstTimeOfMonth(2024, 1)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 3, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeQuarterly, 5, 1, neo4jenum.BilledTypeQuarterly, 3, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 20)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 20)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_1_Multi_Year_Renewal_1_SLI_V1_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 5)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_2_Multi_Year_Renewal_1_SLI_V1_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   24,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 5, 1, sli1StartedAt)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 10)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_1_Multi_Year_Renewal_1_SLI_V2_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2022, 6)
	sli2date := utils.FirstTimeOfMonth(2024, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 5, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 10, 1, neo4jenum.BilledTypeAnnually, 5, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2022, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 37, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2022, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2022, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2022, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2022, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2022, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2022, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2022, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 5)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 10)
}

func Test_Dashboard_ARR_Breakdown_Renewals_Live_Contract_2_Multi_Year_Renewal_1_SLI_V2_Annually(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1date := utils.FirstTimeOfMonth(2023, 6)
	sli2date := utils.FirstTimeOfMonth(2024, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1date,
		LengthInMonths:   24,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 5, 1, sli1date, sli2date)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 10, 1, neo4jenum.BilledTypeAnnually, 5, 1, sli2date, sliId)

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 6)
	endTime := utils.FirstTimeOfMonth(2025, 6)

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_arr_breakdown",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_ARRBreakdown model.DashboardARRBreakdown
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 25, len(dashboardReport.Dashboard_ARRBreakdown.PerMonth))

	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2023, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 6, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 7, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 8, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 9, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 10, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 11, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2024, 12, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 1, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 2, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 3, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 4, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 5, 0)
	assertRenewalsMonthData(t, &dashboardReport.Dashboard_ARRBreakdown, 2025, 6, 20)
}

func assertRenewalsMonthData(t *testing.T, dashboardReport *model.DashboardARRBreakdown, year, month int, expectedRenewals float64) {
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
	require.Equal(t, expectedRenewals, dashboardReport.PerMonth[index].Renewals)
}
