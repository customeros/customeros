package agent_listeners

import (
	"context"

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

type NewWebSessionListener struct {
	events.BaseEventListener
	postgresRepositories *postgres_repository.Repositories
	agentRunnerService   interfaces.AgentRunnerService
}

type IdentifyWebsiteVisitorConfig struct {
	Websites ConfigMultipleValues `json:"websites"`
}

func (c *IdentifyWebsiteVisitorConfig) Validate() bool {
	isValid := true

	websiteConfigured := false
	for _, website := range c.Websites.Value {
		if !utils.IsBlank(website) {
			websiteConfigured = true
			break
		}
	}
	if !websiteConfigured {
		c.Websites.Error = "Add at least 1 website"
		isValid = false
	} else {
		c.Websites.Error = ""
	}

	return isValid
}

// Compile-time interface check for AgentListenerUntyped
var (
	_ interfaces.AgentListenerUntyped = (*NewWebSessionListener)(nil)
	_ interfaces.EventListener        = (*NewWebSessionListener)(nil)
)

func NewNewWebSessionListener(
	logger logger.Logger,
	postgresRepositories *postgres_repository.Repositories,
	agentRunnerService interfaces.AgentRunnerService,
) *NewWebSessionListener {
	return &NewWebSessionListener{
		BaseEventListener: events.NewBaseEventListener(
			logger,
			events.GetEventType[dto.NewWebSession](), // subscribed event
			events.QueueAgents,                       // listening on Agents queue
		),
		postgresRepositories: postgresRepositories,
		agentRunnerService:   agentRunnerService,
	}
}

func (l *NewWebSessionListener) Type() enum.AgentListenerEvent {
	return enum.EventNewWebSession
}

func (l *NewWebSessionListener) Name() string {
	return "New website sessions"
}

func (l *NewWebSessionListener) DefaultConfig() any {
	config := IdentifyWebsiteVisitorConfig{}
	config.Websites.Value = []string{}
	return &config
}

func (l *NewWebSessionListener) ExecutingAgents() []enum.AgentType {
	return []enum.AgentType{
		enum.AgentWebVisitorIdentifier,
	}
}

func (l *NewWebSessionListener) Handle(ctx context.Context, baseEvent any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewWebSessionListener.Handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "baseEvent", baseEvent)

	event, err := l.ValidateBaseEvent(ctx, baseEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return l.handleExecution(ctx, event.Event.Id)
}

func (l *NewWebSessionListener) handleExecution(ctx context.Context, webSessionId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewWebSessionListener.handle")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	activeAgents := l.lookupActiveAgents(ctx)
	if activeAgents == nil || len(activeAgents) == 0 {
		return nil
	}

	startParams := struct {
		WebSessionID string
	}{
		WebSessionID: webSessionId,
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

func (l *NewWebSessionListener) run(ctx context.Context, agent postgres_entity.Agent, startParams any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewWebSessionListener.run")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	initialParams, err := utils.StructToMap(startParams)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	_, err = l.agentRunnerService.Run(ctx, agent, l.Type().String(), initialParams, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (l *NewWebSessionListener) lookupActiveAgents(ctx context.Context) []postgres_entity.Agent {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NewWebSessionListener.lookupActiveAgents")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	agents, err := l.postgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypes(ctx, l.ExecutingAgents())
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	return agents
}
