package agent

import (
	"context"
	"errors"

	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/multierr"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type SupportAgent struct {
	postgresRepositories *postgres_repository.Repositories
	agentCapabilities    *agent_capability.AgentCapabilities
	tagService           interfaces.TagService
}

func NewSupportAgent(
	postgresRepo *postgres_repository.Repositories,
	agentCapabilities *agent_capability.AgentCapabilities,
	tagService interfaces.TagService,
) *SupportAgent {
	return &SupportAgent{
		postgresRepositories: postgresRepo,
		agentCapabilities:    agentCapabilities,
		tagService:           tagService,
	}
}

func (a *SupportAgent) ProcessIntentEvent(ctx context.Context, intentEvent *dto.IntentEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SupportAgent.ProcessIntentEvent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "event", intentEvent)

	// find active agents for tenant
	agents, err := a.postgresRepositories.AgentsRepository.GetActiveAgentsByTypes(ctx, []enum.AgentType{enum.AgentSupport})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if len(agents) == 0 {
		return nil
	}

	// route to agents
	var errs error
	for _, agent := range agents {
		err := a.Run(ctx, agent.ID, intentEvent)
		if err != nil {
			tracing.TraceErr(span, err)
			errs = multierr.Append(errs, err)
		}
	}

	return errs
}

func (a *SupportAgent) Run(ctx context.Context, agentID string, intentEvent *dto.IntentEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SupportAgent.Run")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// validate intentEvent
	if intentEvent.OrganizationID == "" {
		err := errors.New("OrganizationID is not set")
		tracing.TraceErr(span, err)
		return err
	}
	span.LogKV(
		"tenant", common.GetTenantFromContext(ctx),
		"agentID", agentID,
		"organizationID", intentEvent.OrganizationID,
	)

	// lookup capabilities

	// create execution record

	// execute create tag capability
	tag, err := a.tagService.Save(ctx, nil, &neo4j_entity.TagEntity{
		EntityType: model.ORGANIZATION,
		Name:       "Support",
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if tag == nil {
		err := errors.New("no tag found for Support on Organization")
		tracing.TraceErr(span, err)
		return err
	}

	_, err = a.agentCapabilities.ApplyTag.Execute(ctx, agent_capability.ApplyTagInput{
		EntityType: model.ORGANIZATION,
		EntityID:   intentEvent.OrganizationID,
		TagID:      tag.Id,
	}, agent_capability.NoConfig{})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}
