package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

func MapCustomFieldTemplateInputToEntity(input model.CustomFieldTemplateInput) *neo4jentity.CustomFieldTemplateEntity {
	templateEntity := neo4jentity.CustomFieldTemplateEntity{
		Name:      input.Name,
		Type:      input.Type.String(),
		Order:     int64(input.Order),
		Mandatory: utils.IfNotNilBool(input.Mandatory),
		Length:    utils.IntPtrToInt64Ptr(input.Length),
		Min:       utils.IntPtrToInt64Ptr(input.Min),
		Max:       utils.IntPtrToInt64Ptr(input.Max),
	}
	return &templateEntity
}

func MapEntityToCustomFieldTemplate(entity *neo4jentity.CustomFieldTemplateEntity) *model.CustomFieldTemplate {
	fieldType := model.CustomFieldTemplateType(entity.Type)
	if !fieldType.IsValid() {
		fieldType = model.CustomFieldTemplateTypeText
	}
	output := model.CustomFieldTemplate{
		ID:        entity.Id,
		Name:      entity.Name,
		Type:      fieldType,
		Order:     int(entity.Order),
		Mandatory: entity.Mandatory,
		Length:    utils.Int64PtrToIntPtr(entity.Length),
		Min:       utils.Int64PtrToIntPtr(entity.Min),
		Max:       utils.Int64PtrToIntPtr(entity.Max),
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
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

func MapTemplateTypeToFieldDataType(templateType string) *model.CustomFieldDataType {
	switch templateType {
	case model.CustomFieldTemplateTypeText.String():
		return utils.ToPtr(model.CustomFieldDataTypeText)
	case model.CustomFieldTemplateTypeLink.String():
		return utils.ToPtr(model.CustomFieldDataTypeText)
	default:
		return utils.ToPtr(model.CustomFieldDataTypeText)
	}
}
