package entity

import (
	"fmt"
	"time"
)

type EmailEntity struct {
	Id            string
	Email         string `neo4jDb:"property:email;lookupName:EMAIL;supportCaseSensitive:true"`
	RawEmail      string `neo4jDb:"property:rawEmail;lookupName:RAW_EMAIL;supportCaseSensitive:true"`
	Label         string
	Primary       bool
	Source        DataSource
	SourceOfTruth DataSource
	AppSource     string
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Validated      *bool
	IsReachable    *string
	IsValidSyntax  *bool
	CanConnectSMTP *bool
	AcceptsMail    *bool
	HasFullInbox   *bool
	IsCatchAll     *bool
	IsDeliverable  *bool
	IsDisabled     *bool
	Error          *string

	InteractionEventParticipantDetails   InteractionEventParticipantDetails
	InteractionSessionParticipantDetails InteractionSessionParticipantDetails
	DataloaderKey                        string
}

func (email EmailEntity) ToString() string {
	return fmt.Sprintf("id: %s\nemail: %s\nlabel: %s", email.Id, email.Email, email.Label)
}

type EmailEntities []EmailEntity

func (EmailEntity) IsInteractionEventParticipant() {}

func (EmailEntity) InteractionEventParticipantLabel() string {
	return NodeLabel_Email
}

func (EmailEntity) IsInteractionSessionParticipant() {}

func (EmailEntity) InteractionSessionParticipantLabel() string {
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
