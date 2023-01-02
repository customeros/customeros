package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type EntityDefinitionService interface {
	Create(ctx context.Context, entity *entity.EntityDefinitionEntity) (*entity.EntityDefinitionEntity, error)
	FindAll(ctx context.Context, extends *string) (*entity.EntityDefinitionEntities, error)
	FindLinkedWithContact(ctx context.Context, contactId string) (*entity.EntityDefinitionEntity, error)
}

type entityDefinitionService struct {
	repository *repository.Repositories
}

func NewEntityDefinitionService(repository *repository.Repositories) EntityDefinitionService {
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
	return entityDefinition, nil
}

func (s *entityDefinitionService) FindAll(ctx context.Context, extends *string) (*entity.EntityDefinitionEntities, error) {
	session := utils.NewNeo4jReadSession(*s.repository.Drivers.Neo4jDriver)
	defer session.Close()

	var err error
	var entityDefinitionsDbRecords []*db.Record
	if extends == nil {
		entityDefinitionsDbRecords, err = s.repository.EntityDefinitionRepository.FindAllByTenant(session, common.GetContext(ctx).Tenant)
	} else {
		entityDefinitionsDbRecords, err = s.repository.EntityDefinitionRepository.FindAllByTenantAndExtends(session, common.GetContext(ctx).Tenant, *extends)
	}
	if err != nil {
		return nil, err
	}
	entityDefinitionEntities := entity.EntityDefinitionEntities{}
	for _, dbRecord := range entityDefinitionsDbRecords {
		entityDefinitionEntity := s.mapDbNodeToEntityDefinition(dbRecord.Values[0].(dbtype.Node))
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
	return entityDefinitionEntity, nil
}

func (s *entityDefinitionService) mapDbNodeToEntityDefinition(dbNode dbtype.Node) *entity.EntityDefinitionEntity {
	props := utils.GetPropsFromNode(dbNode)
	entityDefinition := entity.EntityDefinitionEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		Extends:   utils.GetStringPropOrNil(props, "extends"),
		Version:   utils.GetIntPropOrMinusOne(props, "version"),
		CreatedAt: utils.GetTimePropOrNow(props, "createdAt"),
	}
	return &entityDefinition
}
