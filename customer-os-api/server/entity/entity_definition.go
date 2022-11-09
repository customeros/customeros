package entity

import (
	"fmt"
)

type EntityDefinitionEntity struct {
	Id           string
	Name         string
	Extends      *string
	Version      int64
	CustomFields []CustomFieldDefinitionEntity
	FieldSets    []FieldSetDefinitionEntity
}

func (entityDefinition EntityDefinitionEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s\nextends: %s", entityDefinition.Id, entityDefinition.Name, entityDefinition.Extends)
}

type EntityDefinitionEntities []EntityDefinitionEntity
