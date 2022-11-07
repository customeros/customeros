package entity

import (
	"fmt"
	"time"
)

type FieldsSetEntity struct {
	Id    string
	Type  string
	Name  string
	Added time.Time
}

func (set FieldsSetEntity) ToString() string {
	return fmt.Sprintf("id: %s\ntype: %s\nname: %s", set.Id, set.Type, set.Name)
}

type FieldsSetEntities []FieldsSetEntity
