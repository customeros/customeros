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

type NewLeadListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   interfaces.AgentRunnerService
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*NewLeadListener)(nil)
	_ interfaces.EventListener        = (*NewLeadListener)(nil)
)

func NewNewLeadListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *NewLeadListener {
	return &NewLeadListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.NewLead](), // subscribed event
			events.QueueAgents,                 // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

func (l *NewLeadListener) Type() enum.AgentListenerEvent {
	return enum.EventNewLead
}

func (l *NewLeadListener) Name() string {
	return "New leads added to the app"
}

func (l *NewLeadListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *NewLeadListener) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentICPQualifier,
	}
}

func (l *NewLeadListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NewLeadListener.Handle")
	defer spans.Finish()

	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return l.handleExecution(ctx, event.Event.EntityId)
}

func (l *NewLeadListener) handleExecution(ctx context.Context, orgID string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NewLeadListener.handleExecution")
	defer spans.Finish()

	activeAgents := l.lookupActiveAgents(ctx)
	if len(activeAgents) == 0 {
		spans.LogKV("result", "No active agents found")
		return nil
	}

	message := struct {
		OrganizationID string
	}{
		OrganizationID: orgID,
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

func (l *NewLeadListener) lookupActiveAgents(ctx context.Context) []postgres_entity.Agent {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NewLeadListener.lookupActiveAgents")
	defer spans.Finish()

	agents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, l.ExecutingAgents())
	if err != nil {
		spans.TraceError(err)
		return nil
	}
	return agents
}
