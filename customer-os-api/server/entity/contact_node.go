package entity

import (
	"fmt"
	"time"
)

type ContactNode struct {
	Id          string
	FirstName   string
	LastName    string
	Label       string
	ContactType string
	CreatedAt   time.Time
	Groups      ContactGroupNodes
}

func (contact ContactNode) ToString() string {
	return fmt.Sprintf("id: %s\nfirstName: %s\nlastName: %s", contact.Id, contact.FirstName, contact.LastName)
}

type ContactNodes []ContactNode
