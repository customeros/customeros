package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapTextCustomFieldInputsToEntities(inputs []*model.TextCustomFieldInput) *entity.TextCustomFieldEntities {
	var result entity.TextCustomFieldEntities
	for _, singleInput := range inputs {
		result = append(result, *MapTextCustomFieldInputToEntity(*singleInput))
	}
	return &result
}

func MapTextCustomFieldInputToEntity(input model.TextCustomFieldInput) *entity.TextCustomFieldEntity {
	textCustomFieldEntity := entity.TextCustomFieldEntity{
		Name:  input.Name,
		Value: input.Value,
	}
	if input.Group != nil {
		textCustomFieldEntity.Group = *input.Group
	}
	return &textCustomFieldEntity
}
