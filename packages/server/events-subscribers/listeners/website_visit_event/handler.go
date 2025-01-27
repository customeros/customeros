package website_visit_event

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
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

	if event.SessionID == "" {
		err := errors.New("SessionID cannot be empty")
		return nil, err
	}

	if event.IPAddress == "" {
		err := errors.New("IP address cannot be empty")
		return nil, err
	}

	return &WebsiteVisitEventHandler{
		dependencies: dependencies,
		event:        event,
	}, nil
}

// Add all subscribed Agents here
func (h *WebsiteVisitEventHandler) subscribedAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentVisitorID,
	}
}

func (h *WebsiteVisitEventHandler) Handle(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteVisitEventHandler.Handle")
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
		// replaced route with generic mapping
		initialParams, err := utils.StructToMap(h.event)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
		err = h.dependencies.CommonServices.AgentRunnerService.Run(ctx, agent, h.event.Type(), initialParams)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (h *WebsiteVisitEventHandler) lookupActiveAgents(ctx context.Context) []postgres_entity.Agents {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteVisitEventHandler.lookupActiveAgents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	agents, err := h.dependencies.PostgresRepositories.AgentsRepository.GetActiveAgentsByTypes(ctx, h.subscribedAgents())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	return agents
}
