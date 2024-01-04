package entity

type CustomFieldDataType string

const (
	CustomFieldDataTypeText     CustomFieldDataType = "TEXT"
	CustomFieldDataTypeBool     CustomFieldDataType = "BOOL"
	CustomFieldDataTypeDatetime CustomFieldDataType = "DATETIME"
	CustomFieldDataTypeInteger  CustomFieldDataType = "INTEGER"
	CustomFieldDataTypeDecimal  CustomFieldDataType = "DECIMAL"
)

func (cfdt CustomFieldDataType) String() string {
	return string(cfdt)
}

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
	case string(CustomFieldDataTypeText):
		return textNodeLabel.String()
	case string(CustomFieldDataTypeInteger):
		return intNodeLabel.String()
	case string(CustomFieldDataTypeDecimal):
		return floatNodeLabel.String()
	case string(CustomFieldDataTypeDatetime):
		return timeNodeLabel.String()
	case string(CustomFieldDataTypeBool):
		return boolNodeLabel.String()
	}
	return ""
}

func PropertyNameForCustomFieldDataType(dataType string) string {
	switch dataType {
	case string(CustomFieldDataTypeText):
		return CustomFieldTextProperty.String()
	case string(CustomFieldDataTypeInteger):
		return CustomFieldIntProperty.String()
	case string(CustomFieldDataTypeDecimal):
		return CustomFieldFloatProperty.String()
	case string(CustomFieldDataTypeDatetime):
		return CustomFieldTimeProperty.String()
	case string(CustomFieldDataTypeBool):
		return CustomFieldBoolProperty.String()
	}
	return ""
}
