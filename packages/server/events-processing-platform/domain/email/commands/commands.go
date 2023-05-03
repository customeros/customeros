package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type CreateEmailCommand struct {
	eventstore.BaseCommand
	Tenant    string
	Email     string
	Source    models.Source
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type UpsertEmailCommand struct {
	eventstore.BaseCommand
	Tenant    string
	RawEmail  string
	Source    models.Source
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type FailEmailValidationCommand struct {
	eventstore.BaseCommand
	Tenant          string
	ValidationError string
}

type EmailValidatedCommand struct {
	eventstore.BaseCommand
	Tenant          string
	RawEmail        string
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
	NormalizedEmail string
}

func NewCreateEmailCommand(aggregateID, tenant, rawEmail, source, sourceOfTruth, appSource string, createdAt, updatedAt *time.Time) *CreateEmailCommand {
	return &CreateEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(aggregateID),
		Tenant:      tenant,
		Email:       rawEmail,
		Source: models.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewUpsertEmailCommand(aggregateID, tenant, rawEmail, source, sourceOfTruth, appSource string, createdAt, updatedAt *time.Time) *UpsertEmailCommand {
	return &UpsertEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(aggregateID),
		Tenant:      tenant,
		RawEmail:    rawEmail,
		Source: models.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewFailEmailValidationCommand(aggregateID, tenant, validationError string) *FailEmailValidationCommand {
	return &FailEmailValidationCommand{
		BaseCommand:     eventstore.NewBaseCommand(aggregateID),
		Tenant:          tenant,
		ValidationError: validationError,
	}
}

func NewEmailValidatedCommand(aggregateID, tenant, rawEmail, validationError, domain, username, normalizedEmail string, acceptsMail, canConnectSmtp, hasFullInbox, isCatchAll, isDisabled, isValidSyntax bool) *EmailValidatedCommand {
	return &EmailValidatedCommand{
		BaseCommand:     eventstore.NewBaseCommand(aggregateID),
		Tenant:          tenant,
		RawEmail:        rawEmail,
		ValidationError: validationError,
		Domain:          domain,
		Username:        username,
		NormalizedEmail: normalizedEmail,
		AcceptsMail:     acceptsMail,
		CanConnectSmtp:  canConnectSmtp,
		HasFullInbox:    hasFullInbox,
		IsCatchAll:      isCatchAll,
		IsDisabled:      isDisabled,
		IsValidSyntax:   isValidSyntax,
	}
}
