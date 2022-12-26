package entity

import (
	"fmt"
)

type ContactGroupEntity struct {
	Id   string
	Name string `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
}

func (contactGroup ContactGroupEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", contactGroup.Id, contactGroup.Name)
}

type ContactGroupEntities []ContactGroupEntity

func (contactGroup ContactGroupEntity) Labels() []string {
	return []string{"ContactGroup"}
}
