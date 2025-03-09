package agent

import (
	"context"
	"fmt"
	"time"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type agentRunnerService struct {
	postgresRepositories       *postgres_repository.Repositories
	opensearchService          interfaces.OpensearchService
	agentCapabilitiesService   *agent_capability.AgentCapabilities
	agentService               interfaces.AgentService
	capabilityExecutionService interfaces.AgentCapabilityExecutionService
}

func NewAgentRunnerService(
	postgresRepositories *postgres_repository.Repositories,
	opensearchService interfaces.OpensearchService,
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
		opensearchService:          opensearchService,
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
		span:          span,
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
		tracing.TraceErr(span, err)
		return err
	}

	allParams := make(map[string]any)
	utils.MergeMapToMap(params.initialParams, allParams)
	untypedExecutors := a.agentCapabilitiesService.GetExecutors()

	for _, capabilityTypeStr := range *&play.Capabilities {
		// build observability metrics
		metrics := a.newObservabilityContainer(params)
		metrics.Capability = capabilityTypeStr

		status, err := a.executeCapability(ctx, metrics, capabilityParams{
			executionID:       params.executionID,
			capabilityTypeStr: capabilityTypeStr,
			agent:             params.agent,
			allParams:         allParams,
			untypedExecutors:  untypedExecutors,
			span:              params.span,
		})
		if !metrics.SkipPublishingObserbility {
			metrics.Status = status.String()
			a.pushObservabilityMetrics(ctx, metrics)
		}

		if err != nil {
			return err
		}
		if status == enum.CapabilityExecutionStop {
			break
		} else if status == enum.CapabilityExecutionSkip {
			continue
		} else if status != enum.CapabilityExecutionCompleted {
			return nil
		}
	}

	return a.postgresRepositories.AgentExecutionRepository.Finish(ctx, params.executionID)
}

func (a *agentRunnerService) executeCapability(ctx context.Context, metrics *dto.AgentExecutionObservability, params capabilityParams) (enum.CapabilityExecutionStatus, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.executeCapability")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	execution, err := a.getExecution(ctx, params.executionID)
	if err != nil {
		tracing.TraceErr(span, err)
		metrics.ErrorMessage = err.Error()
		metrics.Success = false
		return enum.CapabilityExecutionError, err
	}

	// Check if this capability was already completed
	if execution.Checkpoints != nil {
		if checkpoint, exists := execution.Checkpoints[params.capabilityTypeStr]; exists {
			if resultMap, ok := checkpoint.(map[string]any); ok {
				metrics.SkipPublishingObserbility = true
				utils.MergeMapToMap(resultMap, params.allParams)
				return enum.CapabilityExecutionCompleted, nil
			}
		}
	}

	capability, err := a.getCapability(ctx, params)
	if err != nil {
		tracing.TraceErr(span, err)
		metrics.ErrorMessage = err.Error()
		metrics.Success = false
		return enum.CapabilityExecutionError, err
	}
	if capability == nil {
		err = errors.New("capability not found")
		tracing.TraceErr(span, err)
		metrics.ErrorMessage = err.Error()
		metrics.Success = false
		return enum.CapabilityExecutionError, err
	}
	if !capability.Active {
		span.LogFields(log.String("result", "capability not active"))
		return enum.CapabilityExecutionSkip, nil
	}

	executionContainer := interfaces.ExecutionContainer{
		AgentExecutionID: params.executionID,
		Capability:       *capability,
		ExecutionParams:  params.allParams,
		UntypedExecutors: params.untypedExecutors,
	}
	status, output, execErr := a.capabilityExecutionService.Execute(ctx, executionContainer)
	err = a.handleExecutionResult(ctx, metrics, params, status, output, execErr)
	if err != nil {
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

func (a *agentRunnerService) handleExecutionResult(
	ctx context.Context,
	metrics *dto.AgentExecutionObservability,
	params capabilityParams,
	status enum.CapabilityExecutionStatus,
	output map[string]any,
	execErr error,
) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.handleExecutionResult")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("status", status.String()))

	execErrStr := ""
	if execErr != nil {
		execErrStr = execErr.Error()
		metrics.ErrorMessage = execErrStr
	}

	switch status {
	case enum.CapabilityExecutionSkip:
		return nil
	case enum.CapabilityExecutionError:
		metrics.Success = false
		metrics.Retry = false
		if err := a.postgresRepositories.AgentExecutionRepository.Fail(ctx, params.executionID, execErrStr); err != nil {
			return errors.Wrap(err, "unable to update agent execution record")
		}
		return nil

	case enum.CapabilityExecutionRetry:
		metrics.Success = false
		metrics.Retry = true

		// Save state for retry
		stateData := map[string]any{
			"params": params.allParams,
		}
		retryAt, retryErr := a.postgresRepositories.AgentExecutionRepository.ScheduleRetry(ctx, params.executionID, execErr, stateData)
		metrics.RetryAt = retryAt
		if retryErr != nil {
			tracing.TraceErr(span, retryErr)
			// If retry scheduling fails, mark as failed
			if err := a.postgresRepositories.AgentExecutionRepository.Fail(ctx, params.executionID, execErrStr); err != nil {
				return errors.Wrap(err, "unable to update agent execution record")
			}
		}
		return nil

	case enum.CapabilityExecutionCompleted:
		metrics.Success = true
		metrics.CompletedAt = utils.NowPtr()
		metrics.OutputData = output
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
		metrics.Success = true
		metrics.CompletedAt = utils.NowPtr()
		metrics.OutputData = output
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
		err := fmt.Errorf("unexpected capability execution result status {%s}", status.String())
		tracing.TraceErr(span, err)
		return err
	}
}

func (a *agentRunnerService) createAgentExecutionRecord(ctx context.Context, agentID, triggerEventName, traceId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.createAgentExecutionRecord")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogKV("agentID", agentID, "triggerEventType", triggerEventName, "traceId", traceId)

	agent, err := a.postgresRepositories.AgentRepository.GetById(ctx, agentID)
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
	agent, err := a.postgresRepositories.AgentRepository.GetById(ctx, *execution.AgentID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if agent == nil {
		err = a.postgresRepositories.AgentExecutionRepository.Fail(ctx, executionID, "Agent removed")
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		return nil
	}
	if !agent.IsActive {
		_, err = a.postgresRepositories.AgentExecutionRepository.ScheduleRetry(ctx, executionID, errors.New(utils.IfNotNilString(execution.ErrorMessage)), execution.StateData)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	// Resume from stored state
	params := map[string]any{}
	if execution.StateData == nil {
		err := errors.New("StateData is empty, cannot retry")
		tracing.TraceErr(span, err)
		return err
	}

	if paramsVal, ok := execution.StateData["params"]; ok {
		if paramsMap, ok := paramsVal.(map[string]any); ok {
			params = paramsMap
		} else {
			err := errors.New("paramsMap is invalid")
			tracing.TraceErr(span, err)
			return err
		}
	}

	// Resume execution
	_, err = a.Run(ctx, *agent, execution.TriggerEvent, params, utils.StringPtr(executionID))
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return err
}

func (a *agentRunnerService) RerunExecution(ctx context.Context, executionID string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.RerunExecution")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.TagEntity(span, executionID)

	// Get execution details
	execution, err := a.postgresRepositories.AgentExecutionRepository.GetById(ctx, executionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if execution.Status != enum.AgentExecutionRetrying {
		return errors.New("execution is not in retry state")
	}

	// Check if not retrying too early
	if execution.NextRetryAt != nil && execution.NextRetryAt.After(time.Now()) {
		return errors.New("retry attempt too early")
	}

	// Get agent details
	agent, err := a.postgresRepositories.AgentRepository.GetById(ctx, *execution.AgentID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	// if agent was removed stop retrying
	if agent == nil {
		err = a.postgresRepositories.AgentExecutionRepository.Fail(ctx, executionID, "Stop retrying, agent not found")
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to stop retrying"))
			return err
		}
		return nil
	}
	// if agent is not active re-schedule retry
	if !agent.IsActive {
		_, err = a.postgresRepositories.AgentExecutionRepository.ScheduleRetry(ctx, executionID, errors.New(utils.IfNotNilString(execution.ErrorMessage)), execution.StateData)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to schedule retry"))
			return err
		}
		return nil
	}

	// Retry from last known state
	params := map[string]any{}
	if execution.StateData == nil {
		err := errors.New("StateData is empty, cannot retry")
		tracing.TraceErr(span, err)
		return err
	}

	if paramsVal, ok := execution.StateData["params"]; ok {
		if paramsMap, ok := paramsVal.(map[string]any); ok {
			params = paramsMap
		} else {
			err := errors.New("paramsMap is invalid")
			tracing.TraceErr(span, err)
			return err
		}
	}

	// Retry execution
	_, err = a.Run(ctx, *agent, execution.TriggerEvent, params, &executionID)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
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

func (a *agentRunnerService) newObservabilityContainer(executionParams executionParams) *dto.AgentExecutionObservability {
	return &dto.AgentExecutionObservability{
		CapabilityExecutionID: utils.GenerateNanoIdWithPrefix("ace", 16),
		ExecutionID:           executionParams.executionID,
		AgentID:               executionParams.agent.ID,
		AgentType:             executionParams.agent.Type.String(),
		AgentScope:            executionParams.agent.Scope.String(),
		Tenant:                executionParams.agent.Tenant,
		UserID:                executionParams.agent.Owner,
		TriggerEvent:          executionParams.triggerEvent.String(),
		StartedAt:             utils.Now(),
		InputData:             executionParams.initialParams,
		TraceID:               utils.GetTraceIDFromSpan(executionParams.span),
	}
}

func (a *agentRunnerService) pushObservabilityMetrics(ctx context.Context, metrics *dto.AgentExecutionObservability) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentRunnerService.pushObservabilityMetrics")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	index := fmt.Sprintf("agent-%s", utils.CurrentMonth())
	err := a.opensearchService.AgentExecutionObservabilityIndexCheck(ctx, index)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	err = a.opensearchService.UpsertDocument(ctx, index, &metrics.CapabilityExecutionID, metrics)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}
}
