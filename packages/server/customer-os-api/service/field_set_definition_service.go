package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type FieldSetDefinitionService interface {
	FindAll(entityDefinitionId string) (*entity.FieldSetDefinitionEntities, error)
	FindLinkedWithFieldSet(ctx context.Context, fieldSetId string) (*entity.FieldSetDefinitionEntity, error)
}

type fieldSetDefinitionService struct {
	repository *repository.Repositories
}

func NewFieldSetDefinitionService(repository *repository.Repositories) FieldSetDefinitionService {
	return &fieldSetDefinitionService{
		repository: repository,
	}
}

func (s *fieldSetDefinitionService) FindAll(entityDefinitionId string) (*entity.FieldSetDefinitionEntities, error) {
	all, err := s.repository.FieldSetDefinitionRepository.FindAllByEntityDefinitionId(entityDefinitionId)
	if err != nil {
		return nil, err
	}
	fieldSetDefinitionEntities := entity.FieldSetDefinitionEntities{}
	for _, dbRecord := range all.([]*db.Record) {
		fieldSetDefinitionEntities = append(fieldSetDefinitionEntities, *s.mapDbNodeToFieldSetDefinition(dbRecord.Values[0].(dbtype.Node)))
	}
	return &fieldSetDefinitionEntities, nil
}

func (s *fieldSetDefinitionService) FindLinkedWithFieldSet(ctx context.Context, fieldSetId string) (*entity.FieldSetDefinitionEntity, error) {
	queryResult, err := s.repository.FieldSetDefinitionRepository.FindByFieldSetId(fieldSetId)
	if err != nil {
		return nil, err
	}
	if len(queryResult.([]*db.Record)) == 0 {
		return nil, nil
	}
	return s.mapDbNodeToFieldSetDefinition((queryResult.([]*db.Record))[0].Values[0].(dbtype.Node)), nil
}

func (s *fieldSetDefinitionService) mapDbNodeToFieldSetDefinition(dbNode dbtype.Node) *entity.FieldSetDefinitionEntity {
	props := utils.GetPropsFromNode(dbNode)
	fieldSetDefinition := entity.FieldSetDefinitionEntity{
		Id:    utils.GetStringPropOrEmpty(props, "id"),
		Name:  utils.GetStringPropOrEmpty(props, "name"),
		Order: utils.GetIntPropOrMinusOne(props, "order"),
	}
	return &fieldSetDefinition
}
