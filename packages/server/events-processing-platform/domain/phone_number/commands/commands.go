package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type CreatePhoneNumberCommand struct {
	eventstore.BaseCommand
	Tenant      string
	PhoneNumber string
}

type UpsertPhoneNumberCommand struct {
	eventstore.BaseCommand
	Tenant         string
	RawPhoneNumber string
	Source         string
	SourceOfTruth  string
	AppSource      string
	CreatedAt      *time.Time
	UpdatedAt      *time.Time
}

func NewCreatePhoneNumberCommand(aggregateID, tenant, rawPhoneNumber string) *CreatePhoneNumberCommand {
	return &CreatePhoneNumberCommand{
		BaseCommand: eventstore.NewBaseCommand(aggregateID),
		Tenant:      tenant,
		PhoneNumber: rawPhoneNumber,
	}
}

func NewUpsertPhoneNumberCommand(aggregateID, tenant, rawPhoneNumber, source, sourceOfTruth, appSource string, createdAt, updatedAt *time.Time) *UpsertPhoneNumberCommand {
	return &UpsertPhoneNumberCommand{
		BaseCommand:    eventstore.NewBaseCommand(aggregateID),
		Tenant:         tenant,
		RawPhoneNumber: rawPhoneNumber,
		Source:         source,
		SourceOfTruth:  sourceOfTruth,
		AppSource:      appSource,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
