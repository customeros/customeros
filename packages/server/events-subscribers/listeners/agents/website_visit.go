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
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type WebsiteVisitListener struct {
	events.BaseEventListener
	dependencies *model.DependencyContainer
}

func NewWebsiteVisitListener(logger logger.Logger, deps *model.DependencyContainer) interfaces.EventListener {
	return &WebsiteVisitListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.WebsiteVisit](), // subscribed event
			events.QueueAgents,                      // listening on Agents queue
		),
		dependencies: deps,
	}
}

// Add all Agent types subscribed to this event here
func (h *WebsiteVisitListener) subscribedAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentVisitorID,
	}
}

func (l *WebsiteVisitListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteVisitListener.Handle")
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

func (l *WebsiteVisitListener) handle(ctx context.Context, message *dto.WebsiteVisit) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteVisitListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	activeAgents := l.lookupActiveAgents(ctx)
	if activeAgents == nil || len(activeAgents) == 0 {
		return nil
	}

	var errs error
	for _, agent := range activeAgents {
		// replaced route with generic mapping
		initialParams, err := utils.StructToMap(message)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
		err = l.dependencies.CommonServices.AgentRunnerService.Run(ctx, agent, message.Type(), initialParams)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (h *WebsiteVisitListener) lookupActiveAgents(ctx context.Context) []postgres_entity.Agents {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteVisitListener.lookupActiveAgents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	agents, err := h.dependencies.PostgresRepositories.AgentsRepository.GetActiveAgentsByTypes(ctx, h.subscribedAgents())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	return agents
}

func (l *WebsiteVisitListener) validateMessage(ctx context.Context, event *dto.Event) (*dto.WebsiteVisit, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteVisitListener.validateMessage")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	message, ok := event.Event.Data.(*dto.WebsiteVisit)
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

	if message.SessionID == "" {
		err := errors.New("SessionID not set on event")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if message.IPAddress == "" {
		err := errors.New("IPAddress not set on event")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return message, nil
}
