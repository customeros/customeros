package website_visit_event

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/events-subscribers/model"
)

type WebsiteVisitEventHandler struct {
	dependencies *model.DependencyContainer
	event        *data_fields.WebsiteVisitEvent
}

func NewWebsiteVisitEventHandler(dependencies *model.DependencyContainer, event *data_fields.WebsiteVisitEvent) (*WebsiteVisitEventHandler, error) {
	if event == nil {
		err := errors.New("Event cannot be nil")
		return nil, err
	}

	if event.Tenant == "" {
		err := errors.New("Tenant cannot be empty")
		return nil, err
	}

	return &WebsiteVisitEventHandler{
		dependencies: dependencies,
		event:        event,
	}, nil
}

// Add all subscribed Agents here
func (h *WebsiteVisitEventHandler) subscribedAgents(ctx context.Context) []enum.AgentType {
	return []enum.AgentType{
		enum.AgentVisitorID,
	}
}

func (h *WebsiteVisitEventHandler) Handle(ctx context.Context) error {
	span, ctx := tracing.StartTracerSpan(ctx, "WebsiteVisitEventHandler.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant: h.event.Tenant,
	})

	activeAgents := h.lookupActiveAgents(ctx)
	if activeAgents == nil || len(activeAgents) == 0 {
		return nil
	}

	var errs error
	for _, agent := range activeAgents {
		err := h.route(ctx, agent)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (h *WebsiteVisitEventHandler) route(ctx context.Context, agent postgres_entity.Agents) error {
	span, ctx := tracing.StartTracerSpan(ctx, "WebsiteVisitEventHandler.execute")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	if agent.Type == "" {
		err := errors.New("agent.Type is empty")
		tracing.TraceErr(span, err)
		return err
	}

	if agent.ID == "" {
		err := errors.New("agent.ID is empty")
		tracing.TraceErr(span, err)
		return err
	}

	// run agent
	agentType, err := enum.GetAgentType(agent.Type)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	switch agentType {
	case enum.AgentVisitorID:
		return h.dependencies.CommonServices.AgentVisitorIDService.Run(ctx, agent.ID, h.event)

	default:
		err := errors.New("Unsupported agent type")
		span.LogKV("agentType", agentType.String())
		tracing.TraceErr(span, err)
		return err
	}
}

func (h *WebsiteVisitEventHandler) lookupActiveAgents(ctx context.Context) []postgres_entity.Agents {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteVisitEventHandler.lookupActiveAgents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	agents, err := h.dependencies.PostgresRepositories.AgentsRepository.FindAllFromAgentsList(ctx, h.subscribedAgents(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	if agents == nil {
		return nil
	}
	return agents
}
