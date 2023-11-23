package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type ResetRenewalOpportunityOnExpirationCommand struct {
	eventstore.BaseCommand
	AppSource string
}

func NewResetRenewalOpportunityOnExpirationCommand(contractId, tenant, loggedInUserId, appSource string) *ResetRenewalOpportunityOnExpirationCommand {
	return &ResetRenewalOpportunityOnExpirationCommand{
		BaseCommand: eventstore.NewBaseCommand(contractId, tenant, loggedInUserId),
		AppSource:   appSource,
	}
}
