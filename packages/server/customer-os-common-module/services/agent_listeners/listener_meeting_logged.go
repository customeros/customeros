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

type MeetingLoggedListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*MeetingLoggedListener)(nil)
	_ interfaces.EventListener        = (*MeetingLoggedListener)(nil)
)

func NewMeetingLoggedListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
) *MeetingLoggedListener {
	return &MeetingLoggedListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.MeegingLogged](), // subscribed event
			events.QueueAgents,                       // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
	}
}

func (l *MeetingLoggedListener) Type() enum.AgentListenerEvent {
	return enum.EventMeetingLogged
}

func (l *MeetingLoggedListener) Name() string {
	return "Meeting is logged"
}

func (l *MeetingLoggedListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *MeetingLoggedListener) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{}
}

func (l *MeetingLoggedListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MeetingLoggedListener.Handle")
	defer spans.Finish()

	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	data, err := events.DecodeEventData[dto.MeegingLogged](ctx, event)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	if data.AgentExecutionId != "" {
		return l.handleGoalAchieved(ctx, data.AgentExecutionId)
	}

	return nil
}

func (l *MeetingLoggedListener) handleGoalAchieved(ctx context.Context, agentExecutionId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "MeetingLoggedListener.handle")
	defer spans.Finish()

	spans.LogKV("agentExecutionId", agentExecutionId)

	var agentExecution *postgres_entity.AgentExecution
	var err error
	if agentExecutionId != "" {
		agentExecution, err = l.postgresRepositories.AgentExecutionRepository.GetById(ctx, agentExecutionId)
		if err != nil {
			spans.TraceError(err)
			return err
		}
	}
	if agentExecution == nil {
		err = fmt.Errorf("agent execution not found")
		spans.TraceError(err)
		return err
	}

	// update execution with goal achieved
	_, err = l.postgresRepositories.AgentExecutionRepository.Completed(ctx, agentExecution.ID, utils.TruePtr())
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}
