package agent_listeners

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type StartInvoiceRunWithAutopayment struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   interfaces.AgentRunnerService
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*StartInvoiceRunWithAutopayment)(nil)
	_ interfaces.EventListener        = (*StartInvoiceRunWithAutopayment)(nil)
)

func NewStartInvoiceRunWithAutopayment(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *StartInvoiceRunWithAutopayment {
	return &StartInvoiceRunWithAutopayment{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.InvoiceContractWithAutopayment](), // subscribed event
			events.QueueAgents, // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

func (l *StartInvoiceRunWithAutopayment) Type() enum.AgentListenerEvent {
	return enum.EventStartInvoiceRunWithAutopayment
}

func (l *StartInvoiceRunWithAutopayment) Name() string {
	return "Start invoice run"
}

func (l *StartInvoiceRunWithAutopayment) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *StartInvoiceRunWithAutopayment) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentCashflowGuardian,
	}
}

func (l *StartInvoiceRunWithAutopayment) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StartInvoiceRunWithAutopayment.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	data, err := events.DecodeEventData[dto.InvoiceContractWithAutopayment](ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if data.ContractId == "" {
		err := errors.New("Missing contract id")
		tracing.TraceErr(span, err)
		return err
	}

	return l.handleExecution(ctx, data)
}

func (l *StartInvoiceRunWithAutopayment) handleExecution(ctx context.Context, data dto.InvoiceContractWithAutopayment) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StartInvoiceRunWithAutopayment.handleExecution")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	activeAgents := l.lookupActiveAgents(ctx)
	if activeAgents == nil || len(activeAgents) == 0 {
		return nil
	}

	var errs error
	for _, agent := range activeAgents {
		// replaced route with generic mapping
		initialParams, err := utils.StructToMap(data)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
		err = l.agentRunnerService.Run(ctx, agent, l.Type().String(), initialParams)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (l *StartInvoiceRunWithAutopayment) lookupActiveAgents(ctx context.Context) []postgres_entity.Agent {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StartInvoiceRunWithAutopayment.lookupActiveAgents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	agents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, l.ExecutingAgents())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	return agents
}
