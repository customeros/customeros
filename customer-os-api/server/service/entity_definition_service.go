package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type EntityDefinitionService interface {
	Create(ctx context.Context, entity *entity.EntityDefinitionEntity) (*entity.EntityDefinitionEntity, error)
}

type entityDefinitionService struct {
	repository *repository.RepositoryContainer
}

func NewEntityDefinitionService(repository *repository.RepositoryContainer) EntityDefinitionService {
	return &entityDefinitionService{
		repository: repository,
	}
}

func (s *entityDefinitionService) Create(ctx context.Context, entity *entity.EntityDefinitionEntity) (*entity.EntityDefinitionEntity, error) {
	entity.Version = 1
	record, err := s.repository.EntityDefinitionRepository.Create(common.GetContext(ctx).Tenant, entity)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToEntityDefinition(record.(dbtype.Node)), nil
}

func (s *entityDefinitionService) mapDbNodeToEntityDefinition(dbNode dbtype.Node) *entity.EntityDefinitionEntity {
	props := utils.GetPropsFromNode(dbNode)
	entityDefinition := entity.EntityDefinitionEntity{
		Id:      utils.GetStringPropOrEmpty(props, "id"),
		Name:    utils.GetStringPropOrEmpty(props, "name"),
		Extends: utils.GetStringPropOrNil(props, "extends"),
		Version: utils.GetIntPropOrDefault(props, "version"),
	}
	return &entityDefinition
}
