package command

import "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"

type RefreshRenewalSummaryCommand struct {
	eventstore.BaseCommand
}

func NewRefreshRenewalSummaryCommand(tenant, orgId, userId, appSource string) *RefreshRenewalSummaryCommand {
	return &RefreshRenewalSummaryCommand{
		BaseCommand: eventstore.NewBaseCommand(orgId, tenant, userId).WithAppSource(appSource),
	}
}
