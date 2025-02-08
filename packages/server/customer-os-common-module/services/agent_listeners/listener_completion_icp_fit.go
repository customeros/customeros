package agent_listeners

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type IcpFitListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*IcpFitListener)(nil)
	_ interfaces.EventListener        = (*IcpFitListener)(nil)
)

func NewIcpFitListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
) *IcpFitListener {
	return &IcpFitListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.IcpFit](), // subscribed event
			events.QueueAgents,                // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
	}
}

func (l *IcpFitListener) Type() enum.AgentListenerEvent {
	return enum.EventICPFit
}

func (l *IcpFitListener) Name() string {
	return "ICP Fit"
}

func (l *IcpFitListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *IcpFitListener) SubscribedAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentICPQualifier,
	}
}

func (l *IcpFitListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IcpFitListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	data, err := events.DecodeEventData[dto.IcpFit](ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return l.handleGoalAchieved(ctx, data.AgentExecutionId)
}

func (l *IcpFitListener) handleGoalAchieved(ctx context.Context, agentExecutionId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IcpFitListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	span.LogFields(log.String("agentExecutionId", agentExecutionId))

	var agentExecution *postgres_entity.AgentExecution
	var err error
	if agentExecutionId != "" {
		agentExecution, err = l.postgresRepositories.AgentExecutionRepository.GetById(ctx, agentExecutionId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}
	if agentExecution == nil {
		err = fmt.Errorf("agent execution not found")
		tracing.TraceErr(span, err)
		return err
	}

	// update execution with goal achieved
	if agentExecution.Status == enum.AgentExecutionCompleted || agentExecution.Status == enum.AgentExecutionFail {
		err = fmt.Errorf("agent execution already completed or failed")
		tracing.TraceErr(span, err)
		return err
	}
	agentExecution.GoalAchieved = true
	_, err = l.postgresRepositories.AgentExecutionRepository.Update(ctx, agentExecution.ID, utils.NowPtr(), "", true)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
