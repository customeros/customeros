package entity

import (
	"fmt"
	"time"
)

type FieldSetEntity struct {
	Id    string
	Type  string
	Name  string
	Added time.Time
}

func (set FieldSetEntity) ToString() string {
	return fmt.Sprintf("id: %s\ntype: %s\nname: %s", set.Id, set.Type, set.Name)
}

type FieldSetEntities []FieldSetEntity
