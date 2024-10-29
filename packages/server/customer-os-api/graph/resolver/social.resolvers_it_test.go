package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_SocialUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	socialId := neo4jt.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{})

	rawResponse := callGraphQL(t, "social/update_social", map[string]interface{}{"socialId": socialId})

	var socialStruct struct {
		Social_Update model.Social
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &socialStruct)
	require.Nil(t, err)

	updatedSocial := socialStruct.Social_Update

	require.Equal(t, socialId, updatedSocial.ID)
	test.AssertRecentTime(t, updatedSocial.UpdatedAt)
	require.Equal(t, "new url", updatedSocial.URL)

	// Check the number of nodes in the Neo4j database
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Social"))
	require.Equal(t, 1, neo4jtest.GetCountOfNodes(ctx, driver, "Social_"+tenantName))
}

func TestMutationResolver_SocialRemove(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jtest.CreateTenant(ctx, driver, tenantName)
	socialId := neo4jt.CreateSocial(ctx, driver, tenantName, neo4jentity.SocialEntity{})

	rawResponse := callGraphQL(t, "social/remove_social", map[string]interface{}{"socialId": socialId})

	var resultStruct struct {
		Social_Remove model.Result
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &resultStruct)
	require.Nil(t, err)

	require.True(t, resultStruct.Social_Remove.Result)

	// Check the number of nodes in the Neo4j database
	require.Equal(t, 0, neo4jtest.GetCountOfNodes(ctx, driver, "Social"))
	require.Equal(t, 0, neo4jtest.GetCountOfNodes(ctx, driver, "Social_"+tenantName))
}
