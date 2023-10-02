package unit_tests

import (
	"fmt"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/resolver"
	cosHandler "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handler"

	//cosHandler "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handler"

	//cosHandler "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handler"
	srv "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

type WorkspaceInput struct {
	name      string
	provider  string
	appSource string
}

func TestMutationResolver_WorkspaceMergeToTenant_FullyPopulated(t *testing.T) {
	workspaceInput := WorkspaceInput{
		name:      "test name",
		provider:  "test provider",
		appSource: "test app",
	}
	t.Run("should create fully populated workspace correctly", func(t *testing.T) {
		testWorkspaceService := new(MockedWorkspaceService)
		mockedServices := srv.Services{
			WorkspaceService: testWorkspaceService,
		}
		resolvers := resolver.Resolver{Services: &mockedServices}
		q := fmt.Sprintf(`
		 mutation {
		   workspace_MergeToTenant(workspace: {name: "%s", provider: "%s", appSource: "%s"}, tenant: "test") {
		result
		   }
		 }
		`, workspaceInput.name, workspaceInput.provider, workspaceInput.appSource)
		schemaConfig := generated.Config{Resolvers: &resolvers}
		schemaConfig.Directives.HasRole = cosHandler.GetRoleChecker()
		schemaConfig.Directives.HasTenant = cosHandler.GetTenantChecker()
		c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(schemaConfig)))
		rawResponse, err := c.RawPost(q)
		require.Nil(t, err)

		var workspaceStruct struct {
			Workspace_Create model.Workspace
		}

		err = decode.Decode(rawResponse.Data.(map[string]any), &workspaceStruct)
		require.Nil(t, err)
		require.NotNil(t, workspaceStruct)

		workspace := workspaceStruct.Workspace_Create

		require.Equal(t, "", workspace.ID)
		require.Equal(t, workspaceInput.name, workspace.Name)
		require.Equal(t, workspaceInput.provider, workspace.Provider)
		require.Equal(t, workspaceInput.appSource, workspace.AppSource)
	})
}

//func GetRoleChecker() func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []model.Role) (res interface{}, err error) {
//	return func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []model.Role) (res interface{}, err error) {
//		common.CustomContext{Roles: roles}
//		currentRoles := common.GetRolesFromContext(ctx)
//		// Check if the current role is in the list of allowed roles
//		for _, allowedRole := range roles {
//			for _, currentRole := range currentRoles {
//				if currentRole == allowedRole {
//					// If the role is in the list of allowed roles, call the next resolver
//					return next(ctx)
//				}
//			}
//		}
//		// If the role is not in the list of allowed roles, return an error
//		return nil, errors.ErrAccessDenied
//	}
//}
