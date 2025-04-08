package agent_listeners

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

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
	return "Scheduled invoices with auto payment"
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
	spans, ctx := telemetry.StartServiceSpan(ctx, "StartInvoiceRunWithAutopayment.Handle")
	defer spans.Finish()

	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	data, err := events.DecodeEventData[dto.InvoiceContractWithAutopayment](ctx, event)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	if data.ContractId == "" {
		err := errors.New("Missing contract id")
		spans.TraceError(err)
		return err
	}

	return l.handleExecution(ctx, data)
}

func (l *StartInvoiceRunWithAutopayment) handleExecution(ctx context.Context, data dto.InvoiceContractWithAutopayment) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "StartInvoiceRunWithAutopayment.handleExecution")
	defer spans.Finish()

	activeAgents := l.lookupActiveAgents(ctx)
	if activeAgents == nil || len(activeAgents) == 0 {
		return nil
	}

	var errs error
	for _, agent := range activeAgents {
		// replaced route with generic mapping
		initialParams, err := utils.StructToMap(data)
		if err != nil {
			spans.TraceError(err)
			errs = multierr.Append(errs, err)
		}
		_, err = l.agentRunnerService.Run(ctx, agent, l.Type().String(), initialParams, nil)
		if err != nil {
			spans.TraceError(err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (l *StartInvoiceRunWithAutopayment) lookupActiveAgents(ctx context.Context) []postgres_entity.Agent {
	spans, ctx := telemetry.StartServiceSpan(ctx, "StartInvoiceRunWithAutopayment.lookupActiveAgents")
	defer spans.Finish()

	agents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, l.ExecutingAgents())
	if err != nil {
		spans.TraceError(err)
		return nil
	}
	return agents
}
