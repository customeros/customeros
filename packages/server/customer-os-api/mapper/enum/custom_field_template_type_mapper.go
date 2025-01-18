package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func MapTemplateTypeToFieldDataType(templateType string) *model.CustomFieldDataType {
	switch templateType {
	case model.CustomFieldTemplateTypeFreeText.String():
		return utils.ToPtr(model.CustomFieldDataTypeText)
	case model.CustomFieldTemplateTypeSingleSelect.String():
		return utils.ToPtr(model.CustomFieldDataTypeText)
	case model.CustomFieldTemplateTypeNumber.String():
		return utils.ToPtr(model.CustomFieldDataTypeInteger)
	default:
		return utils.ToPtr(model.CustomFieldDataTypeText)
	}
}
