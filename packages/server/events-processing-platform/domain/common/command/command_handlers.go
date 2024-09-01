package command

import (
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/config"
	issuecmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/command_handler"
	jobrolecmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/job_role/commands"
	locationcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/location/command_handler"
	logentrycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/command_handler"
	masterplancmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/master_plan/command_handler"
	opportunitycmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/opportunity/command_handler"
	organizationcmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization/command_handler"
	orgplanevents "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/organization_plan/event_handler"
	phonenumbercmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/phone_number/command_handler"
	usercmdhandler "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/user/command_handler"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventbuffer"
	"github.com/openline-ai/openline-customer-os/packages/server/events/eventstore"
)

type CommandHandlers struct {
	Organization     *organizationcmdhandler.CommandHandlers
	PhoneNumber      *phonenumbercmdhandler.CommandHandlers
	User             *usercmdhandler.CommandHandlers
	Location         *locationcmdhandler.CommandHandlers
	JobRole          *jobrolecmdhandler.CommandHandlers
	LogEntry         *logentrycmdhandler.CommandHandlers
	Issue            *issuecmdhandler.CommandHandlers
	Opportunity      *opportunitycmdhandler.CommandHandlers
	MasterPlan       *masterplancmdhandler.CommandHandlers
	OrganizationPlan *orgplanevents.EventHandlers
}

func NewCommandHandlers(log logger.Logger,
	cfg *config.Config,
	aggregateStore eventstore.AggregateStore,
	ebs *eventbuffer.EventBufferStoreService,
) *CommandHandlers {

	return &CommandHandlers{
		Organization:     organizationcmdhandler.NewCommandHandlers(log, cfg, aggregateStore, ebs),
		PhoneNumber:      phonenumbercmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
		Location:         locationcmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
		User:             usercmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
		JobRole:          jobrolecmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
		LogEntry:         logentrycmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
		Issue:            issuecmdhandler.NewCommandHandlers(log, aggregateStore),
		Opportunity:      opportunitycmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
		MasterPlan:       masterplancmdhandler.NewCommandHandlers(log, cfg, aggregateStore),
		OrganizationPlan: orgplanevents.NewEventHandlers(log, cfg, aggregateStore),
	}
}
