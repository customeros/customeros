package service

import (
	"errors"
	"fmt"
	"github.com/machinebox/graphql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"

	"context"

	"time"
)

type CustomerOsClient interface {
	GetTenantByWorkspace(workspace *model.WorkspaceInput) (*string, error)
	GetTenantByUserEmail(email string) (*string, error)
	MergeTenantToWorkspace(workspace *model.WorkspaceInput, tenant string) (bool, error)
	MergeTenant(tenant *model.TenantInput) (string, error)
	HardDeleteTenant(context context.Context, tenant, username, reqTenant, reqConfirmTenant string) error

	GetPlayer(authId string, provider string) (*model.GetPlayerResponse, error)
	CreatePlayer(tenant, userId, identityId, authId, provider string) error

	GetUserById(tenant, userId string) (*model.UserResponse, error)
	GetUserByEmail(tenant, email string) (*model.UserResponse, error)

	CreateUser(user *model.UserInput, tenant string, roles []model.Role) (*model.UserResponse, error)
	AddUserRole(tenant, userId string, role model.Role) (*model.UserResponse, error)
	AddUserRoles(tenant, userId string, roles []model.Role) (*model.UserResponse, error)
	//CreateContact(tenant, username, firstName, lastname, email string, profilePhotoUrl *string) (string, error)
	CreateContact(tenant, username string, contactInput model.ContactInput) (string, error)
	AddSocialContact(tenant, username, contactId string, socialInput model.SocialInput) (string, error)
	CreateNoteForContact(tenant, username, contactId string, socialInput model.NoteInput) (string, error)
	LinkContactToOrganization(tenant, contactId, organizationId string) (string, error)
	CreateTenantBillingProfile(tenant, username string, input model.TenantBillingProfileInput) (string, error)
	GetOrganizations(tenant, username string) ([]string, int64, error)
	ArchiveOrganizations(tenant, username string, ids []string) (bool, error)
	CreateOrganization(tenant, username string, input model.OrganizationInput) (string, error)
	UpdateOrganization(tenant, username string, input model.OrganizationUpdateInput) (string, error)
	AddSocialOrganization(tenant, username, organizationId string, socialInput model.SocialInput) (string, error)
	UpdateOrganizationOnboardingStatus(tenant, username string, onboardingStatus model.OrganizationUpdateOnboardingStatus) (string, error)

	CreateContract(tenant, username string, input model.ContractInput) (string, error)
	UpdateContract(tenant, username string, input model.ContractUpdateInput) (string, error)
	GetContractById(tenant, contractId string) (*dbtype.Node, error)

	CreateServiceLine(tenant, username string, input interface{}) (string, error)
	GetServiceLine(tenant, serviceLineId string) (*dbtype.Node, error)

	DryRunNextInvoiceForContractInput(tenant, username, contractId string) (string, error)

	CreateMeeting(tenant, username string, input model.MeetingInput) (string, error)

	CreateInteractionSession(tenant, username string, options ...InteractionSessionBuilderOption) (*string, error)
	CreateInteractionEvent(tenant, username string, options ...InteractionEventBuilderOption) (*string, error)
	CreateLogEntry(tenant, username string, organizationId, author, content, contentType string, startedAt time.Time) (*string, error)

	AddContactToOrganization(tenant, username, contactId, organizationId, jobTitle, description string) error

	CreateMasterPlan(tenant, username, name string) (string, error)
	CreateMasterPlanMilestone(tenant, username string, masterPlanMilestoneInput model.MasterPlanMilestoneInput) (string, error)
}

type customerOsClient struct {
	cfg           *config.Config
	graphqlClient *graphql.Client
	driver        *neo4j.DriverWithContext
	database      string
}

func NewCustomerOsClient(cfg *config.Config, driver *neo4j.DriverWithContext) CustomerOsClient {
	return &customerOsClient{
		cfg:           cfg,
		graphqlClient: graphql.NewClient(cfg.CustomerOS.CustomerOsAPI),
		driver:        driver,
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

func (s *customerOsClient) HardDeleteTenant(ctx context.Context, tenant, username, reqTenant, reqConfirmTenant string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomerOsClient.HardDeleteTenant")
	defer span.Finish()

	span.LogFields(tracingLog.String("tenant", tenant), tracingLog.String("username", username), tracingLog.String("reqTenant", reqTenant), tracingLog.String("reqConfirmTenant", reqConfirmTenant))

	graphqlRequest := graphql.NewRequest(
		`
			mutation HardDeleteTenant($reqTenant: String!, $reqConfirmTenant: String!) {
			   tenant_hardDelete(
					tenant: $reqTenant,
					confirmTenant: $reqConfirmTenant)
			}
	`)
	graphqlRequest.Var("reqTenant", reqTenant)
	graphqlRequest.Var("reqConfirmTenant", reqConfirmTenant)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)

	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return err
	}
	defer cancel()

	var graphqlResponse model.TenantHardDeleteResponse
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return err
	}
	return nil
}

func (s *customerOsClient) GetPlayer(authId string, provider string) (*model.GetPlayerResponse, error) {

	graphqlRequest := graphql.NewRequest(
		`
		query GetPlayer ($authId: String!, $provider: String!) {
				player_ByAuthIdProvider(
					  authId: $authId,
					  provider: $provider
				) { 
					id
					users {
						tenant
					}
				   }
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

func (cosService *customerOsClient) CreateContact(tenant, username string, contactInput model.ContactInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateContact($contactInput: ContactInput!) {
				contact_Create(input: $contactInput) {
					id
				}
			}`)

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

func (cosService *customerOsClient) AddSocialContact(tenant, username, contactId string, socialInput model.SocialInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation AddSocialContact($contactId: ID!, $socialInput: SocialInput!) {
				contact_AddSocial(contactId: $contactId, input: $socialInput) {
					id
				}
			}`)

	graphqlRequest.Var("contactId", contactId)
	graphqlRequest.Var("socialInput", socialInput)

	err := cosService.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)

	if err != nil {
		return "", fmt.Errorf("add headers contact_AddSocial: %w", err)
	}

	ctx, cancel, err := cosService.contextWithTimeout()
	if err != nil {
		return "", fmt.Errorf("context contact_AddSocial: %v", err)
	}
	defer cancel()

	var graphqlResponse map[string]map[string]string
	if err := cosService.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("contact_AddSocial: %w", err)
	}
	id := graphqlResponse["contact_AddSocial"]["id"]
	return id, nil
}

func (cosService *customerOsClient) CreateNoteForContact(tenant, username, contactId string, noteInput model.NoteInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateNoteForContact($contactId: ID!, $noteInput: NoteInput!) {
				note_CreateForContact(contactId: $contactId, input: $noteInput) {
					id
				}
			}`)

	graphqlRequest.Var("contactId", contactId)
	graphqlRequest.Var("noteInput", noteInput)

	err := cosService.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)

	if err != nil {
		return "", fmt.Errorf("add headers note_CreateForContact: %w", err)
	}

	ctx, cancel, err := cosService.contextWithTimeout()
	if err != nil {
		return "", fmt.Errorf("context note_CreateForContact: %v", err)
	}
	defer cancel()

	var graphqlResponse map[string]map[string]string
	if err := cosService.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("note_CreateForContact: %w", err)
	}
	id := graphqlResponse["note_CreateForContact"]["id"]
	return id, nil
}

func (cosService *customerOsClient) LinkContactToOrganization(tenant, contactId, organizationId string) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation LinkContactToOrganization($input: ContactOrganizationInput!) {
				contact_AddOrganizationById(input: $input) {
					id
				}
			}`)

	input := model.ContactOrganizationInput{
		ContactId:      contactId,
		OrganizationId: organizationId,
	}
	graphqlRequest.Var("input", input)

	err := cosService.addHeadersToGraphRequest(graphqlRequest, &tenant, nil)

	if err != nil {
		return "", fmt.Errorf("add headers contact_AddOrganizationById: %w", err)
	}

	ctx, cancel, err := cosService.contextWithTimeout()
	if err != nil {
		return "", fmt.Errorf("context contact_AddOrganizationById: %v", err)
	}
	defer cancel()

	var graphqlResponse map[string]map[string]string
	if err := cosService.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("contact_AddOrganizationById: %w", err)
	}
	id := graphqlResponse["contact_AddOrganizationById"]["id"]
	return id, nil
}

func (s *customerOsClient) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *s.driver, utils.WithDatabaseName(s.database))
}

func (s *customerOsClient) CreateTenantBillingProfile(tenant, username string, input model.TenantBillingProfileInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation TenantAddBillingProfile($input: TenantBillingProfileInput!) {
				tenant_AddBillingProfile(input: $input) {
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

	var graphqlResponse model.TenantAddBillingProfileResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("tenantBillingProfile_Create: %w", err)
	}

	return graphqlResponse.TenantBillingProfileAdd.Id, nil
}

func (s *customerOsClient) ArchiveOrganizations(tenant, username string, ids []string) (bool, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation ArchiveOrganizations($ids: [ID!]!) {
  				organization_ArchiveAll(ids: $ids) {
					result
			}
		}`)

	graphqlRequest.Var("ids", ids)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)
	if err != nil {
		return false, err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return false, err
	}
	defer cancel()

	var graphqlResponse model.ArchiveOrganizationResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return false, fmt.Errorf("organization_ArchiveAll: %w", err)
	}

	return graphqlResponse.Result, nil
}

func (s *customerOsClient) CreateOrganization(tenant, username string, input model.OrganizationInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateOrganization($input: OrganizationInput!) {
  				organization_Create(input: $input) {
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

	var graphqlResponse model.CreateOrganizationResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("organization_Create: %w", err)
	}

	return graphqlResponse.OrganizationCreate.Id, nil
}

func (s *customerOsClient) UpdateOrganization(tenant, username string, input model.OrganizationUpdateInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation UpdateOrganization($input: OrganizationUpdateInput!) {
  				organization_Update(input: $input) {
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

	var graphqlResponse model.UpdateOrganizationResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("organization_Update: %w", err)
	}

	return graphqlResponse.OrganizationUpdate.Id, nil
}

func (cosService *customerOsClient) AddSocialOrganization(tenant, username, organizationId string, socialInput model.SocialInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation AddSocialOrganization($organizationId: ID!, $socialInput: SocialInput!) {
				organization_AddSocial(organizationId: $organizationId, input: $socialInput) {
					id
				}
			}`)

	graphqlRequest.Var("organizationId", organizationId)
	graphqlRequest.Var("socialInput", socialInput)

	err := cosService.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)

	if err != nil {
		return "", fmt.Errorf("add headers organization_AddSocial: %w", err)
	}

	ctx, cancel, err := cosService.contextWithTimeout()
	if err != nil {
		return "", fmt.Errorf("context organization_AddSocial: %v", err)
	}
	defer cancel()

	var graphqlResponse map[string]map[string]string
	if err := cosService.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("organization_AddSocial: %w", err)
	}
	id := graphqlResponse["organization_AddSocial"]["id"]
	return id, nil
}

func (s *customerOsClient) GetOrganizations(tenant, username string) ([]string, int64, error) {
	graphqlRequest := graphql.NewRequest(
		`
			query getOrganizations() {
			  organizations(pagination: {limit: 100, page: 1}) {
				totalElements
				content {
                  id
                }
			  }
			}`)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)
	if err != nil {
		return nil, 0, err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, 0, err
	}
	defer cancel()

	var graphqlResponse model.GetOrganizationsResponse
	if err = s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return nil, 0, err
	}
	var ids []string
	for _, org := range graphqlResponse.Organizations.Content {
		ids = append(ids, org.ID)
	}
	return ids, graphqlResponse.Organizations.TotalElements, nil

}

func (s *customerOsClient) UpdateOrganizationOnboardingStatus(tenant, username string, onboardingStatus model.OrganizationUpdateOnboardingStatus) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation UpdateOrganizationOnboardingStatus($organizationId: ID!, $onboardingStatus: OnboardingStatus!, $onboardingComments: String) {
  				organization_UpdateOnboardingStatus(input: {
					organizationId: $organizationId,
					status: $onboardingStatus,
					comments: $onboardingComments,
					}) {
					id
			}
		}`)

	graphqlRequest.Var("organizationId", onboardingStatus.OrganizationId)
	graphqlRequest.Var("onboardingStatus", onboardingStatus.Status)
	graphqlRequest.Var("onboardingComments", onboardingStatus.Comments)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)
	if err != nil {
		return "", err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", err
	}
	defer cancel()

	var graphqlResponse model.UpdateOrganizationResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("organization_UpdateOnboardingStatus: %w", err)
	}

	return graphqlResponse.OrganizationUpdate.Id, nil
}

func (s *customerOsClient) CreateContract(tenant, username string, input model.ContractInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation createContract($input: ContractInput!) {
				contract_Create(input: $input) {
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

	var graphqlResponse model.CreateContractResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("contract_Create: %w", err)
	}

	return graphqlResponse.ContractCreate.Id, nil
}

func (s *customerOsClient) GetContractById(tenant, contractId string) (*dbtype.Node, error) {
	cypher := `MATCH (:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract {id:$id}) RETURN c`
	params := map[string]any{
		"tenant": tenant,
		"id":     contractId,
	}

	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, err
	}
	defer cancel()
	session := s.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)

	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (s *customerOsClient) UpdateContract(tenant, username string, input model.ContractUpdateInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation updateContract($input: ContractUpdateInput!) {
				contract_Update(input: $input) {
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

	var graphqlResponse model.UpdateContractResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("contract_Update: %w", err)
	}

	return graphqlResponse.ContractUpdate.Id, nil
}

func (s *customerOsClient) CreateServiceLine(tenant, username string, input interface{}) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation serviceLineItem($input: ServiceLineItemInput!) {
				contractLineItem_Create(input: $input) {
					metadata {
						id
					}
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

	var graphqlResponse model.CreateServiceLineItemResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("contractLineItem_Create: %w", err)
	}

	return graphqlResponse.ContractLineItemCreate.Metadata.Id, nil
}

func (s *customerOsClient) GetServiceLine(contractId, serviceLineId string) (*dbtype.Node, error) {
	cypher := `MATCH (c:Contract {id:$contractId})-[:HAS_SERVICE]->(sli:ServiceLineItem {id:$serviceLineId}) RETURN sli`
	params := map[string]any{
		"contractId":    contractId,
		"serviceLineId": serviceLineId,
	}

	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, err
	}
	defer cancel()
	session := s.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)

	})
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (s *customerOsClient) DryRunNextInvoiceForContractInput(tenant, username, contractId string) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation invoice_NextDryRunForContract($contractId: ID!) {
				invoice_NextDryRunForContract(contractId: $contractId)
		}`)

	graphqlRequest.Var("contractId", contractId)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)
	if err != nil {
		return "", err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", err
	}
	defer cancel()

	var graphqlResponse map[string]string
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("invoice_NextDryRunForContract_Create: %w", err)
	}

	id := graphqlResponse["invoice_NextDryRunForContract"]
	return id, nil
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

func (s *customerOsClient) CreateMasterPlan(tenant, username, masterPlanName string) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateMasterPlan($masterPlanName: String!) {
				masterPlan_Create(input: {
						name: $masterPlanName
					}) {
					id
					name
			}
		}`)

	graphqlRequest.Var("masterPlanName", masterPlanName)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)
	if err != nil {
		return "", err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", err
	}
	defer cancel()

	var graphqlResponse model.CreateMasterPlanResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("masterPlan_Create: %w", err)
	}

	return graphqlResponse.MasterPlanCreate.Id, nil
}

func (s *customerOsClient) CreateMasterPlanMilestone(tenant, username string, masterPlanMilestoneInput model.MasterPlanMilestoneInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateMasterPlanMilestone($input: MasterPlanMilestoneInput!) {
				masterPlanMilestone_Create(input: $input) {
					id
				  }
				}`)
	graphqlRequest.Var("input", masterPlanMilestoneInput)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)
	if err != nil {
		return "", err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", err
	}
	defer cancel()

	var graphqlResponse model.CreateMasterPlanMilestoneResponse
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("masterPlanMilestone_Create: %w", err)
	}

	return graphqlResponse.MasterPlanMilestoneCreate.Id, nil
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
