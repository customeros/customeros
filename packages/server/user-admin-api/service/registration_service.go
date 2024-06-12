package service

import (
	"context"
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commonUtils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"strings"
)

type RegistrationService interface {
	CreateOrganizationAndContact(ctx context.Context, tenant, email string) (*string, *string, error)
}

type registrationService struct {
	services *Services
}

func NewRegistrationService(services *Services) RegistrationService {
	return &registrationService{
		services: services,
	}
}

func (s *registrationService) CreateOrganizationAndContact(ctx context.Context, tenant, email string) (*string, *string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RegistrationService.CreateOrganizationAndContact")
	defer span.Finish()

	domain := commonUtils.ExtractDomain(email)

	isPersonalEmail := false
	//check if the user is using a personal email provider
	for _, personalEmailProvider := range s.services.Cache.GetPersonalEmailProviders() {
		if strings.Contains(domain, personalEmailProvider) {
			isPersonalEmail = true
			break
		}
	}

	organizationId := ""
	contactId := ""

	if !isPersonalEmail {

		organizationByDomain, err := s.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationWithDomain(ctx, tenant, domain)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, nil, err
		}

		if organizationByDomain == nil {
			prospect := model.OrganizationRelationshipProspect
			lead := model.OrganizationStageLead
			leadSource := "tracking"
			organizationId, err = s.services.CustomerOSApiClient.CreateOrganization(tenant, "", model.OrganizationInput{Relationship: &prospect, Stage: &lead, Domains: []string{domain}, LeadSource: &leadSource})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, nil, err
			}
		} else {
			organizationId = mapper.MapDbNodeToOrganizationEntity(organizationByDomain).ID
		}

		if organizationId == "" {
			tracing.TraceErr(span, errors.New("organization id empty"))
			return nil, nil, errors.New("organization id empty")
		}
		span.LogFields(tracingLog.String("result.organizationId", organizationId))

		contactNode, err := s.services.CommonServices.Neo4jRepositories.ContactReadRepository.GetContactInOrganizationByEmail(ctx, tenant, organizationId, email)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, nil, err
		}

		if contactNode == nil {
			contactInput := model.ContactInput{
				ProfilePhotoURL: nil,
			}

			contactId, err = s.services.CustomerOSApiClient.CreateContact(tenant, "", contactInput)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, nil, err
			}

			_, err = s.services.CustomerOSApiClient.LinkContactToOrganization(tenant, contactId, organizationId)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, nil, err
			}
		} else {
			contactId = mapper.MapDbNodeToContactEntity(contactNode).Id
		}

		if contactId == "" {
			tracing.TraceErr(span, errors.New("contact id empty"))
			return nil, nil, errors.New("contact id empty")
		}
		span.LogFields(tracingLog.String("result.contactId", contactId))
	}

	return &organizationId, &contactId, nil
}
