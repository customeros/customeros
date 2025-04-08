package agent_listeners

import (
	"context"
	"errors"
	"fmt"

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

type TEMPLATEListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   interfaces.AgentRunnerService
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*TEMPLATEListener)(nil)
	_ interfaces.EventListener        = (*TEMPLATEListener)(nil)
)

func NewTEMPLATEListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *TEMPLATEListener {
	return &TEMPLATEListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.TEMPLATE](), // subscribed event
			events.QueueAgents,                  // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

func (l *TEMPLATEListener) Type() enum.AgentListenerEvent {
	return enum.EventTEMPLATE
}

func (l *TEMPLATEListener) Name() string {
	return "TEMPLATE"
}

func (l *TEMPLATEListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *TEMPLATEListener) ExecutingAgents() []enum.AgentType {
	// add execution agents
	return []enum.AgentType{}
}

func (l *TEMPLATEListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TEMPLATEListener.Handle")
	defer spans.Finish()

	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	data, err := events.DecodeEventData[dto.TEMPLATE](ctx, event)
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

func (l *TEMPLATEListener) handleExecution(ctx context.Context, orgID string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TEMPLATEListener.handleExecution")
	defer spans.Finish()

	activeAgents := l.lookupActiveAgents(ctx)
	if len(activeAgents) == 0 {
		err := errors.New("No agent types configured for company identified listener")
		spans.TraceError(err)
		return err
	}

	// replace with start params
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

func (l *TEMPLATEListener) run(ctx context.Context, agent postgres_entity.Agent, startParams any) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TEMPLATEListener.run")
	defer spans.Finish()

	initialParams, err := utils.StructToMap(startParams)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	err = l.agentRunnerService.Run(ctx, agent, l.Type().String(), initialParams)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}

func (l *TEMPLATEListener) handleGoalAchieved(ctx context.Context, agentExecutionId string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TEMPLATEListener.handleGoalAchieved")
	defer spans.Finish()

	spans.LogKV("agentExecutionId", agentExecutionId))

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
	if agentExecution.Status == enum.AgentExecutionCompleted {
		err = fmt.Errorf("agent execution already completed")
		spans.TraceError(err)
		return err
	}
	_, err = l.postgresRepositories.AgentExecutionRepository.Completed(ctx, agentExecution.ID, false)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (l *TEMPLATEListener) lookupActiveAgents(ctx context.Context) []postgres_entity.Agent {
	spans, ctx := telemetry.StartServiceSpan(ctx, "TEMPLATEListener.lookupActiveAgents")
	defer spans.Finish()

	agents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, l.ExecutingAgents())
	if err != nil {
		spans.TraceError(err)
		return nil
	}
	return agents
}
