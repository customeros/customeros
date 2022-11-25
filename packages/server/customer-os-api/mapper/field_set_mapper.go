package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

func MapFieldSetInputsToEntities(inputs []*model.FieldSetInput) *entity.FieldSetEntities {
	if inputs == nil {
		return nil
	}
	var result entity.FieldSetEntities
	for _, singleInput := range inputs {
		result = append(result, *MapFieldSetInputToEntity(singleInput))
	}
	return &result
}

func MapFieldSetInputToEntity(input *model.FieldSetInput) *entity.FieldSetEntity {
	fieldSetEntity := entity.FieldSetEntity{
		Id:           input.ID,
		Name:         input.Name,
		DefinitionId: input.DefinitionID,
		CustomFields: MapCustomFieldInputsToEntities(input.CustomFields),
	}
	return &fieldSetEntity
}

func MapFieldSetUpdateInputToEntity(input *model.FieldSetUpdateInput) *entity.FieldSetEntity {
	fieldSetEntity := entity.FieldSetEntity{
		Id:   utils.StringPtr(input.ID),
		Name: input.Name,
	}
	return &fieldSetEntity
}

func MapEntitiesToFieldSets(fieldSetEntities *entity.FieldSetEntities) []*model.FieldSet {
	var fieldSet []*model.FieldSet
	for _, fieldSetEntity := range *fieldSetEntities {
		fieldSet = append(fieldSet, MapEntityToFieldSet(&fieldSetEntity))
	}
	return fieldSet
}

func MapEntityToFieldSet(fieldSetEntity *entity.FieldSetEntity) *model.FieldSet {
	return &model.FieldSet{
		ID:    *fieldSetEntity.Id,
		Name:  fieldSetEntity.Name,
		Added: fieldSetEntity.Added,
	}
}
