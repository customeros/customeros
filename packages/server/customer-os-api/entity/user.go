package entity

import (
	"fmt"
	"time"
)

type UserEntity struct {
	Id            string
	FirstName     string     `neo4jDb:"property:firstName;lookupName:FIRST_NAME;supportCaseSensitive:true"`
	LastName      string     `neo4jDb:"property:lastName;lookupName:LAST_NAME;supportCaseSensitive:true"`
	CreatedAt     time.Time  `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
	UpdatedAt     time.Time  `neo4jDb:"property:updatedAt;lookupName:UPDATED_AT;supportCaseSensitive:false"`
	Source        DataSource `neo4jDb:"property:source;lookupName:SOURCE;supportCaseSensitive:false"`
	SourceOfTruth DataSource
}

func (User UserEntity) ToString() string {
	return fmt.Sprintf("id: %s\nfirstName: %s\nlastName: %s", User.Id, User.FirstName, User.LastName)
}

type UserEntities []UserEntity

func (user UserEntity) Labels() []string {
	return []string{"User"}
}
