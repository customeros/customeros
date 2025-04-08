package agent_listeners

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
)

type WebVisitorIdentifiedListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*WebVisitorIdentifiedListener)(nil)
	_ interfaces.EventListener        = (*WebVisitorIdentifiedListener)(nil)
)

func NewWebVisitorIdentifiedListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
) *WebVisitorIdentifiedListener {
	return &WebVisitorIdentifiedListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.WebVisitorIdentified](), // subscribed event
			events.QueueAgents, // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
	}
}

func (l *WebVisitorIdentifiedListener) Type() enum.AgentListenerEvent {
	return enum.EventWebVisitorIdentified
}

func (l *WebVisitorIdentifiedListener) Name() string {
	return "Web visitor identified"
}

func (l *WebVisitorIdentifiedListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *WebVisitorIdentifiedListener) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{}
}

func (l *WebVisitorIdentifiedListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebVisitorIdentifiedListener.Handle")
	defer spans.Finish()

	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	data, err := events.DecodeEventData[dto.WebVisitorIdentified](ctx, event)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	if data.AgentExecutionId != "" {
		return l.handleGoalAchieved(ctx, data.AgentExecutionId)
	}
	return nil
}

func (l *WebVisitorIdentifiedListener) handleGoalAchieved(ctx context.Context, agentExecutionId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebVisitorIdentifiedListener.handleGoalAchieved")
	defer spans.Finish()

	spans.LogKV("agentExecutionId", agentExecutionId)

	var agentExecution *postgres_entity.AgentExecution
	var err error
	agentExecution, err = l.postgresRepositories.AgentExecutionRepository.GetById(ctx, agentExecutionId)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	if agentExecution == nil {
		err = fmt.Errorf("agent execution not found")
		spans.TraceError(err)
		return err
	}

	// update execution with goal achieved
	err = l.postgresRepositories.AgentExecutionRepository.GoalAchieved(ctx, agentExecution.ID, true, nil)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}
