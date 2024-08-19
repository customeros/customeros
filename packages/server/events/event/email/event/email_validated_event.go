package event

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/validator"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"github.com/pkg/errors"
	"time"
)

type EmailValidatedEventV2 struct {
	Tenant         string    `json:"tenant" validate:"required"`
	RawEmail       string    `json:"rawEmail" validate:"required"`
	ValidatedAt    time.Time `json:"validatedAt" validate:"required"`
	Email          string    `json:"email"`
	Domain         string    `json:"domain"`
	Username       string    `json:"username"`
	IsValidSyntax  bool      `json:"isValidSyntax"`
	IsRisky        bool      `json:"isRisky"`
	IsFirewalled   bool      `json:"isFirewalled"`
	Provider       string    `json:"provider"`
	Firewall       string    `json:"firewall"`
	IsCatchAll     bool      `json:"isCatchAll"`
	CanConnectSMTP bool      `json:"canConnectSMTP"`
	IsDeliverable  bool      `json:"isDeliverable"`
	IsMailboxFull  bool      `json:"isMailboxFull"`
	IsRoleAccount  bool      `json:"isRoleAccount"`
	IsFreeAccount  bool      `json:"isFreeAccount"`
	SmtpSuccess    bool      `json:"smtpSuccess"`
	ResponseCode   string    `json:"responseCode"`
	ErrorCode      string    `json:"errorCode"`
	Description    string    `json:"description"`
	SmtpResponse   string    `json:"smtpResponse"`
}

func NewEmailValidatedEventV2(aggregate eventstore.Aggregate, tenant, rawEmail, email, domain, username string,
	isValidSyntax, risky, firewalled bool, provider, firewall string,
	isCatchAll, canConnectSMTP, isDeliverable, isMailboxFull, isRoleAccount, isFreeAccount, smtpSuccess bool,
	responseCode, errorCode, description, smtpResponse string) (eventstore.Event, error) {
	eventData := EmailValidatedEventV2{
		Tenant:         tenant,
		RawEmail:       rawEmail,
		Email:          email,
		ValidatedAt:    utils.Now(),
		Domain:         domain,
		Username:       username,
		IsValidSyntax:  isValidSyntax,
		IsRisky:        risky,
		IsFirewalled:   firewalled,
		Provider:       provider,
		Firewall:       firewall,
		IsCatchAll:     isCatchAll,
		CanConnectSMTP: canConnectSMTP,
		IsDeliverable:  isDeliverable,
		IsMailboxFull:  isMailboxFull,
		IsRoleAccount:  isRoleAccount,
		IsFreeAccount:  isFreeAccount,
		SmtpSuccess:    smtpSuccess,
		ResponseCode:   responseCode,
		ErrorCode:      errorCode,
		Description:    description,
		SmtpResponse:   smtpResponse,
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
