package service

import (
	"context"
	"fmt"
	"github.com/customeros/mailsherpa/domaincheck"
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
	MergeDomain(ctx context.Context, domain string) error
	ExtractDomainFromOrganizationWebsite(ctx context.Context, websiteUrl string) string
	IsKnownCompanyHostingUrl(ctx context.Context, website string) bool
	GetAllDomainsForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.DomainEntities, error)
	UpdateDomainPrimaryDetails(ctx context.Context, domain string) error
	GetDomain(ctx context.Context, domain string) (*neo4jentity.DomainEntity, error)
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
	span.LogKV("websiteUrl", websiteUrl)

	if strings.TrimSpace(websiteUrl) == "" {
		return ""
	}

	if s.IsKnownCompanyHostingUrl(ctx, websiteUrl) {
		return ""
	}

	personalEmailProviders := s.services.Cache.GetPersonalEmailProviders()
	if personalEmailProviders == nil || len(personalEmailProviders) == 0 {
		err := fmt.Errorf("personal email providers not loaded")
		tracing.TraceErr(span, err)
		return ""
	}

	domain := utils.ExtractDomain(websiteUrl)

	if s.services.Cache.IsPersonalEmailProvider(domain) {
		return ""
	}

	return ""
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

func (s *domainService) GetAllDomainsForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.DomainEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainService.GetAllDomainsForOrganizations")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationIds", strings.Join(organizationIds, ",")))

	err := common.ValidateTenant(ctx)
	if err != nil {
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

func (s *domainService) UpdateDomainPrimaryDetails(ctx context.Context, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainService.UpdateDomainPrimaryDetails")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, domain)

	// Run primary domain check in a separate goroutine
	go func() {
		// Perform the primary domain check asynchronously
		isPrimary, primaryDomain := domaincheck.PrimaryDomainCheck(domain)

		// Call the saving logic after the primary domain check finishes
		err := s.services.Neo4jRepositories.DomainWriteRepository.SetPrimaryDetails(context.Background(), domain, primaryDomain, isPrimary)
		if err != nil {
			// Log the error in tracing
			tracing.TraceErr(span, errors.Wrap(err, "Error while setting primary details asynchronously"))
		}

		// If the domain is not primary, trigger the domain merge
		if !isPrimary && primaryDomain != "" {
			err = s.MergeDomain(context.Background(), primaryDomain)
			if err != nil {
				// Log the error during domain merging
				tracing.TraceErr(span, errors.Wrap(err, "Error while merging primary domain asynchronously"))
			}
		}
	}()

	// Return early, as the operation is now async
	return nil
}

func (s *domainService) MergeDomain(ctx context.Context, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainService.MergeDomain")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	domain = strings.TrimSpace(domain)
	domain = strings.ToLower(domain)

	if domain == "" {
		tracing.TraceErr(span, fmt.Errorf("domain is empty"))
		return nil
	}

	// create domain db node in neo4j if missing
	err := s.services.Neo4jRepositories.DomainWriteRepository.MergeDomain(ctx, domain, neo4jentity.DataSourceOpenline.String(), common.GetAppSourceFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error while merging domain"))
		return err
	}

	// read domain from neo4j
	domainEntity, err := s.GetDomain(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error while getting domain"))
		return err
	}

	// if domain was already checked for primary skip the check
	if domainEntity.IsPrimary == nil {
		err = s.UpdateDomainPrimaryDetails(ctx, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error while checking and updating domain primary"))
		}
	}

	return nil
}

func (s *domainService) GetDomain(ctx context.Context, domain string) (*neo4jentity.DomainEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainService.GetDomain")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, domain)

	domainDbNode, err := s.services.Neo4jRepositories.DomainReadRepository.GetDomain(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	domainEntity := neo4jmapper.MapDbNodeToDomainEntity(domainDbNode)
	return domainEntity, nil
}
