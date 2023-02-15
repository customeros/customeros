package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"github.com/sirupsen/logrus"
)

type SearchService interface {
	SearchBasic(ctx context.Context, keyword string) (*entity.SearchResultEntities, error)
}

type searchService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewSearchService(repositories *repository.Repositories, services *Services) SearchService {
	return &searchService{
		repositories: repositories,
		services:     services,
	}
}

func (s *searchService) SearchBasic(ctx context.Context, keyword string) (*entity.SearchResultEntities, error) {
	tenant := common.GetTenantFromContext(ctx)
	records, err := s.repositories.SearchRepository.SearchBasic(ctx, tenant, keyword)
	if err != nil {
		return nil, err
	}

	// Create a map of entity types and their corresponding labels
	entityLabels := s.prepareEntityLabels(tenant)

	// Create a map of entity types and their corresponding extract functions
	extractFunctions := s.prepareExtractFunctions()

	result := entity.SearchResultEntities{}
	for _, v := range records {
		labels, err := utils.AnySliceToStringSlice(v.Values[0].([]any))
		if err != nil {
			logrus.Errorf("error while converting labels %v to string slice: %v", v.Values[0].([]any), err)
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
		entity.SearchResultEntityTypeContact:      entity.ContactEntity{}.Labels(tenant),
		entity.SearchResultEntityTypeOrganization: entity.OrganizationEntity{}.Labels(tenant),
		entity.SearchResultEntityTypeEmail:        entity.EmailEntity{}.Labels(tenant),
	}
	return entityLabels
}

func (s *searchService) prepareExtractFunctions() map[entity.SearchResultEntityType]func(dbtype.Node) any {
	extractFunctions := map[entity.SearchResultEntityType]func(dbtype.Node) any{
		entity.SearchResultEntityTypeContact:      s.extractFieldsFromContactNode,
		entity.SearchResultEntityTypeOrganization: s.extractFieldsFromOrganizationNode,
		entity.SearchResultEntityTypeEmail:        s.extractFieldsFromEmailNode,
	}
	return extractFunctions
}

func (s *searchService) extractFieldsFromContactNode(node dbtype.Node) any {
	return s.services.ContactService.mapDbNodeToContactEntity(node)
}

func (s *searchService) extractFieldsFromOrganizationNode(node dbtype.Node) any {
	return s.services.OrganizationService.mapDbNodeToOrganizationEntity(node)
}

func (s *searchService) extractFieldsFromEmailNode(node dbtype.Node) any {
	return s.services.EmailService.mapDbNodeToEmailEntity(node)
}
