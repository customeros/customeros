package service

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/model"

	"time"
)

type CustomerOsClient interface {
	GetUserById(tenant, userId string) (*model.UserResponse, error)
	GetUserByEmail(tenant, email string) (*model.UserResponse, error)

	AddSocialContact(tenant, username, contactId string, socialInput model.SocialInput) (string, error)
	CreateNoteForContact(tenant, username, contactId string, socialInput model.NoteInput) (string, error)

	CreateTenantBillingProfile(tenant, username string, input model.TenantBillingProfileInput) (string, error)
	GetOrganizations(tenant, username string) ([]string, int64, error)
	ArchiveOrganizations(tenant, username string, ids []string) (bool, error)

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
		return nil, fmt.Errorf("GetById: %w", err)
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
