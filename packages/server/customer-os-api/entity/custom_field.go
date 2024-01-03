package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"time"
)

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

type CustomFieldEntity struct {
	Id            *string
	Name          string
	DataType      string
	Value         model.AnyTypeValue
	TemplateId    *string
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CustomFieldEntities []CustomFieldEntity

func (f *CustomFieldEntity) NodeLabel() string {
	switch f.DataType {
	case model.CustomFieldDataTypeText.String():
		return textNodeLabel.String()
	case model.CustomFieldDataTypeInteger.String():
		return intNodeLabel.String()
	case model.CustomFieldDataTypeDecimal.String():
		return floatNodeLabel.String()
	case model.CustomFieldDataTypeDatetime.String():
		return timeNodeLabel.String()
	case model.CustomFieldDataTypeBool.String():
		return boolNodeLabel.String()
	}
	return ""
}

func (f *CustomFieldEntity) PropertyName() string {
	switch f.DataType {
	case model.CustomFieldDataTypeText.String():
		return CustomFieldTextProperty.String()
	case model.CustomFieldDataTypeInteger.String():
		return CustomFieldIntProperty.String()
	case model.CustomFieldDataTypeDecimal.String():
		return CustomFieldFloatProperty.String()
	case model.CustomFieldDataTypeDatetime.String():
		return CustomFieldTimeProperty.String()
	case model.CustomFieldDataTypeBool.String():
		return CustomFieldBoolProperty.String()
	}
	return ""
}

func (f *CustomFieldEntity) ToString() string {
	return fmt.Sprintf("id: %v\nname: %s\nvalue: %s", f.Id, f.Name, f.Value.RealValue())
}

func (f *CustomFieldEntity) AdjustValueByDatatype() {
	switch f.DataType {
	case model.CustomFieldDataTypeText.String():
		if f.Value.Str == nil {
			if f.Value.Time != nil {
				f.Value.TimeToStr()
			} else if f.Value.Bool != nil {
				f.Value.BoolToStr()
			} else if f.Value.Int != nil {
				f.Value.IntToStr()
			} else if f.Value.Float != nil {
				f.Value.FloatToStr()
			}
		}
	case model.CustomFieldDataTypeDatetime.String():
		if f.Value.Time == nil {
			f.Value.StrToTime()
		}
	case model.CustomFieldDataTypeBool.String():
		if f.Value.Bool == nil {
			if f.Value.Str != nil {
				f.Value.StrToBool()
			} else if f.Value.Int != nil {
				f.Value.IntToBool()
			}
		}
	case model.CustomFieldDataTypeInteger.String():
		if f.Value.Int == nil {
			if f.Value.Str != nil {
				f.Value.StrToInt()
			} else if f.Value.Float != nil {
				f.Value.FloatToInt()
			}
		}
	case model.CustomFieldDataTypeDecimal.String():
		if f.Value.Float == nil {
			if f.Value.Str != nil {
				f.Value.StrToFloat()
			} else if f.Value.Int != nil {
				f.Value.IntToFloat()
			}
		}
	}
}

func (f *CustomFieldEntity) Labels() []string {
	return []string{"CustomField", f.NodeLabel()}
}
