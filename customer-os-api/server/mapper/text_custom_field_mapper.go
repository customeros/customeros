package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapTextCustomFieldInputsToEntities(inputs []*model.TextCustomFieldInput) *entity.TextCustomFieldEntities {
	if inputs == nil {
		return nil
	}
	var result entity.TextCustomFieldEntities
	for _, singleInput := range inputs {
		result = append(result, *MapTextCustomFieldInputToEntity(singleInput))
	}
	return &result
}

func MapTextCustomFieldInputToEntity(input *model.TextCustomFieldInput) *entity.TextCustomFieldEntity {
	textCustomFieldEntity := entity.TextCustomFieldEntity{
		Name:  input.Name,
		Value: input.Value,
	}
	return &textCustomFieldEntity
}

func MapTextCustomFieldUpdateInputToEntity(input *model.TextCustomFieldUpdateInput) *entity.TextCustomFieldEntity {
	textCustomFieldEntity := entity.TextCustomFieldEntity{
		Id:    input.ID,
		Name:  input.Name,
		Value: input.Value,
	}
	return &textCustomFieldEntity
}

func MapEntitiesToTextCustomFields(textCustomFieldEntities *entity.TextCustomFieldEntities) []*model.CustomField {
	var textCustomFields []*model.CustomField
	for _, textCustomFieldEntity := range *textCustomFieldEntities {
		textCustomFields = append(textCustomFields, MapEntityToTextCustomField(&textCustomFieldEntity))
	}
	return textCustomFields
}

func MapEntityToTextCustomField(textCustomFieldEntity *entity.TextCustomFieldEntity) *model.CustomField {
	return &model.CustomField{
		ID:    textCustomFieldEntity.Id,
		Name:  textCustomFieldEntity.Name,
		Value: textCustomFieldEntity.Value,
	}
}
