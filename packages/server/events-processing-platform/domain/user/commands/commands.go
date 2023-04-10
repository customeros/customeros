package commands

import (
	common_models "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/models"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UserCoreFields struct {
	Name      string
	FirstName string
	LastName  string
}

type UpsertUserCommand struct {
	eventstore.BaseCommand
	Tenant     string
	CoreFields UserCoreFields
	Source     common_models.Source
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

func UpsertUserCommandToUserDto(command *UpsertUserCommand) *models.UserDto {
	return &models.UserDto{
		ID:        command.AggregateID,
		Tenant:    command.Tenant,
		Name:      command.CoreFields.Name,
		FirstName: command.CoreFields.FirstName,
		LastName:  command.CoreFields.LastName,
		Source:    command.Source,
		CreatedAt: command.CreatedAt,
		UpdatedAt: command.UpdatedAt,
	}
}

func NewUpsertUserCommand(aggregateID, tenant, source, sourceOfTruth, appSource string, coreFields UserCoreFields, createdAt, updatedAt *time.Time) *UpsertUserCommand {
	return &UpsertUserCommand{
		BaseCommand: eventstore.NewBaseCommand(aggregateID),
		Tenant:      tenant,
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

type LinkPhoneNumberCommand struct {
	eventstore.BaseCommand
	Tenant        string
	PhoneNumberId string
	Primary       bool
	Label         string
}

func NewLinkPhoneNumberCommand(aggregateID, tenant, phoneNumberId, label string, primary bool) *LinkPhoneNumberCommand {
	return &LinkPhoneNumberCommand{
		BaseCommand:   eventstore.NewBaseCommand(aggregateID),
		Tenant:        tenant,
		PhoneNumberId: phoneNumberId,
		Primary:       primary,
		Label:         label,
	}
}

type LinkEmailCommand struct {
	eventstore.BaseCommand
	Tenant  string
	EmailId string
	Primary bool
	Label   string
}

func NewLinkEmailCommand(aggregateID, tenant, emailId, label string, primary bool) *LinkEmailCommand {
	return &LinkEmailCommand{
		BaseCommand: eventstore.NewBaseCommand(aggregateID),
		Tenant:      tenant,
		EmailId:     emailId,
		Primary:     primary,
		Label:       label,
	}
}
