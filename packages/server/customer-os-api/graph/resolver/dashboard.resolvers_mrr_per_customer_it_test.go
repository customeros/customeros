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

func TestQueryResolver_Dashboard_MRR_Per_Customer_No_Period_No_Data_In_DB(t *testing.T) {
	ctx := context.TODO()
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
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 12, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, float64(0), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_InvalidPeriod(t *testing.T) {
	ctx := context.TODO()
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
	ctx := context.TODO()
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
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, months, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_InMonth_HiddenOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Hide: true,
	})

	sli1StartedAt := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
	contractId := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, &sli1StartedAt, nil)

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
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_InMonth_ProspectOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: false,
	})

	sli1StartedAt := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
	contractId := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, &sli1StartedAt, nil)

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
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_BeforeMonth(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := time.Date(2023, 6, 30, 23, 59, 59, 999999999, time.UTC)
	contractId := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, &sli1StartedAt, nil)

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
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(2), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_AfterMonth(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC)
	contractId := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, &sli1StartedAt, nil)

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
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_AtBeginningOfMonth(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contractId := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)

	sli1StartedAt := time.Date(2023, 7, 1, 0, 0, 0, 0, time.UTC)
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, &sli1StartedAt, nil)

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
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(2), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_AtEndOfMonth(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := time.Date(2023, 7, 31, 20, 59, 59, 999999999, time.UTC)
	contractId := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, &sli1StartedAt, nil)

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
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(2), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_InMonth_EndedImmediately(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)
	contractId := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, &sli1StartedAt, &sli1StartedAt)

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
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(0), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_SLI_InMonth_EndedAtEndOfMonth(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)
	sli1EndedAt := time.Date(2023, 7, 31, 23, 59, 59, 999999999, time.UTC)
	contractId := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, &sli1StartedAt, &sli1EndedAt)

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
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(2), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_2_SLI_SameMonth_SameOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)
	contract1Id := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, &sli1StartedAt, nil)

	contract2ServiceStartedAt := time.Date(2023, 7, 20, 0, 0, 0, 0, time.UTC)
	contract2Id := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, &contract2ServiceStartedAt, nil)

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

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(4), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_2_SLI_SameMonth_DifferentOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	orgId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)
	contract1Id := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId1)
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, &sli1StartedAt, nil)

	contract2ServiceStartedAt := time.Date(2023, 7, 20, 0, 0, 0, 0, time.UTC)
	contract2Id := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId2)
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, &contract2ServiceStartedAt, nil)

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

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(2), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_2_SLI_DifferentMonths_SameOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)
	contract1Id := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, &sli1StartedAt, nil)

	contract2ServiceStartedAt := time.Date(2023, 9, 20, 0, 0, 0, 0, time.UTC)
	contract2Id := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId)
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, &contract2ServiceStartedAt, nil)

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

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(2), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_2_SLI_DifferentMonths_DifferentOrganization(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	orgId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1Id := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId1)
	sli1StartedAt := time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, &sli1StartedAt, nil)

	contract2Id := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId2)
	sli2StartedAt := time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC)
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, &sli2StartedAt, nil)

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

	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.MrrPerCustomer)
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	for _, month := range dashboardReport.Dashboard_MRRPerCustomer.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, float64(1), month.Value)
	}
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_2_SLI_SameOrganization_Overlaps_2_Months(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1Id := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId1)
	sli1StartedAt := time.Date(2023, 6, 01, 0, 0, 0, 0, time.UTC)
	sli1EndedAt := time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, &sli1StartedAt, &sli1EndedAt)

	contract2Id := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId1)
	sli2StartedAt := time.Date(2023, 7, 01, 0, 0, 0, 0, time.UTC)
	sli2EndedAt := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, &sli2StartedAt, &sli2EndedAt)

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
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 6, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	may := dashboardReport.Dashboard_MRRPerCustomer.PerMonth[0]
	require.Equal(t, 2023, may.Year)
	require.Equal(t, 5, may.Month)
	require.Equal(t, float64(0), may.Value)

	june := dashboardReport.Dashboard_MRRPerCustomer.PerMonth[1]
	require.Equal(t, 2023, june.Year)
	require.Equal(t, 6, june.Month)
	require.Equal(t, float64(2), june.Value)

	july := dashboardReport.Dashboard_MRRPerCustomer.PerMonth[2]
	require.Equal(t, 2023, july.Year)
	require.Equal(t, 7, july.Month)
	require.Equal(t, float64(4), july.Value)

	august := dashboardReport.Dashboard_MRRPerCustomer.PerMonth[3]
	require.Equal(t, 2023, august.Year)
	require.Equal(t, 8, august.Month)
	require.Equal(t, float64(4), august.Value)

	september := dashboardReport.Dashboard_MRRPerCustomer.PerMonth[4]
	require.Equal(t, 2023, september.Year)
	require.Equal(t, 9, september.Month)
	require.Equal(t, float64(2), september.Value)

	october := dashboardReport.Dashboard_MRRPerCustomer.PerMonth[5]
	require.Equal(t, 2023, october.Year)
	require.Equal(t, 10, october.Month)
	require.Equal(t, float64(0), october.Value)
}

func TestQueryResolver_Dashboard_MRR_Per_Customer_2_SLI_DifferentOrganization_Overlaps_2_Months(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})
	orgId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	contract1Id := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId1)
	sli1StartedAt := time.Date(2023, 6, 01, 0, 0, 0, 0, time.UTC)
	sli1EndedAt := time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 12, &sli1StartedAt, &sli1EndedAt)

	contract2Id := insertMRRPerCustomerContractWithOpportunity(ctx, driver, orgId2)
	sli2StartedAt := time.Date(2023, 7, 01, 0, 0, 0, 0, time.UTC)
	sli2EndedAt := time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeAnnually, 12, &sli2StartedAt, &sli2EndedAt)

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
	require.Equal(t, float64(0), dashboardReport.Dashboard_MRRPerCustomer.IncreasePercentage)
	require.Equal(t, 6, len(dashboardReport.Dashboard_MRRPerCustomer.PerMonth))

	may := dashboardReport.Dashboard_MRRPerCustomer.PerMonth[0]
	require.Equal(t, 2023, may.Year)
	require.Equal(t, 5, may.Month)
	require.Equal(t, float64(0), may.Value)

	june := dashboardReport.Dashboard_MRRPerCustomer.PerMonth[1]
	require.Equal(t, 2023, june.Year)
	require.Equal(t, 6, june.Month)
	require.Equal(t, float64(1), june.Value)

	july := dashboardReport.Dashboard_MRRPerCustomer.PerMonth[2]
	require.Equal(t, 2023, july.Year)
	require.Equal(t, 7, july.Month)
	require.Equal(t, float64(2), july.Value)

	august := dashboardReport.Dashboard_MRRPerCustomer.PerMonth[3]
	require.Equal(t, 2023, august.Year)
	require.Equal(t, 8, august.Month)
	require.Equal(t, float64(2), august.Value)

	september := dashboardReport.Dashboard_MRRPerCustomer.PerMonth[4]
	require.Equal(t, 2023, september.Year)
	require.Equal(t, 9, september.Month)
	require.Equal(t, float64(1), september.Value)

	october := dashboardReport.Dashboard_MRRPerCustomer.PerMonth[5]
	require.Equal(t, 2023, october.Year)
	require.Equal(t, 10, october.Month)
	require.Equal(t, float64(0), october.Value)
}

func insertMRRPerCustomerContractWithOpportunity(ctx context.Context, driver *neo4j.DriverWithContext, orgId string) string {
	contractId := neo4jt.CreateContractForOrganization(ctx, driver, tenantName, orgId, entity.ContractEntity{})
	opportunityId := neo4jt.CreateOpportunityForContract(ctx, driver, tenantName, contractId, entity.OpportunityEntity{})
	neo4jt.ActiveRenewalOpportunityForContract(ctx, driver, tenantName, contractId, opportunityId)
	return contractId
}

func insertServiceLineItem(ctx context.Context, driver *neo4j.DriverWithContext, contractId string, billedType entity.BilledType, price float64, startedAt, endedAt *time.Time) {
	neo4jt.CreateServiceLineItemForContract(ctx, driver, tenantName, contractId, entity.ServiceLineItemEntity{
		Billed:    billedType,
		Price:     price,
		Quantity:  2,
		StartedAt: *startedAt,
		EndedAt:   endedAt,
	})
}
