package agent

import (
	"context"
	"time"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type agentRunnerService struct {
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
) interfaces.AgentRunnerService {
	if postgresRepositories == nil {
		panic("postgresRepositories cannot be nil")
	}
	if agentCapabilities == nil {
		panic("agentCapabilities cannot be nil")
	}
	if agentService == nil {
		panic("agentService cannot be nil")
	}

	return &agentRunnerService{
		postgresRepositories:       postgresRepositories,
		agentCapabilitiesService:   agentCapabilities,
		agentService:               agentService,
		capabilityExecutionService: capabilityExecution,
	}
}

type executionParams struct {
	agent         postgres_entity.Agent
	triggerEvent  enum.AgentListenerEvent
	executionID   string
	initialParams map[string]any
	span          opentracing.Span
}

type capabilityParams struct {
	executionID       string
	capabilityTypeStr string
	agent             postgres_entity.Agent
	allParams         map[string]any
	untypedExecutors  map[enum.AgentCapability]interfaces.AgentCapabilityUntyped
	span              opentracing.Span
}

func (a *agentRunnerService) Run(ctx context.Context, agent postgres_entity.Agent, agentEventName string, initialParams map[string]any, existingExecutionId *string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.Run")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("agentID", agent.ID), log.String("agentEventName", agentEventName), log.String("existingExecutionId", utils.IfNotNilString(existingExecutionId)))

	if err := a.validateAgent(ctx, agent); err != nil {
		return "", err
	}

	triggerEvent, err := a.validateAndGetTriggerEvent(ctx, agentEventName)
	if err != nil {
		return "", err
	}

	executionID := utils.IfNotNilString(existingExecutionId)
	if executionID == "" {
		executionID, err = a.setupExecution(ctx, agent, agentEventName)
		if err != nil {
			return "", err
		}
	}

	return executionID, a.processCapabilities(ctx, executionParams{
		agent:         agent,
		triggerEvent:  triggerEvent,
		executionID:   executionID,
		initialParams: initialParams,
	})
}

func (a *agentRunnerService) validateAgent(ctx context.Context, agent postgres_entity.Agent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.validateAgent")
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
	return nil
}

func (a *agentRunnerService) validateAndGetTriggerEvent(ctx context.Context, agentEventName string) (enum.AgentListenerEvent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.validateAndGetTriggerEvent")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	triggerEvent, err := enum.GetAgentListener(agentEventName)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	return triggerEvent, nil
}

func (a *agentRunnerService) setupExecution(ctx context.Context, agent postgres_entity.Agent, agentEventName string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.setupExecution")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	executionID, err := a.createAgentExecutionRecord(ctx, agent.ID, agentEventName, tracing.GetTraceId(span))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "unable to create agent execution record"))
		return "", err
	}
	span.LogFields(log.String("result.executionID", executionID))
	return executionID, nil
}

func (a *agentRunnerService) processCapabilities(ctx context.Context, params executionParams) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.processCapabilities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	play, err := a.postgresRepositories.AgentRegistryRepository.FindPlay(ctx, params.agent.Type, params.triggerEvent)
	if err != nil {
		tracing.TraceErr(params.span, err)
		return err
	}

	allParams := make(map[string]any)
	utils.MergeMapToMap(params.initialParams, allParams)
	untypedExecutors := a.agentCapabilitiesService.GetExecutors()

	for _, capabilityTypeStr := range *&play.Capabilities {
		status, err := a.executeCapability(ctx, capabilityParams{
			executionID:       params.executionID,
			capabilityTypeStr: capabilityTypeStr,
			agent:             params.agent,
			allParams:         allParams,
			untypedExecutors:  untypedExecutors,
			span:              params.span,
		})
		if err != nil {
			return err
		}
		if status == enum.CapabilityExecutionStop {
			break
		}
		if status != enum.CapabilityExecutionCompleted {
			return nil
		}
	}

	return a.postgresRepositories.AgentExecutionRepository.Finish(ctx, params.executionID)
}

func (a *agentRunnerService) executeCapability(ctx context.Context, params capabilityParams) (enum.CapabilityExecutionStatus, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.executeCapability")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	execution, err := a.getExecution(ctx, params.executionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, err
	}

	// Check if this capability was already completed
	if execution.Checkpoints != nil {
		if checkpoint, exists := execution.Checkpoints[params.capabilityTypeStr]; exists {
			if resultMap, ok := checkpoint.(map[string]any); ok {
				utils.MergeMapToMap(resultMap, params.allParams)
				return enum.CapabilityExecutionCompleted, nil
			}
		}
	}

	capability, err := a.getCapability(ctx, params)
	if err != nil {
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, err
	}
	if capability == nil {
		err = errors.New("capability not found")
		tracing.TraceErr(span, err)
		return enum.CapabilityExecutionError, err
	}
	if !capability.Active {
		span.LogFields(log.String("result", "capability not active"))
		return enum.CapabilityExecutionCompleted, nil
	}

	executionContainer := interfaces.ExecutionContainer{
		AgentExecutionID: params.executionID,
		Capability:       *capability,
		ExecutionParams:  params.allParams,
		UntypedExecutors: params.untypedExecutors,
	}
	status, output, execErr := a.capabilityExecutionService.Execute(ctx, executionContainer)
	if err = a.handleExecutionResult(ctx, params, status, output, execErr); err != nil {
		tracing.TraceErr(span, err)
		return status, err
	}

	if status == enum.CapabilityExecutionCompleted {
		utils.MergeMapToMap(output, params.allParams)
	}

	return status, nil
}

func (a *agentRunnerService) getExecution(ctx context.Context, executionID string) (*postgres_entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.getExecution")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	execution, err := a.postgresRepositories.AgentExecutionRepository.GetById(ctx, executionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if execution == nil {
		err = errors.New("execution not found")
		tracing.TraceErr(span, err)
		return nil, err
	}
	return execution, nil
}

func (a *agentRunnerService) getCapability(ctx context.Context, params capabilityParams) (*postgres_entity.Capability, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.getCapability")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	capabilityType, err := enum.GetAgentCapability(params.capabilityTypeStr)
	if err != nil {
		tracing.TraceErr(params.span, err)
		return nil, err
	}

	capability, err := a.postgresRepositories.AgentRepository.FindCapability(ctx, params.agent.ID, capabilityType)
	if err != nil {
		tracing.TraceErr(params.span, err)
		return nil, err
	}
	if capability == nil {
		err = errors.New("cannot identify capability")
		tracing.TraceErr(params.span, err,
			log.String("capabilityType", capabilityType.String()),
			log.String("agentID", params.agent.ID))
		return nil, err
	}
	return capability, nil
}

func (a *agentRunnerService) handleExecutionResult(ctx context.Context, params capabilityParams, status enum.CapabilityExecutionStatus, output map[string]any, execErr error) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.handleExecutionResult")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("status", status.String()))

	execErrStr := ""
	if execErr != nil {
		execErrStr = execErr.Error()
	}

	switch status {
	case enum.CapabilityExecutionError:
		if err := a.postgresRepositories.AgentExecutionRepository.Fail(ctx, params.executionID, execErrStr); err != nil {
			return errors.Wrap(err, "unable to update agent execution record")
		}
		return nil

	case enum.CapabilityExecutionRetry:
		// Save state for retry
		stateData := map[string]any{
			"params": params.allParams,
		}
		retryErr := a.postgresRepositories.AgentExecutionRepository.ScheduleRetry(ctx, params.executionID, execErr, stateData)
		if retryErr != nil {
			tracing.TraceErr(span, retryErr)
			// If retry scheduling fails, mark as failed
			if err := a.postgresRepositories.AgentExecutionRepository.Fail(ctx, params.executionID, execErrStr); err != nil {
				return errors.Wrap(err, "unable to update agent execution record")
			}
		}
		return nil

	case enum.CapabilityExecutionCompleted:
		checkpointData := map[string]any{
			"status":       status.String(),
			"output":       output,
			"completed_at": utils.Now(),
		}
		if err := a.postgresRepositories.AgentExecutionRepository.CompleteStep(ctx, params.executionID, params.capabilityTypeStr, checkpointData); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		return nil

	case enum.CapabilityExecutionStop:
		checkpointData := map[string]any{
			"status":     status.String(),
			"output":     output,
			"stopped_at": utils.Now(),
		}
		if err := a.postgresRepositories.AgentExecutionRepository.CompleteStep(ctx, params.executionID, params.capabilityTypeStr, checkpointData); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		return nil

	case enum.CapabilityExecutionPending:
		// Save current step and state for async resume
		currentStep := output["current_step"].(string)
		stateData := map[string]any{
			"step":   currentStep,
			"params": output,
		}
		return a.postgresRepositories.AgentExecutionRepository.SaveAsyncState(ctx, params.executionID, currentStep, stateData)

	default:
		err := errors.New("unexpected capability execution status")
		tracing.TraceErr(span, err)
		return err
	}
}

func (a *agentRunnerService) createAgentExecutionRecord(ctx context.Context, agentID, triggerEventName, traceId string) (string, error) {
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

	createdRecordId, err := a.agentService.CreateAgentExecutionRecord(ctx, *agent, triggerEventName, traceId)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	span.LogFields(log.String("result.agentExecutionId", createdRecordId))
	return createdRecordId, err
}

func (a *agentRunnerService) ResumeExecution(ctx context.Context, executionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.ResumeExecution")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Get execution details
	execution, err := a.postgresRepositories.AgentExecutionRepository.GetById(ctx, executionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if execution.Status != enum.AgentExecutionPending {
		return errors.New("execution is not in pending state")
	}

	// Get agent details
	agent, err := a.postgresRepositories.AgentRepository.Find(ctx, postgres_entity.Agent{
		ID: *execution.AgentID,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Resume from stored state
	params := map[string]any{}
	if execution.StateData != nil {
		params = execution.StateData["params"].(map[string]any)
	}

	// Resume execution
	_, err = a.Run(ctx, *agent, execution.TriggerEvent, params, utils.StringPtr(executionID))
	return err
}

func (a *agentRunnerService) RetryExecution(ctx context.Context, executionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.RetryExecution")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	// Get execution details
	execution, err := a.postgresRepositories.AgentExecutionRepository.GetById(ctx, executionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if execution.Status != enum.AgentExecutionRetrying {
		return errors.New("execution is not in retry state")
	}

	// Check retry timing
	if execution.NextRetryAt != nil && execution.NextRetryAt.After(time.Now()) {
		return errors.New("retry attempt too early")
	}

	// Get agent details
	agent, err := a.postgresRepositories.AgentRepository.Find(ctx, postgres_entity.Agent{
		ID: *execution.AgentID,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Resume from last known state
	params := map[string]any{}
	if execution.StateData != nil {
		params = execution.StateData["params"].(map[string]any)
	}

	// Retry execution
	_, err = a.Run(ctx, *agent, execution.TriggerEvent, params, utils.StringPtr(executionID))
	return err
}

func (a *agentRunnerService) GetExecutionStatus(ctx context.Context, executionID string) (*postgres_entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.GetExecutionStatus")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	execution, err := a.postgresRepositories.AgentExecutionRepository.GetById(ctx, executionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return execution, nil
}
