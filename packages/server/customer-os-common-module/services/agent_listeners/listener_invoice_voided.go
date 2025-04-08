package agent_listeners

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
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
	spans, ctx := telemetry.StartServiceSpan(ctx, "InvoiceVoidedListener.Handle")
	defer spans.Finish()

	spans.LogObjectAsJson("baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	data, err := events.DecodeEventData[dto.InvoiceVoided](ctx, event)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	if data.InvoiceID == "" {
		err := errors.New("missing invoice id")
		spans.TraceError(err)
		return err
	}

	activeAgents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, l.ExecutingAgents())
	if err != nil {
		spans.TraceError(err)
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
			spans.TraceError(err)
			errs = multierr.Append(errs, err)
			continue
		}
		agentExecutionId, err := l.agentRunnerService.Run(ctx, agent, l.Type().String(), initialParams, nil)
		if err != nil {
			spans.TraceError(err)
			errs = multierr.Append(errs, err)
		}
		agentExecutionIds = append(agentExecutionIds, agentExecutionId)
	}

	if !data.DryRun {
		for _, agentExecutionId := range agentExecutionIds {
			err := l.postgresRepositories.AgentExecutionRepository.GoalAchieved(ctx, agentExecutionId, false, utils.StringPtr(data.InvoiceID))
			if err != nil {
				spans.TraceError(errors.Wrap(err, "failed to mark agent execution as goal achieved"))
			}
		}
	}

	return errs
}
