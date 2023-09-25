package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func MapCustomFieldTemplateInputToEntity(input model.CustomFieldTemplateInput) *entity.CustomFieldTemplateEntity {
	templateEntity := entity.CustomFieldTemplateEntity{
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

func MapEntityToCustomFieldTemplate(entity *entity.CustomFieldTemplateEntity) *model.CustomFieldTemplate {
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

func MapEntitiesToCustomFieldTemplates(entities *entity.CustomFieldTemplateEntities) []*model.CustomFieldTemplate {
	var customFieldTemplates []*model.CustomFieldTemplate
	for _, v := range *entities {
		customFieldTemplates = append(customFieldTemplates, MapEntityToCustomFieldTemplate(&v))
	}
	return customFieldTemplates
}
