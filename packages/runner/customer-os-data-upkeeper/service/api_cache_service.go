package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"sync"
)

type ApiCacheService interface {
	RefreshApiCache()
}

type apiCacheService struct {
	cfg            *config.Config
	log            logger.Logger
	repositories   *repository.Repositories
	commonServices *commonService.Services
}

func NewApiCacheService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories, commonServices *commonService.Services) ApiCacheService {
	return &apiCacheService{
		cfg:            cfg,
		log:            log,
		repositories:   repositories,
		commonServices: commonServices,
	}
}

func (s *apiCacheService) RefreshApiCache() {
	ctx, cancel := utils.GetContextWithTimeout(context.Background(), utils.HalfOfHourDuration)
	defer cancel() // Cancel context on exit

	span, ctx := tracing.StartTracerSpan(ctx, "ApiCacheService.RefreshApiCache")
	defer span.Finish()
	tracing.TagComponentCronJob(span)

	tenantNodeList, err := s.repositories.Neo4jRepositories.TenantReadRepository.GetAll(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	tenants := make([]*neo4jEntity.TenantEntity, len(tenantNodeList))
	for i, tenantNode := range tenantNodeList {
		tenants[i] = neo4jmapper.MapDbNodeToTenantEntity(tenantNode)
	}

	span.LogFields(log.Int("tenant.count", len(tenants)))

	if err := s.processTenants(ctx, tenants, span); err != nil {
		tracing.TraceErr(span, err)
	}
}

func (s *apiCacheService) processTenants(ctx context.Context, tenants []*neo4jEntity.TenantEntity, span opentracing.Span) error {
	semaphore := make(chan struct{}, 20)
	errChan := make(chan error, len(tenants))

	var wg sync.WaitGroup
	for _, tenant := range tenants {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(tenant *neo4jEntity.TenantEntity) {
			defer wg.Done()
			defer func() { <-semaphore }()

			if err := s.processTenant(ctx, tenant, span); err != nil {
				errChan <- fmt.Errorf("error processing tenant %s: %w", tenant.Name, err)
			}
		}(tenant)
	}

	wg.Wait()
	close(errChan)

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred while processing tenants: %v", errs)
	}

	return nil
}

func (s *apiCacheService) processTenant(ctx context.Context, tenant *neo4jEntity.TenantEntity, span opentracing.Span) error {
	organizationCount, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.CountByTenant(ctx, tenant.Name)
	if err != nil {
		return fmt.Errorf("error counting organizations: %w", err)
	}

	span.LogFields(log.Int("tenant."+tenant.Name, int(organizationCount)))

	response, err := s.fetchAllOrganizations(ctx, tenant.Name, int(organizationCount))
	if err != nil {
		return fmt.Errorf("error fetching organizations: %w", err)
	}

	data, err := json.Marshal(response)
	if err != nil {
		span.LogFields(log.String("tenant."+tenant.Name, err.Error()))
		return fmt.Errorf("error marshaling response: %w", err)
	}

	err = s.repositories.PostgresRepositories.ApiCacheRepository.Save(ctx, entity.ApiCache{
		CreatedAt: utils.Now(),
		Tenant:    tenant.Name,
		Type:      "ORGANIZATION",
		Data:      string(data),
	})

	if err != nil {
		span.LogFields(log.String("tenant."+tenant.Name, err.Error()))
		return fmt.Errorf("error saving to API cache: %w", err)
	}

	return nil
}

func (s *apiCacheService) fetchAllOrganizations(ctx context.Context, tenantName string, totalCount int) ([]*commonService.ApiCacheOrganization, error) {
	const limit = 1000
	response := make([]*commonService.ApiCacheOrganization, 0, totalCount)

	for page := 0; page*limit < totalCount; page++ {
		cache, err := s.commonServices.ApiCacheService.GetApiCache(ctx, tenantName, page, limit)
		if err != nil {
			return nil, fmt.Errorf("error fetching API cache for page %d: %w", page, err)
		}
		response = append(response, cache...)
	}

	return response, nil
}
