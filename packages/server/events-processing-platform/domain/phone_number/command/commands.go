package command

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertPhoneNumberCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	RawPhoneNumber  string
	Source          cmnmod.Source
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func NewUpsertPhoneNumberCommand(objectId, tenant, loggedInUserId, rawPhoneNumber string, source cmnmod.Source, createdAt, updatedAt *time.Time) *UpsertPhoneNumberCommand {
	return &UpsertPhoneNumberCommand{
		BaseCommand:    eventstore.NewBaseCommand(objectId, tenant, loggedInUserId),
		RawPhoneNumber: rawPhoneNumber,
		Source:         source,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}
