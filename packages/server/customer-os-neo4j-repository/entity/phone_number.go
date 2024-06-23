package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

type PhoneNumberEntity struct {
	DataLoaderKey
	Id             string
	E164           string
	Validated      *bool
	RawPhoneNumber string
	Source         DataSource
	SourceOfTruth  DataSource
	AppSource      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Label          string
	Primary        bool

	InteractionEventParticipantDetails   InteractionEventParticipantDetails
	InteractionSessionParticipantDetails InteractionSessionParticipantDetails
}

type PhoneNumberEntities []PhoneNumberEntity

func (e PhoneNumberEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (PhoneNumberEntity) IsInteractionEventParticipant() {}

func (PhoneNumberEntity) IsInteractionSessionParticipant() {}

func (PhoneNumberEntity) EntityLabel() string {
	return neo4jutil.NodeLabelPhoneNumber
}

func (e PhoneNumberEntity) Labels(tenant string) []string {
	return []string{
		e.EntityLabel(),
		e.EntityLabel() + "_" + tenant,
	}
}
