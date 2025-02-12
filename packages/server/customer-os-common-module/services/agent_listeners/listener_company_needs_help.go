package agent_listeners

import (
	"context"
	"errors"

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
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type CompanyNeedsHelpListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   interfaces.AgentRunnerService
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*CompanyNeedsHelpListener)(nil)
	_ interfaces.EventListener        = (*CompanyNeedsHelpListener)(nil)
)

func NewCompanyNeedsHelpListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *CompanyNeedsHelpListener {
	return &CompanyNeedsHelpListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.CompanyNeedsHelp](), // subscribed event
			events.QueueAgents,                          // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

func (l *CompanyNeedsHelpListener) Type() enum.AgentListenerEvent {
	return enum.EventCompanyNeedsHelp
}

func (l *CompanyNeedsHelpListener) Name() string {
	return "Companies that need help"
}

func (h *CompanyNeedsHelpListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (h *CompanyNeedsHelpListener) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentSupportSpotter,
	}
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

	return l.handleExecution(ctx, event.Event.EntityId)
}

func (l *CompanyNeedsHelpListener) handleExecution(ctx context.Context, orgID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CompanyNeedsHelpListener.handleExecution")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	activeAgents := l.lookupActiveAgents(ctx)
	if len(activeAgents) == 0 {
		err := errors.New("No agent types configured for company identified listener")
		tracing.TraceErr(span, err)
		return err
	}

	startParams := struct {
		OrganizationID string
	}{
		OrganizationID: orgID,
	}

	var errs error
	for _, agent := range activeAgents {
		err := l.run(ctx, agent, startParams)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (l *CompanyNeedsHelpListener) run(ctx context.Context, agent postgres_entity.Agent, startParams any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CompanyNeedsHelpListener.run")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	initialParams, err := utils.StructToMap(startParams)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	err = l.agentRunnerService.Run(ctx, agent, l.Type().String(), initialParams)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (l *CompanyNeedsHelpListener) lookupActiveAgents(ctx context.Context) []postgres_entity.Agent {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CompanyNeedsHelpListener.lookupActiveAgents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	agents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, l.ExecutingAgents())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	return agents
}
