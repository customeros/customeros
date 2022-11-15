package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type EntityDefinitionService interface {
	Create(ctx context.Context, entity *entity.EntityDefinitionEntity) (*entity.EntityDefinitionEntity, error)
	FindAll(ctx context.Context) (*entity.EntityDefinitionEntities, error)
	FindLinkedWithContact(ctx context.Context, contactId string) (*entity.EntityDefinitionEntity, error)
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
	entityDefinition := s.mapDbNodeToEntityDefinition((record.([]*db.Record)[0]).Values[0].(dbtype.Node))
	s.addDbRelationshipToEntity((record.([]*db.Record)[0]).Values[1].(dbtype.Relationship), entityDefinition)
	return entityDefinition, nil
}

func (s *entityDefinitionService) FindAll(ctx context.Context) (*entity.EntityDefinitionEntities, error) {
	all, err := s.repository.EntityDefinitionRepository.FindAllByTenant(common.GetContext(ctx).Tenant)
	if err != nil {
		return nil, err
	}
	entityDefinitionEntities := entity.EntityDefinitionEntities{}
	for _, dbRecord := range all.([]*db.Record) {
		entityDefinitionEntity := s.mapDbNodeToEntityDefinition(dbRecord.Values[0].(dbtype.Node))
		s.addDbRelationshipToEntity(dbRecord.Values[1].(dbtype.Relationship), entityDefinitionEntity)
		entityDefinitionEntities = append(entityDefinitionEntities, *entityDefinitionEntity)
	}
	return &entityDefinitionEntities, nil
}

func (s *entityDefinitionService) FindLinkedWithContact(ctx context.Context, contactId string) (*entity.EntityDefinitionEntity, error) {
	queryResult, err := s.repository.EntityDefinitionRepository.FindByContactId(common.GetContext(ctx).Tenant, contactId)
	if err != nil {
		return nil, err
	}
	if len(queryResult.([]*db.Record)) == 0 {
		return nil, nil
	}
	entityDefinitionEntity := s.mapDbNodeToEntityDefinition((queryResult.([]*db.Record))[0].Values[0].(dbtype.Node))
	s.addDbRelationshipToEntity((queryResult.([]*db.Record))[0].Values[1].(dbtype.Relationship), entityDefinitionEntity)
	return entityDefinitionEntity, nil
}

func (s *entityDefinitionService) mapDbNodeToEntityDefinition(dbNode dbtype.Node) *entity.EntityDefinitionEntity {
	props := utils.GetPropsFromNode(dbNode)
	entityDefinition := entity.EntityDefinitionEntity{
		Id:      utils.GetStringPropOrEmpty(props, "id"),
		Name:    utils.GetStringPropOrEmpty(props, "name"),
		Extends: utils.GetStringPropOrNil(props, "extends"),
		Version: utils.GetIntPropOrMinusOne(props, "version"),
	}
	return &entityDefinition
}

func (s *entityDefinitionService) addDbRelationshipToEntity(relationship dbtype.Relationship, entityDefinition *entity.EntityDefinitionEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	entityDefinition.Added = utils.GetTimePropOrNow(props, "added")
}
