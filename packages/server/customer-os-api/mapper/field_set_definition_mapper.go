package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapFieldSetDefinitionInputToEntity(input model.FieldSetDefinitionInput) *entity.FieldSetDefinitionEntity {
	definitionEntity := entity.FieldSetDefinitionEntity{
		Name:  input.Name,
		Order: int64(input.Order),
	}
	for _, v := range input.CustomFields {
		definitionEntity.CustomFields = append(definitionEntity.CustomFields, MapCustomFieldDefinitionInputToEntity(*v))
	}
	return &definitionEntity
}

func MapEntityToFieldSetDefinition(entity *entity.FieldSetDefinitionEntity) *model.FieldSetDefinition {
	output := model.FieldSetDefinition{
		ID:    entity.Id,
		Name:  entity.Name,
		Order: int(entity.Order),
	}
	return &output
}

func MapEntitiesToFieldSetDefinitions(entities *entity.FieldSetDefinitionEntities) []*model.FieldSetDefinition {
	var fieldSetDefinitions []*model.FieldSetDefinition
	for _, v := range *entities {
		fieldSetDefinitions = append(fieldSetDefinitions, MapEntityToFieldSetDefinition(&v))
	}
	return fieldSetDefinitions
}
