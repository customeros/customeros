package entity

import "fmt"

type ContactNode struct {
	Id          string
	FirstName   string
	LastName    string
	Label       string
	ContactType string
}

func NewContact(id string, firstName string, lastName string, label string, contactType string) ContactNode {
	contact := ContactNode{Id: id, FirstName: firstName, LastName: lastName, Label: label, ContactType: contactType}
	return contact
}

func (contact ContactNode) ToString() string {
	return fmt.Sprintf("id: %s\nfirstName: %s", contact.Id, contact.FirstName)
}
