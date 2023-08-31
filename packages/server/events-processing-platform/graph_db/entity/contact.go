package entity

import (
	"fmt"
	"time"
)

type ContactEntity struct {
	Id              string
	Prefix          string
	Name            string
	FirstName       string
	LastName        string
	Description     string
	Timezone        string
	ProfilePhotoUrl string
	CreatedAt       *time.Time
	UpdatedAt       time.Time
	Source          DataSource
	SourceOfTruth   DataSource
	AppSource       string

	InteractionEventParticipantDetails   InteractionEventParticipantDetails
	InteractionSessionParticipantDetails InteractionSessionParticipantDetails
	DataloaderKey                        string
}

func (contact ContactEntity) ToString() string {
	return fmt.Sprintf("ContactEntity{Id: %s, Prefix: %s, Name: %s, FirstName: %s, LastName: %s, Description: %s, CreatedAt: %s, UpdatedAt: %s, Source: %s, SourceOfTruth: %s, AppSource: %s}",
		contact.Id, contact.Prefix, contact.Name, contact.FirstName, contact.LastName, contact.Description, contact.CreatedAt, contact.UpdatedAt, contact.Source, contact.SourceOfTruth, contact.AppSource)
}

type ContactEntities []ContactEntity

func (ContactEntity) IsInteractionEventParticipant() {}

func (ContactEntity) InteractionEventParticipantLabel() string {
	return NodeLabel_Contact
}

func (ContactEntity) IsInteractionSessionParticipant() {}

func (ContactEntity) InteractionSessionParticipantLabel() string {
	return NodeLabel_Contact
}

func (ContactEntity) IsMeetingParticipant() {}

func (ContactEntity) MeetingParticipantLabel() string {
	return NodeLabel_Contact
}

func (ContactEntity) IsNotedEntity() {}

func (ContactEntity) NotedEntityLabel() string {
	return NodeLabel_Contact
}

func (contact ContactEntity) GetDataloaderKey() string {
	return contact.DataloaderKey
}

func (ContactEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Contact,
		NodeLabel_Contact + "_" + tenant,
	}
}
