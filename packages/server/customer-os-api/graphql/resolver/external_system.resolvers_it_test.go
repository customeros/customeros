package resolver

import (
	"context"
	"testing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
)

func TestQueryResolver_ExternalSystemInstances(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateExternalSystem(ctx, driver, tenantName, neo4jentity.ExternalSystemEntity{
		ExternalSystemId: enum.SourceStripe,
		Stripe: struct {
			PaymentMethodTypes []string
		}{
			PaymentMethodTypes: []string{"card", "ideal"},
		},
	})
	neo4jtest.CreateExternalSystem(ctx, driver, tenantName, neo4jentity.ExternalSystemEntity{
		ExternalSystemId: enum.SourceHubspot,
	})

	rawResponse, err := c.RawPost(getQuery("external_system/get_external_system_instances"))
	assertRawResponseSuccess(t, rawResponse, err)

	var graphqlResponse struct {
		ExternalSystemInstances []model.ExternalSystemInstance
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &graphqlResponse)
	require.Nil(t, err)
	require.NotNil(t, graphqlResponse)

	require.Equal(t, 2, len(graphqlResponse.ExternalSystemInstances))
	for _, instance := range graphqlResponse.ExternalSystemInstances {
		switch instance.Type {
		case model.ExternalSystemTypeStripe:
			require.Equal(t, []string{"card", "ideal"}, instance.StripeDetails.PaymentMethodTypes)
		case model.ExternalSystemTypeHubspot:
			require.Nil(t, instance.StripeDetails)
		}
	}
}
