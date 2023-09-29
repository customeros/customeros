package command

import (
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertUserCommand struct {
	eventstore.BaseCommand
	IsCreateCommand bool
	DataFields      models.UserDataFields
	Source          common_models.Source
	ExternalSystem  common_models.ExternalSystem
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
}

func NewUpsertUserCommand(objectID, tenant, userId string, source common_models.Source, externalSystem common_models.ExternalSystem, dataFields models.UserDataFields, createdAt, updatedAt *time.Time) *UpsertUserCommand {
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

type LinkPhoneNumberCommand struct {
	eventstore.BaseCommand
	PhoneNumberId string
	Primary       bool
	Label         string
}

// TODO add userId
func NewLinkPhoneNumberCommand(objectID, tenant, phoneNumberId, label string, primary bool) *LinkPhoneNumberCommand {
	return &LinkPhoneNumberCommand{
		BaseCommand:   eventstore.NewBaseCommand(objectID, tenant, ""),
		PhoneNumberId: phoneNumberId,
		Primary:       primary,
		Label:         label,
	}
}

type LinkEmailCommand struct {
	eventstore.BaseCommand
	EmailId string
	Primary bool
	Label   string
}

// TODO add userId
func NewLinkEmailCommand(objectID, tenant, emailId, label string, primary bool) *LinkEmailCommand {
	return &LinkEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant, ""),
		EmailId:     emailId,
		Primary:     primary,
		Label:       label,
	}
}
