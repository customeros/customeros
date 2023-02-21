package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type QueryService interface {
	GetDashboardViewData(ctx context.Context, page int, limit int, searchTerm *string) (*utils.Pagination, error)
}

type queryService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewQueryService(repositories *repository.Repositories, services *Services) QueryService {
	return &queryService{
		repositories: repositories,
		services:     services,
	}
}

func (s *queryService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

func (s *queryService) GetDashboardViewData(ctx context.Context, page int, limit int, searchTerm *string) (*utils.Pagination, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getNeo4jDriver())
	defer session.Close(ctx)

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}

	dbNodes, err := s.repositories.QueryRepository.GetOrganizationsAndContacts(ctx, session, common.GetContext(ctx).Tenant, paginatedResult.GetSkip(), paginatedResult.GetLimit(), searchTerm)
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodes.Count)

	var rows []*entity.DashboardViewResultEntity

	for _, v := range dbNodes.Pairs {
		row := entity.DashboardViewResultEntity{}
		if v.First != nil {
			row.Organization = s.services.OrganizationService.mapDbNodeToOrganizationEntity(*v.First)
		}
		if v.Second != nil {
			row.Contact = s.services.ContactService.mapDbNodeToContactEntity(*v.Second)
		}
		rows = append(rows, &row)
	}

	paginatedResult.SetRows(rows)
	return &paginatedResult, nil
}
