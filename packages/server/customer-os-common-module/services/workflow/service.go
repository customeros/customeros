package workflow

import (
	"errors"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"golang.org/x/net/context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type workflowService struct {
	postgres *postgres_repository.Repositories
}

func NewWorkflowService(postgres *postgres_repository.Repositories) interfaces.WorkflowService {
	return &workflowService{
		postgres: postgres,
	}
}

// Flow
func (w *workflowService) SaveFlow(ctx context.Context, flowRecord postgres_entity.Flows) (*postgres_entity.Flows, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WorkflowService.SaveFlow")
	defer spans.Finish()

	if flowRecord.ID == "" {
		return w.postgres.FlowsRepository.Create(ctx, flowRecord)
	}

	return w.postgres.FlowsRepository.Update(ctx, flowRecord)
}

func (w *workflowService) GetFlowsByTrigger(ctx context.Context, listenerEvent enum.AgentListenerEvent) ([]postgres_entity.Flows, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WorkflowService.GetFlowsByTrigger")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not set in context")
		spans.TraceError(err)
		return nil, err

	}

	query := postgres_entity.Flows{
		AgentID: "",
	}

	allFlows, err := w.postgres.FlowsRepository.FindAll(ctx, query)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return allFlows, nil
}

func (w *workflowService) GetNextStepInFlow(ctx context.Context, flowId string, fromNodeId *string) (*interfaces.FlowNextStep, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EventHandlers.getNextStepInFlow")
	defer spans.Finish()

	if fromNodeId == nil {
		query := postgres_entity.Flows{
			ID: flowId,
		}

		flow, err := w.postgres.FlowsRepository.Find(ctx, query)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}
		if flow == nil {
			err = errors.New("flow does not exist")
			spans.TraceError(err)
			return nil, err
		}

		fromNodeId = nil
	}

	nextStep, err := w.nextStepInFlow(ctx, flowId, *fromNodeId)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to get next step for flow %s: %w", flowId, err)
	}

	// validate the transition to next action
	validAgent, err := w.ValidateTransition(ctx, *fromNodeId, nextStep.ToNodeID)
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("failed to validate transition for flow %s: %w", flowId, err)
	}
	if !validAgent {
		err = fmt.Errorf("not a valid action transition for flow %s", flowId)
		spans.TraceError(err)
		return nil, err
	}
	return nextStep, nil
}

func (w *workflowService) GetFlowsForListenerEvent(ctx context.Context, sourceEvent enum.AgentListenerEvent, eventType string, eventData any) ([]postgres_entity.Flows, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "EventHandlers.getFlowsForEvent")
	defer spans.Finish()

	// check to see if tenant has flows configured for this event
	flows, err := w.GetFlowsByTrigger(ctx, sourceEvent)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	// send to dead events if no flows configured to receive event
	// if err != nil || len(*flows) == 0 {
	// 	err := w.SendToDeadEvents(ctx, sourceEvent, eventType, eventData)
	// 	if err != nil {
	// 		spans.TraceError( err)
	// 		return nil, err
	// 	}
	// 	return nil, err
	// }

	return flows, nil
}

func (w *workflowService) nextStepInFlow(ctx context.Context, flowId, fromNodeId string) (*interfaces.FlowNextStep, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WorkflowService.GetNextStepInFlow")
	defer spans.Finish()

	if flowId == "" {
		err := errors.New("flowID is empty")
		spans.TraceError(err)
		return nil, err
	}

	if fromNodeId == "" {
		err := errors.New("FromNodeID is empty")
		spans.TraceError(err)
		return nil, err
	}

	query := postgres_entity.FlowEdge{
		FlowID:     flowId,
		FromNodeID: fromNodeId,
	}

	edges, err := w.postgres.FlowEdgeRepository.FindAll(ctx, query)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	if len(edges) == 0 { // mark flow as completed
		return &interfaces.FlowNextStep{
			FlowID:     flowId,
			FromNodeID: fromNodeId,
			ToNodeType: enum.NodeFlowEnd,
		}, nil
	}

	// take the first result
	// todo - handle conditions
	edge := (edges)[0]

	nextNodeQuery := postgres_entity.FlowNode{
		ID:     edge.ToNodeID,
		FlowID: edge.FlowID,
	}

	nextNodeDetails, err := w.postgres.FlowNodeRepository.Find(ctx, nextNodeQuery)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	if nextNodeDetails == nil {
		err := errors.New("No flow node record found for ToNode")
		spans.TraceError(err)
		return nil, err
	}

	toNodeType, err := enum.GetFlowNodeType(nextNodeDetails.Type)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	nextStep := interfaces.FlowNextStep{
		FlowID:     edge.FlowID,
		FromNodeID: edge.FromNodeID,
		ToNodeID:   edge.ToNodeID,
		ToNodeType: toNodeType,
	}

	if toNodeType == enum.NodeFlowAgent {
		nextStep.ToNodeAgent = enum.FlowAgent(*nextNodeDetails.AgentID)
	}
	return &nextStep, nil
}

// Flow Execution

func (w *workflowService) GetFlowExecutionRecordById(ctx context.Context, id string) (*postgres_entity.FlowExecution, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WorkflowService.GetFlowExecutionRecordById")
	defer spans.Finish()

	spans.LogKV("id", id)
	tenant := common.GetTenantFromContext(ctx)

	if id == "" {
		err := errors.New("ID is empty")
		spans.TraceError(err)
		return nil, err
	}

	if tenant == "" {
		err := errors.New("Tenant is not set")
		spans.TraceError(err)
		return nil, err
	}

	searchParams := postgres_entity.FlowExecution{
		ID: id,
	}

	result, err := w.postgres.FlowExecutionRepository.Find(ctx, searchParams)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	return result, nil
}

func (w *workflowService) BuildAndSaveFlowExecutionRecord(
	ctx context.Context,
	flowStatus, flowId, flowNodeId, entityId, entityType, currentStep string,
	eventData any,
) (*postgres_entity.FlowExecution, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WorkflowService.BuildAndSaveFlowExecutionRecord")
	defer spans.Finish()

	record := postgres_entity.FlowExecution{
		FlowID:            flowId,
		Status:            enum.FlowExecutionRunning.String(),
		StartedAt:         utils.NowPtr(),
		CurrentStepNodeId: flowNodeId,
		CreatedAt:         utils.Now(),
	}

	switch flowStatus {
	case "INACTIVE":
		reason := enum.FlowBlockedNotActive.String()
		record.Status = enum.FlowExecutionBlocked.String()
		record.BlockedReason = &reason

	case "ARCHIVED":
		reason := enum.FlowBlockedArchived.String()
		record.Status = enum.FlowExecutionBlocked.String()
		record.BlockedReason = &reason

	default:
		record.Status = enum.FlowExecutionRunning.String()
	}

	// save flow execution record
	flowExecutionRecord, err := w.SaveFlowExecutionRecord(ctx, record)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	return flowExecutionRecord, nil
}

func (w *workflowService) SaveFlowExecutionRecord(ctx context.Context, flowExecutionRecord postgres_entity.FlowExecution) (*postgres_entity.FlowExecution, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WorkflowService.SaveFlowExecutionRecord")
	defer spans.Finish()

	if flowExecutionRecord.FlowID == "" {
		err := errors.New("Tenant, FlowID, or EntityID missing")
		spans.LogObjectAsJson("executionRecord", flowExecutionRecord)
		spans.TraceError(err)
		return nil, err
	}

	if flowExecutionRecord.ID == "" {
		return w.postgres.FlowExecutionRepository.Create(ctx, flowExecutionRecord)
	}

	return w.postgres.FlowExecutionRepository.Update(ctx, flowExecutionRecord)
}

// Flow Agent Execution

func (w *workflowService) SaveFlowAgentExecutionRecord(ctx context.Context, flowAgentExecutionRecord postgres_entity.AgentExecution) (*postgres_entity.AgentExecution, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WorkflowService.SaveFlowExecutionAgentRecord")
	defer spans.Finish()

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set in context")
		spans.TraceError(err)
		return nil, err
	}

	if flowAgentExecutionRecord.ID == "" {
		if *flowAgentExecutionRecord.AgentID == "" || flowAgentExecutionRecord.ID == "" {
			spans.LogObjectAsJson("flowAgentExecutionRecord", flowAgentExecutionRecord)
			err := errors.New("Agent or FlowExecutionID missing")
			spans.TraceError(err)
			return nil, err
		}
		return w.postgres.AgentExecutionRepository.Create(ctx, flowAgentExecutionRecord)
	}

	return &postgres_entity.AgentExecution{}, nil
}
