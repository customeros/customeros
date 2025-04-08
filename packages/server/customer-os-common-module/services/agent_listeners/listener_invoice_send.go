package agent_listeners

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type SendInvoiceListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   interfaces.AgentRunnerService
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*SendInvoiceListener)(nil)
	_ interfaces.EventListener        = (*SendInvoiceListener)(nil)
)

func NewSendInvoiceListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *SendInvoiceListener {
	return &SendInvoiceListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.SendInvoice](), // subscribed event
			events.QueueAgents,                     // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

func (l *SendInvoiceListener) Type() enum.AgentListenerEvent {
	return enum.EventSendInvoice
}

func (l *SendInvoiceListener) Name() string {
	return "Invoices that are sent"
}

func (l *SendInvoiceListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *SendInvoiceListener) ExecutingAgents() []enum.AgentType {
	// add execution agents
	return []enum.AgentType{
		enum.AgentCashflowGuardian,
	}
}

func (l *SendInvoiceListener) Handle(ctx context.Context, baseEvent any) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "SendInvoiceListener.Handle")
	defer spans.Finish()

	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	data, err := events.DecodeEventData[dto.SendInvoice](ctx, event)
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

func (l *SendInvoiceListener) handleExecution(ctx context.Context, data dto.SendInvoice) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "SendInvoiceListener.handleExecution")
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
