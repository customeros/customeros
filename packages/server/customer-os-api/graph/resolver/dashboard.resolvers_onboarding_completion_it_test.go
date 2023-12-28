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
	"time"
)

func TestQueryResolver_Dashboard_OnboardingCompletion_No_Period_No_Data_In_DB(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_onboarding_completion_no_period",
		map[string]interface{}{})

	var dashboardReport struct {
		Dashboard_OnboardingCompletion model.DashboardOnboardingCompletion
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0.0, dashboardReport.Dashboard_OnboardingCompletion.CompletionPercentage)
	require.Equal(t, 0.0, dashboardReport.Dashboard_OnboardingCompletion.IncreasePercentage)
	require.Equal(t, 0, len(dashboardReport.Dashboard_OnboardingCompletion.PerMonth))
}

func TestQueryResolver_Dashboard_OnboardingCompletion_InvalidPeriod(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1})

	response := callGraphQLExpectError(t, "dashboard_view/dashboard_onboarding_completion",
		map[string]interface{}{
			"start": "2020-02-01T00:00:00.000Z",
			"end":   "2020-01-01T00:00:00.000Z",
		})

	require.Contains(t, "Failed to get the data for period", response.Message)
}

func TestQueryResolver_Dashboard_OnboardingCompletion_SingleMonth_OneCompleted_OneNotCompleted(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	now := inCurrentMonthExceptFirstAndLastDays()
	hoursAgo2 := now.Add(-2 * time.Hour)
	in1Minute := now.Add(1 * time.Minute)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, now, map[string]string{"status": "SUCCESSFUL"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, hoursAgo2, map[string]string{"status": "LATE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, in1Minute, map[string]string{"status": "STUCK"})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1, "Organization": 1, "Action": 3})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_onboarding_completion",
		map[string]interface{}{
			"start": now,
			"end":   now,
		})

	var dashboardReport struct {
		Dashboard_OnboardingCompletion model.DashboardOnboardingCompletion
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, float64(50), dashboardReport.Dashboard_OnboardingCompletion.CompletionPercentage)
	require.Equal(t, float64(0), dashboardReport.Dashboard_OnboardingCompletion.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_OnboardingCompletion.PerMonth))
	require.Equal(t, 50.0, dashboardReport.Dashboard_OnboardingCompletion.PerMonth[0].Value)
}

func TestQueryResolver_Dashboard_OnboardingCompletion_PreviousMonthNoData(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	now := inCurrentMonthExceptFirstAndLastDays()
	hoursAgo4 := now.Add(-4 * time.Hour)
	monthAgo := now.Add(-30 * 24 * time.Hour)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, now, map[string]string{"status": "DONE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, hoursAgo4, map[string]string{"status": "LATE"})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1, "Organization": 1, "Action": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_onboarding_completion",
		map[string]interface{}{
			"start": monthAgo,
			"end":   now,
		})

	var dashboardReport struct {
		Dashboard_OnboardingCompletion model.DashboardOnboardingCompletion
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 100.0, dashboardReport.Dashboard_OnboardingCompletion.CompletionPercentage)
	require.Equal(t, 0.0, dashboardReport.Dashboard_OnboardingCompletion.IncreasePercentage)
	require.Equal(t, 2, len(dashboardReport.Dashboard_OnboardingCompletion.PerMonth))
	require.Equal(t, 0.0, dashboardReport.Dashboard_OnboardingCompletion.PerMonth[0].Value)
	require.Equal(t, 100.0, dashboardReport.Dashboard_OnboardingCompletion.PerMonth[1].Value)
}

func TestQueryResolver_Dashboard_OnboardingCompletion_PreviousMonthHasDataCurrentMonthNoData(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	now := inCurrentMonthExceptFirstAndLastDays()
	hoursAgo12 := now.Add(-12 * time.Hour)
	inAMonth := now.Add(30 * 24 * time.Hour)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, now, map[string]string{"status": "DONE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, hoursAgo12, map[string]string{"status": "LATE"})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1, "Organization": 1, "Action": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_onboarding_completion",
		map[string]interface{}{
			"start": now,
			"end":   inAMonth,
		})

	var dashboardReport struct {
		Dashboard_OnboardingCompletion model.DashboardOnboardingCompletion
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0.0, dashboardReport.Dashboard_OnboardingCompletion.CompletionPercentage)
	require.Equal(t, 0.0, dashboardReport.Dashboard_OnboardingCompletion.IncreasePercentage)
	require.Equal(t, 2, len(dashboardReport.Dashboard_OnboardingCompletion.PerMonth))
	require.Equal(t, 100.0, dashboardReport.Dashboard_OnboardingCompletion.PerMonth[0].Value)
	require.Equal(t, 0.0, dashboardReport.Dashboard_OnboardingCompletion.PerMonth[1].Value)
}

func TestQueryResolver_Dashboard_OnboardingCompletion_PercentageIncrease(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	now := inCurrentMonthExceptFirstAndLastDays()
	hoursAgo1 := now.Add(-1 * time.Hour)
	monthAgo := now.Add(-30 * 24 * time.Hour)
	monthAgoMinus1Hour := now.Add(-30 * 24 * time.Hour).Add(-1 * time.Hour)
	monthAgoPlus1Hour := now.Add(-30 * 24 * time.Hour).Add(1 * time.Hour)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, hoursAgo1, map[string]string{"status": "STUCK"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, now, map[string]string{"status": "DONE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, monthAgoMinus1Hour, map[string]string{"status": "LATE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, monthAgo, map[string]string{"status": "DONE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, monthAgoPlus1Hour, map[string]string{"status": "STUCK"})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1, "Organization": 1, "Action": 5})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_onboarding_completion",
		map[string]interface{}{
			"start": monthAgo,
			"end":   now,
		})

	var dashboardReport struct {
		Dashboard_OnboardingCompletion model.DashboardOnboardingCompletion
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 100.0, dashboardReport.Dashboard_OnboardingCompletion.CompletionPercentage)
	require.Equal(t, 100.0, dashboardReport.Dashboard_OnboardingCompletion.IncreasePercentage)
	require.Equal(t, 2, len(dashboardReport.Dashboard_OnboardingCompletion.PerMonth))
	require.Equal(t, 50.0, dashboardReport.Dashboard_OnboardingCompletion.PerMonth[0].Value)
	require.Equal(t, 100.0, dashboardReport.Dashboard_OnboardingCompletion.PerMonth[1].Value)
}

func TestQueryResolver_Dashboard_OnboardingCompletion_PercentageDecrease(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	now := inCurrentMonthExceptFirstAndLastDays()
	hoursAgo1 := now.Add(-1 * time.Hour)
	in1Hour := now.Add(1 * time.Hour)
	monthAgo := now.Add(-30 * 24 * time.Hour)
	monthAgoMinus1Hour := now.Add(-30 * 24 * time.Hour).Add(-1 * time.Hour)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, now, map[string]string{"status": "DONE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, hoursAgo1, map[string]string{"status": "LATE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, in1Hour, map[string]string{"status": "LATE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, monthAgo, map[string]string{"status": "DONE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, monthAgoMinus1Hour, map[string]string{"status": "LATE"})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1, "Organization": 1, "Action": 5})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_onboarding_completion",
		map[string]interface{}{
			"start": monthAgo,
			"end":   now,
		})

	var dashboardReport struct {
		Dashboard_OnboardingCompletion model.DashboardOnboardingCompletion
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 50.0, dashboardReport.Dashboard_OnboardingCompletion.CompletionPercentage)
	require.Equal(t, -50.0, dashboardReport.Dashboard_OnboardingCompletion.IncreasePercentage)
	require.Equal(t, 2, len(dashboardReport.Dashboard_OnboardingCompletion.PerMonth))
	require.Equal(t, 100.0, dashboardReport.Dashboard_OnboardingCompletion.PerMonth[0].Value)
	require.Equal(t, 50.0, dashboardReport.Dashboard_OnboardingCompletion.PerMonth[1].Value)
}

func TestQueryResolver_Dashboard_OnboardingCompletion_DoneIsFirstStatus(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	now := inCurrentMonthExceptFirstAndLastDays()
	monthAgo := now.Add(-30 * 24 * time.Hour)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, now, map[string]string{"status": "DONE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, monthAgo, map[string]string{"status": "DONE"})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1, "Organization": 1, "Action": 2})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_onboarding_completion",
		map[string]interface{}{
			"start": monthAgo,
			"end":   now,
		})

	var dashboardReport struct {
		Dashboard_OnboardingCompletion model.DashboardOnboardingCompletion
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0.0, dashboardReport.Dashboard_OnboardingCompletion.CompletionPercentage)
	require.Equal(t, 0.0, dashboardReport.Dashboard_OnboardingCompletion.IncreasePercentage)
	require.Equal(t, 2, len(dashboardReport.Dashboard_OnboardingCompletion.PerMonth))
	require.Equal(t, 100.0, dashboardReport.Dashboard_OnboardingCompletion.PerMonth[0].Value)
	require.Equal(t, 0.0, dashboardReport.Dashboard_OnboardingCompletion.PerMonth[1].Value)
}

func TestQueryResolver_Dashboard_OnboardingCompletion_MultipleOrgs(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	now := inCurrentMonthExceptFirstAndLastDays()
	hoursAgo4 := now.Add(-4 * time.Hour)
	hoursAgo8 := now.Add(-8 * time.Hour)
	in4Hours := now.Add(4 * time.Hour)
	orgId1 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	orgId2 := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId1, entity.ActionOnboardingStatusChanged, hoursAgo4, map[string]string{"status": "STUCK"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId1, entity.ActionOnboardingStatusChanged, now, map[string]string{"status": "DONE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId2, entity.ActionOnboardingStatusChanged, hoursAgo8, map[string]string{"status": "ON_TRACK"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId2, entity.ActionOnboardingStatusChanged, now, map[string]string{"status": "DONE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId2, entity.ActionOnboardingStatusChanged, in4Hours, map[string]string{"status": "NOT_STARTED"})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1, "Organization": 2, "Action": 5})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_onboarding_completion",
		map[string]interface{}{
			"start": now,
			"end":   now,
		})

	var dashboardReport struct {
		Dashboard_OnboardingCompletion model.DashboardOnboardingCompletion
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 1, len(dashboardReport.Dashboard_OnboardingCompletion.PerMonth))
	require.Equal(t, 67.0, dashboardReport.Dashboard_OnboardingCompletion.PerMonth[0].Value)
}

func TestQueryResolver_Dashboard_OnboardingCompletion_NoDoneOnboardings(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	now := inCurrentMonthExceptFirstAndLastDays()
	hoursAgo1 := now.Add(-1 * time.Hour)
	hoursAgo2 := now.Add(-2 * time.Hour)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, now, map[string]string{"status": "NOT_STARTED"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, hoursAgo1, map[string]string{"status": "LATE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, hoursAgo2, map[string]string{"status": "STUCK"})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1, "Organization": 1, "Action": 3})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_onboarding_completion",
		map[string]interface{}{
			"start": now,
			"end":   now,
		})

	var dashboardReport struct {
		Dashboard_OnboardingCompletion model.DashboardOnboardingCompletion
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 0.0, dashboardReport.Dashboard_OnboardingCompletion.CompletionPercentage)
	require.Equal(t, 0.0, dashboardReport.Dashboard_OnboardingCompletion.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_OnboardingCompletion.PerMonth))
	require.Equal(t, 0.0, dashboardReport.Dashboard_OnboardingCompletion.PerMonth[0].Value)
}

func TestQueryResolver_Dashboard_OnboardingCompletion_MultipleDonesInAMonth(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)

	now := inCurrentMonthExceptFirstAndLastDays()
	hoursAgo1 := now.Add(-1 * time.Hour)
	hoursAgo2 := now.Add(-2 * time.Hour)
	hoursAgo3 := now.Add(-3 * time.Hour)
	hoursAgo4 := now.Add(-4 * time.Hour)
	orgId := neo4jt.CreateOrg(ctx, driver, tenantName, entity.OrganizationEntity{})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, now, map[string]string{"status": "NOT_STARTED"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, hoursAgo1, map[string]string{"status": "SUCCESSFUL"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, hoursAgo2, map[string]string{"status": "LATE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, hoursAgo3, map[string]string{"status": "DONE"})
	neo4jt.CreateActionForOrganizationWithProperties(ctx, driver, tenantName, orgId, entity.ActionOnboardingStatusChanged, hoursAgo4, map[string]string{"status": "STUCK"})

	assertNeo4jNodeCount(ctx, t, driver, map[string]int{"Tenant": 1, "Organization": 1, "Action": 5})

	rawResponse := callGraphQL(t, "dashboard_view/dashboard_onboarding_completion",
		map[string]interface{}{
			"start": now,
			"end":   now,
		})

	var dashboardReport struct {
		Dashboard_OnboardingCompletion model.DashboardOnboardingCompletion
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &dashboardReport)
	require.Nil(t, err)

	require.Equal(t, 67.0, dashboardReport.Dashboard_OnboardingCompletion.CompletionPercentage)
	require.Equal(t, 0.0, dashboardReport.Dashboard_OnboardingCompletion.IncreasePercentage)
	require.Equal(t, 1, len(dashboardReport.Dashboard_OnboardingCompletion.PerMonth))
	require.Equal(t, 67.0, dashboardReport.Dashboard_OnboardingCompletion.PerMonth[0].Value)
}

func inCurrentMonthExceptFirstAndLastDays() time.Time {
	now := utils.Now()
	if now.Day() == 1 {
		return now.Add(24 * time.Hour)
	}
	if now.Day() >= 29 {
		return now.Add(-72 * time.Hour)
	}
	return now
}
