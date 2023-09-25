package entity

import "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/models"

type nodeLabel string

const (
	textNodeLabel  nodeLabel = "TextField"
	intNodeLabel   nodeLabel = "IntField"
	floatNodeLabel nodeLabel = "FloatField"
	boolNodeLabel  nodeLabel = "BoolField"
	timeNodeLabel  nodeLabel = "TimeField"
)

func (l nodeLabel) String() string {
	return string(l)
}

type propertyName string

const (
	CustomFieldTextProperty  propertyName = "textValue"
	CustomFieldIntProperty   propertyName = "intValue"
	CustomFieldFloatProperty propertyName = "floatValue"
	CustomFieldBoolProperty  propertyName = "boolValue"
	CustomFieldTimeProperty  propertyName = "timeValue"
)

func (p propertyName) String() string {
	return string(p)
}

func NodeLabelForCustomFieldDataType(dataType string) string {
	switch dataType {
	case string(models.CustomFieldDataTypeText):
		return textNodeLabel.String()
	case string(models.CustomFieldDataTypeInteger):
		return intNodeLabel.String()
	case string(models.CustomFieldDataTypeDecimal):
		return floatNodeLabel.String()
	case string(models.CustomFieldDataTypeDatetime):
		return timeNodeLabel.String()
	case string(models.CustomFieldDataTypeBool):
		return boolNodeLabel.String()
	}
	return ""
}

func PropertyNameForCustomFieldDataType(dataType string) string {
	switch dataType {
	case string(models.CustomFieldDataTypeText):
		return CustomFieldTextProperty.String()
	case string(models.CustomFieldDataTypeInteger):
		return CustomFieldIntProperty.String()
	case string(models.CustomFieldDataTypeDecimal):
		return CustomFieldFloatProperty.String()
	case string(models.CustomFieldDataTypeDatetime):
		return CustomFieldTimeProperty.String()
	case string(models.CustomFieldDataTypeBool):
		return CustomFieldBoolProperty.String()
	}
	return ""
}
