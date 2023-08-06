package entity

import (
	"fmt"
)

type EmailEntity struct {
	Id          string
	Email       string
	RawEmail    string
	IsReachable *string
}

func (email EmailEntity) ToString() string {
	return fmt.Sprintf("id: %s\nemail: %s", email.Id, email.Email)
}

type EmailEntities []EmailEntity
