package agent_listeners

import (
	"context"
	"fmt"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type NewLeadListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   *agent.AgentRunnerService
}

func NewNewLeadListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService *agent.AgentRunnerService,
) *NewLeadListener {
	return &NewLeadListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.NewLead](), // subscribed event
			events.QueueAgents,                 // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

// Add all Agent types subscribed to this event here
func (l *NewLeadListener) subscribedAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentICPQualifier,
	}
}

func (l *NewLeadListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewLeadListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return l.handleExecution(ctx, event.Event.EntityId, event.Event.AgentEventName)
}

func (l *NewLeadListener) handleExecution(ctx context.Context, orgID string, eventName string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewLeadListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	subscribedAgents := l.subscribedAgents()
	if len(subscribedAgents) == 0 {
		err := errors.New("No agent types configured for company stage lead listener")
		tracing.TraceErr(span, err)
		return err
	}

	message := struct {
		OrganizationID string
	}{
		OrganizationID: orgID,
	}

	activeAgents := l.lookupActiveAgents(ctx, subscribedAgents)
	var errs error
	for _, agent := range activeAgents {

		initialParams, err := utils.StructToMap(message)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
		err = l.agentRunnerService.Run(ctx, agent, eventName, initialParams)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (l *NewLeadListener) lookupActiveAgents(ctx context.Context, agentTypes []enum.AgentType) []postgres_entity.Agent {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewLeadListener.lookupActiveAgents")
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
