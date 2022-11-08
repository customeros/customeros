package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapFieldSetInputToEntity(input *model.FieldSetInput) *entity.FieldSetEntity {
	fieldSetEntity := entity.FieldSetEntity{
		Type: input.Type,
		Name: input.Name,
	}
	return &fieldSetEntity
}

func MapFieldSetUpdateInputToEntity(input *model.FieldSetUpdateInput) *entity.FieldSetEntity {
	fieldSetEntity := entity.FieldSetEntity{
		Id:   input.ID,
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
		ID:    fieldSetEntity.Id,
		Type:  fieldSetEntity.Type,
		Name:  fieldSetEntity.Name,
		Added: fieldSetEntity.Added,
	}
}
