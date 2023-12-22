package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_Dashboard_Retention_Rate_No_Period_No_Data_In_DB(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate_no_period",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "0", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 12, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
		require.Equal(t, 0, month.RenewCount)
		require.Equal(t, 0, month.ChurnCount)
	}
}

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

func TestQueryResolver_Dashboard_Retention_Rate_PeriodIntervals(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	assert_Dashboard_Retention_Rate_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-01-31T00:00:00.000Z", 1)
	assert_Dashboard_Retention_Rate_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-01-01T00:00:00.000Z", 1)
	assert_Dashboard_Retention_Rate_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-02-01T00:00:00.000Z", 2)
	assert_Dashboard_Retention_Rate_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2020-02-28T00:00:00.000Z", 2)
	assert_Dashboard_Retention_Rate_PeriodIntervals(t, "2020-01-01T00:00:00.000Z", "2029-12-01T00:00:00.000Z", 120)
}

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
	require.Equal(t, "0", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, months, len(dashboardReport.Dashboard_RetentionRate.PerMonth))
}

func TestQueryResolver_Dashboard_Retention_Rate_Hidden_Organization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		Hide:       true,
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "0", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 0, month.RenewCount)
		require.Equal(t, 0, month.ChurnCount)
	}
}

func TestQueryResolver_Dashboard_Retention_Rate_Prospect_Organization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: false,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeAnnually, 12, 2, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "0", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 0, month.RenewCount)
		require.Equal(t, 0, month.ChurnCount)
	}
}

func TestQueryResolver_Dashboard_Retention_Rate_1_Renewal_1_SLI_V1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "+100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 1, month.RenewCount)
		require.Equal(t, 0, month.ChurnCount)
	}
}

func TestQueryResolver_Dashboard_Retention_Rate_1_Renewal_1_SLI_V2(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := neo4jt.MiddleTimeOfMonth(2023, 6)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	sliId := insertServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 1, sli1StartedAt, sli1EndedAt)
	insertServiceLineItemWithParent(ctx, driver, contractId, entity.BilledTypeMonthly, 2, 1, entity.BilledTypeMonthly, 1, 1, sli1EndedAt, sliId)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "+100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 1, month.RenewCount)
		require.Equal(t, 0, month.ChurnCount)
	}
}

func TestQueryResolver_Dashboard_Retention_Rate_1_Renewal_2_SLI_V1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})

	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "+100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 1, month.RenewCount)
		require.Equal(t, 0, month.ChurnCount)
	}
}

func TestQueryResolver_Dashboard_Retention_Rate_1_Renewal_2_SLI_V2(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := neo4jt.MiddleTimeOfMonth(2023, 6)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})

	sli1Id := insertServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 1, sli1StartedAt, sli1EndedAt)
	insertServiceLineItemWithParent(ctx, driver, contractId, entity.BilledTypeMonthly, 2, 1, entity.BilledTypeMonthly, 1, 1, sli1EndedAt, sli1Id)

	sli2Id := insertServiceLineItemEnded(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 1, sli1StartedAt, sli1EndedAt)
	insertServiceLineItemWithParent(ctx, driver, contractId, entity.BilledTypeMonthly, 2, 1, entity.BilledTypeMonthly, 1, 1, sli1EndedAt, sli2Id)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 4})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "+100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 1, month.RenewCount)
		require.Equal(t, 0, month.ChurnCount)
	}
}

func TestQueryResolver_Dashboard_Retention_Rate_2_Renewals_1_SLI_V1(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	sli2StartedAt := neo4jt.FirstTimeOfMonth(2023, 5)
	contract2Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract2Id, entity.BilledTypeMonthly, 1, 1, sli2StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "+100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 2, month.RenewCount)
		require.Equal(t, 0, month.ChurnCount)
	}
}

func TestQueryResolver_Dashboard_Retention_Rate_2_Renewals_1_SLI_V2(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := neo4jt.MiddleTimeOfMonth(2023, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	sli1Id := insertServiceLineItemEnded(ctx, driver, contract1Id, entity.BilledTypeMonthly, 1, 1, sli1StartedAt, sli1EndedAt)
	insertServiceLineItemWithParent(ctx, driver, contract1Id, entity.BilledTypeMonthly, 2, 1, entity.BilledTypeMonthly, 1, 1, sli1StartedAt, sli1Id)

	sli2StartedAt := neo4jt.FirstTimeOfMonth(2023, 5)
	sli2EndedAt := neo4jt.MiddleTimeOfMonth(2023, 5)
	contract2Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	sli2Id := insertServiceLineItemEnded(ctx, driver, contract2Id, entity.BilledTypeMonthly, 1, 1, sli2StartedAt, sli2EndedAt)
	insertServiceLineItemWithParent(ctx, driver, contract2Id, entity.BilledTypeMonthly, 2, 1, entity.BilledTypeMonthly, 1, 1, sli2StartedAt, sli2Id)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 4})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "+100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 2, month.RenewCount)
		require.Equal(t, 0, month.ChurnCount)
	}
}

func TestQueryResolver_Dashboard_Retention_Rate_Renewals_Monthly_Contract_Monthly_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "+100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 1, month.RenewCount)
		require.Equal(t, 0, month.ChurnCount)
	}
}

func TestQueryResolver_Dashboard_Retention_Rate_Renewals_Monthly_Contract_Quarterly_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeQuarterly, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := neo4jt.FirstTimeOfMonth(2023, 7)
	endTime := neo4jt.FirstTimeOfMonth(2023, 12)
	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "+100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 6, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 7, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 8, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 9, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 10, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 11, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 12, 1, 0)
}

func TestQueryResolver_Dashboard_Retention_Rate_Renewals_Monthly_Contract_Annually_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := neo4jt.FirstTimeOfMonth(2023, 7)
	endTime := neo4jt.FirstTimeOfMonth(2024, 7)
	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "-100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 13, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 7, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 8, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 9, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 10, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 11, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 12, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 1, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 2, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 3, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 4, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 5, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 6, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 7, 0, 0)
}

func TestQueryResolver_Dashboard_Retention_Rate_Renewals_Quarterly_Contract_Monthly_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleQuarterlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := neo4jt.FirstTimeOfMonth(2023, 7)
	endTime := neo4jt.FirstTimeOfMonth(2024, 7)
	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "-100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 13, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 7, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 8, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 9, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 10, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 11, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 12, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 1, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 2, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 3, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 4, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 5, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 6, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 7, 0, 0)
}

func TestQueryResolver_Dashboard_Retention_Rate_Renewals_Quarterly_Contract_Quarterly_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleQuarterlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeQuarterly, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := neo4jt.FirstTimeOfMonth(2023, 7)
	endTime := neo4jt.FirstTimeOfMonth(2024, 7)
	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "-100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 13, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 7, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 8, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 9, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 10, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 11, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 12, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 1, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 2, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 3, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 4, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 5, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 6, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 7, 0, 0)
}

func TestQueryResolver_Dashboard_Retention_Rate_Renewals_Quarterly_Contract_Annually_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleQuarterlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := neo4jt.FirstTimeOfMonth(2023, 7)
	endTime := neo4jt.FirstTimeOfMonth(2024, 7)
	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "-100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 13, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 7, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 8, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 9, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 10, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 11, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 12, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 1, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 2, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 3, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 4, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 5, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 6, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 7, 0, 0)
}

func TestQueryResolver_Dashboard_Retention_Rate_Renewals_Annually_Contract_Monthly_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleAnnualRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := neo4jt.FirstTimeOfMonth(2023, 7)
	endTime := neo4jt.FirstTimeOfMonth(2024, 7)
	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "-100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 13, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 7, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 8, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 9, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 10, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 11, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 12, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 1, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 2, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 3, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 4, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 5, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 6, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 7, 0, 0)
}

func TestQueryResolver_Dashboard_Retention_Rate_Renewals_Annually_Contract_Quarterly_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleAnnualRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeQuarterly, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := neo4jt.FirstTimeOfMonth(2023, 7)
	endTime := neo4jt.FirstTimeOfMonth(2024, 7)
	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "-100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 13, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 7, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 8, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 9, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 10, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 11, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 12, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 1, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 2, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 3, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 4, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 5, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 6, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 7, 0, 0)
}

func TestQueryResolver_Dashboard_Retention_Rate_Renewals_Annually_Contract_Annually_SLI(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleAnnualRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeAnnually, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := neo4jt.FirstTimeOfMonth(2023, 7)
	endTime := neo4jt.FirstTimeOfMonth(2024, 7)
	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "-100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 13, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 7, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 8, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 9, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 10, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 11, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 12, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 1, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 2, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 3, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 4, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 5, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 6, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 7, 0, 0)
}

func TestQueryResolver_Dashboard_Retention_Rate_Renewals_1_Multi_Year_Contract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleAnnualRenewal,
		RenewalPeriods:   utils.Ptr[int64](1),
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := neo4jt.FirstTimeOfMonth(2023, 7)
	endTime := neo4jt.FirstTimeOfMonth(2025, 7)
	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "-100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 25, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 7, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 8, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 9, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 10, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 11, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 12, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 1, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 2, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 3, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 4, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 5, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 6, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 7, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 8, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 9, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 10, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 11, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 12, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2025, 1, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2025, 2, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2025, 3, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2025, 4, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2025, 5, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2025, 6, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2025, 7, 0, 0)
}

func TestQueryResolver_Dashboard_Retention_Rate_Renewals_2_Multi_Year_Contract(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		RenewalCycle:     entity.RenewalCycleAnnualRenewal,
		RenewalPeriods:   utils.Ptr[int64](2),
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := neo4jt.FirstTimeOfMonth(2023, 7)
	endTime := neo4jt.FirstTimeOfMonth(2025, 7)
	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": startTime.Format(format),
			"end":   endTime.Format(format),
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "-100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 25, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 7, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 8, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 9, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 10, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 11, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 12, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 1, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 2, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 3, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 4, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 5, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 6, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 7, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 8, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 9, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 10, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 11, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2024, 12, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2025, 1, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2025, 2, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2025, 3, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2025, 4, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2025, 5, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2025, 6, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2025, 7, 0, 0)
}

func TestQueryResolver_Dashboard_Retention_Rate_Churned_Before_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 5)
	sli1EndedAt := neo4jt.MiddleTimeOfMonth(2023, 6)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1EndedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "0", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 0, month.RenewCount)
		require.Equal(t, 0, month.ChurnCount)
	}
}

func TestQueryResolver_Dashboard_Retention_Rate_Churned_After_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 5)
	sli1EndedAt := neo4jt.MiddleTimeOfMonth(2023, 8)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1EndedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2023-07-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "+100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
		require.Equal(t, 2023, month.Year)
		require.Equal(t, 7, month.Month)
		require.Equal(t, 1, month.RenewCount)
		require.Equal(t, 0, month.ChurnCount)
	}
}

func TestQueryResolver_Dashboard_Retention_Rate_Churned_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sli1StartedAt := neo4jt.FirstTimeOfMonth(2023, 5)
	sli1EndedAt := neo4jt.MiddleTimeOfMonth(2023, 7)
	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1EndedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 1, sli1StartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2023-05-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "-100", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 3, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 5, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 6, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 7, 0, 1)
}

func TestQueryResolver_Dashboard_Retention_Rate_1_Renewal_1_Churned_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateDefaultUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{
		IsCustomer: true,
	})

	sliStartedAt := neo4jt.FirstTimeOfMonth(2023, 6)
	sliEndedAt := neo4jt.FirstTimeOfMonth(2023, 7)
	contract1Id := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusLive,
		ServiceStartedAt: &sliStartedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contract1Id, entity.BilledTypeMonthly, 1, 1, sliStartedAt)

	contractId := insertContractWithOpportunity(ctx, driver, orgId, entity.ContractEntity{
		ContractStatus:   entity.ContractStatusEnded,
		ServiceStartedAt: &sliStartedAt,
		EndedAt:          &sliEndedAt,
		RenewalCycle:     entity.RenewalCycleMonthlyRenewal,
	}, entity.OpportunityEntity{})
	insertServiceLineItem(ctx, driver, contractId, entity.BilledTypeMonthly, 1, 1, sliStartedAt)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2023-06-01T00:00:00.000Z",
			"end":   "2023-07-01T00:00:00.000Z",
		})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(50), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, "+50", dashboardReport.Dashboard_RetentionRate.IncreasePercentage)
	require.Equal(t, 2, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 6, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 7, 1, 1)
}

func assertRetentionRateMonthData(t *testing.T, dashboardReport *model.DashboardRetentionRate, year, month, renewExpected, churnExpected int) {
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
