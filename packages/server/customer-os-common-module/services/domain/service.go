package domain

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/opentracing/opentracing-go/log"
	"strings"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/customeros/mailsherpa/domaincheck"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/caches"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type domainService struct {
	log      logger.Logger
	cache    *caches.Cache
	postgres *postgres_repository.Repositories
	neo4j    *neo4j_repository.Repositories
	events   *events.EventsService
}

func NewDomainService(log logger.Logger, cache *caches.Cache, postgres *postgres_repository.Repositories, neo4j *neo4j_repository.Repositories, events *events.EventsService) interfaces.DomainService {
	return &domainService{
		log:      log,
		cache:    cache,
		postgres: postgres,
		neo4j:    neo4j,
		events:   events,
	}
}

func (s *domainService) GetPrimaryDomainForOrganizationWebsite(ctx context.Context, websiteUrl string) string {
	spans, ctx := telemetry.StartServiceSpan(ctx, "DomainService.GetPrimaryDomainForOrganizationWebsite")
	defer spans.Finish()
	spans.LogKV("websiteUrl", websiteUrl)

	websiteUrl = strings.ToLower(websiteUrl)

	if strings.TrimSpace(websiteUrl) == "" {
		return ""
	}

	if s.IsKnownCompanyHostingUrl(ctx, websiteUrl) {
		spans.LogFields(log.Bool("isKnownCompanyHostingUrl", true))
		return ""
	}

	isPrimary, primaryDomain := domaincheck.PrimaryDomainCheck(websiteUrl)
	spans.LogFields(log.Bool("isPrimary", isPrimary), log.String("primaryDomain", primaryDomain))

	if primaryDomain == "" {
		return ""
	}

	if !utils.IsValidDomain(primaryDomain) {
		primaryDomain = utils.ExtractDomain(primaryDomain)
	}

	if !s.IsAcceptedDomainForOrganization(ctx, primaryDomain) {
		return ""
	}

	spans.LogKV("result.primaryDomain", primaryDomain)

	return primaryDomain
}

func (s *domainService) IsKnownCompanyHostingUrl(ctx context.Context, website string) bool {
	spans, ctx := telemetry.StartServiceSpan(ctx, "DomainService.IsKnownCompanyHostingUrl")
	defer spans.Finish()
	spans.LogKV("website", website)

	website = strings.ToLower(website)

	urlPatterns := s.getKnownOrganizationHostingUrlPatterns(ctx)

	// check if url start with pattern or contains pattern prefixed with . or /
	for _, pattern := range urlPatterns {
		if pattern == "" {
			continue
		}
		if strings.HasPrefix(website, pattern) || strings.Contains(website, "."+pattern) || strings.Contains(website, "/"+pattern) {
			spans.LogKV("result.pattern", pattern)
			spans.LogFields(log.Bool("result", true))
			return true
		}
	}
	spans.LogFields(log.Bool("result", false))
	return false
}

func (s *domainService) getKnownOrganizationHostingUrlPatterns(ctx context.Context) []string {
	spans, ctx := telemetry.StartServiceSpan(ctx, "DomainService.getKnownOrganizationHostingUrlPatterns")
	defer spans.Finish()

	urlPatterns := s.cache.GetOrganizationWebsiteHostingUrlPatters()
	if len(urlPatterns) == 0 {
		dbUrlPatterns, err := s.postgres.OranizationWebsiteHostingPlatformRepository.GetAllUrlPatterns(ctx)
		if err != nil {
			spans.TraceError(err)
			s.log.Errorf("Error while getting known organization hosting url patterns: %v", err)
			return []string{}
		}
		for _, pattern := range dbUrlPatterns {
			if pattern == "" || !strings.Contains(pattern, ".") || !strings.Contains(pattern, "/") {
				// Not a valid pattern, continue. All patterns should have . or /
				continue
			}
			urlPatterns = append(urlPatterns, pattern)
		}
		s.cache.SetOrganizationWebsiteHostingUrlPatters(urlPatterns)
	}
	spans.LogKV("result.count", len(urlPatterns))
	return urlPatterns
}

func (s *domainService) GetAllDomainsForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.DomainEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "DomainService.GetAllDomainsForOrganizations")
	defer spans.Finish()

	spans.LogKV("organizationIds", strings.Join(organizationIds, ","))

	err := common.ValidateTenant(ctx)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	domainsDbResponse, err := s.neo4j.DomainReadRepository.GetForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
	if err != nil {
		spans.TraceError(err)
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
	spans, ctx := telemetry.StartServiceSpan(ctx, "DomainService.UpdateDomainPrimaryDetails")
	defer spans.Finish()

	spans.TagEntity(domain)

	domain = strings.ToLower(strings.TrimSpace(domain))

	// check if domain exists
	domainNode, err := s.neo4j.DomainReadRepository.GetDomain(ctx, nil, domain)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "Error while getting domain"))
		return err
	}
	if domainNode == nil {
		err = errors.New("Domain not found: " + domain)
		spans.TraceError(err)
		return err
	}

	accessible, isPrimary, primaryDomain := s.CheckDomainWithMailsherpa(ctx, domain)
	spans.LogFields(log.Bool("result.mailsherpa.isPrimary", isPrimary), log.String("result.mailsherpa.primaryDomain", primaryDomain), log.Bool("result.mailsherpa.accessible", accessible))
	if primaryDomain != "" {
		if !utils.IsValidDomain(primaryDomain) {
			primaryDomain = ""
		}
		_, isPrimaryDomainPrimary, _ := s.CheckDomainWithMailsherpa(ctx, primaryDomain)
		if !isPrimaryDomainPrimary {
			primaryDomain = ""
		}
	}

	err = s.neo4j.DomainWriteRepository.SetPrimaryDetails(ctx, domain, primaryDomain, isPrimary, accessible)
	if err != nil {
		// Log the error in tracing
		spans.TraceError(errors.Wrap(err, "Error while setting primary details asynchronously"))
	}

	_ = s.events.Publisher.PublishFanoutEvent(ctx, domain, model.DOMAIN, dto.UpdateDomain{Primary: isPrimary, PrimaryDomain: primaryDomain, Accessible: accessible})

	// If the domain is not primary, trigger the domain merge
	if !isPrimary && primaryDomain != "" {
		err = s.MergeDomain(ctx, nil, primaryDomain)
		if err != nil {
			// Log the error during domain merging
			spans.TraceError(errors.Wrap(err, "Error while merging primary domain"))
		}
	}
	return nil
}

func (s *domainService) MergeDomain(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, domain string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "DomainService.MergeDomain")
	defer spans.Finish()

	spans.LogKV("domain", domain)

	domain = strings.ToLower(strings.TrimSpace(domain))

	if domain == "" {
		return nil
	}

	if !utils.IsValidDomain(domain) {
		err := errors.New("Invalid domain: " + domain)
		spans.TraceError(err)
		return err
	}

	_, err := utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		// create domain db node in neo4j if missing
		domainJustCreated, err := s.neo4j.DomainWriteRepository.MergeDomain(ctx, txWithPostCommit.Tx, domain, neo4jentity.DataSourceOpenline.String(), common.GetAppSourceFromContext(ctx))
		if err != nil {
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			if domainJustCreated {
				_ = s.events.Publisher.PublishFanoutEvent(ctx, domain, model.DOMAIN, dto.CreateDomain{Domain: domain, Source: neo4jentity.DataSourceOpenline.String()})
			}
			return nil
		})

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			if domainJustCreated {
				// read domain from neo4j
				domainEntity, err := s.GetDomain(ctx, domain)
				if err != nil {
					spans.TraceError(errors.Wrap(err, "Error while getting domain"))
					return nil
				}

				// if domain was already checked for primary skip the check
				if domainEntity.IsPrimary == nil {
					err = s.UpdateDomainPrimaryDetails(ctx, domain)
					if err != nil {
						spans.TraceError(errors.Wrap(err, "Error while checking and updating domain primary"))
					}
				}
			}
			return nil
		})

		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (s *domainService) GetDomain(ctx context.Context, domain string) (*neo4jentity.DomainEntity, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "DomainService.GetDomain")
	defer spans.Finish()

	spans.TagEntity(domain)

	domainDbNode, err := s.neo4j.DomainReadRepository.GetDomain(ctx, nil, domain)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	domainEntity := neo4jmapper.MapDbNodeToDomainEntity(domainDbNode)
	return domainEntity, nil
}

func (s *domainService) IsAcceptedDomainForOrganization(ctx context.Context, domain string) bool {
	spans, ctx := telemetry.StartServiceSpan(ctx, "DomainService.IsAcceptedDomainForOrganization")
	defer spans.Finish()
	spans.TagEntity(domain)

	personalEmailProviders := s.cache.GetPersonalEmailProviders()
	if len(personalEmailProviders) == 0 {
		// get personal email providers from personal_email_providers table
		personalEmailProviderEntities, err := s.postgres.PersonalEmailProviderRepository.GetPersonalEmailProviders(ctx)
		if err != nil {
			spans.TraceError(err)
			return false
		}
		// convert to slice of strings
		personalEmailProviders = make([]string, 0, len(personalEmailProviderEntities))
		for _, v := range personalEmailProviderEntities {
			personalEmailProviders = append(personalEmailProviders, v.ProviderDomain)
		}
		// set personal email providers in cache
		s.cache.SetPersonalEmailProviders(personalEmailProviders)
	}

	if s.cache.IsPersonalEmailProvider(domain) {
		return false
	}

	return true
}

func (s *domainService) CheckDomainWithMailsherpa(ctx context.Context, domain string) (bool, bool, string) {
	if domain == "" {
		return false, false, ""
	}
	isPrimary, primaryDomain := domaincheck.PrimaryDomainCheck(domain)

	accessible := true
	if !isPrimary && primaryDomain == "" {
		// check if current domain is an exception case
		var err error
		accessible, err = s.postgres.DomainPrimaryExceptionRepository.Exists(ctx, domain)
		if err != nil {
			s.log.Errorf("Error while checking domain primary exception: %v", err.Error())
		}
		if accessible {
			primaryDomain = domain
		}
	}
	return accessible, isPrimary, primaryDomain
}
