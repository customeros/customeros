package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"strings"
)

type DomainService interface {
	ExtractDomainFromOrganizationWebsite(ctx context.Context, websiteUrl string) string
	IsKnownCompanyHostingUrl(ctx context.Context, website string) bool
	GetDomainsForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.DomainEntities, error)
}

type domainService struct {
	log      logger.Logger
	services *Services
}

func NewDomainService(log logger.Logger, services *Services) DomainService {
	return &domainService{
		log:      log,
		services: services,
	}
}

func (s *domainService) ExtractDomainFromOrganizationWebsite(ctx context.Context, websiteUrl string) string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainService.ExtractDomainFromOrganizationWebsite")
	defer span.Finish()

	if strings.TrimSpace(websiteUrl) == "" {
		return ""
	}

	if s.IsKnownCompanyHostingUrl(ctx, websiteUrl) {
		return ""
	}

	return utils.ExtractDomain(websiteUrl)
}

func (s *domainService) IsKnownCompanyHostingUrl(ctx context.Context, website string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainService.IsKnownCompanyHostingUrl")
	defer span.Finish()

	urlPatterns := s.getKnownOrganizationHostingUrlPatterns(ctx)
	for _, pattern := range urlPatterns {
		if strings.Contains(website, pattern) {
			return true
		}
	}
	return false
}

func (s *domainService) getKnownOrganizationHostingUrlPatterns(ctx context.Context) []string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainService.getKnownOrganizationHostingUrlPatterns")
	defer span.Finish()

	urlPatterns := s.services.Cache.GetOrganizationWebsiteHostingUrlPatters()
	var err error
	if len(urlPatterns) == 0 {
		urlPatterns, err = s.services.PostgresRepositories.OranizationWebsiteHostingPlatformRepository.GetAllUrlPatterns(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error while getting known organization hosting url patterns: %v", err)
			return []string{}
		}
		s.services.Cache.SetOrganizationWebsiteHostingUrlPatters(urlPatterns)
	}
	return urlPatterns
}

func (s *domainService) GetDomainsForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.DomainEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainService.GetDomainsForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationIds", strings.Join(organizationIds, ",")))

	if common.GetTenantFromContext(ctx) == "" {
		err := errors.New("missing tenant on context")
		tracing.TraceErr(span, err)
		return nil, err
	}

	domainsDbResponse, err := s.services.Neo4jRepositories.DomainReadRepository.GetForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	domainEntities := neo4jentity.DomainEntities{}
	for _, v := range domainsDbResponse {
		domainEntity := neo4jmapper.MapDbNodeToDomainEntity(v.Node)
		domainEntity.DataloaderKey = v.LinkedNodeId
		domainEntities = append(domainEntities, *domainEntity)
	}
	return &domainEntities, nil
}
