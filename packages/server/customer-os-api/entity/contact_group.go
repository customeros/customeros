package entity

import (
	"fmt"
)

type ContactGroupEntity struct {
	Id   string
	Name string `neo4jDb:"property:name;exposedName:NAME;supportCaseSensitive:true"`
}

func (contactGroupEntity ContactGroupEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", contactGroupEntity.Id, contactGroupEntity.Name)
}

type ContactGroupEntities []ContactGroupEntity
