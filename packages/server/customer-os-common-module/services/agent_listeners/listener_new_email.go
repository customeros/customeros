package agent_listeners

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type NewEmailListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   interfaces.AgentRunnerService
}

type NewEmailConfig struct {
	Emails ConfigMultipleValuesWithObject `json:"emails"`
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*NewEmailListener)(nil)
	_ interfaces.EventListener        = (*NewEmailListener)(nil)
)

func NewNewEmailListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *NewEmailListener {
	return &NewEmailListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.NewEmail](), // subscribed event
			events.QueueAgents,                  // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

func (l *NewEmailListener) Type() enum.AgentListenerEvent {
	return enum.EventNewEmail
}

func (l *NewEmailListener) Name() string {
	return "New email received"
}

func (l *NewEmailListener) DefaultConfig() any {
	n := &NewEmailConfig{}
	n.Emails = ConfigMultipleValuesWithObject{}
	n.Emails.Value = []interface{}{}
	return n
}

func (l *NewEmailListener) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentEmailKeeper,
	}
}

func (l *NewEmailListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NewEmailListener.Handle")
	defer spans.Finish()

	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return l.handleExecution(ctx, event.Event.EntityId)
}

func (l *NewEmailListener) handleExecution(ctx context.Context, rawEmailId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NewEmailListener.handleExecution")
	defer spans.Finish()

	activeAgents := l.lookupActiveAgents(ctx)
	if len(activeAgents) == 0 {
		spans.LogKV("result", "No active agents found")
		return nil
	}

	message := struct {
		RawEmailId string
	}{
		RawEmailId: rawEmailId,
	}

	var errs error
	for _, agent := range activeAgents {

		initialParams, err := utils.StructToMap(message)
		if err != nil {
			spans.TraceError(err)
			errs = multierr.Append(errs, err)
		}
		_, err = l.agentRunnerService.Run(ctx, agent, l.Type().String(), initialParams, nil)
		if err != nil {
			spans.TraceError(err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (l *NewEmailListener) lookupActiveAgents(ctx context.Context) []postgres_entity.Agent {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NewEmailListener.lookupActiveAgents")
	defer spans.Finish()

	agents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, l.ExecutingAgents())
	if err != nil {
		spans.TraceError(err)
		return nil
	}
	return agents
}
