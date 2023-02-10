package entity

import (
	"fmt"
	"time"
)

type FieldSetTemplateEntity struct {
	Id           string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string
	Order        int64
	CustomFields []*CustomFieldTemplateEntity
}

func (template FieldSetTemplateEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", template.Id, template.Name)
}

type FieldSetTemplateEntities []FieldSetTemplateEntity

func (template FieldSetTemplateEntity) Labels() []string {
	return []string{"FieldSetTemplate"}
}
