package entity

import (
	"fmt"
)

type ContactGroupEntity struct {
	Id            string
	Name          string     `neo4jDb:"property:name;lookupName:NAME;supportCaseSensitive:true"`
	Source        DataSource `neo4jDb:"property:source;lookupName:SOURCE;supportCaseSensitive:false"`
	SourceOfTruth DataSource `neo4jDb:"property:sourceOfTruth;lookupName:SOURCE_OF_TRUTH;supportCaseSensitive:false"`
}

func (contactGroup ContactGroupEntity) ToString() string {
	return fmt.Sprintf("id: %s\nname: %s", contactGroup.Id, contactGroup.Name)
}

type ContactGroupEntities []ContactGroupEntity

func (contactGroup ContactGroupEntity) Labels() []string {
	return []string{"ContactGroup"}
}
