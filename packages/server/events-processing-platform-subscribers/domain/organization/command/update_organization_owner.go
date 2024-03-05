package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
)

type UpdateOrganizationOwnerCommand struct {
	eventstore.BaseCommand
	OwnerUserId    string `json:"userId" validate:"required"` // who became owner
	OrganizationId string `json:"organizationId" validate:"required"`
	ActorUserId    string `json:"actorUserId"` // who set the owner
}

func NewUpdateOrganizationOwnerCommand(organizationId, tenant, loggedInUserId, ownerUserId, appSource string) *UpdateOrganizationOwnerCommand {
	return &UpdateOrganizationOwnerCommand{
		BaseCommand:    eventstore.NewBaseCommand(organizationId, tenant, loggedInUserId).WithAppSource(appSource),
		OrganizationId: organizationId,
		OwnerUserId:    ownerUserId,
		ActorUserId:    loggedInUserId,
	}
}
