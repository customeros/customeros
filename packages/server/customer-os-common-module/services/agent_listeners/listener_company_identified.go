package agent_listeners

import (
	"context"
	"errors"
	"fmt"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type CompanyIdentifiedListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   interfaces.AgentRunnerService
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*CompanyIdentifiedListener)(nil)
	_ interfaces.EventListener        = (*CompanyIdentifiedListener)(nil)
)

func NewCompanyIdentifiedListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *CompanyIdentifiedListener {
	return &CompanyIdentifiedListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.CompanyIdentified](), // subscribed event
			events.QueueAgents,                           // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

func (l *CompanyIdentifiedListener) Type() enum.AgentListenerEvent {
	return enum.EventCompanyIdentified
}

func (l *CompanyIdentifiedListener) Name() string {
	return "Companies identified from web visit"
}

func (l *CompanyIdentifiedListener) DefaultConfig() any {
	return &postgres_entity.NoConfig{}
}

func (l *CompanyIdentifiedListener) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentSupportSpotter,
	}
}

func (l *CompanyIdentifiedListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CompanyIdentifiedListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	data, err := events.DecodeEventData[dto.CompanyIdentified](ctx, event)
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

	err = l.handleExecution(ctx, data.AgentExecutionId)
	if err != nil {
		tracing.TraceErr(span, err)
		errs = multierr.Append(errs, err)
	}

	return errs
}

func (l *CompanyIdentifiedListener) handleExecution(ctx context.Context, orgID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CompanyIdentifiedListener.handleExecution")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	activeAgents := l.lookupActiveAgents(ctx)
	if len(activeAgents) == 0 {
		err := errors.New("No agent types configured for company identified listener")
		tracing.TraceErr(span, err)
		return err
	}

	// get web session
	lookbackInMins := 60
	session, err := l.postgresRepositories.WebSessionRepository.FindSession(ctx, postgres_entity.WebSession{
		Tenant:         common.GetTenantFromContext(ctx),
		OrganizationId: &orgID,
		IsActive:       false,
	}, &lookbackInMins)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if session == nil {
		return nil
	}

	message := struct {
		OrganizationID  string   `json:"organizationId"`
		UniquePageViews []string `json:"uniquePageViews"`
	}{
		OrganizationID:  orgID,
		UniquePageViews: session.UniquePageViews,
	}

	var errs error
	for _, agent := range activeAgents {

		initialParams, err := utils.StructToMap(message)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
		_, err = l.agentRunnerService.Run(ctx, agent, l.Type().String(), initialParams)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (l *CompanyIdentifiedListener) handleGoalAchieved(ctx context.Context, agentExecutionId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CompanyIdentifiedListener.handleGoalAchieved")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	span.LogFields(log.String("agentExecutionId", agentExecutionId))

	var agentExecution *postgres_entity.AgentExecution
	var err error

	agentExecution, err = l.postgresRepositories.AgentExecutionRepository.GetById(ctx, agentExecutionId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err

	}
	if agentExecution == nil {
		err = fmt.Errorf("agent execution not found")
		tracing.TraceErr(span, err)
		return err
	}

	// update execution with goal achieved
	_, err = l.postgresRepositories.AgentExecutionRepository.Completed(ctx, agentExecution.ID, utils.TruePtr())
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (l *CompanyIdentifiedListener) lookupActiveAgents(ctx context.Context) []postgres_entity.Agent {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CompanyIdentifiedListener.lookupActiveAgents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	agents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, l.ExecutingAgents())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	return agents
}
