package service

import (
	"errors"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"

	"context"

	"time"
)

type CustomerOsClient interface {
	GetTenantByWorkspace(workspace *model.WorkspaceInput) (*string, error)
	MergeTenantToWorkspace(workspace *model.WorkspaceInput, tenant string) (bool, error)
	CreateUser(user *model.UserInput, tenant string, roles []Role) (string, error)
	MergeTenant(tenant *model.TenantInput) (string, error)
	IsPlayer(authId string, provider string) (string, error)
}

type customerOsClient struct {
	cfg           *config.Config
	graphqlClient *graphql.Client
}

type Role string

const (
	ROLE_OWNER Role = "OWNER"
	ROLE_USER  Role = "USER"
)

func NewCustomerOsClient(cfg *config.Config, graphqlClient *graphql.Client) CustomerOsClient {
	return &customerOsClient{
		cfg:           cfg,
		graphqlClient: graphqlClient,
	}
}

func (s *customerOsClient) GetTenantByWorkspace(workspace *model.WorkspaceInput) (*string, error) {
	if workspace == nil {
		return nil, errors.New("GetTenantByWorkspace: workspace is nil")
	}
	graphqlRequest := graphql.NewRequest(
		`
		query GetTenantByWorkspace ($name: String!, $provider: String!) {
				tenant_ByWorkspace(workspace: {
			  name: $name,
			  provider: $provider
			}) 
		}
	`)
	graphqlRequest.Var("name", workspace.Name)
	graphqlRequest.Var("provider", workspace.Provider)

	err := s.addHeadersToGraphRequest(graphqlRequest, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, err
	}
	defer cancel()

	var graphqlResponse model.GetTenantByWorkspaceResponse
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return nil, err
	}
	return graphqlResponse.Workspace, nil
}

func (s *customerOsClient) MergeTenantToWorkspace(workspace *model.WorkspaceInput, tenant string) (bool, error) {
	if workspace == nil {
		return false, errors.New("MergeTenantToWorkspace: workspace is nil")
	}
	graphqlRequest := graphql.NewRequest(
		`
			mutation MergeWorkspaceToTenant($tenant: String!, $name: String!, $provider: String!, $appSource: String) {
			   workspace_MergeToTenant(tenant: $tenant, 
										workspace: {name: $name,
										provider: $provider,
										appSource: $appSource}) {
				result
			  }
			}
	`)
	graphqlRequest.Var("tenant", tenant)
	graphqlRequest.Var("name", workspace.Name)
	graphqlRequest.Var("provider", workspace.Provider)
	graphqlRequest.Var("appSource", workspace.AppSource)

	err := s.addHeadersToGraphRequest(graphqlRequest, nil)
	if err != nil {
		return false, err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return false, err
	}
	defer cancel()

	var graphqlResponse model.MergeTenantToWorkspaceResponse
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return false, err
	}
	return graphqlResponse.Workspace.Result, nil
}

func (s *customerOsClient) CreateUser(user *model.UserInput, tenant string, roles []Role) (string, error) {
	if user == nil {
		return "", errors.New("CreateUser: user is nil")
	}
	graphqlRequest := graphql.NewRequest(
		`
			mutation AddUser($user: UserInput!) {
			   user_Create(
					input: $user) {
				id,
				firstName,
				lastName,
			  }
			}
	`)
	graphqlRequest.Var("tenant", tenant)
	graphqlRequest.Var("user", *user)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant)
	if err != nil {
		return "", err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", err
	}
	defer cancel()

	var graphqlResponse model.CreateUserResponse
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", err
	}

	for _, role := range roles {
		_, err = s.AddRole(graphqlResponse.User.ID, tenant, role)
		if err != nil {
			return "", err
		}
	}
	return graphqlResponse.User.ID, nil
}

func (s *customerOsClient) AddRole(userId string, tenant string, role Role) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`
			mutation AddRole($userId: ID!, $role: Role!) {
			   user_AddRole(
					id: $userId,
					role: $role) {
				id
			  }
			}
	`)
	graphqlRequest.Var("userId", userId)
	graphqlRequest.Var("role", role)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant)
	if err != nil {
		return "", err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", err
	}
	defer cancel()

	var graphqlResponse model.AddRoleResponse
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", err
	}
	return graphqlResponse.Role.ID, nil
}

func (s *customerOsClient) MergeTenant(tenant *model.TenantInput) (string, error) {
	if tenant == nil {
		return "", errors.New("MergeTenant: tenant is nil")
	}
	graphqlRequest := graphql.NewRequest(
		`
			mutation CreateTenant($tenant: TenantInput!) {
			   tenant_Merge(
					tenant: $tenant) 
			}
	`)
	graphqlRequest.Var("tenant", *tenant)

	err := s.addHeadersToGraphRequest(graphqlRequest, nil)
	if err != nil {
		return "", err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", err
	}
	defer cancel()

	var graphqlResponse model.CreateTenantResponse
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", err
	}
	return graphqlResponse.Tenant, nil
}

func (s *customerOsClient) addHeadersToGraphRequest(req *graphql.Request, tenant *string) error {
	req.Header.Add("X-Openline-API-KEY", s.cfg.CustomerOS.CustomerOsAPIKey)

	if tenant != nil {
		req.Header.Add("X-Openline-TENANT", *tenant)
	}

	return nil
}

func (s *customerOsClient) contextWithTimeout() (context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	return ctx, cancel, nil
}

func (s *customerOsClient) IsPlayer(authId string, provider string) (string, error) {

	graphqlRequest := graphql.NewRequest(
		`
		query GetPlayer ($authId: String!, $provider: String!) {
				player_ByAuthIdProvider(
					  authId: $authId,
					  provider: $provider
				) { id }
		}
	`)
	graphqlRequest.Var("authId", authId)
	graphqlRequest.Var("provider", provider)

	err := s.addHeadersToGraphRequest(graphqlRequest, nil)
	if err != nil {
		return "", err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", err
	}
	defer cancel()

	var graphqlResponse model.GetPlayerResponse
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", err
	}
	return graphqlResponse.Id, nil
}
