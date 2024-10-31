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
	CreateOrganization(tenant, username string, input model.OrganizationInput) (string, error)
	CreateContact(tenant, username string, contactInput model.ContactInput) (string, error)

	MergeEmailToContact(tenant, contactId string, emailInput model.EmailInput) (string, error)

	LinkContactToOrganization(tenant, contactId, organizationId string) (string, error)
}

func NewCustomerOsClient(customerOSApiPath, customerOSApiKey string) *customerOSApiClient {
	return &customerOSApiClient{
		customerOSApiKey: customerOSApiKey,
		graphqlClient:    graphql.NewClient(customerOSApiPath),
	}
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
		`mutation LinkContactWithOrganization($input: ContactOrganizationInput!) {
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
