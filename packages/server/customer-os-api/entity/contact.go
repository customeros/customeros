package entity

import (
	"fmt"
	"time"
)

type ContactEntity struct {
	Id        string
	Title     string
	FirstName string
	LastName  string
	Label     string
	Notes     string
	CreatedAt time.Time
}

func (contact ContactEntity) ToString() string {
	return fmt.Sprintf("id: %s\nfirstName: %s\nlastName: %s", contact.Id, contact.FirstName, contact.LastName)
}

type ContactEntities []ContactEntity
