package entity

import (
	"fmt"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"time"
)

type EmailEntity struct {
	Id            string
	Email         string `neo4jDb:"property:email;lookupName:EMAIL;supportCaseSensitive:true"`
	RawEmail      string `neo4jDb:"property:rawEmail;lookupName:RAW_EMAIL;supportCaseSensitive:true"`
	Label         string
	Primary       bool
	Source        neo4jentity.DataSource
	SourceOfTruth neo4jentity.DataSource
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

func (EmailEntity) IsInteractionSessionParticipant() {}

func (EmailEntity) ParticipantLabel() string {
	return neo4jentity.NodeLabel_Email
}

func (EmailEntity) IsMeetingParticipant() {}

func (EmailEntity) MeetingParticipantLabel() string {
	return neo4jentity.NodeLabel_Email
}

func (email EmailEntity) GetDataloaderKey() string {
	return email.DataloaderKey
}

func (EmailEntity) Labels(tenant string) []string {
	return []string{
		neo4jentity.NodeLabel_Email,
		neo4jentity.NodeLabel_Email + "_" + tenant,
	}
}
