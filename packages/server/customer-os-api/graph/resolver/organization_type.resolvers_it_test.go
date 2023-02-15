package resolver

import (
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"testing"
)

func TestQueryResolver_OrganizationTypes(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, "other")
	organizationTypeId1 := neo4jt.CreateOrganizationType(ctx, driver, tenantName, "first")
	organizationTypeId2 := neo4jt.CreateOrganizationType(ctx, driver, tenantName, "second")
	neo4jt.CreateOrganizationType(ctx, driver, "other", "organization type for other tenant")

	require.Equal(t, 3, neo4jt.GetCountOfNodes(ctx, driver, "OrganizationType"))

	rawResponse, err := c.RawPost(getQuery("get_organization_types"))
	assertRawResponseSuccess(t, rawResponse, err)

	var organizationType struct {
		OrganizationTypes []model.OrganizationType
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationType)
	require.Nil(t, err)
	require.NotNil(t, organizationType)
	require.Equal(t, 2, len(organizationType.OrganizationTypes))
	require.Equal(t, organizationTypeId1, organizationType.OrganizationTypes[0].ID)
	require.Equal(t, "first", organizationType.OrganizationTypes[0].Name)
	require.Equal(t, organizationTypeId2, organizationType.OrganizationTypes[1].ID)
	require.Equal(t, "second", organizationType.OrganizationTypes[1].Name)
}

func TestMutationResolver_OrganizationTypeCreate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, "otherTenantName")

	rawResponse, err := c.RawPost(getQuery("create_organization_type"))
	assertRawResponseSuccess(t, rawResponse, err)

	var organizationType struct {
		OrganizationType_Create model.OrganizationType
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationType)
	require.Nil(t, err)
	require.NotNil(t, organizationType)
	require.NotNil(t, organizationType.OrganizationType_Create.ID)
	require.Equal(t, "the organization type", organizationType.OrganizationType_Create.Name)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "OrganizationType"))
	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "OrganizationType_"+tenantName))

	assertNeo4jLabels(ctx, t, driver, []string{"Tenant", "OrganizationType", "OrganizationType_" + tenantName})
}

func TestMutationResolver_OrganizationTypeUpdate(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	organizationTypeId := neo4jt.CreateOrganizationType(ctx, driver, tenantName, "original type")

	rawResponse, err := c.RawPost(getQuery("update_organization_type"),
		client.Var("organizationTypeId", organizationTypeId))
	assertRawResponseSuccess(t, rawResponse, err)

	var organizationType struct {
		OrganizationType_Update model.OrganizationType
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &organizationType)
	require.Nil(t, err)
	require.NotNil(t, organizationType)
	require.Equal(t, organizationTypeId, organizationType.OrganizationType_Update.ID)
	require.Equal(t, "updated type", organizationType.OrganizationType_Update.Name)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "OrganizationType"))
}

func TestMutationResolver_OrganizationTypeDelete(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	organizationTypeId := neo4jt.CreateOrganizationType(ctx, driver, tenantName, "the type")

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "OrganizationType"))

	rawResponse, err := c.RawPost(getQuery("delete_organization_type"),
		client.Var("organizationTypeId", organizationTypeId))
	assertRawResponseSuccess(t, rawResponse, err)

	var result struct {
		OrganizationType_Delete model.Result
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &result)
	require.Nil(t, err)
	require.NotNil(t, result)
	require.Equal(t, true, result.OrganizationType_Delete.Result)

	require.Equal(t, 0, neo4jt.GetCountOfNodes(ctx, driver, "OrganizationType"))
}
