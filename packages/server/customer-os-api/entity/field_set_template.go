package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type FieldSetTemplateEntity struct {
	Id           string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string
	Order        int64
	CustomFields []*neo4jentity.CustomFieldTemplateEntity
}

func (template FieldSetTemplateEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", template.Id, template.Name)
}

type FieldSetTemplateEntities []FieldSetTemplateEntity

func (template FieldSetTemplateEntity) Labels() []string {
	return []string{"FieldSetTemplate"}
}
