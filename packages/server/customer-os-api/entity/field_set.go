package entity

import (
	"fmt"
	"time"
)

type FieldSetEntity struct {
	Id           *string
	Name         string
	Added        time.Time
	DefinitionId *string
	CustomFields *CustomFieldEntities
}

func (set FieldSetEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", set.Id, set.Name)
}

type FieldSetEntities []FieldSetEntity
