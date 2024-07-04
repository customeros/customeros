package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type apiCacheService struct {
	repositories *neo4jRepository.Repositories
	services     *Services
}

type ApiCacheService interface {
	GetApiCache(ctx context.Context, tenant string, page, limit int) ([]*map[string]interface{}, error)
	GetPatchesForApiCache(ctx context.Context, tenant string) ([]*map[string]interface{}, error)
}

func (s *apiCacheService) GetApiCache(ctx context.Context, tenant string, page, limit int) ([]*map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ApiCacheService.GetApiCache")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("tenant", tenant), log.Int("page", page), log.Int("limit", limit))

	data, err := s.repositories.OrganizationReadRepository.GetForApiCache(ctx, tenant, page*limit, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	response := make([]*map[string]interface{}, 0)

	for _, row := range data {
		response = append(response, mapNeo4jRowResultToCacheRowResult(row))
	}

	return response, nil
}

func (s *apiCacheService) GetPatchesForApiCache(ctx context.Context, tenant string) ([]*map[string]interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ApiCacheService.GetPatchesForApiCache")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("tenant", tenant))

	now := time.Now().UTC()
	lastPatchTimestamp := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-1, 15, 0, 0, time.UTC)

	data, err := s.repositories.OrganizationReadRepository.GetPatchesForApiCache(ctx, tenant, lastPatchTimestamp)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	response := make([]*map[string]interface{}, 0)

	for _, row := range data {
		response = append(response, mapNeo4jRowResultToCacheRowResult(row))
	}

	return response, nil
}

func mapNeo4jRowResultToCacheRowResult(row map[string]interface{}) *map[string]interface{} {
	organizationNode := row["organization"].(dbtype.Node)
	cacheRowResult := utils.GetPropsFromNode(organizationNode)

	cacheRowResult["metadata"] = map[string]interface{}{
		"id":               utils.GetStringPropOrEmpty(cacheRowResult, "id"),
		"created":          utils.GetTimePropOrNow(cacheRowResult, "createdAt"),
		"lastUpdated":      utils.GetTimePropOrNow(cacheRowResult, "updatedAt"),
		"source":           utils.GetStringPropOrEmpty(cacheRowResult, "source"),
		"sourceOfTruth":    utils.GetStringPropOrEmpty(cacheRowResult, "sourceOfTruth"),
		"appSource":        utils.GetStringPropOrEmpty(cacheRowResult, "appSource"),
		"aggregateVersion": utils.GetStringPropOrEmpty(cacheRowResult, "aggregateVersion"),
	}

	contactList := make([]interface{}, 0)
	for _, dataId := range row["contactList"].([]interface{}) {
		contactList = append(contactList, map[string]interface{}{
			"metadata": map[string]interface{}{
				"id": dataId,
			},
		})
	}
	if len(contactList) > 0 {
		cacheRowResult["contacts"] = map[string]interface{}{
			"content": contactList,
		}
	}

	socialList := make([]interface{}, 0)
	for _, dataId := range row["socialList"].([]interface{}) {
		socialList = append(socialList, map[string]interface{}{
			"metadata": map[string]interface{}{
				"id": dataId,
			},
		})
	}
	if len(socialList) > 0 {
		cacheRowResult["socialMedia"] = socialList
	}

	tagList := make([]interface{}, 0)
	for _, dataId := range row["tagList"].([]interface{}) {
		tagList = append(tagList, map[string]interface{}{
			"metadata": map[string]interface{}{
				"id": dataId,
			},
		})
	}
	if len(tagList) > 0 {
		cacheRowResult["tags"] = tagList
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
		cacheRowResult["subsidiaries"] = subsidiaryList
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
		cacheRowResult["parentCompanies"] = parentList
	}

	if row["ownerId"] != nil && row["ownerId"] != "" {
		cacheRowResult["owner"] = map[string]interface{}{
			"id": row["ownerId"],
		}
	}

	return &cacheRowResult
}

func NewApiCacheService(repositories *neo4jRepository.Repositories, services *Services) ApiCacheService {
	return &apiCacheService{
		repositories: repositories,
		services:     services,
	}
}
