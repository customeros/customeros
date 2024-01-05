package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	neo4jt "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/test/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestMutationResolver_WorkspaceMergeToTenant(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenant(ctx, driver, "other")
	neo4jt.CreateWorkspace(ctx, driver, "testworkspace2", "testprovider", "other")

	rawResponse, err := cAdmin.RawPost(getQuery("workspace/merge_to_tenant"),
		client.Var("name", "testworkspace"),
		client.Var("tenant", tenantName),
		client.Var("provider", "testprovider"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var workspaceResponse struct {
		Workspace_MergeToTenant struct {
			Result bool
		}
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &workspaceResponse)
	require.Nil(t, err)
	require.NotNil(t, workspaceResponse)

	//verify the mapping happened
	rawResponse3, err := cAdmin.RawPost(getQuery("tenant/get_by_workspace"),
		client.Var("name", "testworkspace"),
		client.Var("provider", "testprovider"),
	)
	assertRawResponseSuccess(t, rawResponse3, err)
	log.Println(rawResponse3.Data)

	require.True(t, workspaceResponse.Workspace_MergeToTenant.Result)

	var tenantResponse struct {
		Tenant_ByWorkspace *string
	}

	err = decode.Decode(rawResponse3.Data.(map[string]any), &tenantResponse)
	require.Nil(t, err)
	require.NotNil(t, tenantResponse)
	require.Equal(t, tenantName, *tenantResponse.Tenant_ByWorkspace)

	// testworkspace2 is already mapped to a workspace, this operation should return false
	rawResponse2, err := cAdmin.RawPost(getQuery("workspace/merge_to_tenant"),
		client.Var("name", "testworkspace2"),
		client.Var("tenant", tenantName),
		client.Var("provider", "testprovider"),
	)
	assertRawResponseSuccess(t, rawResponse2, err)

	var workspaceResponse2 struct {
		Workspace_MergeToTenant struct {
			Result bool
		}
	}

	err = decode.Decode(rawResponse2.Data.(map[string]any), &workspaceResponse2)
	require.Nil(t, err)
	require.NotNil(t, workspaceResponse2)

	log.Println(rawResponse2.Data)
	require.False(t, workspaceResponse2.Workspace_MergeToTenant.Result)

	//verify the mapping didn't get modified
	rawResponse4, err := cAdmin.RawPost(getQuery("tenant/get_by_workspace"),
		client.Var("name", "testworkspace2"),
		client.Var("provider", "testprovider"),
	)
	assertRawResponseSuccess(t, rawResponse4, err)
	log.Println(rawResponse4.Data)

	require.True(t, workspaceResponse.Workspace_MergeToTenant.Result)

	var tenantResponse2 struct {
		Tenant_ByWorkspace *string
	}

	err = decode.Decode(rawResponse4.Data.(map[string]any), &tenantResponse2)
	require.Nil(t, err)
	require.NotNil(t, tenantResponse2)
	require.Equal(t, "other", *tenantResponse2.Tenant_ByWorkspace)

}

func TestMutationResolver_WorkspaceMerge(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenant(ctx, driver, "other")
	neo4jt.CreateWorkspace(ctx, driver, "testworkspace2", "testprovider", "other")

	rawResponse, err := cOwner.RawPost(getQuery("workspace/merge"),
		client.Var("name", "testworkspace"),
		client.Var("provider", "testprovider"),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var workspaceResponse struct {
		Workspace_Merge struct {
			Result bool
		}
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &workspaceResponse)
	require.Nil(t, err)
	require.NotNil(t, workspaceResponse)
	require.True(t, workspaceResponse.Workspace_Merge.Result)

	//verify the mapping happened
	rawResponse3, err := cAdmin.RawPost(getQuery("tenant/get_by_workspace"),
		client.Var("name", "testworkspace"),
		client.Var("provider", "testprovider"),
	)
	assertRawResponseSuccess(t, rawResponse3, err)
	log.Println(rawResponse3.Data)

	var tenantResponse struct {
		Tenant_ByWorkspace *string
	}

	err = decode.Decode(rawResponse3.Data.(map[string]any), &tenantResponse)
	require.Nil(t, err)
	require.NotNil(t, tenantResponse)
	require.Equal(t, tenantName, *tenantResponse.Tenant_ByWorkspace)

	// testworkspace2 is already mapped to a workspace, this operation should return false
	rawResponse2, err := cOwner.RawPost(getQuery("workspace/merge"),
		client.Var("name", "testworkspace2"),
		client.Var("provider", "testprovider"),
	)
	assertRawResponseSuccess(t, rawResponse2, err)

	var workspaceResponse2 struct {
		Workspace_MergeToTenant struct {
			Result bool
		}
	}

	err = decode.Decode(rawResponse2.Data.(map[string]any), &workspaceResponse2)
	require.Nil(t, err)
	require.NotNil(t, workspaceResponse2)

	log.Println(rawResponse2.Data)
	require.False(t, workspaceResponse2.Workspace_MergeToTenant.Result)

	//verify the mapping didn't get modified
	rawResponse4, err := cAdmin.RawPost(getQuery("tenant/get_by_workspace"),
		client.Var("name", "testworkspace2"),
		client.Var("provider", "testprovider"),
	)
	assertRawResponseSuccess(t, rawResponse4, err)
	log.Println(rawResponse4.Data)

	var tenantResponse2 struct {
		Tenant_ByWorkspace *string
	}

	err = decode.Decode(rawResponse4.Data.(map[string]any), &tenantResponse2)
	require.Nil(t, err)
	require.NotNil(t, tenantResponse2)
	require.Equal(t, "other", *tenantResponse2.Tenant_ByWorkspace)

}

func TestMutationResolver_WorkspaceMergeToTenantAccessControlled(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)
	neo4jtest.CreateTenant(ctx, driver, tenantName)
	neo4jtest.CreateTenant(ctx, driver, "other")
	neo4jt.CreateWorkspace(ctx, driver, "testworkspace2", "testprovider", "other")

	rawResponse, err := c.RawPost(getQuery("workspace/merge_to_tenant"),
		client.Var("name", "testworkspace"),
		client.Var("tenant", tenantName),
		client.Var("provider", "testprovider"),
	)

	require.Nil(t, err)
	require.NotNil(t, rawResponse.Errors)

}
