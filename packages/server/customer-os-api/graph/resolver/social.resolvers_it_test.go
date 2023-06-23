package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_SocialUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	neo4jt.CreateTenant(ctx, driver, tenantName)
	socialId := neo4jt.CreateSocial(ctx, driver, tenantName, entity.SocialEntity{})

	rawResponse := callGraphQL(t, "social/update_social", map[string]interface{}{"socialId": socialId})

	var socialStruct struct {
		Social_Update model.Social
	}

	err := decode.Decode(rawResponse.Data.(map[string]any), &socialStruct)
	require.Nil(t, err)

	updatedSocial := socialStruct.Social_Update

	require.Equal(t, socialId, updatedSocial.ID)
	test.AssertRecentTime(t, updatedSocial.UpdatedAt)
	require.Equal(t, model.DataSourceOpenline, updatedSocial.SourceOfTruth)
	require.Equal(t, "new name", *updatedSocial.PlatformName)
	require.Equal(t, "new url", updatedSocial.URL)

	// Check the number of nodes in the Neo4j database
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Social"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Social_"+tenantName))
}
