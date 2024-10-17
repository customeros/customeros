package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_EntityTemplates_FilterExtendsProperty(t *testing.T) {
	ctx := context.Background()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)

	neo4jt.CreateEntityTemplate(ctx, driver, tenantName, "")
	id2 := neo4jt.CreateEntityTemplate(ctx, driver, tenantName, model.EntityTemplateExtensionContact.String())
	id3 := neo4jt.CreateEntityTemplate(ctx, driver, tenantName, model.EntityTemplateExtensionContact.String())

	rawResponse, err := c.RawPost(getQuery("get_entity_templates_filter_by_extends"),
		client.Var("extends", model.EntityTemplateExtensionContact.String()))
	assertRawResponseSuccess(t, rawResponse, err)

	var entityTemplate struct {
		EntityTemplates []model.EntityTemplate
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &entityTemplate)
	require.Nil(t, err)
	require.NotNil(t, entityTemplate.EntityTemplates)
	require.Equal(t, 2, len(entityTemplate.EntityTemplates))
	require.Equal(t, "CONTACT", entityTemplate.EntityTemplates[0].Extends.String())
	require.Equal(t, "CONTACT", entityTemplate.EntityTemplates[1].Extends.String())
	require.ElementsMatch(t, []string{id2, id3}, []string{entityTemplate.EntityTemplates[0].ID, entityTemplate.EntityTemplates[1].ID})

	require.Equal(t, 3, neo4jtest.GetCountOfNodes(ctx, driver, "EntityTemplate"))
}
