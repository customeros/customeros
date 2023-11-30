package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type RolloutRenewalOpportunityOnExpirationCommand struct {
	eventstore.BaseCommand
	AppSource string
}

func NewRolloutRenewalOpportunityOnExpirationCommand(contractId, tenant, loggedInUserId, appSource string) *RolloutRenewalOpportunityOnExpirationCommand {
	return &RolloutRenewalOpportunityOnExpirationCommand{
		BaseCommand: eventstore.NewBaseCommand(contractId, tenant, loggedInUserId),
		AppSource:   appSource,
	}
}
