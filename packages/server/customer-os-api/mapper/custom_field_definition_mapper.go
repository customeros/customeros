package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapCustomFieldDefinitionInputToEntity(input model.CustomFieldDefinitionInput) *entity.CustomFieldDefinitionEntity {
	definitionEntity := entity.CustomFieldDefinitionEntity{
		Name:      input.Name,
		Type:      input.Type.String(),
		Order:     int64(input.Order),
		Mandatory: input.Mandatory,
		Length:    utils.IntPtrToInt64Ptr(input.Length),
		Min:       utils.IntPtrToInt64Ptr(input.Min),
		Max:       utils.IntPtrToInt64Ptr(input.Max),
	}
	return &definitionEntity
}

func MapEntityToCustomFieldDefinition(entity *entity.CustomFieldDefinitionEntity) *model.CustomFieldDefinition {
	fieldType := model.CustomFieldDefinitionType(entity.Type)
	if !fieldType.IsValid() {
		fieldType = model.CustomFieldDefinitionTypeText
	}
	output := model.CustomFieldDefinition{
		ID:        entity.Id,
		Name:      entity.Name,
		Type:      fieldType,
		Order:     int(entity.Order),
		Mandatory: entity.Mandatory,
		Length:    utils.Int64PtrToIntPtr(entity.Length),
		Min:       utils.Int64PtrToIntPtr(entity.Min),
		Max:       utils.Int64PtrToIntPtr(entity.Max),
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
