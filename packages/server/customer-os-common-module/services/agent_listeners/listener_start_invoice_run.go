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

type StartInvoiceRun struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   interfaces.AgentRunnerService
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*StartInvoiceRun)(nil)
	_ interfaces.EventListener        = (*StartInvoiceRun)(nil)
)

func NewStartInvoiceRun(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *StartInvoiceRun {
	return &StartInvoiceRun{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.InvoiceContract](), // subscribed event
			events.QueueAgents,                         // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

func (l *StartInvoiceRun) Type() enum.AgentListenerEvent {
	return enum.EventStartInvoiceRun
}

func (l *StartInvoiceRun) Name() string {
	return "Start invoice run"
}

func (l *StartInvoiceRun) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *StartInvoiceRun) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentCashflowGuardian,
	}
}

func (l *StartInvoiceRun) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StartInvoiceRun.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	data, err := events.DecodeEventData[dto.InvoiceContract](ctx, event)
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

func (l *StartInvoiceRun) handleExecution(ctx context.Context, data dto.InvoiceContract) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StartInvoiceRun.handleExecution")
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

func (l *StartInvoiceRun) lookupActiveAgents(ctx context.Context) []postgres_entity.Agent {
	span, ctx := opentracing.StartSpanFromContext(ctx, "StartInvoiceRun.lookupActiveAgents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	agents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, l.ExecutingAgents())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	return agents
}
