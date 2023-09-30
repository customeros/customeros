package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain"
	contact_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/commands"
	email_command_handler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command_handler"
	interaction_event_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/commands"
	job_role_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands"
	location_commands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/commands"
	log_entry_command_handler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command_handler"
	org_command_handler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	phone_num_command_handler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command_handler"
	user_command_handler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
)

func CreateCommands(log logger.Logger, cfg *config.Config, aggregateStore eventstore.AggregateStore, repositories *repository.Repositories) *domain.Commands {
	return &domain.Commands{
		ContactCommands:          contact_commands.NewContactCommands(log, cfg, aggregateStore),
		OrganizationCommands:     org_command_handler.NewOrganizationCommands(log, cfg, aggregateStore, repositories),
		InteractionEventCommands: interaction_event_commands.NewInteractionEventCommands(log, cfg, aggregateStore),
		PhoneNumberCommands:      phone_num_command_handler.NewPhoneNumberCommands(log, cfg, aggregateStore),
		LocationCommands:         location_commands.NewLocationCommands(log, cfg, aggregateStore),
		EmailCommands:            email_command_handler.NewEmailCommands(log, cfg, aggregateStore),
		UserCommands:             user_command_handler.NewUserCommands(log, cfg, aggregateStore),
		JobRoleCommands:          job_role_commands.NewJobRoleCommands(log, cfg, aggregateStore),
		LogEntryCommands:         log_entry_command_handler.NewLogEntryCommands(log, cfg, aggregateStore),
	}
}
