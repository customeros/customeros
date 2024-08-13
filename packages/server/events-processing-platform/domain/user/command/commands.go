package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events/event/common"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
	"time"
)

type UpsertUserCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	DataFields      models.UserDataFields
	Source          cmnmod.Source
	ExternalSystem  cmnmod.ExternalSystem
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func NewUpsertUserCommand(objectID, tenant, userId string, source cmnmod.Source, externalSystem cmnmod.ExternalSystem, dataFields models.UserDataFields, createdAt, updatedAt *time.Time) *UpsertUserCommand {
	return &UpsertUserCommand{
		BaseCommand:    eventstore.NewBaseCommand(objectID, tenant, userId),
		DataFields:     dataFields,
		Source:         source,
		ExternalSystem: externalSystem,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
}

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
