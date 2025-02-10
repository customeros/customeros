package agent_listeners

import (
	"context"
	"fmt"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type CompanyNeedsHelpListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*CompanyNeedsHelpListener)(nil)
	_ interfaces.EventListener        = (*CompanyNeedsHelpListener)(nil)
)

func NewCompanyNeedsHelpListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
) *CompanyNeedsHelpListener {
	return &CompanyNeedsHelpListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.CompanyNeedsHelp](), // subscribed event
			events.QueueAgents,                          // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
	}
}

func (l *CompanyNeedsHelpListener) Type() enum.AgentListenerEvent {
	return enum.EventNewSupportVisit
}

func (l *CompanyNeedsHelpListener) Name() string {
	return "New Support Visit"
}

func (h *CompanyNeedsHelpListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (h *CompanyNeedsHelpListener) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{}
}

func (l *CompanyNeedsHelpListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CompanyNeedsHelpListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	data, err := events.DecodeEventData[dto.CompanyNeedsHelp](ctx, event)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	var errs error
	if data.AgentExecutionId != "" {
		err = l.handleGoalAchieved(ctx, data.AgentExecutionId)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (l *CompanyNeedsHelpListener) handleGoalAchieved(ctx context.Context, agentExecutionId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CompanyNeedsHelpListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	span.LogKV("agentExecutionId", agentExecutionId)

	var agentExecution *postgres_entity.AgentExecution
	var err error
	if agentExecutionId != "" {
		agentExecution, err = l.postgresRepositories.AgentExecutionRepository.GetById(ctx, agentExecutionId)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}
	if agentExecution == nil {
		err = fmt.Errorf("agent execution not found")
		tracing.TraceErr(span, err)
		return err
	}

	// update execution with goal achieved
	if agentExecution.Status == enum.AgentExecutionCompleted {
		err = fmt.Errorf("agent execution already completed")
		tracing.TraceErr(span, err)
		return err
	}
	_, err = l.postgresRepositories.AgentExecutionRepository.Completed(ctx, agentExecution.ID, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
