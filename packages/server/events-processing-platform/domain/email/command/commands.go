package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/events"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"time"
)

type UpsertEmailCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	RawEmail        string `json:"rawEmail"`
	Source          events.Source
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func NewUpsertEmailCommand(objectId, tenant, loggedInUserId, rawEmail string, source events.Source, createdAt, updatedAt *time.Time) *UpsertEmailCommand {
	return &UpsertEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(objectId, tenant, loggedInUserId),
		RawEmail:    rawEmail,
		Source:      source,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
