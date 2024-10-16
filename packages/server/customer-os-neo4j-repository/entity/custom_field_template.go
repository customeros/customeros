package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"time"
)

type CustomFieldTemplateEntity struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Type      string // TODO convert to enum
	Order     int64
	Mandatory bool
	Length    *int64
	Min       *int64
	Max       *int64
}

type CustomFieldTemplateEntities []CustomFieldTemplateEntity

func (cft CustomFieldTemplateEntity) EntityLabel() []string {
	return []string{model.NodeLabelCustomFieldTemplate}
}
