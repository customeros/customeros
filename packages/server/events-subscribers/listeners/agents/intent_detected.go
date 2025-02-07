package agent_listeners

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

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
		enum.AgentWebVisitorIdentifier,
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

	data, err := events.DecodeEventData[dto.IntentDetected](ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return l.handle(ctx, data)
}

func (l *IntentDetectedListener) handle(ctx context.Context, message dto.IntentDetected) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IntentDetectedListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	var agentTypes []enum.AgentType
	switch message.IntentType {
	case enum.IntentSupportRequired:
		agentTypes = append(agentTypes, enum.AgentSupportSpotter)
	case enum.IntentIcpCheck:
		agentTypes = append(agentTypes, enum.AgentICPQualifier)
	case enum.IntentGenerateCycleInvoice:
		agentTypes = append(agentTypes, enum.AgentCashflowGuardian)
	default:
		err := errors.New("IntentType not supported")
		tracing.TraceErr(span, err)
		return err
	}

	if len(agentTypes) == 0 {
		err := errors.New("No agent types configured for intent")
		tracing.TraceErr(span, err)
		return err
	}

	activeAgents := l.lookupActiveAgents(ctx, agentTypes)
	var errs error
	for _, agent := range activeAgents {
		initialParams, err := utils.StructToMap(message)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
		err = l.dependencies.CommonServices.AgentRunnerService.Run(ctx, agent, message.IntentType.String(), initialParams)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (l *IntentDetectedListener) lookupActiveAgents(ctx context.Context, agentTypes []enum.AgentType) []postgres_entity.Agent {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Listeners.lookupActiveAgents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	span.LogFields(log.String("agentTypes", fmt.Sprintf("%v", agentTypes)))

	agents, err := l.dependencies.PostgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, agentTypes)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	return agents
}
