package agent_listeners

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
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

	return l.handleGoalAchieved(ctx, event.Event.EntityId, event.Event.AgentEventName)
}

func (l *IcpFitListener) handleGoalAchieved(ctx context.Context, orgId, eventName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IcpFitListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	subscribedAgents := l.SubscribedAgents()
	if len(subscribedAgents) == 0 {
		err := errors.New("No agent types configured for company stage lead listener")
		tracing.TraceErr(span, err)
		return err
	}

	activeAgents := l.lookupActiveAgents(ctx, subscribedAgents)
	var errs error
	for _, agent := range activeAgents {
		fmt.Println(agent)
		// todo mark execution as goal achieved
	}

	return errs
}

func (l *IcpFitListener) lookupActiveAgents(ctx context.Context, agentTypes []enum.AgentType) []postgres_entity.Agent {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IcpFitListener.lookupActiveAgents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	span.LogFields(log.String("agentTypes", fmt.Sprintf("%v", agentTypes)))

	agents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, agentTypes)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	return agents
}
