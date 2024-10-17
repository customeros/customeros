package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"time"
)

type CustomFieldTemplateProperty string

const (
	CustomFieldTemplatePropertyId          CustomFieldTemplateProperty = "id"
	CustomFieldTemplatePropertyName        CustomFieldTemplateProperty = "name"
	CustomFieldTemplatePropertyEntityType  CustomFieldTemplateProperty = "entityType"
	CustomFieldTemplatePropertyType        CustomFieldTemplateProperty = "type"
	CustomFieldTemplatePropertyValidValues CustomFieldTemplateProperty = "validValues"
	CustomFieldTemplatePropertyOrder       CustomFieldTemplateProperty = "order"
	CustomFieldTemplatePropertyRequired    CustomFieldTemplateProperty = "required"
	CustomFieldTemplatePropertyLength      CustomFieldTemplateProperty = "length"
	CustomFieldTemplatePropertyMin         CustomFieldTemplateProperty = "min"
	CustomFieldTemplatePropertyMax         CustomFieldTemplateProperty = "max"
	CustomFieldTemplatePropertyCreatedAt   CustomFieldTemplateProperty = "createdAt"
	CustomFieldTemplatePropertyUpdatedAt   CustomFieldTemplateProperty = "updatedAt"
)

type CustomFieldTemplateEntity struct {
	DataLoaderKey
	Id          string
	Name        string
	EntityType  model.EntityType
	Type        string
	ValidValues []string
	Order       *int64
	Required    *bool
	Length      *int64
	Min         *int64
	Max         *int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CustomFieldTemplateEntities []CustomFieldTemplateEntity

func (cft CustomFieldTemplateEntity) EntityLabel() []string {
	return []string{model.NodeLabelCustomFieldTemplate}
}
