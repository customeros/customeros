package commands

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain"
	commentcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/command_handler"
	contactcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/command_handler"
	emailcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command_handler"
	iecmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command_handler"
	issuecmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/command_handler"
	jobrolecmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands"
	locationcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/command_handler"
	logentrycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command_handler"
	organizationcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	phonenumbercmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command_handler"
	usercmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
)

func InitCommandHandlers(log logger.Logger,
	cfg *config.Config,
	aggregateStore eventstore.AggregateStore,
	repositories *repository.Repositories) *domain.CommandHandlers {

	return &domain.CommandHandlers{
		Contact:          contactcmdhandler.NewContactCommandHandlers(log, aggregateStore),
		Organization:     organizationcmdhandler.NewOrganizationCommands(log, cfg, aggregateStore, repositories),
		InteractionEvent: iecmdhandler.NewInteractionEventCommandHandlers(log, aggregateStore),
		PhoneNumber:      phonenumbercmdhandler.NewPhoneNumberCommands(log, cfg, aggregateStore),
		Location:         locationcmdhandler.NewLocationCommands(log, cfg, aggregateStore),
		Email:            emailcmdhandler.NewEmailCommandHandlers(log, cfg, aggregateStore),
		User:             usercmdhandler.NewUserCommands(log, cfg, aggregateStore),
		JobRole:          jobrolecmdhandler.NewJobRoleCommands(log, cfg, aggregateStore),
		LogEntry:         logentrycmdhandler.NewLogEntryCommands(log, aggregateStore),
		Issue:            issuecmdhandler.NewIssueCommandHandlers(log, aggregateStore),
		Comment:          commentcmdhandler.NewCommentCommandHandlers(log, aggregateStore),
	}
}
