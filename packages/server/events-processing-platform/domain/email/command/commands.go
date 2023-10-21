package command

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertEmailCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	RawEmail        string `json:"rawEmail" validate:"required"`
	Source          cmnmod.Source
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

type FailedEmailValidationCommand struct {
	eventstore.BaseCommand
	ValidationError string
}

type EmailValidatedCommand struct {
	eventstore.BaseCommand
	RawEmail        string
	IsReachable     string
	ValidationError string
	AcceptsMail     bool
	CanConnectSmtp  bool
	HasFullInbox    bool
	IsCatchAll      bool
	IsDeliverable   bool
	IsDisabled      bool
	Domain          string
	IsValidSyntax   bool
	Username        string
	EmailAddress    string
}

func NewUpsertEmailCommand(objectId, tenant, loggedInUserId, rawEmail string, source cmnmod.Source, createdAt, updatedAt *time.Time) *UpsertEmailCommand {
	return &UpsertEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(objectId, tenant, loggedInUserId),
		RawEmail:    rawEmail,
		Source:      source,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

func NewFailedEmailValidationCommand(objectID, tenant, validationError string) *FailedEmailValidationCommand {
	return &FailedEmailValidationCommand{
		BaseCommand:     eventstore.NewBaseCommand(objectID, tenant, ""),
		ValidationError: validationError,
	}
}

func NewEmailValidatedCommand(objectID, tenant, rawEmail, isReachable, validationError, domain, username, emailAddress string, acceptsMail, canConnectSmtp, hasFullInbox, isCatchAll, isDisabled, isValidSyntax bool) *EmailValidatedCommand {
	return &EmailValidatedCommand{
		BaseCommand:     eventstore.NewBaseCommand(objectID, tenant, ""),
		IsReachable:     isReachable,
		RawEmail:        rawEmail,
		ValidationError: validationError,
		Domain:          domain,
		Username:        username,
		EmailAddress:    emailAddress,
		AcceptsMail:     acceptsMail,
		CanConnectSmtp:  canConnectSmtp,
		HasFullInbox:    hasFullInbox,
		IsCatchAll:      isCatchAll,
		IsDisabled:      isDisabled,
		IsValidSyntax:   isValidSyntax,
	}
}
