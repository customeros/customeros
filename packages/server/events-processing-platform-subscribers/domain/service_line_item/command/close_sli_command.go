package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type CloseServiceLineItemCommand struct {
	eventstore.BaseCommand
	UpdatedAt     *time.Time
	EndedAt       *time.Time
	SliIsCanceled bool
}

func NewCloseServiceLineItemCommand(serviceLineItemId, tenant, loggedInUserId, appSource string, sliIsCanceled bool, endedAt, updatedAt *time.Time) *CloseServiceLineItemCommand {
	return &CloseServiceLineItemCommand{
		BaseCommand:   eventstore.NewBaseCommand(serviceLineItemId, tenant, loggedInUserId).WithAppSource(appSource),
		UpdatedAt:     updatedAt,
		EndedAt:       endedAt,
		SliIsCanceled: sliIsCanceled,
	}
}
