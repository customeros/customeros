package entity

import (
	"fmt"
)

type UserEntity struct {
	Id        string
	FirstName string `neo4jDb:"property:firstName;lookupName:FIRST_NAME;supportCaseSensitive:true"`
	LastName  string `neo4jDb:"property:lastName;lookupName:LAST_NAME;supportCaseSensitive:true"`
}

func (User UserEntity) ToString() string {
	return fmt.Sprintf("id: %s\nfirstName: %s\nlastName: %s", User.Id, User.FirstName, User.LastName)
}
