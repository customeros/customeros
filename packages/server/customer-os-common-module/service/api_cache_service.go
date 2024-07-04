package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
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
	GetApiCache(ctx context.Context, tenant string, page, limit int) ([]*ApiCacheOrganization, error)
	GetPatchesForApiCache(ctx context.Context, tenant string) ([]*ApiCacheOrganization, error)
}

func (s *apiCacheService) GetApiCache(ctx context.Context, tenant string, page, limit int) ([]*ApiCacheOrganization, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ApiCacheService.GetApiCache")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("tenant", tenant), log.Int("page", page), log.Int("limit", limit))

	data, err := s.repositories.OrganizationReadRepository.GetForApiCache(ctx, tenant, page*limit, limit)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	response := make([]*ApiCacheOrganization, 0)

	for _, row := range data {
		response = append(response, mapNeo4jRowResultToCacheRowResult(row))
	}

	return response, nil
}

func (s *apiCacheService) GetPatchesForApiCache(ctx context.Context, tenant string) ([]*ApiCacheOrganization, error) {
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

	response := make([]*ApiCacheOrganization, 0)

	for _, row := range data {
		response = append(response, mapNeo4jRowResultToCacheRowResult(row))
	}

	return response, nil
}

func mapNeo4jRowResultToCacheRowResult(row map[string]interface{}) *ApiCacheOrganization {
	organizationNode := row["organization"].(dbtype.Node)

	response := ApiCacheOrganization{
		Organization:    neo4jmapper.MapDbNodeToOrganizationEntity(&organizationNode),
		Contacts:        make([]*string, 0),
		SocialMedia:     make([]*string, 0),
		Tags:            make([]*string, 0),
		ParentCompanies: make([]*string, 0),
		Subsidiaries:    make([]*string, 0),
		Owner:           nil,
	}

	for _, dataId := range row["contactList"].([]interface{}) {
		s := dataId.(string)
		response.Contacts = append(response.Contacts, &s)
	}

	for _, dataId := range row["socialList"].([]interface{}) {
		s := dataId.(string)
		response.SocialMedia = append(response.SocialMedia, &s)
	}

	for _, dataId := range row["tagList"].([]interface{}) {
		s := dataId.(string)
		response.Tags = append(response.Tags, &s)
	}

	for _, dataId := range row["subsidiaryList"].([]interface{}) {
		s := dataId.(string)
		response.Subsidiaries = append(response.Subsidiaries, &s)
	}

	for _, dataId := range row["parentList"].([]interface{}) {
		s := dataId.(string)
		response.ParentCompanies = append(response.ParentCompanies, &s)
	}

	if row["ownerId"] != nil && row["ownerId"] != "" {
		s := row["ownerId"].(string)
		response.Owner = &s
	}

	return &response
}

type ApiCacheOrganization struct {
	Organization    *neo4jEntity.OrganizationEntity
	Contacts        []*string
	SocialMedia     []*string
	Tags            []*string
	Subsidiaries    []*string
	ParentCompanies []*string
	Owner           *string
}

func NewApiCacheService(repositories *neo4jRepository.Repositories, services *Services) ApiCacheService {
	return &apiCacheService{
		repositories: repositories,
		services:     services,
	}
}
