package agent_listeners

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

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

type InvoiceVoidedListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   interfaces.AgentRunnerService
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*InvoiceVoidedListener)(nil)
	_ interfaces.EventListener        = (*InvoiceVoidedListener)(nil)
)

func NewInvoiceVoidedListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *InvoiceVoidedListener {
	return &InvoiceVoidedListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.InvoiceVoided](), // subscribed event
			events.QueueAgents,                       // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

func (l *InvoiceVoidedListener) Type() enum.AgentListenerEvent {
	return enum.EventInvoiceVoided
}

func (l *InvoiceVoidedListener) Name() string {
	return "Invoices that are voided"
}

func (l *InvoiceVoidedListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *InvoiceVoidedListener) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentCashflowGuardian,
	}
}

func (l *InvoiceVoidedListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InvoiceVoidedListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	data, err := events.DecodeEventData[dto.InvoiceVoided](ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if data.InvoiceID == "" {
		err := errors.New("missing invoice id")
		tracing.TraceErr(span, err)
		return err
	}

	activeAgents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, l.ExecutingAgents())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	if len(activeAgents) == 0 {
		return nil
	}

	var errs error
	var agentExecutionIds []string
	for _, agent := range activeAgents {
		initialParams, err := utils.StructToMap(data)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
			continue
		}
		agentExecutionId, err := l.agentRunnerService.Run(ctx, agent, l.Type().String(), initialParams)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
		agentExecutionIds = append(agentExecutionIds, agentExecutionId)
	}

	if !data.DryRun {
		for _, agentExecutionId := range agentExecutionIds {
			err := l.postgresRepositories.AgentExecutionRepository.GoalAchieved(ctx, agentExecutionId, false)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "failed to mark agent execution as goal achieved"))
			}
		}
	}

	return errs
}
