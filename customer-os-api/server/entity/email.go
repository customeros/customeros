package entity

import (
	"fmt"
)

type EmailEntity struct {
	Email string
	Label string
}

func (email EmailEntity) ToString() string {
	return fmt.Sprintf("email: %s\nlabel: %s", email.Email, email.Label)
}

type EmailEntities []EmailEntity
