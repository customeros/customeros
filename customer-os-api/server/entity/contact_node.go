package entity

import "fmt"

type ContactNode struct {
	Id          string
	FirstName   string
	LastName    string
	Label       string
	ContactType string
}

func (contact ContactNode) ToString() string {
	return fmt.Sprintf("id: %s\nfirstName: %s\nlastName: %s", contact.Id, contact.FirstName, contact.LastName)
}

type ContactNodes []ContactNode
