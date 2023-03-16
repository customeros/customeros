package entity

import (
	"fmt"
	"time"
)

type EmailEntity struct {
	Id            string
	Email         string
	RawEmail      string
	Validated     bool
	Label         string
	Primary       bool
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	CreatedAt     time.Time
	UpdatedAt     time.Time

	InteractionEventParticipantDetails InteractionEventParticipantDetails
	DataloaderKey                      string
}

func (email EmailEntity) ToString() string {
	return fmt.Sprintf("id: %s\nemail: %s\nlabel: %s", email.Id, email.Email, email.Label)
}

type EmailEntities []EmailEntity

func (EmailEntity) IsInteractionEventParticipant() {}

func (EmailEntity) InteractionEventParticipantLabel() string {
	return NodeLabel_Email
}

func (email EmailEntity) GetDataloaderKey() string {
	return email.DataloaderKey
}

func (EmailEntity) Labels(tenant string) []string {
	return []string{
		NodeLabel_Email,
		NodeLabel_Email + "_" + tenant,
	}
}
