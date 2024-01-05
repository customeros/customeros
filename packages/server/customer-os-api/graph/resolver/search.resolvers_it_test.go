package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_GCliSearch(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateFullTextBasicSearchIndexes(ctx, driver, tenantName)

	neo4jt.CreateCountry(ctx, driver, "US", "USA", "United States", "1")

	neo4jt.CreateContactWith(ctx, driver, tenantName, "c", "1")

	neo4jt.CreateState(ctx, driver, "USA", "Alabama", "AL")
	neo4jt.CreateState(ctx, driver, "USA", "Louisiana", "LA")

	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Country"))
	require.Equal(t, 2, neo4jtest.GetCountOfNodes(ctx, driver, "State"))

	rawResponse, err := c.RawPost(getQuery("search/gcli_search"),
		client.Var("keyword", "AL"),
		client.Var("limit", "1"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	gcliSearchResult := rawResponse.Data.(map[string]interface{})["gcli_Search"]
	require.NotNil(t, gcliSearchResult)
	require.Equal(t, 1, len(gcliSearchResult.([]interface{})))

	require.Equal(t, "STATE", gcliSearchResult.([]interface{})[0].(map[string]interface{})["type"])
	require.Equal(t, "Alabama", gcliSearchResult.([]interface{})[0].(map[string]interface{})["display"])
}
