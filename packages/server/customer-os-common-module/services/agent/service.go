package agent

import (
	"context"
	"errors"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type agentService struct {
	postgresRepositories *postgres_repository.Repositories
}

func NewAgentService(postgresRepositories *postgres_repository.Repositories) interfaces.AgentService {
	return &agentService{
		postgresRepositories: postgresRepositories,
	}
}

func (a *agentService) CreateAgent(ctx context.Context, agentType enum.AgentType) (*postgres_entity.Agents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.CreateAgent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not found in context")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// get config from registry
	masterAgent, err := a.postgresRepositories.AgentRegistryRepository.Find(ctx, agentType)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if masterAgent == nil {
		err := errors.New("no agent found in agent repository")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// build new agent
	agent := postgres_entity.Agents{
		Type:         masterAgent.ID,
		Tenant:       tenant,
		Name:         masterAgent.Name,
		Capabilities: masterAgent.Capabilities,
		Goal:         masterAgent.Goal,
		IsActive:     false,
		VisibleInUI:  true,
		Icon:         masterAgent.Icon,
		Color:        utils.GetRandomColor(),
	}

	// create agent instance in database
	newAgent, err := a.postgresRepositories.AgentsRepository.Create(ctx, agent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return newAgent, nil
}

func (a *agentService) CreateAgentExecutionRecord(ctx context.Context, agent postgres_entity.Agents, triggerEvent string) (string, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.CreateAgentExecutionRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agentExecutionRecord := postgres_entity.AgentExecution{
		AgentID:      &agent.ID,
		TriggerEvent: triggerEvent,
		Status:       enum.AgentExecutionRunning.String(),
		StartedAt:    utils.NowPtr(),
	}

	createdRecord, err := a.postgresRepositories.AgentExecutionRepository.Create(ctx, agentExecutionRecord)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if createdRecord == nil {
		err := errors.New("unable to create agent execution record")
		tracing.TraceErr(span, err)
		return "", err
	}

	return createdRecord.ID, nil
}

func (a *agentService) SaveAgentExecutionCompleted(ctx context.Context, executionID string, goalAchieved bool) error {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.SaveAgentExecutionSuccess")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	executionRecord, err := a.postgresRepositories.AgentExecutionRepository.Find(ctx, postgres_entity.AgentExecution{
		ID: executionID,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	executionRecord.Status = enum.AgentExecutionCompleted.String()
	executionRecord.CompletedAt = utils.NowPtr()
	executionRecord.GoalAchieved = goalAchieved
	// update

	return nil
}

func (a *agentService) SaveAgentExecutionError(ctx context.Context, executionID, errorMessage string) error {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.SaveAgentExecutionError")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	executionRecord, err := a.postgresRepositories.AgentExecutionRepository.Find(ctx, postgres_entity.AgentExecution{
		ID: executionID,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	executionRecord.Status = enum.AgentExecutionFail.String()
	executionRecord.ErrorMessage = &errorMessage
	// update

	return nil
}
