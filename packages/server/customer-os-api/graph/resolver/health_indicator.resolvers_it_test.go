package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_HealthIndicators(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	healthIndicator1 := neo4jt.CreateHealthIndicator(ctx, driver, tenantName, "Green", 10)
	healthIndicator2 := neo4jt.CreateHealthIndicator(ctx, driver, tenantName, "Red", 20)

	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "HealthIndicator"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "HealthIndicator_"+tenantName))

	rawResponse := callGraphQL(t, "health_indicator/get_health_indicators", map[string]interface{}{})

	var healthIndicatorStruct struct {
		HealthIndicators []model.HealthIndicator
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &healthIndicatorStruct)
	healthIndicators := healthIndicatorStruct.HealthIndicators
	require.Nil(t, err)
	require.Equal(t, 2, len(healthIndicators))
	require.Equal(t, healthIndicator1, healthIndicators[0].ID)
	require.Equal(t, "Green", healthIndicators[0].Name)
	require.Equal(t, int64(10), healthIndicators[0].Order)
	require.Equal(t, healthIndicator2, healthIndicators[1].ID)
	require.Equal(t, "Red", healthIndicators[1].Name)
	require.Equal(t, int64(20), healthIndicators[1].Order)
}
