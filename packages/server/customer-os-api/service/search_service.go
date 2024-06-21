package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
)

type SearchService interface {
	GCliSearch(ctx context.Context, keyword string, limit *int) (*entity.SearchResultEntities, error)
}

type searchService struct {
	log          logger.Logger
	repositories *repository.Repositories
	services     *Services
}

func NewSearchService(log logger.Logger, repositories *repository.Repositories, services *Services) SearchService {
	return &searchService{
		log:          log,
		repositories: repositories,
		services:     services,
	}
}

func (s *searchService) GCliSearch(ctx context.Context, keyword string, limit *int) (*entity.SearchResultEntities, error) {
	tenant := common.GetTenantFromContext(ctx)
	if limit == nil {
		limit = utils.IntPtr(10)
	}
	records, err := s.repositories.SearchRepository.GCliSearch(ctx, tenant, keyword, *limit)
	if err != nil {
		return nil, err
	}

	// Create a map of entity types and their corresponding labels
	entityLabels := s.prepareEntityLabels(tenant)

	// Create a map of entity types and their corresponding extract functions
	extractFunctions := s.prepareExtractFunctions()

	result := make(entity.SearchResultEntities, 0, len(records))
	for _, v := range records {
		labels, err := utils.AnySliceToStringSlice(v.Values[0].([]any))
		if err != nil {
			s.log.Errorf("(%s) error while converting labels {%v} to string slice: {%v}", utils.GetFunctionName(), v.Values[0].([]any), err.Error())
			continue
		}
		resultEntity := entity.SearchResultEntity{
			Labels: labels,
			Score:  v.Values[1].(float64),
		}
		for entityType, labels := range entityLabels {
			if utils.ContainsAll(resultEntity.Labels, labels) {
				resultEntity.EntityType = entityType
				resultEntity.Node = extractFunctions[entityType](v.Values[2].(dbtype.Node))
				break
			}
		}
		result = append(result, resultEntity)
	}
	return &result, nil
}

func (s *searchService) prepareEntityLabels(tenant string) map[entity.SearchResultEntityType][]string {
	entityLabels := map[entity.SearchResultEntityType][]string{
		entity.SearchResultEntityTypeContact:      neo4jentity.ContactEntity{}.Labels(tenant),
		entity.SearchResultEntityTypeOrganization: neo4jentity.OrganizationEntity{}.Labels(tenant),
		entity.SearchResultEntityTypeEmail:        entity.EmailEntity{}.Labels(tenant),
		entity.SearchResultEntityTypeState:        neo4jentity.StateEntity{}.Labels(),
	}
	return entityLabels
}

func (s *searchService) prepareExtractFunctions() map[entity.SearchResultEntityType]func(dbtype.Node) any {
	extractFunctions := map[entity.SearchResultEntityType]func(dbtype.Node) any{
		entity.SearchResultEntityTypeContact:      s.extractFieldsFromContactNode,
		entity.SearchResultEntityTypeOrganization: s.extractFieldsFromOrganizationNode,
		entity.SearchResultEntityTypeEmail:        s.extractFieldsFromEmailNode,
		entity.SearchResultEntityTypeState:        s.extractFieldsFromStateNode,
	}
	return extractFunctions
}

func (s *searchService) extractFieldsFromContactNode(node dbtype.Node) any {
	return neo4jmapper.MapDbNodeToContactEntity(&node)
}

func (s *searchService) extractFieldsFromOrganizationNode(node dbtype.Node) any {
	return neo4jmapper.MapDbNodeToOrganizationEntity(&node)
}

func (s *searchService) extractFieldsFromEmailNode(node dbtype.Node) any {
	return s.services.EmailService.mapDbNodeToEmailEntity(node)
}

func (s *searchService) extractFieldsFromStateNode(node dbtype.Node) any {
	return mapper.MapDbNodeToStateEntity(node)
}
