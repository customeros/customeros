package entity

import (
	"fmt"
	"time"
)

type CustomFieldTemplateEntity struct {
	Id        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Type      string
	Order     int64
	Mandatory bool
	Length    *int64
	Min       *int64
	Max       *int64
}

func (template CustomFieldTemplateEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s\ntype: %s", template.Id, template.Name, template.Type)
}

type CustomFieldTemplateEntities []CustomFieldTemplateEntity

func (template CustomFieldTemplateEntity) Labels() []string {
	return []string{"CustomFieldTemplate"}
}
