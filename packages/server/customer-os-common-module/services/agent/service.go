package agent

import (
	"context"
	"encoding/json"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
	"github.com/google/uuid"
	"github.com/pkg/errors"

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
	events               *events.EventsService
}

func NewAgentService(postgresRepositories *postgresrepository.Repositories, events *events.EventsService) interfaces.AgentService {
	return &agentService{
		postgresRepositories: postgresRepositories,
		events:               events,
	}
}

func (a *agentService) GetAgentById(ctx context.Context, agentID string) (*postgresentity.Agents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.GetAgentById")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.GetAllAgentsByTenant")
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
	tracing.TagEntity(span, newAgent.ID)

	err = a.events.Publisher.PublishEvent(ctx, newAgent.ID, model.AGENT, dto.CreateAgent{
		Active:       newAgent.IsActive,
		Name:         newAgent.Name,
		Type:         newAgent.Type.String(),
		Icon:         newAgent.Icon,
		Color:        newAgent.Color,
		Capabilities: newAgent.GetCapabilitiesConfigAsString(),
		Goal:         newAgent.Goal,
		Status:       newAgent.Status,
		FlowID:       newAgent.FlowID,
		VisibleInUI:  newAgent.VisibleInUI,
		RegistryID:   newAgent.RegistryID,
	})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateAgent"))
	}

	return newAgent, nil
}

func (a *agentService) UpdateAgent(ctx context.Context, agentId string, agentFields data_fields.AgentFields, capabilitiesConfig *postgresentity.CapabilitiesConfig) (*postgresentity.Agents, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.UpdateAgent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, agentId)
	tracing.LogObjectAsJson(span, "agentFields", agentFields)

	agentEntity, err := a.postgresRepositories.AgentsRepository.GetById(ctx, agentId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if agentEntity == nil {
		err := errors.New("agent not found")
		tracing.TraceErr(span, err)
		return nil, err
	}
	if agentFields.Name != nil {
		agentEntity.Name = *agentFields.Name
	}
	if agentFields.Goal != nil {
		agentEntity.Goal = *agentFields.Goal
	}
	if agentFields.Icon != nil {
		agentEntity.Icon = *agentFields.Icon
	}
	if agentFields.Color != nil {
		agentEntity.Color = *agentFields.Color
	}
	if agentFields.Active != nil {
		agentEntity.IsActive = *agentFields.Active
	}
	if agentFields.VisibleInUI != nil {
		agentEntity.VisibleInUI = *agentFields.VisibleInUI
	}
	if agentFields.FlowID != nil {
		agentEntity.FlowID = *agentFields.FlowID
	}
	if capabilitiesConfig != nil {
		agentEntity.CapabilitiesConfig = *capabilitiesConfig
	}

	// set default and non-updatable fields
	updatedAgent, err := a.postgresRepositories.AgentsRepository.Update(ctx, *agentEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	eventFields := dto.UpdateAgent{agentFields, capabilitiesConfig}
	err = a.events.Publisher.PublishEvent(ctx, agentId, model.AGENT, eventFields)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateAgent"))
	}

	return updatedAgent, nil
}

func (a *agentService) CreateAgentExecutionRecord(ctx context.Context, agent postgresentity.Agents, triggerEvent string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.CreateAgentExecutionRecord")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.SaveAgentExecutionSuccess")
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.SaveAgentExecutionError")
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
