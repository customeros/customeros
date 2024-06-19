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

func NewUpsertEmailCommand(objectId, tenant, loggedInUserId, rawEmail string, source cmnmod.Source, createdAt, updatedAt *time.Time) *UpsertEmailCommand {
	return &UpsertEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(objectId, tenant, loggedInUserId),
		RawEmail:    rawEmail,
		Source:      source,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}
