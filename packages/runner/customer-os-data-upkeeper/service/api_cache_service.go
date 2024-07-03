package service

import (
	"context"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/repository"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/tracing"
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
	cfg          *config.Config
	log          logger.Logger
	repositories *repository.Repositories
}

func NewApiCacheService(cfg *config.Config, log logger.Logger, repositories *repository.Repositories) ApiCacheService {
	return &apiCacheService{
		cfg:          cfg,
		log:          log,
		repositories: repositories,
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

	tenants = tenants[:1] // TODO: remove this line
	tenants[0].Name = "gasposco"

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

			response := make([]map[string]interface{}, 0)
			for page*limit < int(organizationCount) {
				data, err := s.repositories.Neo4jRepositories.OrganizationReadRepository.GetForApiCache(ctx, tenant.Name, page*limit, limit)
				if err != nil {
					return
				}

				for _, row := range data {

					organizationNode := row["organization"].(dbtype.Node)
					organizationProps := utils.GetPropsFromNode(organizationNode)

					organizationProps["metadata"] = map[string]interface{}{
						"id":               utils.GetStringPropOrEmpty(organizationProps, "id"),
						"created":          utils.GetTimePropOrNow(organizationProps, "createdAt"),
						"lastUpdated":      utils.GetTimePropOrNow(organizationProps, "updatedAt"),
						"source":           utils.GetStringPropOrEmpty(organizationProps, "source"),
						"sourceOfTruth":    utils.GetStringPropOrEmpty(organizationProps, "sourceOfTruth"),
						"appSource":        utils.GetStringPropOrEmpty(organizationProps, "appSource"),
						"aggregateVersion": utils.GetStringPropOrEmpty(organizationProps, "aggregateVersion"),
					}

					contactList := mapNeo4jArrayToGraphArray(row["contactList"].([]interface{}))
					if len(contactList) > 0 {
						organizationProps["contacts"] = map[string]interface{}{
							"content": contactList,
						}
					}

					socialList := mapNeo4jArrayToGraphArray(row["socialList"].([]interface{}))
					if len(socialList) > 0 {
						organizationProps["socialMedia"] = socialList
					}

					tagList := mapNeo4jArrayToGraphArray(row["tagList"].([]interface{}))
					if len(tagList) > 0 {
						organizationProps["tags"] = tagList
					}

					subsidiaryList := make([]interface{}, 0)
					for _, dataId := range row["subsidiaryList"].([]interface{}) {
						subsidiaryList = append(subsidiaryList, map[string]interface{}{
							"organization": map[string]interface{}{
								"id": dataId,
							},
						})
					}
					if len(subsidiaryList) > 0 {
						organizationProps["subsidiaries"] = subsidiaryList
					}

					parentList := make([]interface{}, 0)
					for _, dataId := range row["parentList"].([]interface{}) {
						parentList = append(parentList, map[string]interface{}{
							"organization": map[string]interface{}{
								"id": dataId,
							},
						})
					}
					if len(parentList) > 0 {
						organizationProps["parentCompanies"] = parentList
					}

					if row["ownerId"] != nil && row["ownerId"] != "" {
						organizationProps["owner"] = map[string]interface{}{
							"id": row["ownerId"],
						}
					}

					response = append(response, organizationProps)
				}

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
