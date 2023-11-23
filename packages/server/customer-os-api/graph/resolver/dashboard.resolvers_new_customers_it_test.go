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

func TestQueryResolver_Dashboard_New_Customers_No_Period_No_Data_In_DB(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_new_customers_no_period",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_NewCustomers model.DashboardNewCustomers
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.ThisMonthCount)
	require.Equal(t, float64(0), dashboardReport.Dashboard_NewCustomers.ThisMonthIncreasePercentage)
	require.Equal(t, 12, len(dashboardReport.Dashboard_NewCustomers.PerMonth))

	for _, month := range dashboardReport.Dashboard_NewCustomers.PerMonth {
		require.Equal(t, 0, month.Count)
	}
}

func TestQueryResolver_Dashboard_New_Customers_InvalidPeriod(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	response := callGraphQLExpectError(t, "dashboard_view/dashboard_new_customers",
		map[string]interface{}{
			"start": "2020-02-01T00:00:00.000Z",
			"end":   "2020-01-01T00:00:00.000Z",
		})

	require.Contains(t, "Failed to get the data for period", response.Message)
}

func TestQueryResolver_Dashboard_New_Customers_PeriodIntervals(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	assert_Dashboard_New_Customers_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-01-31T00:00:00.000Z", 1)
	assert_Dashboard_New_Customers_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-01-01T00:00:00.000Z", 1)
	assert_Dashboard_New_Customers_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-02-01T00:00:00.000Z", 2)
	assert_Dashboard_New_Customers_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-02-28T00:00:00.000Z", 2)
	assert_Dashboard_New_Customers_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2029-12-01T00:00:00.000Z", 120)
}

func assert_Dashboard_New_Customers_PeriodIntervals(t *testing.T, start, end string, months int) {
	rawResponse := callGraphQL(t, "dashboard_view/dashboard_new_customers",
		map[string]interface{}{
			"start": start,
			"end":   end,
		})

	var dashboardReport struct {
		Dashboard_NewCustomers model.DashboardNewCustomers
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.ThisMonthCount)
	require.Equal(t, months, len(dashboardReport.Dashboard_NewCustomers.PerMonth))
}

func TestQueryResolver_Dashboard_New_Customers_ContractSignedBeforeMonth(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 6, 30, 23, 59, 59, 999, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_new_customers",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_NewCustomers model.DashboardNewCustomers
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.ThisMonthCount)
	require.Equal(t, float64(0), dashboardReport.Dashboard_NewCustomers.ThisMonthIncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_NewCustomers.PerMonth))

	for _, month := range dashboardReport.Dashboard_NewCustomers.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 0, month.Count)
	}
}

func TestQueryResolver_Dashboard_New_Customers_ContractSignedAfterMonth(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_new_customers",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_NewCustomers model.DashboardNewCustomers
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.ThisMonthCount)
	require.Equal(t, float64(0), dashboardReport.Dashboard_NewCustomers.ThisMonthIncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_NewCustomers.PerMonth))

	for _, month := range dashboardReport.Dashboard_NewCustomers.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 0, month.Count)
	}
}

func TestQueryResolver_Dashboard_New_Customers_ContractSignedAtBeginningOfMonth(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_new_customers",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_NewCustomers model.DashboardNewCustomers
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, dashboardReport.Dashboard_NewCustomers.ThisMonthCount)
	require.Equal(t, float64(0), dashboardReport.Dashboard_NewCustomers.ThisMonthIncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_NewCustomers.PerMonth))

	for _, month := range dashboardReport.Dashboard_NewCustomers.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 1, month.Count)
	}
}

func TestQueryResolver_Dashboard_New_Customers_ContractSignedAtEndOfMonth(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 7, 31, 20, 59, 59, 999, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_new_customers",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_NewCustomers model.DashboardNewCustomers
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, dashboardReport.Dashboard_NewCustomers.ThisMonthCount)
	require.Equal(t, float64(0), dashboardReport.Dashboard_NewCustomers.ThisMonthIncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_NewCustomers.PerMonth))

	for _, month := range dashboardReport.Dashboard_NewCustomers.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 1, month.Count)
	}
}

func TestQueryResolver_Dashboard_New_Customers_ContractSignedInMonth_EndedImmediately(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)
	contract1EndedAt := time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1EndedAt,
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_new_customers",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_NewCustomers model.DashboardNewCustomers
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.ThisMonthCount)
	require.Equal(t, float64(0), dashboardReport.Dashboard_NewCustomers.ThisMonthIncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_NewCustomers.PerMonth))

	for _, month := range dashboardReport.Dashboard_NewCustomers.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 0, month.Count)
	}
}

func TestQueryResolver_Dashboard_New_Customers_ContractSignedInMonth_EndedAtEndOfMonth(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)
	contract1EndedAt := time.Date(2023, 7, 31, 23, 59, 59, 999, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1EndedAt,
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_new_customers",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_NewCustomers model.DashboardNewCustomers
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.ThisMonthCount)
	require.Equal(t, float64(0), dashboardReport.Dashboard_NewCustomers.ThisMonthIncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_NewCustomers.PerMonth))

	for _, month := range dashboardReport.Dashboard_NewCustomers.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 0, month.Count)
	}
}

func TestQueryResolver_Dashboard_New_Customers_ContractSignedInMonth_EndedNextMonth(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)
	contract1EndedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1EndedAt,
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_new_customers",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_NewCustomers model.DashboardNewCustomers
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, dashboardReport.Dashboard_NewCustomers.ThisMonthCount)
	require.Equal(t, float64(0), dashboardReport.Dashboard_NewCustomers.ThisMonthIncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_NewCustomers.PerMonth))

	for _, month := range dashboardReport.Dashboard_NewCustomers.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 1, month.Count)
	}
}

func TestQueryResolver_Dashboard_New_Customers_MultipleContractsSignedInMonth_SameOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)
	contract1EndedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1EndedAt,
	})

	contract2ServiceStartedAt := time.Date(2023, 7, 20, 0, 0, 0, 0, time.UTC)
	contract2EndedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2EndedAt,
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_new_customers",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_NewCustomers model.DashboardNewCustomers
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, dashboardReport.Dashboard_NewCustomers.ThisMonthCount)
	require.Equal(t, float64(0), dashboardReport.Dashboard_NewCustomers.ThisMonthIncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_NewCustomers.PerMonth))

	for _, month := range dashboardReport.Dashboard_NewCustomers.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 1, month.Count)
	}
}

func TestQueryResolver_Dashboard_New_Customers_MultipleContractsSignedInDifferentMonths_SameOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)
	contract1EndedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1EndedAt,
	})

	contract2ServiceStartedAt := time.Date(2023, 9, 20, 0, 0, 0, 0, time.UTC)
	contract2EndedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2EndedAt,
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_new_customers",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_NewCustomers model.DashboardNewCustomers
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, dashboardReport.Dashboard_NewCustomers.ThisMonthCount)
	require.Equal(t, float64(0), dashboardReport.Dashboard_NewCustomers.ThisMonthIncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_NewCustomers.PerMonth))

	for _, month := range dashboardReport.Dashboard_NewCustomers.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 1, month.Count)
	}
}

func TestQueryResolver_Dashboard_New_Customers_MultipleContractsSignedInMonth_DifferentOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	orgId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	contract1ServiceStartedAt := time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)
	contract1EndedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId1, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1EndedAt,
	})

	contract2ServiceStartedAt := time.Date(2023, 7, 20, 0, 0, 0, 0, time.UTC)
	contract2EndedAt := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId2, entity.ContractEntity{
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2EndedAt,
	})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_new_customers",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_NewCustomers model.DashboardNewCustomers
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 2, dashboardReport.Dashboard_NewCustomers.ThisMonthCount)
	require.Equal(t, float64(0), dashboardReport.Dashboard_NewCustomers.ThisMonthIncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_NewCustomers.PerMonth))

	for _, month := range dashboardReport.Dashboard_NewCustomers.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 2, month.Count)
	}
}

func TestQueryResolver_Dashboard_New_Customers_GeneralCount1(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	// Contract 1: Signed on 01.04.2023 with a termination date in 2024
	contract1ServiceStartedAt := time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)
	contract1EndedAt := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract1ServiceStartedAt,
		EndedAt:          &contract1EndedAt,
	})

	// Contract 2: Signed on 01.04.2023 with a termination date in 2024
	contract2ServiceStartedAt := time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC)
	contract2EndedAt := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
	neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
		ServiceStartedAt: &contract2ServiceStartedAt,
		EndedAt:          &contract2EndedAt,
	})

	// Contract 3-4: Signed in 06.2023 with termination date in November/December 2023
	for i := 0; i < 2; i++ {
		contractServiceStartedAt := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
		contractEndedAt := time.Date(2023, 11, 30, 23, 59, 59, 0, time.UTC)
		neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
			ServiceStartedAt: &contractServiceStartedAt,
			EndedAt:          &contractEndedAt,
		})
	}

	// Contract 5-6: Signed in 06.2023 with termination date in December 2023
	for i := 0; i < 2; i++ {
		contractServiceStartedAt := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
		contractEndedAt := time.Date(2023, 12, 30, 23, 59, 59, 0, time.UTC)
		neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
			ServiceStartedAt: &contractServiceStartedAt,
			EndedAt:          &contractEndedAt,
		})
	}

	// Contract 7-18: Signed in 09.2023 with a termination date in 2024
	for i := 0; i < 12; i++ {
		contractServiceStartedAt := time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
		contractEndedAt := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
		neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
			ServiceStartedAt: &contractServiceStartedAt,
			EndedAt:          &contractEndedAt,
		})
	}

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_new_customers",
		map[string]interface{}{
			"start": "2023-01-15T00:00:00.000Z",
			"end":   "2023-12-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_NewCustomers model.DashboardNewCustomers
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.ThisMonthCount)
	require.Equal(t, float64(0), dashboardReport.Dashboard_NewCustomers.ThisMonthIncreasePercentage)
	require.Equal(t, 12, len(dashboardReport.Dashboard_NewCustomers.PerMonth))

	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.PerMonth[0].Count)
	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.PerMonth[1].Count)
	require.Equal(t, 1, dashboardReport.Dashboard_NewCustomers.PerMonth[2].Count)
	require.Equal(t, 1, dashboardReport.Dashboard_NewCustomers.PerMonth[3].Count)
	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.PerMonth[4].Count)
	require.Equal(t, 4, dashboardReport.Dashboard_NewCustomers.PerMonth[5].Count)
	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.PerMonth[6].Count)
	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.PerMonth[7].Count)
	require.Equal(t, 12, dashboardReport.Dashboard_NewCustomers.PerMonth[8].Count)
	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.PerMonth[9].Count)
	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.PerMonth[10].Count)
	require.Equal(t, 0, dashboardReport.Dashboard_NewCustomers.PerMonth[11].Count)
}

func TestQueryResolver_Dashboard_New_Customers_GeneralCount2(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})

	for i := 1; i <= 12; i++ {
		for j := 1; j <= i*1; j++ {
			contractServiceStartedAt := time.Date(2023, time.Month(i), 1, 0, 0, 0, 0, time.UTC)
			contractEndedAt := time.Date(2024, 01, 31, 23, 59, 59, 0, time.UTC)
			neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{
				ServiceStartedAt: &contractServiceStartedAt,
				EndedAt:          &contractEndedAt,
			})
		}
	}

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_new_customers",
		map[string]interface{}{
			"start": "2023-01-15T00:00:00.000Z",
			"end":   "2023-12-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_NewCustomers model.DashboardNewCustomers
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 12, dashboardReport.Dashboard_NewCustomers.ThisMonthCount)
	require.Equal(t, float64(0), dashboardReport.Dashboard_NewCustomers.ThisMonthIncreasePercentage)
	require.Equal(t, 12, len(dashboardReport.Dashboard_NewCustomers.PerMonth))

	require.Equal(t, 1, dashboardReport.Dashboard_NewCustomers.PerMonth[0].Count)
	require.Equal(t, 2, dashboardReport.Dashboard_NewCustomers.PerMonth[1].Count)
	require.Equal(t, 3, dashboardReport.Dashboard_NewCustomers.PerMonth[2].Count)
	require.Equal(t, 4, dashboardReport.Dashboard_NewCustomers.PerMonth[3].Count)
	require.Equal(t, 5, dashboardReport.Dashboard_NewCustomers.PerMonth[4].Count)
	require.Equal(t, 6, dashboardReport.Dashboard_NewCustomers.PerMonth[5].Count)
	require.Equal(t, 7, dashboardReport.Dashboard_NewCustomers.PerMonth[6].Count)
	require.Equal(t, 8, dashboardReport.Dashboard_NewCustomers.PerMonth[7].Count)
	require.Equal(t, 9, dashboardReport.Dashboard_NewCustomers.PerMonth[8].Count)
	require.Equal(t, 10, dashboardReport.Dashboard_NewCustomers.PerMonth[9].Count)
	require.Equal(t, 11, dashboardReport.Dashboard_NewCustomers.PerMonth[10].Count)
	require.Equal(t, 12, dashboardReport.Dashboard_NewCustomers.PerMonth[11].Count)
}
