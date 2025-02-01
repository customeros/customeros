package agent

import (
	"context"
	"encoding/json"

	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgresrepository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/events"
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

func (a *agentService) GetAgentById(ctx context.Context, agentID string) (*postgresentity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.GetAgentById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agent, err := a.postgresRepositories.AgentRepository.GetById(ctx, agentID)
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

func (a *agentService) GetAllAgentsByTenant(ctx context.Context) ([]*postgresentity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.GetAllAgentsByTenant")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agents, err := a.postgresRepositories.AgentRepository.GetAll(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return agents, nil
}

func (a *agentService) CreateAgent(ctx context.Context, agentType enum.AgentType) (*postgresentity.Agent, error) {
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
	agent := postgresentity.Agent{
		Type:        agentType,
		Tenant:      tenant,
		Name:        agentRegistry.Name,
		Goal:        agentRegistry.Goal,
		IsActive:    false,
		VisibleInUI: true,
		Icon:        agentRegistry.Icon,
		Color:       agentRegistry.Color,
		RegistryID:  agentRegistry.ID,
	}
	if agent.Color == "" {
		agent.Color = utils.GetRandomColor()
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

	// validate capabilities
	a.ValidateCapabilities(ctx, &agent)

	// create agent instance in database
	newAgent, err := a.postgresRepositories.AgentRepository.Create(ctx, agent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	tracing.TagEntity(span, newAgent.ID)

	err = a.events.Publisher.PublishFanoutEvent(ctx, newAgent.ID, model.AGENT, dto.CreateAgent{
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

func (a *agentService) UpdateAgent(ctx context.Context, agentId string, agentFields data_fields.AgentFields, capabilitiesConfig *postgresentity.CapabilitiesConfig) (*postgresentity.Agent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.UpdateAgent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, agentId)
	tracing.LogObjectAsJson(span, "agentFields", agentFields)
	tracing.LogObjectAsJson(span, "capabilitiesConfig", capabilitiesConfig)

	agentEntity, err := a.postgresRepositories.AgentRepository.GetById(ctx, agentId)
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
	validateCapabilities := false
	if capabilitiesConfig != nil && len(capabilitiesConfig.Capabilities) > 0 {
		agentEntity.CapabilitiesConfig = *capabilitiesConfig
		validateCapabilities = true
	}

	if validateCapabilities {
		a.ValidateCapabilities(ctx, agentEntity)
	}

	updatedAgent, err := a.postgresRepositories.AgentRepository.Update(ctx, *agentEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	eventFields := dto.UpdateAgent{AgentFields: agentFields, Capabilities: capabilitiesConfig}
	err = a.events.Publisher.PublishFanoutEvent(ctx, agentId, model.AGENT, eventFields)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateAgent"))
	}

	return updatedAgent, nil
}

func (a *agentService) ValidateCapabilities(ctx context.Context, agentEntity *postgresentity.Agent) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.ValidateCapabilities")
	defer span.Finish()

	capabilities := agentEntity.CapabilitiesConfig.Capabilities // []Capability
	allCapabilitiesValid := true

	for i := range capabilities {
		capability := &capabilities[i] // pointer so we can update error
		if !capability.Active {
			// Not active => auto valid
			continue
		}

		// decode JSON config into typed struct (example)
		typedConfig, decodeErr := decodeConfigForCapability(capability.Type, capability.Config)
		if decodeErr != nil {
			tracing.TraceErr(span, decodeErr)
			continue
		}
		if typedConfig == nil {
			continue
		}

		// see if typedConfig implements ConfigValidator
		if validator, ok := typedConfig.(agent_capability.ConfigValidator); ok {
			if !validator.Validate() {
				allCapabilitiesValid = false
			}
			newConfigJson, _ := json.Marshal(typedConfig)
			capability.Config = string(newConfigJson)
		}
	}

	agentEntity.CapabilitiesConfig.Capabilities = capabilities
	agentEntity.Configured = allCapabilitiesValid
}

func decodeConfigForCapability(capType enum.AgentCapabilityType, configJSON string) (any, error) {
	c := agent_capability.GetCapabilityConfigStruct(capType)
	if c == nil {
		return nil, nil
	}

	err := json.Unmarshal([]byte(configJSON), c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (a *agentService) CreateAgentExecutionRecord(ctx context.Context, agent postgresentity.Agent, triggerEvent, traceId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentService.CreateAgentExecutionRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	agentExecutionRecord := postgresentity.AgentExecution{
		Tenant:       agent.Tenant,
		AgentID:      &agent.ID,
		TriggerEvent: triggerEvent,
		Status:       enum.AgentExecutionRunning.String(),
		StartedAt:    utils.NowPtr(),
		TraceId:      traceId,
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
