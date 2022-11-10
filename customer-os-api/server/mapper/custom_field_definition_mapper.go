package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

func MapCustomFieldDefinitionInputToEntity(input model.CustomFieldDefinitionInput) *entity.CustomFieldDefinitionEntity {
	definitionEntity := entity.CustomFieldDefinitionEntity{
		Name:      input.Name,
		Type:      input.Type.String(),
		Order:     int64(input.Order),
		Mandatory: input.Mandatory,
		Length:    utils.ToInt64Ptr(input.Length),
		Min:       utils.ToInt64Ptr(input.Min),
		Max:       utils.ToInt64Ptr(input.Max),
	}
	return &definitionEntity
}

func MapEntityToCustomFieldDefinition(entity *entity.CustomFieldDefinitionEntity) *model.CustomFieldDefinition {
	fieldType := model.CustomFieldType(entity.Type)
	if !fieldType.IsValid() {
		fieldType = model.CustomFieldTypeText
	}
	output := model.CustomFieldDefinition{
		ID:        entity.Id,
		Name:      entity.Name,
		Type:      fieldType,
		Order:     int(entity.Order),
		Mandatory: entity.Mandatory,
		Length:    utils.ToIntPtr(entity.Length),
		Min:       utils.ToIntPtr(entity.Min),
		Max:       utils.ToIntPtr(entity.Max),
	}
	return &output
}

func MapEntitiesToCustomFieldDefinitions(entities *entity.CustomFieldDefinitionEntities) []*model.CustomFieldDefinition {
	var customFieldDefinitions []*model.CustomFieldDefinition
	for _, v := range *entities {
		customFieldDefinitions = append(customFieldDefinitions, MapEntityToCustomFieldDefinition(&v))
	}
	return customFieldDefinitions
}
