package entity

import (
	"fmt"
	"time"
)

type FieldSetEntity struct {
	Id            *string
	Name          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	TemplateId    *string
	Source        DataSource
	SourceOfTruth DataSource
}

func (set FieldSetEntity) ToString() string {
	return fmt.Sprintf("id: %v\nname: %s", set.Id, set.Name)
}

type FieldSetEntities []FieldSetEntity

func (set FieldSetEntity) Labels() []string {
	return []string{"FieldSet"}
}
