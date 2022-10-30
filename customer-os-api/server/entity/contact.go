package entity

import (
	"fmt"
	"time"
)

type ContactEntity struct {
	Id               string
	FirstName        string
	LastName         string
	Label            string
	Company          string
	ContactType      string
	CreatedAt        time.Time
	Groups           ContactGroupEntities
	TextCustomFields TextCustomFieldEntities
}

func (contact ContactEntity) ToString() string {
	return fmt.Sprintf("id: %s\nfirstName: %s\nlastName: %s", contact.Id, contact.FirstName, contact.LastName)
}

type ContactNodes []ContactEntity
