package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func MapEntityTemplateInputToEntity(input model.EntityTemplateInput) *entity.EntityTemplateEntity {
	templateEntity := entity.EntityTemplateEntity{
		Name: input.Name,
	}
	if input.Extends != nil {
		extends := input.Extends.String()
		templateEntity.Extends = &extends
	}
	for _, v := range input.FieldSets {
		templateEntity.FieldSets = append(templateEntity.FieldSets, MapFieldSetTemplateInputToEntity(*v))
	}
	for _, v := range input.CustomFields {
		templateEntity.CustomFields = append(templateEntity.CustomFields, MapCustomFieldTemplateInputToEntity(*v))
	}
	return &templateEntity
}

func MapEntityToEntityTemplate(entity *entity.EntityTemplateEntity) *model.EntityTemplate {
	output := model.EntityTemplate{
		ID:        entity.Id,
		Name:      entity.Name,
		Version:   int(entity.Version),
		CreatedAt: entity.CreatedAt,
	}
	if entity.Extends != nil {
		extends := model.EntityTemplateExtension(*entity.Extends)
		if extends.IsValid() {
			output.Extends = &extends
		}
	}
	return &output
}

func MapEntitiesToEntityTemplates(entities *entity.EntityTemplateEntities) []*model.EntityTemplate {
	var entityTemplates []*model.EntityTemplate
	for _, v := range *entities {
		entityTemplates = append(entityTemplates, MapEntityToEntityTemplate(&v))
	}
	return entityTemplates
}
