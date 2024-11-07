package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type EmailValidatedEventV2 struct {
	Tenant            string    `json:"tenant" validate:"required"`
	RawEmail          string    `json:"rawEmail" validate:"required"`
	ValidatedAt       time.Time `json:"validatedAt" validate:"required"`
	Email             string    `json:"email"`
	Domain            string    `json:"domain"`
	Username          string    `json:"username"`
	IsValidSyntax     bool      `json:"isValidSyntax"`
	IsRisky           bool      `json:"isRisky"`
	IsFirewalled      bool      `json:"isFirewalled"`
	Provider          string    `json:"provider"`
	Firewall          string    `json:"firewall"`
	IsCatchAll        bool      `json:"isCatchAll"`
	Deliverable       string    `json:"deliverable"`
	IsMailboxFull     bool      `json:"isMailboxFull"`
	IsRoleAccount     bool      `json:"isRoleAccount"`
	IsSystemGenerated bool      `json:"isSystemGenerated"`
	IsFreeAccount     bool      `json:"isFreeAccount"`
	SmtpSuccess       bool      `json:"smtpSuccess"`
	ResponseCode      string    `json:"responseCode"`
	ErrorCode         string    `json:"errorCode"`
	Description       string    `json:"description"`
	IsPrimaryDomain   bool      `json:"isPrimaryDomain"`
	PrimaryDomain     string    `json:"primaryDomain"`
	AlternateEmail    string    `json:"alternateEmail"`
	RetryValidation   bool      `json:"retryValidation"`
}

func NewEmailValidatedEventV2(aggregate eventstore.Aggregate, tenant, rawEmail, email, domain, username string,
	isValidSyntax, risky, firewalled bool, provider, firewall, deliverable string,
	isCatchAll, isMailboxFull, isRoleAccount, isSystemGenerated, isFreeAccount, smtpSuccess bool,
	responseCode, errorCode, description string, isPrimaryDomain bool, primaryDomain, alternateEmail string, retryValidation bool) (eventstore.Event, error) {
	eventData := EmailValidatedEventV2{
		Tenant:            tenant,
		RawEmail:          rawEmail,
		Email:             email,
		ValidatedAt:       utils.Now(),
		Domain:            domain,
		Username:          username,
		IsValidSyntax:     isValidSyntax,
		IsRisky:           risky,
		IsFirewalled:      firewalled,
		Provider:          provider,
		Firewall:          firewall,
		IsCatchAll:        isCatchAll,
		Deliverable:       deliverable,
		IsMailboxFull:     isMailboxFull,
		IsRoleAccount:     isRoleAccount,
		IsSystemGenerated: isSystemGenerated,
		IsFreeAccount:     isFreeAccount,
		SmtpSuccess:       smtpSuccess,
		ResponseCode:      responseCode,
		ErrorCode:         errorCode,
		Description:       description,
		IsPrimaryDomain:   isPrimaryDomain,
		PrimaryDomain:     primaryDomain,
		AlternateEmail:    alternateEmail,
		RetryValidation:   retryValidation,
	}

	if err := validator.GetValidator().Struct(eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "failed to validate EmailValidatedEvent")
	}

	event := eventstore.NewBaseEvent(aggregate, EmailValidatedV2)
	if err := event.SetJsonData(&eventData); err != nil {
		return eventstore.Event{}, errors.Wrap(err, "error setting json data for EmailValidatedEvent")
	}
	return event, nil
}
