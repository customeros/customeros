package resolver

import (
	"github.com/99designs/gqlgen/client"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestQueryResolver_GCliSearch(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateFullTextBasicSearchIndexes(ctx, driver, tenantName)

	neo4jt.CreateCountry(ctx, driver, "USA", "United States")

	neo4jt.CreateState(ctx, driver, "USA", "Alabama", "AL")
	neo4jt.CreateState(ctx, driver, "USA", "Louisiana", "LA")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Country"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "State"))

	rawResponse, err := c.RawPost(getQuery("search/gcli_search"),
		client.Var("keyword", "AL"),
		client.Var("limit", "1"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	gcliSearchResult := rawResponse.Data.(map[string]interface{})["gcli_Search"]
	require.NotNil(t, gcliSearchResult)
	require.Equal(t, 1, len(gcliSearchResult.([]interface{})))
	require.NotNil(t, gcliSearchResult.([]interface{})[0].(map[string]interface{})["score"])

	organization := gcliSearchResult.([]interface{})[0].(map[string]interface{})["result"]
	require.Equal(t, "State", organization.(map[string]interface{})["__typename"])
	require.Equal(t, "AL", organization.(map[string]interface{})["code"])
	require.Equal(t, "Alabama", organization.(map[string]interface{})["name"])
}
