package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain"
	contactCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/commands"
	emailCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/commands"
	interactionEventCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/commands"
	jobRoleCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands"
	locationCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/commands"
	logentrcmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command_handler"
	orgcmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	phoneNumberCommands "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/commands"
	usercmdhnd "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
)

func CreateCommands(log logger.Logger, cfg *config.Config, aggregateStore eventstore.AggregateStore, repositories *repository.Repositories) *domain.Commands {
	return &domain.Commands{
		ContactCommands:          contactCommands.NewContactCommands(log, cfg, aggregateStore),
		OrganizationCommands:     orgcmdhnd.NewOrganizationCommands(log, cfg, aggregateStore, repositories),
		InteractionEventCommands: interactionEventCommands.NewInteractionEventCommands(log, cfg, aggregateStore),
		PhoneNumberCommands:      phoneNumberCommands.NewPhoneNumberCommands(log, cfg, aggregateStore),
		LocationCommands:         locationCommands.NewLocationCommands(log, cfg, aggregateStore),
		EmailCommands:            emailCommands.NewEmailCommands(log, cfg, aggregateStore),
		UserCommands:             usercmdhnd.NewUserCommands(log, cfg, aggregateStore),
		JobRoleCommands:          jobRoleCommands.NewJobRoleCommands(log, cfg, aggregateStore),
		LogEntryCommands:         logentrcmdhnd.NewLogEntryCommands(log, cfg, aggregateStore),
	}
}
