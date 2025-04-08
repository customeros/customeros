package agent_listeners

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type PastDueInvoiceListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   interfaces.AgentRunnerService
}

type PastDueInvoiceConfig struct {
	OverdueDays ConfigSingleIntValue `json:"overdueDays"`
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*PastDueInvoiceListener)(nil)
	_ interfaces.EventListener        = (*PastDueInvoiceListener)(nil)
)

func NewPastDueInvoiceListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *PastDueInvoiceListener {
	return &PastDueInvoiceListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.PastDueInvoice](), // subscribed event
			events.QueueAgents,                        // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

func (l *PastDueInvoiceListener) Type() enum.AgentListenerEvent {
	return enum.EventInvoicePastDue
}

func (l *PastDueInvoiceListener) Name() string {
	return "Invoices that are past their due date"
}

func (l *PastDueInvoiceListener) DefaultConfig() any {
	config := PastDueInvoiceConfig{
		OverdueDays: ConfigSingleIntValue{
			Value: 15,
		},
	}
	return &config
}

func (l *PastDueInvoiceListener) ExecutingAgents() []enum.AgentType {
	// add execution agents
	return []enum.AgentType{
		enum.AgentCashflowGuardian,
	}
}

func (l *PastDueInvoiceListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "PastDueInvoiceListener.Handle")
	defer spans.Finish()

	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	data, err := events.DecodeEventData[dto.PastDueInvoice](ctx, event)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	if data.InvoiceId == "" {
		err := errors.New("missing invoice id")
		spans.TraceError(err)
		return err
	}

	return l.handleExecution(ctx, data)
}

func (l *PastDueInvoiceListener) handleExecution(ctx context.Context, data dto.PastDueInvoice) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "PastDueInvoiceListener.handleExecution")
	defer spans.Finish()

	activeAgents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, l.ExecutingAgents())
	if err != nil {
		spans.TraceError(err)
		return nil
	}
	if len(activeAgents) == 0 {
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
