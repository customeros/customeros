package website_visit_event

import (
	"context"

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

func NewWebsiteVisitEventHandler(dependencies *model.DependencyContainer, event *data_fields.WebsiteVisitEvent) *WebsiteVisitEventHandler {
	return &WebsiteVisitEventHandler{
		dependencies: dependencies,
		event:        event,
	}
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

	activeAgents := h.lookupActiveAgents(ctx)
	if activeAgents == nil {
		return nil
	}

	var errs error
	for _, agent := range activeAgents {
		err := h.execute(ctx, agent)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (h *WebsiteVisitEventHandler) execute(ctx context.Context, agent postgres_entity.Agents) error {
	span, ctx := tracing.StartTracerSpan(ctx, "WebsiteVisitEventHandler.execute")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	executionID, err := h.createAgentExecutionRecord(ctx, agent)
	if err != nil {
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
		// send
	default:
		err := errors.New("Unsupported agent type")
		span.LogKV("agentType", agentType.String())
		tracing.TraceErr(span, err)
	}

	// update execution record
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

func (h *WebsiteVisitEventHandler) createAgentExecutionRecord(ctx context.Context, agent postgres_entity.Agents) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteVisitEventHandler.createAgentExecutionRecord")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	agentExecutionRecord := postgres_entity.AgentExecution{
		AgentID:      &agent.ID,
		TriggerEvent: h.event.Type(),
		Status:       enum.AgentExecutionRunning.String(),
		StartedAt:    utils.NowPtr(),
	}

	createdRecord, err := h.dependencies.PostgresRepositories.AgentExecutionRepository.Create(ctx, agentExecutionRecord)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if createdRecord == nil {
		err := errors.New("unable to create agent execution record")
		tracing.TraceErr(span, err)
		return "", err
	}

	return createdRecord.ID, nil
}

func (h *WebsiteVisitEventHandler) agentExecutionSuccess(ctx context.Context, executionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteVisitEventHandler.agentExecutionSuccess")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	agentExecutionRecord.Status = enum.AgentExecutionSuccess.String()
	agentExecutionRecord.CompletedAt = utils.NowPtr()
	execution, err = dependencies.PostgresRepositories.AgentExecutionRepository.Update(ctx, agentExecutionRecord)
	if err != nil {
		tracing.TraceErr(span, err)
		loopErr = multierr.Append(loopErr, err)
	}

	return nil
}

func (h *WebsiteVisitEventHandler) agentExecutionError(ctx context.Context, executionID, errorMessage string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteVisitEventHandler.agentExecutionError")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	errorMessage := loopErr.Error()
	agentExecutionRecord.ErrorMessage = &errorMessage
	agentExecutionRecord.Status = enum.AgentExecutionFail.String()

	execution, err = dependencies.PostgresRepositories.AgentExecutionRepository.Update(ctx, agentExecutionRecord)
	if err != nil {
		tracing.TraceErr(span, err)
		loopErr = multierr.Append(loopErr, err)
	}
	return
}
