package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
)

func MapCustomFieldTemplateInputToSaveFields(input model.CustomFieldTemplateInput) neo4jrepository.CustomFieldTemplateSaveFields {
	fields := neo4jrepository.CustomFieldTemplateSaveFields{}
	if input.Name != nil {
		fields.Name = *input.Name
		fields.UpdateName = true
	}
	if input.EntityType != nil {
		fields.EntityType = commonmodel.DecodeEntityType(input.EntityType.String())
	}
	if input.Type != nil {
		fields.Type = input.Type.String()
		fields.UpdateType = true
	}
	if input.Order != nil {
		fields.Order = input.Order
		fields.UpdateOrder = true
	}
	if input.Required != nil {
		fields.Required = input.Required
		fields.UpdateRequired = true
	}
	if input.Length != nil {
		fields.Length = input.Length
		fields.UpdateLength = true
	}
	if input.Min != nil {
		fields.Min = input.Min
		fields.UpdateMin = true
	}
	if input.Max != nil {
		fields.Max = input.Max
		fields.UpdateMax = true
	}
	if input.ValidValues != nil {
		fields.ValidValues = input.ValidValues
		fields.UpdateValidValues = true
	}

	return fields
}

func MapEntityToCustomFieldTemplate(entity *neo4jentity.CustomFieldTemplateEntity) *model.CustomFieldTemplate {
	output := model.CustomFieldTemplate{
		ID:         entity.Id,
		Name:       entity.Name,
		EntityType: model.EntityType(entity.EntityType.String()),
		Type:       model.CustomFieldTemplateType(entity.Type),
		Order:      entity.Order,
		Required:   entity.Required,
		Length:     entity.Length,
		Min:        entity.Min,
		Max:        entity.Max,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}
	return &output
}

func MapEntitiesToCustomFieldTemplates(entities *neo4jentity.CustomFieldTemplateEntities) []*model.CustomFieldTemplate {
	var customFieldTemplates []*model.CustomFieldTemplate
	for _, v := range *entities {
		customFieldTemplates = append(customFieldTemplates, MapEntityToCustomFieldTemplate(&v))
	}
	return customFieldTemplates
}
