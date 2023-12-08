package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type WebScrapeOrganizationCommand struct {
	eventstore.BaseCommand
	Website string
}

func NewWebScrapeOrganizationCommand(tenant, orgId, loggedInUserId, appSource, website string) *WebScrapeOrganizationCommand {
	return &WebScrapeOrganizationCommand{
		BaseCommand: eventstore.NewBaseCommand(orgId, tenant, loggedInUserId).WithAppSource(appSource),
		Website:     website,
	}
}
