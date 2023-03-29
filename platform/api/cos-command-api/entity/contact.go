package entity

import (
	"fmt"
	"time"
)

type ContactEntity struct {
	Id        string
	Title     string `neo4jDb:"property:title;lookupName:TITLE;supportCaseSensitive:false"`
	FirstName string `neo4jDb:"property:firstName;lookupName:FIRST_NAME;supportCaseSensitive:true"`
	LastName  string `neo4jDb:"property:lastName;lookupName:LAST_NAME;supportCaseSensitive:true"`
	Label     string `neo4jDb:"property:label;lookupName:LABEL;supportCaseSensitive:true"`
	// TODO make non pointer and use different property for input
	CreatedAt     *time.Time `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	Source        DataSource `neo4jDb:"property:source;lookupName:SOURCE;supportCaseSensitive:false"`
	SourceOfTruth DataSource `neo4jDb:"property:sourceOfTruth;lookupName:SOURCE_OF_TRUTH;supportCaseSensitive:false"`
}

func (contact ContactEntity) ToString() string {
	return fmt.Sprintf("id: %s\nfirstName: %s\nlastName: %s", contact.Id, contact.FirstName, contact.LastName)
}

type ContactEntities []ContactEntity

func (contact ContactEntity) Labels() []string {
	return []string{"Contact"}
}
