package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"strings"
)

type DomainService interface {
	ExtractDomainFromOrganizationWebsite(ctx context.Context, websiteUrl string) string
	IsKnownCompanyHostingUrl(ctx context.Context, website string) bool
}

type domainService struct {
	log      logger.Logger
	services *Services
	cache    *caches.Cache
}

func NewDomainService(log logger.Logger, services *Services, cache *caches.Cache) DomainService {
	return &domainService{
		log:      log,
		services: services,
		cache:    cache,
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

	urlPatterns := s.cache.GetOrganizationWebsiteHostingUrlPatters()
	var err error
	if len(urlPatterns) == 0 {
		urlPatterns, err = s.services.PostgresRepositories.OranizationWebsiteHostingPlatformRepository.GetAllUrlPatterns(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error while getting known organization hosting url patterns: %v", err)
			return []string{}
		}
		s.cache.SetOrganizationWebsiteHostingUrlPatters(urlPatterns)
	}
	return urlPatterns
}
