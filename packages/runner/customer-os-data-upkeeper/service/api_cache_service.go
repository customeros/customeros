package service

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/tracing"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go/log"
	"sync"
	"time"
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

	tenantNodeList, err := s.repositories.Neo4jRepositories.TenantReadRepository.GetAll(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	tenants := make([]*neo4jEntity.TenantEntity, len(tenantNodeList))
	for i, tenantNode := range tenantNodeList {
		tenants[i] = neo4jmapper.MapDbNodeToTenantEntity(tenantNode)
	}

	now := time.Now().UTC()

	span.LogFields(log.Int("tenant.count", len(tenants)))

	var wg sync.WaitGroup
	wg.Add(len(tenants))

	for _, tenant := range tenants {

		go func(tenant neo4jEntity.TenantEntity) {
			defer wg.Done()

			organizationCount, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.CountByTenant(ctx, tenant.Name)
			if err != nil {
				return
			}

			span.LogFields(log.Int("tenant."+tenant.Name, int(organizationCount)))

			page := 0
			limit := 1000

			response := make([]*commonService.ApiCacheOrganization, 0)
			for page*limit < int(organizationCount) {
				cache, err := s.commonServices.ApiCacheService.GetApiCache(ctx, tenant.Name, page, limit)
				if err != nil {
					return
				}

				response = append(response, cache...)
				page++
			}

			data, err := json.Marshal(response)
			if err != nil {
				tracing.TraceErr(span, err)
				span.LogFields(log.String("tenant."+tenant.Name, err.Error()))
				return
			}
			jsonStr := string(data)

			err = s.repositories.PostgresRepositories.ApiCacheRepository.Save(ctx, entity.ApiCache{
				CreatedAt: now,
				Tenant:    tenant.Name,
				Type:      "ORGANIZATION",
				Data:      jsonStr,
			})

			if err != nil {
				tracing.TraceErr(span, err)
				span.LogFields(log.String("tenant."+tenant.Name, err.Error()))
				return
			}

		}(*tenant)
	}

	wg.Wait()
}

func mapNeo4jArrayToGraphArray(data []interface{}) []interface{} {
	list := make([]interface{}, 0)

	for _, dataId := range data {
		list = append(list, map[string]interface{}{
			"metadata": map[string]interface{}{
				"id": dataId,
			},
		})
	}

	return list
}
