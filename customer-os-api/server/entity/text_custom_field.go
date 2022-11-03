package entity

import (
	"fmt"
)

type TextCustomFieldEntity struct {
	Id    string
	Group string
	Name  string
	Value string
}

func (field TextCustomFieldEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s\nvalue: %s", field.Id, field.Name, field.Value)
}

type TextCustomFieldEntities []TextCustomFieldEntity
