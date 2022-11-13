package entity

import (
	"fmt"
)

type TextCustomFieldEntity struct {
	Id           string
	Name         string
	Value        string
	DefinitionId *string
}

func (field TextCustomFieldEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s\nvalue: %s", field.Id, field.Name, field.Value)
}

type TextCustomFieldEntities []TextCustomFieldEntity
