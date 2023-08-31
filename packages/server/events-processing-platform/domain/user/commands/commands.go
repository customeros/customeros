package commands

import (
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpsertUserCommand struct {
	eventstore.BaseCommand
	CoreFields models.UserCoreFields
	Source     common_models.Source
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func UpsertUserCommandToUserDto(command *UpsertUserCommand) *models.UserFields {
	return &models.UserFields{
		ID:             command.ObjectID,
		Tenant:         command.Tenant,
		UserCoreFields: command.CoreFields,
		Source:         command.Source,
		CreatedAt:      command.CreatedAt,
		UpdatedAt:      command.UpdatedAt,
	}
}

func NewUpsertUserCommand(objectID, tenant, source, sourceOfTruth, appSource string, coreFields models.UserCoreFields, createdAt, updatedAt *time.Time) *UpsertUserCommand {
	return &UpsertUserCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant),
		CoreFields:  coreFields,
		Source: common_models.Source{
			Source:        source,
			SourceOfTruth: sourceOfTruth,
			AppSource:     appSource,
		},
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type LinkJobRoleCommand struct {
	eventstore.BaseCommand
	JobRoleId string
}

func NewLinkJobRoleCommand(objectID, tenant, jobRoleId string) *LinkJobRoleCommand {
	return &LinkJobRoleCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant),
		JobRoleId:   jobRoleId,
	}
}

type LinkPhoneNumberCommand struct {
	eventstore.BaseCommand
	PhoneNumberId string
	Primary       bool
	Label         string
}

func NewLinkPhoneNumberCommand(objectID, tenant, phoneNumberId, label string, primary bool) *LinkPhoneNumberCommand {
	return &LinkPhoneNumberCommand{
		BaseCommand:   eventstore.NewBaseCommand(objectID, tenant),
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

func NewLinkEmailCommand(objectID, tenant, emailId, label string, primary bool) *LinkEmailCommand {
	return &LinkEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(objectID, tenant),
		EmailId:     emailId,
		Primary:     primary,
		Label:       label,
	}
}
