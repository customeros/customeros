package service

import (
	"context"
	"errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commonUtils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"strings"
)

type RegistrationService interface {
	CreateOrganizationAndContact(ctx context.Context, tenant, email string, allowPersonalEmail bool, leadSource string) (*string, *string, error)
}

type registrationService struct {
	services *Services
}

func NewRegistrationService(services *Services) RegistrationService {
	return &registrationService{
		services: services,
	}
}

func (s *registrationService) CreateOrganizationAndContact(ctx context.Context, tenant, email string, allowPersonalEmail bool, leadSource string) (*string, *string, error) {
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

	if !isPersonalEmail || allowPersonalEmail {
		organizationByDomain, err := s.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(ctx, nil, tenant, domain)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, nil, err
		}

		if organizationByDomain == nil {
			organizationId, err = s.services.CommonServices.OrganizationService.Save(ctx, nil, nil, data_fields.OrganizationFields{
				Domains:      []string{domain},
				Name:         commonUtils.StringPtr(domain),
				Relationship: commonUtils.ToPtr(enum.OrganizationRelationshipProspect),
				Stage:        commonUtils.ToPtr(enum.Trial),
				LeadSource:   commonUtils.StringPtr(leadSource),
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, nil, err
			}
			if organizationId == "" {
				e := errors.New("organization id empty")
				tracing.TraceErr(span, e)
				return nil, nil, e
			}
		} else {
			organizationId = mapper.MapDbNodeToOrganizationEntity(organizationByDomain).ID
		}

		if organizationId == "" {
			tracing.TraceErr(span, err)
			return nil, nil, err
		}
		span.LogFields(tracingLog.String("result.organizationId", organizationId))

		contactNode, err := s.services.CommonServices.Neo4jRepositories.ContactReadRepository.GetContactInOrganizationByEmail(ctx, tenant, organizationId, email)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, nil, err
		}

		if contactNode == nil {
			contactId, err = s.services.CommonServices.ContactService.CreateContactByEmail(ctx, nil, email)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, nil, err
			}

			err = s.services.CommonServices.ContactService.LinkContactWithOrganization(ctx, nil, contactId, organizationId, "", "",
				neo4jentity.DataSourceOpenline.String(), false, nil, nil)
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
