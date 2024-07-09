package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	"time"
)

type EmailProperty string

const (
	EmailPropertyEmail          EmailProperty = "email"
	EmailPropertyRawEmail       EmailProperty = "rawEmail"
	EmailPropertyIsDisposable   EmailProperty = "isDisposable"
	EmailPropertyIsRoleAccount  EmailProperty = "isRoleAccount"
	EmailPropertyIsValidSyntax  EmailProperty = "isValidSyntax"
	EmailPropertyCanConnectSMTP EmailProperty = "canConnectSMTP"
	EmailPropertyAcceptsMail    EmailProperty = "acceptsMail"
	EmailPropertyHasFullInbox   EmailProperty = "hasFullInbox"
	EmailPropertyIsCatchAll     EmailProperty = "isCatchAll"
	EmailPropertyIsDeliverable  EmailProperty = "isDeliverable"
	EmailPropertyIsDisabled     EmailProperty = "isDisabled"
	EmailPropertyError          EmailProperty = "error"
	EmailPropertyValidated      EmailProperty = "validated"
	EmailPropertyIsReachable    EmailProperty = "isReachable"
)

type EmailEntity struct {
	DataLoaderKey
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
	IsDisposable   *bool
	IsRoleAccount  *bool // group inbox indicator

	InteractionEventParticipantDetails   InteractionEventParticipantDetails
	InteractionSessionParticipantDetails InteractionSessionParticipantDetails
}

type EmailEntities []EmailEntity

func (e EmailEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (EmailEntity) IsInteractionEventParticipant() {}

func (EmailEntity) IsInteractionSessionParticipant() {}

func (EmailEntity) IsMeetingParticipant() {}

func (EmailEntity) EntityLabel() string {
	return neo4jutil.NodeLabelEmail
}

func (e EmailEntity) Labels(tenant string) []string {
	return []string{
		e.EntityLabel(),
		e.EntityLabel() + "_" + tenant,
	}
}
