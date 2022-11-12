package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapCustomFieldInputsToEntities(inputs []*model.CustomFieldInput) *entity.CustomFieldEntities {
	if inputs == nil {
		return nil
	}
	var result entity.CustomFieldEntities
	for _, singleInput := range inputs {
		result = append(result, *MapCustomFieldInputToEntity(singleInput))
	}
	return &result
}

func MapCustomFieldInputToEntity(input *model.CustomFieldInput) *entity.CustomFieldEntity {
	customFieldEntity := entity.CustomFieldEntity{
		Name:         input.Name,
		Value:        input.Value,
		DataType:     input.Datatype.String(),
		DefinitionId: input.DefinitionID,
	}
	customFieldEntity.AdjustValueByDatatype()
	return &customFieldEntity
}

func MapTextCustomFieldUpdateInputToEntity(input *model.CustomFieldUpdateInput) *entity.CustomFieldEntity {
	textCustomFieldEntity := entity.CustomFieldEntity{
		Id:   input.ID,
		Name: input.Name,
		//TODO alexb implement update custom fields,
		//Value: input.Value,
	}
	return &textCustomFieldEntity
}

func MapEntitiesToTextCustomFields(textCustomFieldEntities *entity.CustomFieldEntities) []*model.CustomField {
	var textCustomFields []*model.CustomField
	for _, textCustomFieldEntity := range *textCustomFieldEntities {
		textCustomFields = append(textCustomFields, MapEntityToTextCustomField(&textCustomFieldEntity))
	}
	return textCustomFields
}

func MapEntityToTextCustomField(textCustomFieldEntity *entity.CustomFieldEntity) *model.CustomField {
	return &model.CustomField{
		ID:   textCustomFieldEntity.Id,
		Name: textCustomFieldEntity.Name,
		//TODO alexb implement
		//Value: textCustomFieldEntity.Value,
	}
}
