package resolver

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMutationResolver_TagCreate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateTenant(driver, "otherTenant")

	rawResponse, err := c.RawPost(getQuery("create_tag"))
	assertRawResponseSuccess(t, rawResponse, err)

	var tag struct {
		Tag_Create model.Tag
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tag)

	createdTag := tag.Tag_Create
	require.Nil(t, err)
	require.NotNil(t, createdTag)
	require.NotNil(t, createdTag.ID)
	require.NotNil(t, createdTag.CreatedAt)
	require.NotNil(t, createdTag.UpdatedAt)
	require.Equal(t, "the tag", createdTag.Name)
	require.Equal(t, "test", createdTag.AppSource)
	require.Equal(t, model.DataSourceOpenline, createdTag.Source)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tag"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tag_"+tenantName))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "TAG_BELONGS_TO_TENANT"))

	assertNeo4jLabels(t, driver, []string{"Tenant", "Tag", "Tag_" + tenantName})
}
