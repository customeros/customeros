package mapper

import (
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

func MapFieldsSetInputToEntity(input *model.FieldsSetInput) *entity.FieldsSetEntity {
	fieldsSetEntity := entity.FieldsSetEntity{
		Type: input.Type,
		Name: input.Name,
	}
	return &fieldsSetEntity
}

func MapFieldsSetUpdateInputToEntity(input *model.FieldsSetUpdateInput) *entity.FieldsSetEntity {
	fieldsSetEntity := entity.FieldsSetEntity{
		Id:   input.ID,
		Name: input.Name,
	}
	return &fieldsSetEntity
}

func MapEntitiesToFieldsSets(fieldsSetEntities *entity.FieldsSetEntities) []*model.FieldsSet {
	var fieldsSet []*model.FieldsSet
	for _, fieldsSetEntity := range *fieldsSetEntities {
		fieldsSet = append(fieldsSet, MapEntityToFieldsSet(&fieldsSetEntity))
	}
	return fieldsSet
}

func MapEntityToFieldsSet(fieldsSetEntity *entity.FieldsSetEntity) *model.FieldsSet {
	return &model.FieldsSet{
		ID:    fieldsSetEntity.Id,
		Type:  fieldsSetEntity.Type,
		Name:  fieldsSetEntity.Name,
		Added: fieldsSetEntity.Added,
	}
}
