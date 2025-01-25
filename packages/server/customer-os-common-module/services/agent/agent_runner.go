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
	"go.uber.org/multierr"
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

	var errs error
	for _, capability := range agent.CapabilitiesConfig.Capabilities {
		capType := capability.Type

		capExecutor, err := a.agentCapabilitiesService.GetExecutor(capType)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to get capability executor"))
			return err
		}

		input := capExecutor.GetInput()
		err = utils.MapToStruct(allParams, input)
		if err != nil {
			tracing.TraceErr(span, errors.Wrapf(err, "unable to map input for capability : %s", capType))
			return err
		}

		config := capExecutor.GetConfig()
		if capability.Values != "" {
			data := []byte(capability.Values)
			err = json.Unmarshal(data, config)
			if err != nil {
				tracing.TraceErr(span, errors.Wrapf(err, "unable to unmarshal config for capability : %s", capType))
				return err
			}
		}

		output, err := capExecutor.ExecuteUntyped(ctx, input, config)
		if err != nil {
			tracing.TraceErr(span, err)
			if capability.Optional {
				errs = multierr.Append(errs, err)
				continue
			}
			return fmt.Errorf("capability %s execution failed: %w", capType, err)
		}

		// Merge the output back into allParams for subsequent capabilities
		if output != nil {
			outputMap, err := utils.StructToMap(output)
			if err != nil {
				if capability.Optional {
					fmt.Printf("Optional capability %s output conversion failed: %v. Skipping.\n", capType, err)
					errs = multierr.Append(errs, err)
					continue
				}
				return fmt.Errorf("output conversion failed for capability %s: %w", capType, err)
			}

			utils.MergeMapToMap(outputMap, allParams)
		}
	}

	// update agentExecutionRecord
	_, err = a.postgresRepositories.AgentExecutionRepository.Update(
		ctx,
		executionID,
		utils.NowPtr(),
		nil,
		true,
	)
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
