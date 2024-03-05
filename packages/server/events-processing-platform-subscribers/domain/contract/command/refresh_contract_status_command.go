package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type RefreshContractStatusCommand struct {
	eventstore.BaseCommand
	AppSource string
}

func NewRefreshContractStatusCommand(contractId, tenant, loggedInUserId, appSource string) *RefreshContractStatusCommand {
	return &RefreshContractStatusCommand{
		BaseCommand: eventstore.NewBaseCommand(contractId, tenant, loggedInUserId),
		AppSource:   appSource,
	}
}
