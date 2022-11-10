package service

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type FieldSetDefinitionService interface {
	FindAll(entityDefinitionId string) (*entity.FieldSetDefinitionEntities, error)
}

type fieldSetDefinitionService struct {
	repository *repository.RepositoryContainer
}

func NewFieldSetDefinitionService(repository *repository.RepositoryContainer) FieldSetDefinitionService {
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

func (s *fieldSetDefinitionService) mapDbNodeToFieldSetDefinition(dbNode dbtype.Node) *entity.FieldSetDefinitionEntity {
	props := utils.GetPropsFromNode(dbNode)
	fieldSetDefinition := entity.FieldSetDefinitionEntity{
		Id:    utils.GetStringPropOrEmpty(props, "id"),
		Name:  utils.GetStringPropOrEmpty(props, "name"),
		Order: utils.GetIntPropOrDefault(props, "order"),
	}
	return &fieldSetDefinition
}
