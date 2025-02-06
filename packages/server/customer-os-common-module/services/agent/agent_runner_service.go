package agent

import (
	"context"
	"github.com/opentracing/opentracing-go/log"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type AgentRunnerService struct {
	postgresRepositories       *postgres_repository.Repositories
	agentCapabilitiesService   *agent_capability.AgentCapabilities
	agentService               interfaces.AgentService
	capabilityExecutionService interfaces.AgentCapabilityExecutionService
}

func NewAgentRunnerService(
	postgresRepositories *postgres_repository.Repositories,
	agentCapabilities *agent_capability.AgentCapabilities,
	agentService interfaces.AgentService,
	capabilityExecution interfaces.AgentCapabilityExecutionService,
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
		postgresRepositories:       postgresRepositories,
		agentCapabilitiesService:   agentCapabilities,
		agentService:               agentService,
		capabilityExecutionService: capabilityExecution,
	}
}

func (a *AgentRunnerService) Run(ctx context.Context, agent postgres_entity.Agent, agentEventName string, initialParams map[string]any) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.Run")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// validation
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

	triggerEvent, err := enum.GetAgentListener(agentEventName)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// create execution record
	executionID, err := a.createAgentExecutionRecord(ctx, agent.ID, agentEventName, tracing.GetTraceId(span))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to create agent execution record"))
		return err
	}

	// Create a copy of initialParams to avoid mutating the original map
	allParams := make(map[string]any)
	utils.MergeMapToMap(initialParams, allParams)

	// get capability executors
	untypedExecutors := a.agentCapabilitiesService.GetExecutors()

	// get agent Type and lookup the play for the trigger event
	play, err := a.postgresRepositories.AgentRegistryRepository.FindPlay(ctx, agent.Type, triggerEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	var capErr error
	for _, capabilityTypeStr := range *&play.Capabilities {

		capabilityType, err := enum.GetAgentCapability(capabilityTypeStr)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		capability, err := a.postgresRepositories.AgentRepository.FindCapability(ctx, agent.ID, capabilityType)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if capability == nil {
			err = errors.New("cannot identify capability")
			tracing.TraceErr(span, err, log.String("capabilityType", capabilityType.String()), log.String("agentID", agent.ID))
			return err
		}

		if !capability.Active {
			continue
		}

		executionContainer := interfaces.ExecutionContainer{
			AgentExecutionID: executionID,
			Capability:       *capability,
			ExecutionParams:  allParams,
			UntypedExecutors: untypedExecutors,
		}

		output, err := a.capabilityExecutionService.Execute(ctx, executionContainer)
		if err != nil {
			tracing.TraceErr(span, err)
			capErr = err
			break
		}

		utils.MergeMapToMap(output, allParams)
	}

	if capErr != nil {
		tracing.TraceErr(span, capErr)
		_, dbErr := a.postgresRepositories.AgentExecutionRepository.Update(ctx, executionID, nil, utils.StringPtr(capErr.Error()), false)
		if dbErr != nil {
			tracing.TraceErr(span, errors.Wrap(dbErr, "unable to update agent execution record"))
			return dbErr
		}
		return nil
	}

	// update agentExecutionRecord
	_, err = a.postgresRepositories.AgentExecutionRepository.Update(ctx, executionID, utils.NowPtr(), nil, true)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (a *AgentRunnerService) createAgentExecutionRecord(ctx context.Context, agentID, triggerEventName, traceId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.createAgentExecutionRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("agentID", agentID, "triggerEventType", triggerEventName, "traceId", traceId)

	agent, err := a.postgresRepositories.AgentRepository.Find(ctx, postgres_entity.Agent{
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

	return a.agentService.CreateAgentExecutionRecord(ctx, *agent, triggerEventName, traceId)
}
