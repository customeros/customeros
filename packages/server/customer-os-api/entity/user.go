package entity

import (
	"fmt"
	"time"
)

type UserEntity struct {
	Id        string
	FirstName string    `neo4jDb:"property:firstName;lookupName:FIRST_NAME;supportCaseSensitive:true"`
	LastName  string    `neo4jDb:"property:lastName;lookupName:LAST_NAME;supportCaseSensitive:true"`
	Email     string    `neo4jDb:"property:email;lookupName:EMAIL;supportCaseSensitive:true"`
	CreatedAt time.Time `neo4jDb:"property:createdAt;lookupName:CREATED_AT;supportCaseSensitive:false"`
}

func (User UserEntity) ToString() string {
	return fmt.Sprintf("id: %s\nfirstName: %s\nlastName: %s", User.Id, User.FirstName, User.LastName)
}

type UserEntities []UserEntity
