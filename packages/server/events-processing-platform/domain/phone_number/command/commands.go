package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"time"
)

type UpsertPhoneNumberCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	RawPhoneNumber  string
	Source          common.Source
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func NewUpsertPhoneNumberCommand(objectId, tenant, loggedInUserId, rawPhoneNumber string, source common.Source, createdAt, updatedAt *time.Time) *UpsertPhoneNumberCommand {
	return &UpsertPhoneNumberCommand{
		BaseCommand:    eventstore.NewBaseCommand(objectId, tenant, loggedInUserId),
		RawPhoneNumber: rawPhoneNumber,
		Source:         source,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
