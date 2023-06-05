package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type QueryService interface {
	GetDashboardViewContactsData(ctx context.Context, page int, limit int, where *model.Filter, sort *model.SortBy) (*utils.Pagination, error)
	GetDashboardViewOrganizationsData(ctx context.Context, page int, limit int, where *model.Filter, sort *model.SortBy, ownerId string) (*utils.Pagination, error)
}

type queryService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewQueryService(log logger.Logger, repositories *repository.Repositories, services *Services) QueryService {
	return &queryService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (s *queryService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *queryService) GetDashboardViewContactsData(ctx context.Context, page int, limit int, where *model.Filter, sort *model.SortBy) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QueryService.GetDashboardViewContactsData")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	span.LogFields(log.Int("page", page), log.Int("limit", limit))
	if where != nil {
		span.LogFields(log.Object("filter", *where))
	}
	if sort != nil {
		span.LogFields(log.Object("sort", *sort))
	}

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}

	dbNodes, err := s.repositories.QueryRepository.GetDashboardViewContactsData(ctx, common.GetContext(ctx).Tenant, paginatedResult.GetSkip(), paginatedResult.GetLimit(), where, sort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodes.Count)

	contactEntities := entity.ContactEntities{}

	for _, v := range dbNodes.Nodes {
		contactEntities = append(contactEntities, *s.services.ContactService.mapDbNodeToContactEntity(*v))
	}

	paginatedResult.Rows = &contactEntities
	return &paginatedResult, nil
}

func (s *queryService) GetDashboardViewOrganizationsData(ctx context.Context, page int, limit int, where *model.Filter, sort *model.SortBy, ownerId string) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "QueryService.GetDashboardViewOrganizationsData")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Int("page", page), log.Int("limit", limit), log.String("ownerId", ownerId))
	if where != nil {
		span.LogFields(log.Object("filter", *where))
	}
	if sort != nil {
		span.LogFields(log.Object("sort", *sort))
	}

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}

	dbNodes, err := s.repositories.QueryRepository.GetDashboardViewOrganizationData(ctx, common.GetContext(ctx).Tenant, ownerId, paginatedResult.GetSkip(), paginatedResult.GetLimit(), where, sort)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodes.Count)

	organizationEntities := entity.OrganizationEntities{}

	for _, v := range dbNodes.Nodes {
		organizationEntities = append(organizationEntities, *s.services.OrganizationService.mapDbNodeToOrganizationEntity(*v))
	}

	paginatedResult.SetRows(&organizationEntities)
	return &paginatedResult, nil
}
