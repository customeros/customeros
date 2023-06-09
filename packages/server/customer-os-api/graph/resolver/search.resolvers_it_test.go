package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_GCliCache(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateCountryWith(ctx, driver, "1", "USA", "United States")
	neo4jt.CreateState(ctx, driver, "USA", "Alabama", "AL")
	neo4jt.CreateState(ctx, driver, "USA", "Louisiana", "LA")

	neo4jt.CreateContactWith(ctx, driver, tenantName, "a", "1")
	neo4jt.CreateContactWith(ctx, driver, tenantName, "ab", "2")
	neo4jt.CreateContactWith(ctx, driver, tenantName, "abc", "3")
	neo4jt.CreateContactWith(ctx, driver, tenantName, "abcd", "4")
	neo4jt.CreateContactWith(ctx, driver, tenantName, "b", "1")

	neo4jt.CreateOrganization(ctx, driver, tenantName, "a")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "ab")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "abc")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "abcd")
	neo4jt.CreateOrganization(ctx, driver, tenantName, "b")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Country"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "State"))
	require.Equal(t, 5, neo4jt.GetCountOfNodes(ctx, driver, "Contact"))
	require.Equal(t, 5, neo4jt.GetCountOfNodes(ctx, driver, "Organization"))

	rawResponse, err := c.RawPost(getQuery("search/gcli_cache"))
	assertRawResponseSuccess(t, rawResponse, err)

	gcliCacheResult := rawResponse.Data.(map[string]interface{})["gcli_Cache"]
	require.NotNil(t, gcliCacheResult)
	require.Equal(t, 10, len(gcliCacheResult.([]interface{})))

	require.Equal(t, "STATE", gcliCacheResult.([]interface{})[0].(map[string]interface{})["type"])
	require.Equal(t, "STATE", gcliCacheResult.([]interface{})[1].(map[string]interface{})["type"])
	require.Equal(t, "CONTACT", gcliCacheResult.([]interface{})[2].(map[string]interface{})["type"])
	require.Equal(t, "CONTACT", gcliCacheResult.([]interface{})[3].(map[string]interface{})["type"])
	require.Equal(t, "CONTACT", gcliCacheResult.([]interface{})[4].(map[string]interface{})["type"])
	require.Equal(t, "CONTACT", gcliCacheResult.([]interface{})[5].(map[string]interface{})["type"])
	require.Equal(t, "ORGANIZATION", gcliCacheResult.([]interface{})[6].(map[string]interface{})["type"])
	require.Equal(t, "ORGANIZATION", gcliCacheResult.([]interface{})[7].(map[string]interface{})["type"])
	require.Equal(t, "ORGANIZATION", gcliCacheResult.([]interface{})[8].(map[string]interface{})["type"])
	require.Equal(t, "ORGANIZATION", gcliCacheResult.([]interface{})[9].(map[string]interface{})["type"])
}

func TestQueryResolver_GCliSearch(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateFullTextBasicSearchIndexes(ctx, driver, tenantName)

	neo4jt.CreateCountry(ctx, driver, "USA", "United States")

	neo4jt.CreateContactWith(ctx, driver, tenantName, "c", "1")

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

	require.Equal(t, "STATE", gcliSearchResult.([]interface{})[0].(map[string]interface{})["type"])
	require.Equal(t, "Alabama", gcliSearchResult.([]interface{})[0].(map[string]interface{})["display"])
}
