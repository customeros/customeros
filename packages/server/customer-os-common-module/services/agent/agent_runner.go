package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type AgentRunnerService struct {
	postgresRepositories     *postgres_repository.Repositories
	agentCapabilitiesService *agent_capability.AgentCapabilities
	agentService             interfaces.AgentService
}

func NewAgentRunnerService(
	postgresRepositories *postgres_repository.Repositories,
	agentCapabilities *agent_capability.AgentCapabilities,
	agentService interfaces.AgentService,
) *AgentRunnerService {
	if postgresRepositories == nil {
		panic("postgresRepositories cannot be nil")
	}
	if agentCapabilities == nil {
		panic("agentCapabilities cannot be nil")
	}
	if agentService == nil {
		panic("agentService cannot be nil")
	}

	return &AgentRunnerService{
		postgresRepositories:     postgresRepositories,
		agentCapabilitiesService: agentCapabilities,
		agentService:             agentService,
	}
}

func (a *AgentRunnerService) Run(ctx context.Context, agent postgres_entity.Agents, eventType string, initialParams map[string]any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.Run")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	if agent.ID == "" {
		err := errors.New("agent ID cannot be empty")
		tracing.TraceErr(span, err)
		return err
	}
	tracing.TagEntity(span, agent.ID)
	tracing.LogObjectAsJson(span, "agent", agent)

	if !agent.IsActive {
		err := errors.New("agent is not active")
		tracing.TraceErr(span, err)
		return err
	}

	// create execution record
	executionID, err := a.createAgentExecutionRecord(ctx, agent.ID, eventType)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to create agent execution record"))
		return err
	}

	// Create a copy of initialParams to avoid mutating the original map
	allParams := make(map[string]any)
	utils.MergeMapToMap(initialParams, allParams)

	var capErr error
	for _, capability := range agent.CapabilitiesConfig.Capabilities {
		if !capability.Active {
			continue
		}

		capType := capability.Type
		var capExecutor interfaces.AgentCapabilityUntyped
		capExecutor, capErr = a.agentCapabilitiesService.GetExecutor(capType)
		if capErr != nil {
			break
		}
		if capExecutor == nil {
			capErr = fmt.Errorf("missing executor for capability type: %s", capType)
			break
		}

		input := capExecutor.GetInput()
		if input == nil {
			capErr = fmt.Errorf("GetInput returned nil for capability type: %s", capType)
			break
		}
		capErr = utils.MapToStruct(allParams, input)
		if capErr != nil {
			break
		}

		config := capExecutor.GetConfig()
		if config == nil {
			capErr = fmt.Errorf("GetConfig returned nil for capability type: %s", capType)
			break
		}
		if capability.Config != "" {
			data := []byte(capability.Config)
			capErr = json.Unmarshal(data, config)
			if capErr != nil {
				break
			}
		}

		var output any
		output, capErr = capExecutor.ExecuteUntyped(ctx, input, config)
		if capErr != nil {
			break
		}

		// Merge the output back into allParams for subsequent capabilities
		if output != nil {
			var outputMap map[string]any
			outputMap, capErr = utils.StructToMap(output)
			if capErr != nil {
				break
			}
			utils.MergeMapToMap(outputMap, allParams)
		}
	}

	if capErr != nil {
		tracing.TraceErr(span, capErr)
		_, dbErr := a.postgresRepositories.AgentExecutionRepository.Update(ctx, executionID, nil, utils.StringPtr(capErr.Error()), false)
		if dbErr != nil {
			tracing.TraceErr(span, errors.Wrap(dbErr, "unable to update agent execution record"))
			return dbErr
		}
	}

	// update agentExecutionRecord
	_, err = a.postgresRepositories.AgentExecutionRepository.Update(ctx, executionID, utils.NowPtr(), nil, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (a *AgentRunnerService) createAgentExecutionRecord(ctx context.Context, agentID, triggerEventType string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.createAgentExecutionRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("agentID", agentID, "triggerEventType", triggerEventType)

	agent, err := a.postgresRepositories.AgentsRepository.Find(ctx, postgres_entity.Agents{
		ID: agentID,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	if agent == nil {
		err := errors.New("agent not found")
		tracing.TraceErr(span, err)
		return "", err
	}

	return a.agentService.CreateAgentExecutionRecord(ctx, *agent, triggerEventType)
}
