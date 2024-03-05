package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type DeleteServiceLineItemCommand struct {
	eventstore.BaseCommand
}

func NewDeleteServiceLineItemCommand(serviceLineItemId, tenant, loggedInUserId, appSource string) *DeleteServiceLineItemCommand {
	return &DeleteServiceLineItemCommand{
		BaseCommand: eventstore.NewBaseCommand(serviceLineItemId, tenant, loggedInUserId).WithAppSource(appSource),
	}
}
