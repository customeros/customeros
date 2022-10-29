package entity

import (
	"fmt"
)

type TextCustomFieldEntity struct {
	Group string
	Name  string
	Value string
}

func (field TextCustomFieldEntity) ToString() string {
	return fmt.Sprintf("name: %s\nvalue: %s", field.Name, field.Value)
}

type TextCustomFieldEntities []TextCustomFieldEntity
