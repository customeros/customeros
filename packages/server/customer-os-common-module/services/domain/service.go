package domain

import (
	"context"
	"fmt"
	"strings"

	"github.com/customeros/mailsherpa/domaincheck"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neoRepo "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/interfaces"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/services/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type domainService struct {
	log      logger.Logger
	cache    *caches.Cache
	postgres *repository.Repositories
	neo4j    *neoRepo.Repositories
	events   *events.EventsService
}

func NewDomainService(log logger.Logger, cache *caches.Cache, postgres *repository.Repositories, neo4j *neoRepo.Repositories, events *events.EventsService) interfaces.DomainService {
	return &domainService{
		log:      log,
		cache:    cache,
		postgres: postgres,
		neo4j:    neo4j,
		events:   events,
	}
}

func (s *domainService) GetPrimaryDomainForOrganizationWebsite(ctx context.Context, websiteUrl string) (string, string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainService.GetPrimaryDomainForOrganizationWebsite")
	defer span.Finish()
	span.LogKV("websiteUrl", websiteUrl)
	returnedWebsiteUrl := websiteUrl

	websiteUrl = strings.ToLower(websiteUrl)

	if strings.TrimSpace(websiteUrl) == "" {
		return "", ""
	}

	if s.IsKnownCompanyHostingUrl(ctx, websiteUrl) {
		span.LogFields(log.Bool("isKnownCompanyHostingUrl", true))
		return "", returnedWebsiteUrl
	}

	isPrimary, primaryDomain := domaincheck.PrimaryDomainCheck(websiteUrl)
	span.LogFields(log.Bool("isPrimary", isPrimary), log.String("primaryDomain", primaryDomain))
	if !isPrimary && primaryDomain != "" {
		returnedWebsiteUrl = primaryDomain
	}

	if primaryDomain == "" {
		return "", returnedWebsiteUrl
	}

	// TODO: this to be moved into linking org with domain
	if !s.IsAcceptedDomainForOrganization(ctx, primaryDomain) {
		return "", returnedWebsiteUrl
	}

	span.LogKV("result.primaryDomain", primaryDomain, "result.returnedWebsiteUrl", returnedWebsiteUrl)

	return primaryDomain, returnedWebsiteUrl
}

func (s *domainService) IsKnownCompanyHostingUrl(ctx context.Context, website string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainService.IsKnownCompanyHostingUrl")
	defer span.Finish()
	span.LogKV("website", website)

	website = strings.ToLower(website)

	urlPatterns := s.getKnownOrganizationHostingUrlPatterns(ctx)

	// check if url start with pattern or contains pattern prefixed with . or /
	for _, pattern := range urlPatterns {
		if pattern == "" {
			continue
		}
		if strings.HasPrefix(website, pattern) || strings.Contains(website, "."+pattern) || strings.Contains(website, "/"+pattern) {
			span.LogFields(log.String("result.pattern", pattern))
			span.LogFields(log.Bool("result", true))
			return true
		}
	}
	span.LogFields(log.Bool("result", false))
	return false
}

func (s *domainService) getKnownOrganizationHostingUrlPatterns(ctx context.Context) []string {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainService.getKnownOrganizationHostingUrlPatterns")
	defer span.Finish()

	urlPatterns := s.cache.GetOrganizationWebsiteHostingUrlPatters()
	if len(urlPatterns) == 0 {
		dbUrlPatterns, err := s.postgres.OranizationWebsiteHostingPlatformRepository.GetAllUrlPatterns(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
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
	span.LogFields(log.Int("result.count", len(urlPatterns)))
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

	domainsDbResponse, err := s.neo4j.DomainReadRepository.GetForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
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
		err := s.neo4j.DomainWriteRepository.SetPrimaryDetails(context.Background(), domain, primaryDomain, isPrimary)
		if err != nil {
			// Log the error in tracing
			tracing.TraceErr(span, errors.Wrap(err, "Error while setting primary details asynchronously"))
		}

		// If the domain is not primary, trigger the domain merge
		if !isPrimary && primaryDomain != "" {
			err = s.MergeDomain(context.Background(), nil, primaryDomain)
			if err != nil {
				// Log the error during domain merging
				tracing.TraceErr(span, errors.Wrap(err, "Error while merging primary domain asynchronously"))
			}
		}
	}()

	// Return early, as the operation is now async
	return nil
}

func (s *domainService) MergeDomain(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, domain string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainService.MergeDomain")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("domain", domain)

	domain = strings.ToLower(strings.TrimSpace(domain))

	if domain == "" {
		return nil
	}

	if !utils.IsValidDomain(domain) {
		err := errors.New("Invalid domain: " + domain)
		tracing.TraceErr(span, err)
		return err
	}

	_, err := utils.ExecuteWriteInTransactionWithPostCommitActions(ctx, s.neo4j.Neo4jDriver, s.neo4j.Database, txWithPostCommit, func(txWithPostCommit *utils.TxWithPostCommit) (any, error) {
		// create domain db node in neo4j if missing
		domainJustCreated, err := s.neo4j.DomainWriteRepository.MergeDomain(ctx, txWithPostCommit.Tx, domain, neo4jentity.DataSourceOpenline.String(), common.GetAppSourceFromContext(ctx))
		if err != nil {
			return nil, err
		}

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			err = s.events.Publisher.PublishEvent(ctx, domain, model.DOMAIN, dto.CreateDomain{Domain: domain, Source: neo4jentity.DataSourceOpenline.String()})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateDomain"))
			}
			return nil
		})

		txWithPostCommit.AddPostCommitAction(func(ctx context.Context) error {
			if domainJustCreated {
				// read domain from neo4j
				domainEntity, err := s.GetDomain(ctx, domain)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "Error while getting domain"))
					return nil
				}

				// if domain was already checked for primary skip the check
				if domainEntity.IsPrimary == nil {
					err = s.UpdateDomainPrimaryDetails(ctx, domain)
					if err != nil {
						tracing.TraceErr(span, errors.Wrap(err, "Error while checking and updating domain primary"))
					}
				}
			}
			return nil
		})

		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *domainService) GetDomain(ctx context.Context, domain string) (*neo4jentity.DomainEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainService.GetDomain")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, domain)

	domainDbNode, err := s.neo4j.DomainReadRepository.GetDomain(ctx, nil, domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	domainEntity := neo4jmapper.MapDbNodeToDomainEntity(domainDbNode)
	return domainEntity, nil
}

func (s *domainService) IsAcceptedDomainForOrganization(ctx context.Context, domain string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "DomainService.IsAcceptedDomainForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, domain)

	personalEmailProviders := s.cache.GetPersonalEmailProviders()
	if personalEmailProviders == nil || len(personalEmailProviders) == 0 {
		err := fmt.Errorf("personal email providers not loaded")
		tracing.TraceErr(span, err)
		return false
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
		accessible = false
	}
	return accessible, isPrimary, primaryDomain
}
