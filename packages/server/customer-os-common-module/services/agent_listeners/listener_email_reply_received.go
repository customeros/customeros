package agent_listeners

import (
	"context"
	"errors"
	"fmt"
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

type EmailReplyReceivedListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   interfaces.AgentRunnerService
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*EmailReplyReceivedListener)(nil)
	_ interfaces.EventListener        = (*EmailReplyReceivedListener)(nil)
)

func NewEmailReplyReceivedListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *EmailReplyReceivedListener {
	return &EmailReplyReceivedListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.EmailReplyReceived](), // subscribed event
			events.QueueAgents,                            // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

func (l *EmailReplyReceivedListener) Type() enum.AgentListenerEvent {
	return enum.EventEmailReplyReceived
}

func (l *EmailReplyReceivedListener) Name() string {
	return "Replies to emails"
}

func (l *EmailReplyReceivedListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *EmailReplyReceivedListener) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{}
}

func (l *EmailReplyReceivedListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EmailReplyReceivedListener.Handle")
	defer spans.Finish()

	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	data, err := events.DecodeEventData[dto.CompanyIdentified](ctx, event)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	var errs error
	if data.AgentExecutionId != "" {
		err = l.handleGoalAchieved(ctx, data.AgentExecutionId)
		if err != nil {
			spans.TraceError(err)
			errs = multierr.Append(errs, err)
		}
	}

	err = l.handleExecution(ctx, data.AgentExecutionId)
	if err != nil {
		spans.TraceError(err)
		errs = multierr.Append(errs, err)
	}

	return errs
}

func (l *EmailReplyReceivedListener) handleExecution(ctx context.Context, orgID string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EmailReplyReceivedListener.handleExecution")
	defer spans.Finish()

	activeAgents := l.lookupActiveAgents(ctx)
	if len(activeAgents) == 0 {
		err := errors.New("No agent types configured for company identified listener")
		spans.TraceError(err)
		return err
	}

	var startParams any

	var errs error
	for _, agent := range activeAgents {
		err := l.run(ctx, agent, startParams)
		if err != nil {
			spans.TraceError(err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (l *EmailReplyReceivedListener) run(ctx context.Context, agent postgres_entity.Agent, startParams any) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EmailReplyReceivedListener.run")
	defer spans.Finish()

	initialParams, err := utils.StructToMap(startParams)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	_, err = l.agentRunnerService.Run(ctx, agent, l.Type().String(), initialParams, nil)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}

func (l *EmailReplyReceivedListener) handleGoalAchieved(ctx context.Context, agentExecutionId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EmailReplyReceivedListener.handleGoalAchieved")
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
	_, err = l.postgresRepositories.AgentExecutionRepository.Completed(ctx, agentExecution.ID, utils.TruePtr())
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (l *EmailReplyReceivedListener) lookupActiveAgents(ctx context.Context) []postgres_entity.Agent {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EmailReplyReceivedListener.lookupActiveAgents")
	defer spans.Finish()

	agents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, l.ExecutingAgents())
	if err != nil {
		spans.TraceError(err)
		return nil
	}
	return agents
}
