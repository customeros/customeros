package entity

import (
	"fmt"
	"time"
)

type PhoneNumberEntity struct {
	Id             string
	E164           string
	Validated      *bool
	RawPhoneNumber string
	Label          string
	Primary        bool
	Source         DataSource
	SourceOfTruth  DataSource
	AppSource      string
	CreatedAt      time.Time
	UpdatedAt      time.Time ``

	InteractionEventParticipantDetails   InteractionEventParticipantDetails
	InteractionSessionParticipantDetails InteractionSessionParticipantDetails
	DataloaderKey                        string
}

func (phone PhoneNumberEntity) ToString() string {
	return fmt.Sprintf("id: %s\ne164: %s\nlabel: %s", phone.Id, phone.E164, phone.Label)
}

type PhoneNumberEntities []PhoneNumberEntity

func (PhoneNumberEntity) IsInteractionEventParticipant() {}

func (PhoneNumberEntity) IsInteractionSessionParticipant() {}

func (PhoneNumberEntity) ParticipantLabel() string {
	return NodeLabel_PhoneNumber
}

func (phoneNumber PhoneNumberEntity) GetDataloaderKey() string {
	return phoneNumber.DataloaderKey
}

func (PhoneNumberEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_PhoneNumber,
		NodeLabel_PhoneNumber + "_" + tenant,
	}
}
