package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type RequestNextCycleDateCommand struct {
	eventstore.BaseCommand
	AppSource string
}

func NewRequestNextCycleDateCommand(tenant, contractId, loggedInUserId, appSource string) *RequestNextCycleDateCommand {
	return &RequestNextCycleDateCommand{
		BaseCommand: eventstore.NewBaseCommand(contractId, tenant, loggedInUserId),
		AppSource:   appSource,
	}
}
