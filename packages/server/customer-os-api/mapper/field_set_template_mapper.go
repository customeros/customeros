package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapFieldSetTemplateInputToEntity(input model.FieldSetTemplateInput) *entity.FieldSetTemplateEntity {
	templateEntity := entity.FieldSetTemplateEntity{
		Name:  input.Name,
		Order: int64(input.Order),
	}
	for _, v := range input.CustomFields {
		templateEntity.CustomFields = append(templateEntity.CustomFields, MapCustomFieldTemplateInputToEntity(*v))
	}
	return &templateEntity
}

func MapEntityToFieldSetTemplate(entity *entity.FieldSetTemplateEntity) *model.FieldSetTemplate {
	output := model.FieldSetTemplate{
		ID:        entity.Id,
		CreatedAt: entity.CreatedAt,
		Name:      entity.Name,
		Order:     int(entity.Order),
	}
	return &output
}

func MapEntitiesToFieldSetTemplates(entities *entity.FieldSetTemplateEntities) []*model.FieldSetTemplate {
	var fieldSetTemplates []*model.FieldSetTemplate
	for _, v := range *entities {
		fieldSetTemplates = append(fieldSetTemplates, MapEntityToFieldSetTemplate(&v))
	}
	return fieldSetTemplates
}
