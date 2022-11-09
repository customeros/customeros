package entity

import (
	"fmt"
)

type FieldSetDefinitionEntity struct {
	Id           string
	Name         string
	CustomFields []CustomFieldDefinitionEntity
}

func (definition FieldSetDefinitionEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", definition.Id, definition.Name)
}

type FieldSetDefinitionEntities []FieldSetDefinitionEntity
