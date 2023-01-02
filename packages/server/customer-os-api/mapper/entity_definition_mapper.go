package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapEntityDefinitionInputToEntity(input model.EntityDefinitionInput) *entity.EntityDefinitionEntity {
	definitionEntity := entity.EntityDefinitionEntity{
		Name: input.Name,
	}
	if input.Extends != nil {
		extends := input.Extends.String()
		definitionEntity.Extends = &extends
	}
	for _, v := range input.FieldSets {
		definitionEntity.FieldSets = append(definitionEntity.FieldSets, MapFieldSetDefinitionInputToEntity(*v))
	}
	for _, v := range input.CustomFields {
		definitionEntity.CustomFields = append(definitionEntity.CustomFields, MapCustomFieldDefinitionInputToEntity(*v))
	}
	return &definitionEntity
}

func MapEntityToEntityDefinition(entity *entity.EntityDefinitionEntity) *model.EntityDefinition {
	output := model.EntityDefinition{
		ID:        entity.Id,
		Name:      entity.Name,
		Version:   int(entity.Version),
		CreatedAt: entity.CreatedAt,
	}
	if entity.Extends != nil {
		extends := model.EntityDefinitionExtension(*entity.Extends)
		if extends.IsValid() {
			output.Extends = &extends
		}
	}
	return &output
}

func MapEntitiesToEntityDefinitions(entities *entity.EntityDefinitionEntities) []*model.EntityDefinition {
	var entityDefinitions []*model.EntityDefinition
	for _, v := range *entities {
		entityDefinitions = append(entityDefinitions, MapEntityToEntityDefinition(&v))
	}
	return entityDefinitions
}
