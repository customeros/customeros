package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type LinkJobRoleCommand struct {
	eventstore.BaseCommand
	JobRoleId string
}

// TODO add userId
func NewLinkJobRoleCommand(objectID, tenant, jobRoleId string) *LinkJobRoleCommand {
	return &LinkJobRoleCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant, ""),
		JobRoleId:   jobRoleId,
	}
}

type AddRoleCommand struct {
	eventstore.BaseCommand
	Role string `json:"role" validate:"required"`
}

func NewAddRole(objectID, tenant, userId, role string) *AddRoleCommand {
	return &AddRoleCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant, userId),
		Role:        role,
	}
}

type RemoveRoleCommand struct {
	eventstore.BaseCommand
	Role string
}

func NewRemoveRole(objectID, tenant, userId, role string) *RemoveRoleCommand {
	return &RemoveRoleCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant, userId),
		Role:        role,
	}
}
