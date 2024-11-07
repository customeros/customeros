package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"time"
)

type EmailProperty string

const (
	EmailPropertyEmail                 EmailProperty = "email"
	EmailPropertyRawEmail              EmailProperty = "rawEmail"
	EmailPropertyIsRoleAccount         EmailProperty = "isRoleAccount"
	EmailPropertyIsSystemGenerated     EmailProperty = "isSystemGenerated"
	EmailPropertyIsValidSyntax         EmailProperty = "isValidSyntax"
	EmailPropertyIsCatchAll            EmailProperty = "isCatchAll"
	EmailPropertyDeliverable           EmailProperty = "deliverable"
	EmailPropertyValidatedAt           EmailProperty = "techValidatedAt"
	EmailPropertyValidationRequestedAt EmailProperty = "techValidationRequestedAt"
	EmailPropertyUsername              EmailProperty = "username"
	EmailPropertyIsRisky               EmailProperty = "isRisky"
	EmailPropertyIsFirewalled          EmailProperty = "isFirewalled"
	EmailPropertyProvider              EmailProperty = "provider"
	EmailPropertyFirewall              EmailProperty = "firewall"
	EmailPropertyIsMailboxFull         EmailProperty = "isMailboxFull"
	EmailPropertyIsFreeAccount         EmailProperty = "isFreeAccount"
	EmailPropertySmtpSuccess           EmailProperty = "smtpSuccess"
	EmailPropertyResponseCode          EmailProperty = "verifyResponseCode"
	EmailPropertyErrorCode             EmailProperty = "verifyErrorCode"
	EmailPropertyDescription           EmailProperty = "verifyDescription"
	EmailPropertyIsPrimaryDomain       EmailProperty = "isPrimaryDomain"
	EmailPropertyPrimaryDomain         EmailProperty = "primaryDomain"
	EmailPropertyAlternateEmail        EmailProperty = "alternateEmail"
	EmailPropertyWork                  EmailProperty = "work"
	EmailPropertyRetryValidation       EmailProperty = "retryValidation"
)

type EmailEntity struct {
	DataLoaderKey
	Id        string
	Email     string `neo4jDb:"property:email;lookupName:EMAIL;supportCaseSensitive:true"`
	RawEmail  string `neo4jDb:"property:rawEmail;lookupName:RAW_EMAIL;supportCaseSensitive:true"`
	Work      *bool
	Primary   bool
	Source    DataSource
	CreatedAt time.Time
	UpdatedAt time.Time

	IsValidSyntax     *bool
	Username          *string
	IsRisky           *bool
	IsFirewalled      *bool
	Provider          *string
	Firewall          *string
	IsCatchAll        *bool
	Deliverable       *string
	IsMailboxFull     *bool
	IsRoleAccount     *bool
	IsSystemGenerated *bool
	IsFreeAccount     *bool
	SmtpSuccess       *bool
	ResponseCode      *string
	ErrorCode         *string
	Description       *string
	IsPrimaryDomain   *bool
	PrimaryDomain     *string
	AlternateEmail    *string
	RetryValidation   *bool

	EmailInternalFields EmailInternalFields

	InteractionEventParticipantDetails   InteractionEventParticipantDetails
	InteractionSessionParticipantDetails InteractionSessionParticipantDetails
}

type EmailInternalFields struct {
	ValidatedAt           *time.Time
	ValidationRequestedAt *time.Time
}

type EmailEntities []EmailEntity

func (e EmailEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (EmailEntity) IsInteractionEventParticipant() {}

func (EmailEntity) IsInteractionSessionParticipant() {}

func (EmailEntity) IsMeetingParticipant() {}

func (EmailEntity) EntityLabel() string {
	return model.NodeLabelEmail
}

func (e EmailEntity) Labels(tenant string) []string {
	return []string{
		e.EntityLabel(),
		e.EntityLabel() + "_" + tenant,
	}
}
