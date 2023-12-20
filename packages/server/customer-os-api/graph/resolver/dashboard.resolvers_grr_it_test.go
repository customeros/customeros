package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_Dashboard_GRR_1_Contract_1_SLI_V1(t *testing.T) {
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
