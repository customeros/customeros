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
)

func TestQueryResolver_Dashboard_Retention_Rate_No_Period_No_Data_In_DB(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_retention_rate_no_period",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_RetentionRate model.DashboardRetentionRate
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.RetentionRate)
	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
	require.Equal(t, 12, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	for _, month := range dashboardReport.Dashboard_RetentionRate.PerMonth {
		require.Equal(t, 0, month.RenewCount)
		require.Equal(t, 0, month.ChurnCount)
	}
}

func TestQueryResolver_Dashboard_Retention_Rate_InvalidPeriod(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	response := callGraphQLExpectError(t, "dashboard_view/dashboard_retention_rate",
		map[string]interface{}{
			"start": "2020-02-01T00:00:00.000Z",
			"end":   "2020-01-01T00:00:00.000Z",
		})

	require.Contains(t, response.Message, "Failed to get the data for period")
}

func TestQueryResolver_Dashboard_Retention_Rate_PeriodIntervals(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

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
	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
	require.Equal(t, months, len(dashboardReport.Dashboard_RetentionRate.PerMonth))
}

func TestQueryResolver_Dashboard_Retention_Rate_Hidden_Organization(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Hide:         true,
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Prospect,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeAnnually, 12, 2, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.MiddleTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	sliId := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt, sli1EndedAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 2, 1, neo4jenum.BilledTypeMonthly, 1, 1, sli1EndedAt, sliId)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})

	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.MiddleTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})

	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt, sli1EndedAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 2, 1, neo4jenum.BilledTypeMonthly, 1, 1, sli1EndedAt, sli1Id)

	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt, sli1EndedAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 2, 1, neo4jenum.BilledTypeMonthly, 1, 1, sli1EndedAt, sli2Id)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 4})

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
	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	sli2StartedAt := utils.FirstTimeOfMonth(2023, 5)
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeMonthly, 1, 1, sli2StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	sli1EndedAt := utils.MiddleTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	sli1Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt, sli1EndedAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeMonthly, 2, 1, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt, sli1Id)

	sli2StartedAt := utils.FirstTimeOfMonth(2023, 5)
	sli2EndedAt := utils.MiddleTimeOfMonth(2023, 5)
	contract2Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	sli2Id := neo4jtest.InsertServiceLineItemEnded(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeMonthly, 1, 1, sli2StartedAt, sli2EndedAt)
	neo4jtest.InsertServiceLineItemWithParent(ctx, driver, tenantName, contract2Id, neo4jenum.BilledTypeMonthly, 2, 1, neo4jenum.BilledTypeMonthly, 1, 1, sli2StartedAt, sli2Id)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 4})

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
	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeQuarterly, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2023, 12)
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
	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2024, 7)
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
	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   3,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2024, 7)
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
	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   3,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeQuarterly, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2024, 7)
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
	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   3,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2024, 7)
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
	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2024, 7)
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
	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeQuarterly, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2024, 7)
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
	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeAnnually, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2024, 7)
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
	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   12,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2025, 7)
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
	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 6)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sli1StartedAt,
		LengthInMonths:   24,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

	format := "2006-01-02T15:04:05.000Z"
	startTime := utils.FirstTimeOfMonth(2023, 7)
	endTime := utils.FirstTimeOfMonth(2025, 7)
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
	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 5)
	sli1EndedAt := utils.MiddleTimeOfMonth(2023, 6)
	contractId := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1EndedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, float64(0), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 5)
	sli1EndedAt := utils.MiddleTimeOfMonth(2023, 8)
	contractId := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1EndedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, float64(100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sli1StartedAt := utils.FirstTimeOfMonth(2023, 5)
	sli1EndedAt := utils.MiddleTimeOfMonth(2023, 7)
	contractId := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sli1StartedAt,
		EndedAt:          &sli1EndedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 1, sli1StartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 1})

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
	require.Equal(t, float64(-100), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
	require.Equal(t, 3, len(dashboardReport.Dashboard_RetentionRate.PerMonth))

	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 5, 0, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 6, 1, 0)
	assertRetentionRateMonthData(t, &dashboardReport.Dashboard_RetentionRate, 2023, 7, 0, 1)
}

func TestQueryResolver_Dashboard_Retention_Rate_1_Renewal_1_Churned_In_Month(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateUserWithId(ctx, driver, tenantName, testUserId)

	orgId := neo4jtest.CreateOrganization(ctx, driver, tenantName, neo4jentity.OrganizationEntity{
		Relationship: neo4jenum.Customer,
	})

	sliStartedAt := utils.FirstTimeOfMonth(2023, 6)
	sliEndedAt := utils.FirstTimeOfMonth(2023, 7)
	contract1Id := neo4jtest.InsertContractWithActiveRenewalOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusLive,
		ServiceStartedAt: &sliStartedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contract1Id, neo4jenum.BilledTypeMonthly, 1, 1, sliStartedAt)

	contractId := neo4jtest.InsertContractWithOpportunity(ctx, driver, tenantName, orgId, neo4jentity.ContractEntity{
		ContractStatus:   neo4jenum.ContractStatusEnded,
		ServiceStartedAt: &sliStartedAt,
		EndedAt:          &sliEndedAt,
		LengthInMonths:   1,
	}, neo4jentity.OpportunityEntity{})
	neo4jtest.InsertServiceLineItem(ctx, driver, tenantName, contractId, neo4jenum.BilledTypeMonthly, 1, 1, sliStartedAt)

	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Organization": 1})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Contract": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"Opportunity": 2})
	neo4jtest.AssertNeo4jNodeCount(ctx, t, driver, map[string]int{"ServiceLineItem": 2})

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
	require.Equal(t, float64(50), dashboardReport.Dashboard_RetentionRate.IncreasePercentageValue)
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
