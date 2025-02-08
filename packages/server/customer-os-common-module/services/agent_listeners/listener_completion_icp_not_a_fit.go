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

type IcpNotAFitListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*IcpNotAFitListener)(nil)
	_ interfaces.EventListener        = (*IcpNotAFitListener)(nil)
)

func NewIcpNotAFitListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
) *IcpNotAFitListener {
	return &IcpNotAFitListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.IcpNotAFit](), // subscribed event
			events.QueueAgents,                    // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
	}
}

func (l *IcpNotAFitListener) Type() enum.AgentListenerEvent {
	return enum.EventICPNotAFit
}

func (l *IcpNotAFitListener) Name() string {
	return "ICP Not A Fit"
}

func (l *IcpNotAFitListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *IcpNotAFitListener) SubscribedAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentICPQualifier,
	}
}

func (l *IcpNotAFitListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IcpNotAFitListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	agentExecutionId := ""
	data, ok := event.Event.Data.(map[string]interface{})
	if ok {
		agentExecutionId = data["agentExecutionId"].(string)
	}

	return l.handleGoalAchieved(ctx, agentExecutionId)
}

func (l *IcpNotAFitListener) handleGoalAchieved(ctx context.Context, agentExecutionId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IcpNotAFitListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

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
