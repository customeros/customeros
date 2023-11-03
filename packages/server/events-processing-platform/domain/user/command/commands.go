package command

import (
	cmnmod "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
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

type AddPlayerInfoCommand struct {
	eventstore.BaseCommand
	Provider   string `json:"provider" validate:"required"`
	AuthId     string `json:"authId" validate:"required"`
	IdentityId string
	Source     cmnmod.Source
	Timestamp  *time.Time
}

func NewAddPlayerInfoCommand(objectID, tenant, userId string, source cmnmod.Source, provider, authId, identityId string, timestamp *time.Time) *AddPlayerInfoCommand {
	return &AddPlayerInfoCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant, userId),
		Source:      source,
		Timestamp:   timestamp,
		Provider:    provider,
		AuthId:      authId,
		IdentityId:  identityId,
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

type LinkPhoneNumberCommand struct {
	eventstore.BaseCommand
	PhoneNumberId string `json:"phoneNumberId" validate:"required"`
	Primary       bool
	Label         string
}

func NewLinkPhoneNumberCommand(objectID, tenant, userId, phoneNumberId, label string, primary bool) *LinkPhoneNumberCommand {
	return &LinkPhoneNumberCommand{
		BaseCommand:   eventstore.NewBaseCommand(objectID, tenant, userId),
		PhoneNumberId: phoneNumberId,
		Primary:       primary,
		Label:         label,
	}
}

type LinkEmailCommand struct {
	eventstore.BaseCommand
	EmailId   string `json:"emailId" validate:"required"`
	Primary   bool
	Label     string
	AppSource string
}

func NewLinkEmailCommand(objectID, tenant, loggedInUserId, emailId, label, appSource string, primary bool) *LinkEmailCommand {
	return &LinkEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant, loggedInUserId),
		EmailId:     emailId,
		Primary:     primary,
		Label:       label,
		AppSource:   appSource,
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
