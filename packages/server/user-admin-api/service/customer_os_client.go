package service

import (
	"errors"
	"fmt"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"

	"context"

	"time"
)

type CustomerOsClient interface {
	GetTenantByWorkspace(workspace *model.WorkspaceInput) (*string, error)
	GetTenantByUserEmail(email string) (*string, error)
	MergeTenantToWorkspace(workspace *model.WorkspaceInput, tenant string) (bool, error)
	MergeTenant(tenant *model.TenantInput) (string, error)

	GetPlayer(authId string, provider string) (*model.GetPlayerResponse, error)
	CreatePlayer(tenant, userId, identityId, authId, provider string) error

	GetUserById(tenant, userId string) (*model.UserResponse, error)
	GetUserByEmail(tenant, email string) (*model.UserResponse, error)

	CreateUser(user *model.UserInput, tenant string, roles []model.Role) (*model.UserResponse, error)
	AddUserRole(tenant, userId string, role model.Role) (*model.UserResponse, error)
	AddUserRoles(tenant, userId string, roles []model.Role) (*model.UserResponse, error)
	CreateContact(tenant, username, firstName, lastname, email string) (string, error)
	CreateOrganization(tenant, username, organizationName, domain string) (string, error)
	CreateMeeting(tenant, username string, input model.MeetingInput) (string, error)

	CreateInteractionSession(tenant, username string, options ...InteractionSessionBuilderOption) (*string, error)
	CreateInteractionEvent(tenant, username string, options ...InteractionEventBuilderOption) (*string, error)
	CreateLogEntry(tenant, username string, organizationId, author, content, contentType string, startedAt time.Time) (*string, error)

	AddContactToOrganization(tenant, username, contactId, organizationId, jobTitle, description string) error
}

type customerOsClient struct {
	cfg           *config.Config
	graphqlClient *graphql.Client
}

func NewCustomerOsClient(cfg *config.Config) CustomerOsClient {
	return &customerOsClient{
		cfg:           cfg,
		graphqlClient: graphql.NewClient(cfg.CustomerOS.CustomerOsAPI),
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

	err := s.addHeadersToGraphRequest(graphqlRequest, nil, nil)
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
	return graphqlResponse.Name, nil
}

func (s *customerOsClient) GetTenantByUserEmail(email string) (*string, error) {
	if email == "" {
		return nil, errors.New("GetTenantByUserEmail: email is empty")
	}
	graphqlRequest := graphql.NewRequest(
		`
		query GetTenantByEmail ($email: String!) {
				tenant_ByEmail(email: $email) 
		}
	`)
	graphqlRequest.Var("email", email)

	err := s.addHeadersToGraphRequest(graphqlRequest, nil, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, err
	}
	defer cancel()

	var graphqlResponse struct {
		Tenant_ByEmail *string
	}
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return nil, err
	}
	return graphqlResponse.Tenant_ByEmail, nil
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

	err := s.addHeadersToGraphRequest(graphqlRequest, nil, nil)
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

func (s *customerOsClient) AddUserRole(tenant, userId string, role model.Role) (*model.UserResponse, error) {
	graphqlRequest := graphql.NewRequest(
		`
			mutation user_AddRole($id: ID!, $role: Role!) {
			  user_AddRole(id: $id, role: $role) {
				id
				roles
			  }
			}
	`)
	graphqlRequest.Var("id", userId)
	graphqlRequest.Var("role", role)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, err
	}
	defer cancel()

	var userAddRoleResponse model.UserAddRoleResponse
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &userAddRoleResponse); err != nil {
		return nil, err
	}
	return &model.UserResponse{
		ID:    userAddRoleResponse.UserAddRole.ID,
		Roles: userAddRoleResponse.UserAddRole.Roles,
	}, nil
}

func (s *customerOsClient) AddUserRoles(tenant, userId string, roles []model.Role) (*model.UserResponse, error) {
	for _, role := range roles {
		_, err := s.AddUserRole(tenant, userId, role)
		if err != nil {
			return nil, err
		}
	}

	return s.GetUserById(tenant, userId)
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

	err := s.addHeadersToGraphRequest(graphqlRequest, nil, nil)
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

func (s *customerOsClient) GetPlayer(authId string, provider string) (*model.GetPlayerResponse, error) {

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

	err := s.addHeadersToGraphRequest(graphqlRequest, nil, nil)

	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, err
	}
	defer cancel()

	var graphqlResponse model.GetPlayerResponse
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		if err.Error() == "graphql: Failed to get player by authId and provider" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &graphqlResponse, nil
}

func (s *customerOsClient) CreatePlayer(tenant, userId, identityId, authId, provider string) error {

	graphqlRequest := graphql.NewRequest(
		`
		mutation MergePlayer ($userId: ID!, $identityId: String!, $authId: String!, $provider: String!) {
				player_Merge(userId: $userId, input: {
					  identityId: $identityId,
					  authId: $authId,
					  provider: $provider
				}) { result }
		}
	`)
	graphqlRequest.Var("userId", userId)
	graphqlRequest.Var("identityId", identityId)
	graphqlRequest.Var("authId", authId)
	graphqlRequest.Var("provider", provider)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, nil)

	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return err
	}
	defer cancel()

	var graphqlResponse model.Result
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return err
	}
	return nil
}

func (s *customerOsClient) CreateUser(user *model.UserInput, tenant string, roles []model.Role) (*model.UserResponse, error) {
	graphqlRequest := graphql.NewRequest(
		`
			mutation AddUser($user: UserInput!) {
			   user_Create(
					input: $user) {
				id
				roles
			  }
			}
	`)
	graphqlRequest.Var("tenant", tenant)
	graphqlRequest.Var("user", *user)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, err
	}
	defer cancel()

	var graphqlResponse model.CreateUserResponse
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return nil, err
	}

	for _, role := range roles {
		_, err = s.AddUserRole(tenant, graphqlResponse.User.ID, role)
		if err != nil {
			return nil, err
		}
	}
	return &model.UserResponse{
		ID:    graphqlResponse.User.ID,
		Roles: graphqlResponse.User.Roles,
	}, nil
}

func (cosService *customerOsClient) CreateContact(tenant, username, firstName, lastname, email string) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateContact($contactInput: ContactInput!) {
				contact_Create(input: $contactInput) {
					id
				}
			}`)

	contactInput := model.ContactInput{
		FirstName: &firstName,
		LastName:  &lastname,
		Email: &model.EmailInput{
			Email: email,
		},
	}
	graphqlRequest.Var("contactInput", contactInput)

	err := cosService.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)

	if err != nil {
		return "", fmt.Errorf("add headers contact_Create: %w", err)
	}

	ctx, cancel, err := cosService.contextWithTimeout()
	if err != nil {
		return "", fmt.Errorf("context contact_Create: %v", err)
	}
	defer cancel()

	var graphqlResponse map[string]map[string]string
	if err := cosService.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("contact_Create: %w", err)
	}
	id := graphqlResponse["contact_Create"]["id"]
	return id, nil
}

func (s *customerOsClient) CreateOrganization(tenant, username, organizationName, domain string) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateOrganization($organizationName: String!, $domain: String!) {
  				organization_Create(input: {
					name: $organizationName,
					domains: [$domain],
					}) {
					id
			}
		}`)

	graphqlRequest.Var("organizationName", organizationName)
	graphqlRequest.Var("domain", domain)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)
	if err != nil {
		return "", err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", err
	}
	defer cancel()

	var graphqlResponse model.CreateOrganizationResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("organization_Create: %w", err)
	}

	return graphqlResponse.OrganizationCreate.Id, nil
}

func (s *customerOsClient) AddContactToOrganization(tenant, username, contactId, organizationId, jobTitle, description string) error {
	graphqlRequest := graphql.NewRequest(
		`mutation AddOrganizationToContact($contactId: ID!, $organizationId: ID!, $jobTitle: String, $description: String) {
			  jobRole_Create(contactId : $contactId, input: {organizationId: $organizationId, jobTitle: $jobTitle, description: $description}) {
				id
			  }
			}`)

	graphqlRequest.Var("contactId", contactId)
	graphqlRequest.Var("organizationId", organizationId)
	graphqlRequest.Var("jobTitle", jobTitle)
	graphqlRequest.Var("description", description)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)
	if err != nil {
		return err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return err
	}
	defer cancel()

	var graphqlResponse map[string]map[string]string
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return fmt.Errorf("contact_AddOrganizationById: %w", err)
	}
	return nil
}

func (s *customerOsClient) GetUserById(tenant, userId string) (*model.UserResponse, error) {
	graphqlRequest := graphql.NewRequest(
		`
			query GetUserById($id: ID!) {
			  user(id: $id) {
				id
				roles
			  }
			}`)

	graphqlRequest.Var("id", userId)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, nil)

	if err != nil {
		return nil, fmt.Errorf("user_ByEmail: %w", err)
	}

	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, fmt.Errorf("user_ByEmail: %v", err)
	}
	defer cancel()

	var getUserResponse model.GetUserByIdResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &getUserResponse); err != nil {
		if err.Error() == "graphql: User with id "+userId+" not identified" {
			return nil, nil
		} else {
			return nil, fmt.Errorf("GetUserById: %w", err)
		}
	}
	return &model.UserResponse{
		ID:    getUserResponse.User.ID,
		Roles: getUserResponse.User.Roles,
	}, nil
}

func (s *customerOsClient) GetUserByEmail(tenant, email string) (*model.UserResponse, error) {
	graphqlRequest := graphql.NewRequest(
		`
			query GetUserByEmail($email: String!) {
			  user_ByEmail(email: $email) {
				id
				roles
			  }
			}`)

	graphqlRequest.Var("email", email)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, nil)

	if err != nil {
		return nil, fmt.Errorf("user_ByEmail: %w", err)
	}

	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, fmt.Errorf("user_ByEmail: %v", err)
	}
	defer cancel()

	var getUserResponse model.GetUserByEmailResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &getUserResponse); err != nil {
		if err.Error() == "graphql: User with email "+email+" not identified" {
			return nil, nil
		} else {
			return nil, fmt.Errorf("user_ByEmail: %w", err)
		}
	}
	return &model.UserResponse{
		ID:    getUserResponse.UserByEmail.ID,
		Roles: getUserResponse.UserByEmail.Roles,
	}, nil
}

func (s *customerOsClient) CreateMeeting(tenant, username string, input model.MeetingInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateMeeting($input: MeetingInput!) {
  				meeting_Create(meeting: $input) {
					id
			}
		}`)

	graphqlRequest.Var("input", input)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)
	if err != nil {
		return "", err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", err
	}
	defer cancel()

	var graphqlResponse model.CreateMeetingResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("meeting_Create: %w", err)
	}

	return graphqlResponse.MeetingCreate.Id, nil
}

func (s *customerOsClient) CreateInteractionSession(tenant, username string, options ...InteractionSessionBuilderOption) (*string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateInteractionSession($sessionIdentifier: String, $channel: String, $name: String!, $type: String, $status: String!, $appSource: String!, $attendedBy: [InteractionSessionParticipantInput!]) {
				interactionSession_Create(
				session: {
					sessionIdentifier: $sessionIdentifier
        			channel: $channel
        			name: $name
        			status: $status
					type: $type
        			appSource: $appSource
                    attendedBy: $attendedBy
    			}
  			) {
				id
   			}
		}
	`)

	params := InteractionSessionBuilderOptions{}
	for _, opt := range options {
		opt(&params)
	}

	graphqlRequest.Var("sessionIdentifier", params.sessionIdentifier)
	graphqlRequest.Var("channel", params.channel)
	graphqlRequest.Var("name", params.name)
	graphqlRequest.Var("status", params.status)
	graphqlRequest.Var("appSource", params.appSource)
	graphqlRequest.Var("attendedBy", params.attendedBy)
	graphqlRequest.Var("type", params.sessionType)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)
	if err != nil {
		return nil, fmt.Errorf("interactionSession_Create: %w", err)
	}

	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, fmt.Errorf("interactionSession_Create: %v", err)
	}
	defer cancel()

	var graphqlResponse map[string]map[string]string
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return nil, fmt.Errorf("CreateInteractionSession: %w", err)
	}
	id := graphqlResponse["interactionSession_Create"]["id"]
	return &id, nil

}

func (s *customerOsClient) CreateInteractionEvent(tenant, username string, options ...InteractionEventBuilderOption) (*string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateInteractionEvent(
				$sessionId: ID, 
				$meetingId: ID,
				$eventIdentifier: String,
				$externalId: String,
				$externalSystemId: String,
				$channel: String,
				$channelData: String,
				$sentBy: [InteractionEventParticipantInput!]!, 
				$sentTo: [InteractionEventParticipantInput!]!, 
				$appSource: String!, 
				$repliesTo: ID, 
				$content: String, 
				$contentType: String
				$eventType: String,
				$createdAt: Time) {
  					interactionEvent_Create(
    					event: {interactionSession: $sessionId, 
								meetingId: $meetingId,
								eventIdentifier: $eventIdentifier,
								externalId: $externalId,
								externalSystemId: $externalSystemId,
								channel: $channel, 
								channelData: $channelData,
								sentBy: $sentBy, 
								sentTo: $sentTo, 
								appSource: $appSource, 
								repliesTo: $repliesTo, 
								content: $content, 
								contentType: $contentType
								eventType: $eventType,
								createdAt: $createdAt}
  					) {
						id
					  }
					}`)

	params := InteractionEventBuilderOptions{
		sentTo: []model.InteractionEventParticipantInput{},
		sentBy: []model.InteractionEventParticipantInput{},
	}
	for _, opt := range options {
		opt(&params)
	}

	graphqlRequest.Var("sessionId", params.sessionId)
	graphqlRequest.Var("eventIdentifier", params.eventIdentifier)
	graphqlRequest.Var("externalId", params.externalId)
	graphqlRequest.Var("externalSystemId", params.externalSystemId)
	graphqlRequest.Var("content", params.content)
	graphqlRequest.Var("contentType", params.contentType)
	graphqlRequest.Var("channelData", params.channelData)
	graphqlRequest.Var("channel", params.channel)
	graphqlRequest.Var("eventType", params.eventType)
	graphqlRequest.Var("sentBy", params.sentBy)
	graphqlRequest.Var("sentTo", params.sentTo)
	graphqlRequest.Var("appSource", params.appSource)
	graphqlRequest.Var("meetingId", params.meetingId)
	graphqlRequest.Var("createdAt", params.createdAt)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)

	if err != nil {
		return nil, fmt.Errorf("error while adding headers to graph request: %w", err)
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, fmt.Errorf("GetInteractionEvent: %w", err)
	}
	defer cancel()

	var graphqlResponse map[string]map[string]string
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return nil, fmt.Errorf("CreateInteractionSession: %w", err)
	}
	id := graphqlResponse["interactionEvent_Create"]["id"]
	return &id, nil
}

func (s *customerOsClient) CreateLogEntry(tenant, username string, organizationId, author, content, contentType string, startedAt time.Time) (*string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateLogEntry($organizationId: ID!, $content: String, $contentType: String, $startedAt: Time) {
			  logEntry_CreateForOrganization(
				organizationId: $organizationId
				input: {content: $content, contentType: $contentType, startedAt: $startedAt}
			  )
			}`)

	graphqlRequest.Var("organizationId", organizationId)
	graphqlRequest.Var("content", content)
	graphqlRequest.Var("contentType", contentType)
	graphqlRequest.Var("contentType", contentType)
	graphqlRequest.Var("startedAt", startedAt)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, &author)

	if err != nil {
		return nil, fmt.Errorf("error while adding headers to graph request: %w", err)
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, fmt.Errorf("GetInteractionEvent: %w", err)
	}
	defer cancel()

	var graphqlResponse map[string]string
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return nil, fmt.Errorf("Error logEntry_CreateForOrganization: %w", err)
	}
	id := graphqlResponse["logEntry_CreateForOrganization"]
	return &id, nil
}

func (s *customerOsClient) addHeadersToGraphRequest(req *graphql.Request, tenant, username *string) error {
	req.Header.Add("X-Openline-API-KEY", s.cfg.CustomerOS.CustomerOsAPIKey)

	if tenant != nil {
		req.Header.Add("X-Openline-TENANT", *tenant)
	}

	if username != nil {
		req.Header.Add("X-Openline-USERNAME", *username)
	}

	return nil
}

func (s *customerOsClient) contextWithTimeout() (context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	return ctx, cancel, nil
}
