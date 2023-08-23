package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"log"
	"strings"
	"testing"
)

func TestMutationResolver_TenantMerge(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, "other")

	rawResponse, err := cAdmin.RawPost(getQuery("tenant/merge_tenant"),
		client.Var("name", "testtenant"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantResponse struct {
		Tenant_Merge *string `json:"tenant_Merge"`
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tenantResponse)
	require.Nil(t, err)
	require.NotNil(t, tenantResponse)

	log.Println(rawResponse.Data)
	require.NotNil(t, tenantResponse.Tenant_Merge)
	require.Equal(t, "testtenant", *tenantResponse.Tenant_Merge)

	rawResponse2, err := cAdmin.RawPost(getQuery("tenant/merge_tenant"),
		client.Var("name", "testtenant"),
	)
	assertRawResponseSuccess(t, rawResponse2, err)

	var tenantResponse2 struct {
		Tenant_Merge *string `json:"tenant_Merge"`
	}

	err = decode.Decode(rawResponse2.Data.(map[string]any), &tenantResponse2)
	require.Nil(t, err)
	require.NotNil(t, tenantResponse2)

	log.Println(rawResponse2.Data)
	require.NotNil(t, tenantResponse2.Tenant_Merge)
	require.NotEqual(t, "testtenant", *tenantResponse2.Tenant_Merge)
	require.True(t, strings.HasPrefix(*tenantResponse2.Tenant_Merge, "testtenant"))

}

func TestMutationResolver_TenantMerge_AccessControlled(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	rawResponse, err := c.RawPost(getQuery("tenant/merge_tenant"),
		client.Var("name", "testtenant"),
	)

	require.Nil(t, err)
	require.NotNil(t, rawResponse.Errors)

}

func TestMutationResolver_TenantMerge_CheckDefaultData(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateOrganizationRelationship(ctx, driver, "R1")
	neo4jt.CreateOrganizationRelationship(ctx, driver, "R2")

	newTenantName := "test_tenant"
	rawResponse, err := cAdmin.RawPost(getQuery("tenant/merge_tenant"),
		client.Var("name", newTenantName),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	require.Equal(t, 1, neo4jt.GetCountOfNodes(ctx, driver, "Tenant"))
	require.Equal(t, 2, neo4jt.GetCountOfNodes(ctx, driver, "OrganizationRelationship"))
	require.Equal(t, 16, neo4jt.GetCountOfNodes(ctx, driver, "OrganizationRelationshipStage"))
	require.Equal(t, 16, neo4jt.GetCountOfNodes(ctx, driver, "OrganizationRelationshipStage_"+newTenantName))
	require.Equal(t, 16, neo4jt.GetCountOfRelationships(ctx, driver, "HAS_STAGE"))
	require.Equal(t, 16, neo4jt.GetCountOfRelationships(ctx, driver, "STAGE_BELONGS_TO_TENANT"))
}

func TestMutationResolver_GetByWorkspace(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jt.CreateTenant(ctx, driver, tenantName)
	neo4jt.CreateTenant(ctx, driver, "other")
	neo4jt.CreateWorkspace(ctx, driver, "testworkspace", "testprovider", tenantName)

	rawResponse, err := cAdmin.RawPost(getQuery("tenant/get_by_workspace"),
		client.Var("name", "testworkspace"),
		client.Var("provider", "testprovider"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantResponse struct {
		Tenant_ByWorkspace *string
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tenantResponse)
	require.Nil(t, err)
	require.NotNil(t, tenantResponse)

	log.Println(rawResponse.Data)
	require.NotNil(t, tenantResponse.Tenant_ByWorkspace)
	require.Equal(t, tenantName, *tenantResponse.Tenant_ByWorkspace)

	rawResponse2, err := cAdmin.RawPost(getQuery("tenant/get_by_workspace"),
		client.Var("name", "testworkspace2"),
		client.Var("provider", "testprovider"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var tenantResponse2 struct {
		Tenant_ByWorkspace *string
	}
	log.Println(rawResponse2.Data)

	err = decode.Decode(rawResponse2.Data.(map[string]any), &tenantResponse2)
	require.Nil(t, err)
	require.NotNil(t, tenantResponse2)
	require.Nil(t, tenantResponse2.Tenant_ByWorkspace)

}
