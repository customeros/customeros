package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type CustomFieldDefinitionService interface {
	FindAllForEntityDefinition(entityDefinitionId string) (*entity.CustomFieldDefinitionEntities, error)
	FindAllForFieldSetDefinition(fieldSetDefinitionId string) (*entity.CustomFieldDefinitionEntities, error)
	FindLinkedWithCustomField(ctx context.Context, customFieldId string) (*entity.CustomFieldDefinitionEntity, error)
}

type customFieldDefinitionService struct {
	repository *repository.Repositories
}

func NewCustomFieldDefinitionService(repository *repository.Repositories) CustomFieldDefinitionService {
	return &customFieldDefinitionService{
		repository: repository,
	}
}

func (s *customFieldDefinitionService) FindAllForEntityDefinition(entityDefinitionId string) (*entity.CustomFieldDefinitionEntities, error) {
	all, err := s.repository.CustomFieldDefinitionRepository.FindAllByEntityDefinitionId(entityDefinitionId)
	if err != nil {
		return nil, err
	}
	customFieldDefinitionEntities := entity.CustomFieldDefinitionEntities{}
	for _, dbRecord := range all.([]*db.Record) {
		customFieldDefinitionEntities = append(customFieldDefinitionEntities, *s.mapDbNodeToCustomFieldDefinition(dbRecord.Values[0].(dbtype.Node)))
	}
	return &customFieldDefinitionEntities, nil
}

func (s *customFieldDefinitionService) FindAllForFieldSetDefinition(fieldSetDefinitionId string) (*entity.CustomFieldDefinitionEntities, error) {
	all, err := s.repository.CustomFieldDefinitionRepository.FindAllByEntityFieldSetDefinitionId(fieldSetDefinitionId)
	if err != nil {
		return nil, err
	}
	customFieldDefinitionEntities := entity.CustomFieldDefinitionEntities{}
	for _, dbRecord := range all.([]*db.Record) {
		customFieldDefinitionEntities = append(customFieldDefinitionEntities, *s.mapDbNodeToCustomFieldDefinition(dbRecord.Values[0].(dbtype.Node)))
	}
	return &customFieldDefinitionEntities, nil
}

func (s *customFieldDefinitionService) FindLinkedWithCustomField(ctx context.Context, customFieldId string) (*entity.CustomFieldDefinitionEntity, error) {
	queryResult, err := s.repository.CustomFieldDefinitionRepository.FindByCustomFieldId(customFieldId)
	if err != nil {
		return nil, err
	}
	if len(queryResult.([]*db.Record)) == 0 {
		return nil, nil
	}
	return s.mapDbNodeToCustomFieldDefinition((queryResult.([]*db.Record))[0].Values[0].(dbtype.Node)), nil
}

func (s *customFieldDefinitionService) mapDbNodeToCustomFieldDefinition(dbNode dbtype.Node) *entity.CustomFieldDefinitionEntity {
	props := utils.GetPropsFromNode(dbNode)
	customFieldDefinition := entity.CustomFieldDefinitionEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		Order:     utils.GetIntPropOrMinusOne(props, "order"),
		Mandatory: utils.GetBoolPropOrFalse(props, "mandatory"),
		Type:      utils.GetStringPropOrEmpty(props, "type"),
		Length:    utils.GetIntPropOrNil(props, "length"),
		Min:       utils.GetIntPropOrNil(props, "min"),
		Max:       utils.GetIntPropOrNil(props, "max"),
	}
	return &customFieldDefinition
}
