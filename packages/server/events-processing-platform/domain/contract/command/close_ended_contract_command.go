package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type CloseEndedContractCommand struct {
	eventstore.BaseCommand
	AppSource string
}

func NewCloseEndedContractCommand(contractId, tenant, loggedInUserId, appSource string) *CloseEndedContractCommand {
	return &CloseEndedContractCommand{
		BaseCommand: eventstore.NewBaseCommand(contractId, tenant, loggedInUserId),
		AppSource:   appSource,
	}
}
