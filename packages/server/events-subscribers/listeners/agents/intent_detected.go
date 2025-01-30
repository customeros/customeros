package agent_listeners

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type IntentDetectedListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewIntentDetectedListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &IntentDetectedListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.IntentDetected](), // subscribed event
			events.QueueAgents,                        // listening on Agents queue
		),
		dependencies: deps,
	}
}

// Add all Agent types subscribed to this event here
func (h *IntentDetectedListener) subscribedAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentVisitorID,
	}
}

func (l *IntentDetectedListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IntentDetectedListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	message, err := l.validateMessage(ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	ctx = common.SetTenantInContext(ctx, message.Tenant)

	return l.handle(ctx, message)
}

func (l *IntentDetectedListener) handle(ctx context.Context, message *dto.IntentDetected) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IntentDetectedListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	activeAgents := l.lookupActiveAgents(ctx)
	if activeAgents == nil || len(activeAgents) == 0 {
		return nil
	}

	switch message.IntentType {
	case enum.IntentSupportRequired:
		return l.dependencies.CommonServices.SupportAgent.ProcessIntentEvent(ctx, message)

	default:
		err := errors.New("IntentType not supported")
		tracing.TraceErr(span, err)
		return err
	}
}

func (h *IntentDetectedListener) lookupActiveAgents(ctx context.Context) []postgres_entity.Agents {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IntentDetectedListener.lookupActiveAgents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	agents, err := h.dependencies.PostgresRepositories.AgentsRepository.GetActiveAgentsByTypes(ctx, h.subscribedAgents())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	return agents
}

func (l *IntentDetectedListener) validateMessage(ctx context.Context, event *dto.Event) (*dto.IntentDetected, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IntentDetectedListener.validateMessage")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message, ok := event.Event.Data.(*dto.IntentDetected)
	if !ok {
		err := fmt.Errorf("expected WebsiteVisit, got %T", event.Event.Data)
		tracing.TraceErr(span, err)
		return nil, err
	}

	if message.Tenant == "" {
		err := errors.New("Tenant not set on event")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return message, nil
}
