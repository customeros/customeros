package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	commentcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/command_handler"
	contactcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contact/command_handler"
	contractcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/contract/command_handler"
	emailcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/email/command_handler"
	iecmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/command_handler"
	iscmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/command_handler"
	issuecmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/command_handler"
	jobrolecmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands"
	locationcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/command_handler"
	logentrycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command_handler"
	opportunitycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	organizationcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	phonenumbercmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command_handler"
	servicelineitemcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/command_handler"
	usercmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/repository"
)

type CommandHandlers struct {
	Contact            *contactcmdhandler.CommandHandlers
	Organization       *organizationcmdhandler.CommandHandlers
	PhoneNumber        *phonenumbercmdhandler.CommandHandlers
	Email              *emailcmdhandler.CommandHandlers
	User               *usercmdhandler.CommandHandlers
	Location           *locationcmdhandler.CommandHandlers
	JobRole            *jobrolecmdhandler.CommandHandlers
	InteractionEvent   *iecmdhandler.CommandHandlers
	InteractionSession *iscmdhandler.CommandHandlers
	LogEntry           *logentrycmdhandler.CommandHandlers
	Issue              *issuecmdhandler.CommandHandlers
	Comment            *commentcmdhandler.CommandHandlers
	Opportunity        *opportunitycmdhandler.CommandHandlers
	Contract           *contractcmdhandler.CommandHandlers
	ServiceLineItem    *servicelineitemcmdhandler.CommandHandlers
}

func NewCommandHandlers(log logger.Logger,
	cfg *config.Config,
	aggregateStore eventstore.AggregateStore,
	repositories *repository.Repositories) *CommandHandlers {

	return &CommandHandlers{
		Contact:            contactcmdhandler.NewCommandHandlers(log, aggregateStore),
		Organization:       organizationcmdhandler.NewCommandHandlers(log, cfg, aggregateStore, repositories),
		InteractionEvent:   iecmdhandler.NewCommandHandlers(log, aggregateStore),
		InteractionSession: iscmdhandler.NewCommandHandlers(log, aggregateStore),
		PhoneNumber:        phonenumbercmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
		Location:           locationcmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
		Email:              emailcmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
		User:               usercmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
		JobRole:            jobrolecmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
		LogEntry:           logentrycmdhandler.NewCommandHandlers(log, aggregateStore),
		Issue:              issuecmdhandler.NewCommandHandlers(log, aggregateStore),
		Comment:            commentcmdhandler.NewCommandHandlers(log, aggregateStore),
		Opportunity:        opportunitycmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
		Contract:           contractcmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
		ServiceLineItem:    servicelineitemcmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
	}
}
