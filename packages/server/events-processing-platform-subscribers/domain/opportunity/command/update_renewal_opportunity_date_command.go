package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpdateRenewalOpportunityNextCycleDateCommand struct {
	eventstore.BaseCommand
	AppSource  string
	ContractId string
	UpdatedAt  *time.Time
	RenewedAt  *time.Time
}

func NewUpdateRenewalOpportunityNextCycleDateCommand(opportunityId, tenant, loggedInUserId, appSource string, updatedAt, renewedAt *time.Time) *UpdateRenewalOpportunityNextCycleDateCommand {
	return &UpdateRenewalOpportunityNextCycleDateCommand{
		BaseCommand: eventstore.NewBaseCommand(opportunityId, tenant, loggedInUserId),
		AppSource:   appSource,
		UpdatedAt:   updatedAt,
		RenewedAt:   renewedAt,
	}
}
