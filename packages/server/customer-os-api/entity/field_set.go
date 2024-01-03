package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type FieldSetEntity struct {
	Id            *string
	Name          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	TemplateId    *string
	CustomFields  *CustomFieldEntities
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
}

func (set FieldSetEntity) ToString() string {
	return fmt.Sprintf("id: %v\nname: %s", set.Id, set.Name)
}

type FieldSetEntities []FieldSetEntity

func (set FieldSetEntity) Labels() []string {
	return []string{"FieldSet"}
}
