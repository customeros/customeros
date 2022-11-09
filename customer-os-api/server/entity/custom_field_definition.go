package entity

import (
	"fmt"
)

type CustomFieldDefinitionEntity struct {
	Id   string
	Name string
	Type string
}

func (definition CustomFieldDefinitionEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s\ntype: %s", definition.Id, definition.Name, definition.Type)
}

type CustomFieldDefinitionEntities []CustomFieldDefinitionEntity
