package agent_listeners

import (
	"context"
	"errors"
	"fmt"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
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
	return "EmailReplyReceived"
}

func (l *EmailReplyReceivedListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *EmailReplyReceivedListener) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{}
}

func (l *EmailReplyReceivedListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailReplyReceivedListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	data, err := events.DecodeEventData[dto.CompanyIdentified](ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	var errs error
	if data.AgentExecutionId != "" {
		err = l.handleGoalAchieved(ctx, data.AgentExecutionId)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	err = l.handleExecution(ctx, data.AgentExecutionId)
	if err != nil {
		tracing.TraceErr(span, err)
		errs = multierr.Append(errs, err)
	}

	return errs
}

func (l *EmailReplyReceivedListener) handleExecution(ctx context.Context, orgID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailReplyReceivedListener.handleExecution")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	activeAgents := l.lookupActiveAgents(ctx)
	if len(activeAgents) == 0 {
		err := errors.New("No agent types configured for company identified listener")
		tracing.TraceErr(span, err)
		return err
	}

	var startParams any

	var errs error
	for _, agent := range activeAgents {
		err := l.run(ctx, agent, startParams)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (l *EmailReplyReceivedListener) run(ctx context.Context, agent postgres_entity.Agent, startParams any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailReplyReceivedListener.run")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	initialParams, err := utils.StructToMap(startParams)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	err = l.agentRunnerService.Run(ctx, agent, l.Type().String(), initialParams)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (l *EmailReplyReceivedListener) handleGoalAchieved(ctx context.Context, agentExecutionId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailReplyReceivedListener.handleGoalAchieved")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	span.LogFields(log.String("agentExecutionId", agentExecutionId))

	var agentExecution *postgres_entity.AgentExecution
	var err error

	agentExecution, err = l.postgresRepositories.AgentExecutionRepository.GetById(ctx, agentExecutionId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err

	}
	if agentExecution == nil {
		err = fmt.Errorf("agent execution not found")
		tracing.TraceErr(span, err)
		return err
	}

	// update execution with goal achieved
	_, err = l.postgresRepositories.AgentExecutionRepository.Completed(ctx, agentExecution.ID, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (l *EmailReplyReceivedListener) lookupActiveAgents(ctx context.Context) []postgres_entity.Agent {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EmailReplyReceivedListener.lookupActiveAgents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	agents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, l.ExecutingAgents())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	return agents
}
