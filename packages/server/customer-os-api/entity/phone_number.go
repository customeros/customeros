package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type PhoneNumberEntity struct {
	Id             string
	E164           string
	Validated      *bool
	RawPhoneNumber string
	Label          string
	Primary        bool
	Source         neo4jentity.DataSource
	SourceOfTruth  neo4jentity.DataSource
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
	return neo4jentity.NodeLabel_PhoneNumber
}

func (phoneNumber PhoneNumberEntity) GetDataloaderKey() string {
	return phoneNumber.DataloaderKey
}

func (PhoneNumberEntity) Labels(tenant string) []string {
	return []string{
		neo4jentity.NodeLabel_PhoneNumber,
		neo4jentity.NodeLabel_PhoneNumber + "_" + tenant,
	}
}
