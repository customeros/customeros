package agent

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	"github.com/google/uuid"

	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgresrepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type agentService struct {
	postgresRepositories *postgresrepository.Repositories
}

func NewAgentService(postgresRepositories *postgresrepository.Repositories) interfaces.AgentService {
	return &agentService{
		postgresRepositories: postgresRepositories,
	}
}

func (a *agentService) GetAgentById(ctx context.Context, agentID string) (*postgresentity.Agents, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.GetAgentById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agent, err := a.postgresRepositories.AgentsRepository.GetById(ctx, agentID)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if agent == nil {
		err := errors.New("agent not found")
		tracing.TraceErr(span, err)
		return nil, err
	}

	return agent, nil
}

func (a *agentService) GetAllAgentsByTenant(ctx context.Context) ([]*postgresentity.Agents, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.GetAllAgentsByTenant")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agents, err := a.postgresRepositories.AgentsRepository.GetAll(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return agents, nil
}

func (a *agentService) CreateAgent(ctx context.Context, agentType enum.AgentType) (*postgresentity.Agents, error) {
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
	agentRegistry, err := a.postgresRepositories.AgentRegistryRepository.Find(ctx, agentType)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if agentRegistry == nil {
		err := errors.New("no agent found in agent repository")
		tracing.TraceErr(span, err)
		return nil, err
	}

	// build new agent
	agent := postgresentity.Agents{
		Type:        agentType,
		Tenant:      tenant,
		Name:        agentRegistry.Name,
		Goal:        agentRegistry.Goal,
		IsActive:    false,
		VisibleInUI: true,
		Icon:        agentRegistry.Icon,
		Color:       utils.GetRandomColor(),
		RegistryID:  agentRegistry.ID,
	}

	// build capabilities from registry
	var agentCapabilities []postgresentity.Capability
	for _, masterCapability := range agentRegistry.CapabilitiesConfig.Capabilities {
		agentCapability := postgresentity.Capability{
			ID:     uuid.New().String(),
			Name:   masterCapability.Name,
			Type:   masterCapability.Type,
			Error:  "",
			Active: masterCapability.Active,
		}
		if agentCapability.Name == "" {
			agentCapability.Name = masterCapability.Type.GetName()
		}
		if masterCapability.Config != "" {
			agentCapability.Config = masterCapability.Config
		} else {
			config := agent_capability.GetCapabilityConfigStruct(agentCapability.Type)
			if config != nil {
				configBytes, err := json.Marshal(config)
				if err != nil {
					tracing.TraceErr(span, err)
					return nil, err
				}
				agentCapability.Config = string(configBytes)
			}
		}
		agentCapabilities = append(agentCapabilities, agentCapability)
	}
	agent.CapabilitiesConfig = postgresentity.CapabilitiesConfig{
		Capabilities: agentCapabilities,
	}

	// create agent instance in database
	newAgent, err := a.postgresRepositories.AgentsRepository.Create(ctx, agent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return newAgent, nil
}

func (a *agentService) UpdateAgent(ctx context.Context, agentFields postgresentity.Agents) (*postgresentity.Agents, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.UpdateAgent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	existingAgentEntity, err := a.postgresRepositories.AgentsRepository.GetById(ctx, agentFields.ID)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// set default and non-updatable fields
	agentFields.UpdatedAt = utils.NowPtr()
	agentFields.CreatedAt = existingAgentEntity.CreatedAt
	agentFields.Type = existingAgentEntity.Type
	agentFields.Tenant = existingAgentEntity.Tenant
	agentFields.RegistryID = existingAgentEntity.RegistryID
	agentFields.Status = existingAgentEntity.Status
	agentFields.ErrorMessage = existingAgentEntity.ErrorMessage

	updatedAgent, err := a.postgresRepositories.AgentsRepository.Update(ctx, agentFields)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return updatedAgent, nil
}

func (a *agentService) CreateAgentExecutionRecord(ctx context.Context, agent postgresentity.Agents, triggerEvent string) (string, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "AgentService.CreateAgentExecutionRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agentExecutionRecord := postgresentity.AgentExecution{
		Tenant:       agent.Tenant,
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

	executionRecord, err := a.postgresRepositories.AgentExecutionRepository.Find(ctx, postgresentity.AgentExecution{
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

	executionRecord, err := a.postgresRepositories.AgentExecutionRepository.Find(ctx, postgresentity.AgentExecution{
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
