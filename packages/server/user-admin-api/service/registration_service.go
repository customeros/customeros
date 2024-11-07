package service

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	constants "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commonUtils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
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
		organizationByDomain, err := s.services.CommonServices.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByDomain(ctx, tenant, domain)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, nil, err
		}

		if organizationByDomain == nil {
			orgId, err := s.services.CommonServices.OrganizationService.Save(ctx, nil, tenant, nil, &repository.OrganizationSaveFields{
				Domains:            []string{domain},
				Name:               domain,
				Relationship:       enum.Prospect,
				Stage:              enum.Trial,
				LeadSource:         leadSource,
				UpdateName:         true,
				UpdateRelationship: true,
				UpdateStage:        true,
				UpdateLeadSource:   true,
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, nil, err
			}
			if orgId == nil {
				e := errors.New("organization id empty")
				tracing.TraceErr(span, e)
				return nil, nil, e
			}
			organizationId = *orgId
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
			contactId, err = s.services.CommonServices.ContactService.SaveContact(ctx, nil, repository.ContactFields{}, "", neo4jmodel.ExternalSystem{})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, nil, err
			}

			_, err := s.services.CommonServices.EmailService.Merge(ctx, tenant,
				commonservice.EmailFields{
					Email:     email,
					Source:    neo4jentity.DataSourceOpenline,
					AppSource: constants.AppSourceUserAdminApi,
				}, &commonservice.LinkWith{
					Type: commonModel.CONTACT,
					Id:   contactId,
				})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, nil, err
			}

			err = s.services.CommonServices.ContactService.LinkContactWithOrganization(ctx, contactId, organizationId, "", "",
				neo4jentity.DataSourceOpenline.String(), false, nil, nil)
			if err != nil {
				tracing.TraceErr(span, err)
				graphql.AddErrorf(ctx, "Failed to add organization %s to contact %s", organizationId, contactId)
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
