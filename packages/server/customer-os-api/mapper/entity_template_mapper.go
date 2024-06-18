package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapEntityTemplateInputToEntity(input model.EntityTemplateInput) *neo4jentity.EntityTemplateEntity {
	templateEntity := neo4jentity.EntityTemplateEntity{
		Name: input.Name,
	}
	if input.Extends != nil {
		templateEntity.Extends = input.Extends.String()
	}
	for _, v := range input.CustomFieldTemplateInputs {
		templateEntity.CustomFields = append(templateEntity.CustomFields, MapCustomFieldTemplateInputToEntity(*v))
	}
	return &templateEntity
}

func MapEntityToEntityTemplate(entity *neo4jentity.EntityTemplateEntity) *model.EntityTemplate {
	output := model.EntityTemplate{
		ID:          entity.Id,
		Name:        entity.Name,
		Version:     int(entity.Version),
		Created:     entity.CreatedAt,
		LastUpdated: entity.UpdatedAt,
	}
	if entity.Extends != "" {
		extends := model.EntityTemplateExtension(entity.Extends)
		if extends.IsValid() {
			output.Extends = &extends
		}
	}
	return &output
}

func MapEntitiesToEntityTemplates(entities *neo4jentity.EntityTemplateEntities) []*model.EntityTemplate {
	var entityTemplates []*model.EntityTemplate
	for _, v := range *entities {
		entityTemplates = append(entityTemplates, MapEntityToEntityTemplate(&v))
	}
	return entityTemplates
}
