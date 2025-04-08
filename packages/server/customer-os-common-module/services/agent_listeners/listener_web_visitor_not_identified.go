package agent_listeners

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
)

type WebVisitorNotIdentifiedListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*WebVisitorNotIdentifiedListener)(nil)
	_ interfaces.EventListener        = (*WebVisitorNotIdentifiedListener)(nil)
)

func NewWebVisitorNotIdentifiedListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
) *WebVisitorNotIdentifiedListener {
	return &WebVisitorNotIdentifiedListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.WebVisitorNotIdentified](), // subscribed event
			events.QueueAgents, // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
	}
}

func (l *WebVisitorNotIdentifiedListener) Type() enum.AgentListenerEvent {
	return enum.EventWebVisitorNotIdentified
}

func (l *WebVisitorNotIdentifiedListener) Name() string {
	return "Web visitor not identified"
}

func (l *WebVisitorNotIdentifiedListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *WebVisitorNotIdentifiedListener) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{}
}

func (l *WebVisitorNotIdentifiedListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebVisitorNotIdentifiedListener.Handle")
	defer spans.Finish()

	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	data, err := events.DecodeEventData[dto.WebVisitorNotIdentified](ctx, event)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	if data.AgentExecutionId != "" {
		return l.handleGoalAchieved(ctx, data.AgentExecutionId)
	}
	return nil
}

func (l *WebVisitorNotIdentifiedListener) handleGoalAchieved(ctx context.Context, agentExecutionId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebVisitorNotIdentifiedListener.handleGoalAchieved")
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
	_, err = l.postgresRepositories.AgentExecutionRepository.Completed(ctx, agentExecution.ID, utils.FalsePtr())
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}
