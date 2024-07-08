package service

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/model"

	"time"
)

type customerOSApiClient struct {
	customerOSApiKey string
	graphqlClient    *graphql.Client
}

type CustomerOSApiClient interface {
	CreateExternalSystem(tenant, username *string, input model.ExternalSystemInput) (string, error)
	CreateOrganization(tenant, username string, input model.OrganizationInput) (string, error)
	CreateContact(tenant, username string, contactInput model.ContactInput) (string, error)

	MergeEmailToContact(tenant, contactId string, emailInput model.EmailInput) (string, error)

	AddSocialToContact(tenant, contactId string, socialInput model.SocialInput) (string, error)

	AddEmailToUser(tenant, userId string, email model.EmailInput) (string, error)
	RemoveEmailFromUser(tenant, userId string, email string) (string, error)

	LinkContactToOrganization(tenant, contactId, organizationId string) (string, error)

	GetInteractionSessionForInteractionEvent(tenant, user *string, interactionEventId string) (*model.InteractionSession, error)
}

func NewCustomerOsClient(customerOSApiPath, customerOSApiKey string) *customerOSApiClient {
	return &customerOSApiClient{
		customerOSApiKey: customerOSApiKey,
		graphqlClient:    graphql.NewClient(customerOSApiPath),
	}
}

func (s *customerOSApiClient) CreateExternalSystem(tenant, username *string, input model.ExternalSystemInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation ExternalSystemCreate($input: ExternalSystemInput!) {
				externalSystem_Create(input: $input)
				}`)
	graphqlRequest.Var("input", input)

	err := s.addHeadersToGraphRequest(graphqlRequest, tenant, username)
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
		return "", fmt.Errorf("externalSystem_Create: %w", err)
	}
	return graphqlResponse["externalSystem_Create"], nil
}

func (s *customerOSApiClient) CreateOrganization(tenant, username string, input model.OrganizationInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateOrganization($input: OrganizationInput!) {
  				organization_Create(input: $input) {
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

	var graphqlResponse map[string]map[string]map[string]string
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("organization_Create: %w", err)
	}
	id := graphqlResponse["organization_Create"]["metadata"]["id"]
	return id, nil
}

func (s *customerOSApiClient) CreateContact(tenant, username string, contactInput model.ContactInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation CreateContact($contactInput: ContactInput!) {
				contact_Create(input: $contactInput)
			}`)

	graphqlRequest.Var("contactInput", contactInput)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, &username)

	if err != nil {
		return "", fmt.Errorf("add headers contact_Create: %w", err)
	}

	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", fmt.Errorf("context contact_Create: %v", err)
	}
	defer cancel()

	var graphqlResponse map[string]string
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("contact_Create: %w", err)
	}
	id := graphqlResponse["contact_Create"]
	return id, nil
}

func (s *customerOSApiClient) LinkContactToOrganization(tenant, contactId, organizationId string) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation LinkContactToOrganization($input: ContactOrganizationInput!) {
				contact_AddOrganizationById(input: $input) {
					id
				}
			}`)

	input := model.ContactOrganizationInput{
		ContactID:      contactId,
		OrganizationID: organizationId,
	}
	graphqlRequest.Var("input", input)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, nil)

	if err != nil {
		return "", fmt.Errorf("add headers contact_AddOrganizationById: %w", err)
	}

	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", fmt.Errorf("context contact_AddOrganizationById: %v", err)
	}
	defer cancel()

	var graphqlResponse map[string]map[string]string
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("contact_AddOrganizationById: %w", err)
	}
	id := graphqlResponse["contact_AddOrganizationById"]["id"]
	return id, nil
}

func (s *customerOSApiClient) MergeEmailToContact(tenant, contactId string, emailInput model.EmailInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation MergeEmailToContact($contactId : ID!, $input: EmailInput!) {
				emailMergeToContact(contactId: $contactId, input: $input) {
					id
				}
			}`)

	graphqlRequest.Var("contactId", contactId)
	graphqlRequest.Var("input", emailInput)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, nil)

	if err != nil {
		return "", fmt.Errorf("add headers emailMergeToContact: %w", err)
	}

	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", fmt.Errorf("context emailMergeToContact: %v", err)
	}
	defer cancel()

	var graphqlResponse map[string]map[string]string
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("emailMergeToContact: %w", err)
	}
	id := graphqlResponse["emailMergeToContact"]["id"]
	return id, nil
}

func (s *customerOSApiClient) AddSocialToContact(tenant, contactId string, socialInput model.SocialInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation AddSocialToContact($contactId : ID!, $input: SocialInput!) {
				contact_AddSocial(contactId: $contactId, input: $input) {
					id
				}
			}`)

	graphqlRequest.Var("contactId", contactId)
	graphqlRequest.Var("input", socialInput)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, nil)

	if err != nil {
		return "", fmt.Errorf("add headers AddSocialToContact: %w", err)
	}

	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", fmt.Errorf("context AddSocialToContact: %v", err)
	}
	defer cancel()

	var graphqlResponse map[string]map[string]string
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("contact_AddSocial: %w", err)
	}
	id := graphqlResponse["contact_AddSocial"]["id"]
	return id, nil
}

func (s *customerOSApiClient) AddEmailToUser(tenant, userId string, email model.EmailInput) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation emailMergeToUser($userId : ID!, $input: EmailInput!) {
				emailMergeToUser(userId: $userId, input: $input) {
					id
				}
			}`)

	graphqlRequest.Var("userId", userId)
	graphqlRequest.Var("input", email)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, nil)

	if err != nil {
		return "", fmt.Errorf("add headers emailMergeToUser: %w", err)
	}

	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", fmt.Errorf("context emailMergeToUser: %v", err)
	}
	defer cancel()

	var graphqlResponse map[string]map[string]string
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("emailMergeToUser: %w", err)
	}
	id := graphqlResponse["emailMergeToUser"]["id"]
	return id, nil
}

func (s *customerOSApiClient) RemoveEmailFromUser(tenant, userId string, email string) (string, error) {
	graphqlRequest := graphql.NewRequest(
		`mutation emailRemoveFromUser($userId : ID!, $email: String!) {
				emailRemoveFromUser(userId: $userId, email: $email) {
					result
				}
			}`)

	graphqlRequest.Var("userId", userId)
	graphqlRequest.Var("email", email)

	err := s.addHeadersToGraphRequest(graphqlRequest, &tenant, nil)

	if err != nil {
		return "", fmt.Errorf("add headers emailRemoveFromUser: %w", err)
	}

	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return "", fmt.Errorf("context emailRemoveFromUser: %v", err)
	}
	defer cancel()

	var graphqlResponse struct{ bool }
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		return "", fmt.Errorf("emailRemoveFromUser: %w", err)
	}
	return "", nil
}

func (s *customerOSApiClient) GetInteractionSessionForInteractionEvent(tenant, user *string, interactionEventId string) (*model.InteractionSession, error) {
	graphqlRequest := graphql.NewRequest(
		`query GetInteractionSession($eventIdentifier: String!) {
  					interactionSession_ByEventIdentifier(eventIdentifier: $eventIdentifier) {
       					id
						sessionIdentifier
				}
			}`)

	graphqlRequest.Var("eventIdentifier", interactionEventId)

	err := s.addHeadersToGraphRequest(graphqlRequest, tenant, user)
	if err != nil {
		return nil, err
	}
	ctx, cancel, err := s.contextWithTimeout()
	if err != nil {
		return nil, err
	}
	defer cancel()

	var graphqlResponse map[string]model.InteractionSession
	if err := s.graphqlClient.Run(ctx, graphqlRequest, &graphqlResponse); err != nil {
		if err.Error() != "graphql: InteractionSession with EventIdentifier "+interactionEventId+" not found" {
			return nil, fmt.Errorf("GetInteractionSession: %w", err)
		} else {
			return nil, nil
		}
	}
	session := graphqlResponse["interactionSession_ByEventIdentifier"]
	return &session, nil
}

func (s *customerOSApiClient) addHeadersToGraphRequest(req *graphql.Request, tenant, username *string) error {
	req.Header.Add("X-Openline-API-KEY", s.customerOSApiKey)

	if tenant != nil && *tenant != "" {
		req.Header.Add("X-Openline-TENANT", *tenant)
	}

	if username != nil && *username != "" {
		req.Header.Add("X-Openline-USERNAME", *username)
	}

	return nil
}

func (s *customerOSApiClient) contextWithTimeout() (context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	return ctx, cancel, nil
}
