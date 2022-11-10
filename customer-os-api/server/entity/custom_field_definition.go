package entity

import (
	"fmt"
)

type CustomFieldDefinitionEntity struct {
	Id        string
	Name      string
	Type      string
	Order     int64
	Mandatory bool
	Length    *int64
	Min       *int64
	Max       *int64
}

func (definition CustomFieldDefinitionEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s\ntype: %s", definition.Id, definition.Name, definition.Type)
}

type CustomFieldDefinitionEntities []CustomFieldDefinitionEntity
