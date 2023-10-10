package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain"
	contactcmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/command_handler"
	emailcmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command_handler"
	iecmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/commands"
	jobrolecmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands"
	locationcmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/command_handler"
	logentrycmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command_handler"
	orgcmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	phonecmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command_handler"
	usercmd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
)

func CreateCommands(log logger.Logger, cfg *config.Config, aggregateStore eventstore.AggregateStore, repositories *repository.Repositories) *domain.Commands {
	return &domain.Commands{
		ContactCommands:          contactcmd.NewContactCommands(log, cfg, aggregateStore),
		OrganizationCommands:     orgcmd.NewOrganizationCommands(log, cfg, aggregateStore, repositories),
		InteractionEventCommands: iecmd.NewInteractionEventCommands(log, cfg, aggregateStore),
		PhoneNumberCommands:      phonecmd.NewPhoneNumberCommands(log, cfg, aggregateStore),
		LocationCommands:         locationcmd.NewLocationCommands(log, cfg, aggregateStore),
		EmailCommands:            emailcmd.NewEmailCommands(log, cfg, aggregateStore),
		UserCommands:             usercmd.NewUserCommands(log, cfg, aggregateStore),
		JobRoleCommands:          jobrolecmd.NewJobRoleCommands(log, cfg, aggregateStore),
		LogEntryCommands:         logentrycmd.NewLogEntryCommands(log, cfg, aggregateStore),
	}
}
