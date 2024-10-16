package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type EntityTemplateEntity struct {
	Id           string
	Name         string
	Extends      *string
	Version      int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CustomFields []*neo4jentity.CustomFieldTemplateEntity
	FieldSets    []*FieldSetTemplateEntity
}

func (template EntityTemplateEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s\nextends: %s", template.Id, template.Name, *template.Extends)
}

type EntityTemplateEntities []EntityTemplateEntity

func (template EntityTemplateEntity) Labels() []string {
	return []string{"EntityTemplate"}
}
