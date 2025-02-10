package agent_listeners

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
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
	return "Web Visitor Identified"
}

func (l *WebVisitorIdentifiedListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *WebVisitorIdentifiedListener) SubscribedAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentWebVisitorIdentifier,
	}
}

func (l *WebVisitorIdentifiedListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebVisitorIdentifiedListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	// TODO implement
	return nil

	//return l.handleGoalAchieved(ctx, data.AgentExecutionId)
}

//func (l *WebVisitorIdentifiedListener) handleGoalAchieved(ctx context.Context, agentExecutionId string) error {
//	span, ctx := opentracing.StartSpanFromContext(ctx, "WebVisitorIdentifiedListener.handleGoalAchieved")
//	defer span.Finish()
//	tracing.SetDefaultListenerSpanTags(ctx, span)
//	span.LogFields(log.String("agentExecutionId", agentExecutionId))
//
//	var agentExecution *postgres_entity.AgentExecution
//	var err error
//	if agentExecutionId != "" {
//		agentExecution, err = l.postgresRepositories.AgentExecutionRepository.GetById(ctx, agentExecutionId)
//		if err != nil {
//			tracing.TraceErr(span, err)
//			return err
//		}
//	}
//	if agentExecution == nil {
//		err = fmt.Errorf("agent execution not found")
//		tracing.TraceErr(span, err)
//		return err
//	}
//
//	// update execution with goal achieved
//	if agentExecution.Status == enum.AgentExecutionCompleted {
//		err = fmt.Errorf("agent execution already completed or failed")
//		tracing.TraceErr(span, err)
//		return err
//	}
//	agentExecution.GoalAchieved = true
//	_, err = l.postgresRepositories.AgentExecutionRepository.Completed(ctx, agentExecution.ID, true)
//	if err != nil {
//		tracing.TraceErr(span, err)
//		return err
//	}
//
//	return nil
//}
