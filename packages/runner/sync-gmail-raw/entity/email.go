package entity

import (
	"fmt"
)

type EmailEntity struct {
	Id       string
	Email    string `neo4jDb:"property:email;lookupName:EMAIL;supportCaseSensitive:true"`
	RawEmail string `neo4jDb:"property:rawEmail;lookupName:RAW_EMAIL;supportCaseSensitive:true"`
}

func (email EmailEntity) ToString() string {
	return fmt.Sprintf("id: %s\nemail: %s", email.Id, email.Email)
}

type EmailEntities []EmailEntity
