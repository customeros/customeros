package command

import "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"

type RefreshArrCommand struct {
	eventstore.BaseCommand
	AppSource string
}

func NewRefreshArrCommand(tenant, orgId, userId, appSource string) *RefreshArrCommand {
	return &RefreshArrCommand{
		BaseCommand: eventstore.NewBaseCommand(orgId, tenant, userId),
		AppSource:   appSource,
	}
}
