package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type CloseLooseOpportunityCommand struct {
	eventstore.BaseCommand
	AppSource string
	UpdatedAt *time.Time
	ClosedAt  *time.Time
}

func NewCloseLooseOpportunityCommand(opportunityId, tenant, loggedInUserId, appSource string, updatedAt, closedAt *time.Time) *CloseLooseOpportunityCommand {
	return &CloseLooseOpportunityCommand{
		BaseCommand: eventstore.NewBaseCommand(opportunityId, tenant, loggedInUserId),
		UpdatedAt:   updatedAt,
		ClosedAt:    closedAt,
		AppSource:   appSource,
	}
}
