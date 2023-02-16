package models

import (
	"fmt"
)

type Contact struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	Uuid      string `json:"uuid" bson:"uuid,omitempty"`
	FirstName string `json:"firstName" bson:"firstName,omitempty"`
	LastName  string `json:"lastName" bson:"lastName,omitempty"`
}

func (contact *Contact) String() string {
	return fmt.Sprintf("ID: {%s}, uuid: {%s}, firstName: {%s}, lastName: {%s}", contact.ID, contact.Uuid, contact.FirstName, contact.LastName)
}

func NewContact() *Contact {
	return &Contact{}
}

//func ContactToProto(contact *Contact, id string) *contactService.Contact {
//	return &contactService.Contact{
//		ID:        id,
//		Uuid:      contact.uuid,
//		FirstName: contact.FirstName,
//		LastName:  contact.LastName,
//	}
//}
