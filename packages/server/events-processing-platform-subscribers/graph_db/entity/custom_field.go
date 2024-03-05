package entity

import "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/model"

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
	case string(model.CustomFieldDataTypeText):
		return textNodeLabel.String()
	case string(model.CustomFieldDataTypeInteger):
		return intNodeLabel.String()
	case string(model.CustomFieldDataTypeDecimal):
		return floatNodeLabel.String()
	case string(model.CustomFieldDataTypeDatetime):
		return timeNodeLabel.String()
	case string(model.CustomFieldDataTypeBool):
		return boolNodeLabel.String()
	}
	return ""
}

func PropertyNameForCustomFieldDataType(dataType string) string {
	switch dataType {
	case string(model.CustomFieldDataTypeText):
		return CustomFieldTextProperty.String()
	case string(model.CustomFieldDataTypeInteger):
		return CustomFieldIntProperty.String()
	case string(model.CustomFieldDataTypeDecimal):
		return CustomFieldFloatProperty.String()
	case string(model.CustomFieldDataTypeDatetime):
		return CustomFieldTimeProperty.String()
	case string(model.CustomFieldDataTypeBool):
		return CustomFieldBoolProperty.String()
	}
	return ""
}
