package resolver

import (
	"github.com/99designs/gqlgen/client"
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

	rawResponse, err := c.RawPost(getQuery("tag/create_tag"))
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

func TestMutationResolver_TagUpdate(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	tagId := neo4jt.CreateTag(driver, tenantName, "original tag")

	rawResponse, err := c.RawPost(getQuery("tag/update_tag"),
		client.Var("tagId", tagId),
		client.Var("tagName", "new tag name"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var tag struct {
		Tag_Update model.Tag
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tag)
	updatedTag := tag.Tag_Update
	require.Nil(t, err)
	require.NotNil(t, updatedTag)
	require.NotNil(t, updatedTag.UpdatedAt)
	require.NotEqual(t, updatedTag.UpdatedAt, updatedTag.CreatedAt)
	require.Equal(t, tagId, updatedTag.ID)
	require.Equal(t, "new tag name", updatedTag.Name)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tag"))
}

func TestMutationResolver_TagDelete(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	tagId := neo4jt.CreateTag(driver, tenantName, "original tag")
	contactId := neo4jt.CreateDefaultContact(driver, tenantName)
	neo4jt.TagContact(driver, contactId, tagId)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Tag"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "TAGGED"))
	require.Equal(t, 1, neo4jt.GetCountOfRelationships(driver, "TAG_BELONGS_TO_TENANT"))

	rawResponse, err := c.RawPost(getQuery("tag/delete_tag"),
		client.Var("tagId", tagId),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var result struct {
		Tag_Delete model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &result)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, true, result.Tag_Delete.Result)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(driver, "Contact"))
	require.Equal(t, 0, neo4jt.GetCountOfNodes(driver, "Tag"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(driver, "TAGGED"))
	require.Equal(t, 0, neo4jt.GetCountOfRelationships(driver, "TAG_BELONGS_TO_TENANT"))
}

func TestQueryResolver_Tags(t *testing.T) {
	defer tearDownTestCase()(t)
	neo4jt.CreateTenant(driver, tenantName)
	neo4jt.CreateTenant(driver, "other")
	tagId1 := neo4jt.CreateTag(driver, tenantName, "tag B")
	tagId2 := neo4jt.CreateTag(driver, tenantName, "tag A")
	neo4jt.CreateTag(driver, "other", "contact type for other tenant")

	require.Equal(t, 3, neo4jt.GetCountOfNodes(driver, "Tag"))

	rawResponse, err := c.RawPost(getQuery("tag/get_tags"))
	assertRawResponseSuccess(t, rawResponse, err)

	var tagStruct struct {
		Tags []model.Tag
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tagStruct)
	tags := tagStruct.Tags
	require.Nil(t, err)
	require.Equal(t, 2, len(tags))
	require.Equal(t, tagId2, tags[0].ID)
	require.Equal(t, "tag A", tags[0].Name)
	require.Equal(t, tagId1, tags[1].ID)
	require.Equal(t, "tag B", tags[1].Name)
}
