package workflow

import (
	"errors"
	"fmt"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.SaveFlow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if flowRecord.ID == "" {
		return w.postgres.FlowsRepository.Create(ctx, flowRecord)
	}

	return w.postgres.FlowsRepository.Update(ctx, flowRecord)
}

func (w *workflowService) GetFlowsByTrigger(ctx context.Context, listenerEvent enum.FlowListenerEvent) ([]postgres_entity.Flows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.GetFlowsByTrigger")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set in context")
		tracing.TraceErr(span, err)
		return nil, err

	}

	query := postgres_entity.Flows{
		AgentID: "",
	}

	allFlows, err := w.postgres.FlowsRepository.FindAll(ctx, query)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return allFlows, nil
}

func (w *workflowService) GetNextStepInFlow(ctx context.Context, flowId string, fromNodeId *string) (*interfaces.FlowNextStep, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.getNextStepInFlow")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	if fromNodeId == nil {
		query := postgres_entity.Flows{
			ID: flowId,
		}

		flow, err := w.postgres.FlowsRepository.Find(ctx, query)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if flow == nil {
			err = errors.New("flow does not exist")
			tracing.TraceErr(span, err)
			return nil, err
		}

		fromNodeId = nil
	}

	nextStep, err := w.nextStepInFlow(ctx, flowId, *fromNodeId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get next step for flow %s: %w", flowId, err)
	}

	// validate the transition to next action
	validAgent, err := w.ValidateTransition(ctx, *fromNodeId, nextStep.ToNodeID)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to validate transition for flow %s: %w", flowId, err)
	}
	if !validAgent {
		err = fmt.Errorf("not a valid action transition for flow %s", flowId)
		tracing.TraceErr(span, err)
		return nil, err
	}
	return nextStep, nil
}

func (w *workflowService) GetFlowsForListenerEvent(ctx context.Context, sourceEvent enum.FlowListenerEvent, eventType string, eventData any) ([]postgres_entity.Flows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.getFlowsForEvent")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	// check to see if tenant has flows configured for this event
	flows, err := w.GetFlowsByTrigger(ctx, sourceEvent)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	// send to dead events if no flows configured to receive event
	// if err != nil || len(*flows) == 0 {
	// 	err := w.SendToDeadEvents(ctx, sourceEvent, eventType, eventData)
	// 	if err != nil {
	// 		tracing.TraceErr(span, err)
	// 		return nil, err
	// 	}
	// 	return nil, err
	// }

	return flows, nil
}

func (w *workflowService) nextStepInFlow(ctx context.Context, flowId, fromNodeId string) (*interfaces.FlowNextStep, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.GetNextStepInFlow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if flowId == "" {
		err := errors.New("flowID is empty")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if fromNodeId == "" {
		err := errors.New("FromNodeID is empty")
		tracing.TraceErr(span, err)
		return nil, err
	}

	query := postgres_entity.FlowEdge{
		FlowID:     flowId,
		FromNodeID: fromNodeId,
	}

	edges, err := w.postgres.FlowEdgeRepository.FindAll(ctx, query)
	if err != nil {
		tracing.TraceErr(span, err)
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
		tracing.TraceErr(span, err)
		return nil, err
	}

	if nextNodeDetails == nil {
		err := errors.New("No flow node record found for ToNode")
		tracing.TraceErr(span, err)
		return nil, err
	}

	toNodeType, err := enum.GetFlowNodeType(nextNodeDetails.Type)
	if err != nil {
		tracing.TraceErr(span, err)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.GetFlowExecutionRecordById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)
	span.LogKV("id", id)
	tenant := common.GetTenantFromContext(ctx)

	if id == "" {
		err := errors.New("ID is empty")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if tenant == "" {
		err := errors.New("Tenant is not set")
		tracing.TraceErr(span, err)
		return nil, err
	}

	searchParams := postgres_entity.FlowExecution{
		ID: id,
	}

	result, err := w.postgres.FlowExecutionRepository.Find(ctx, searchParams)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return result, nil
}

func (w *workflowService) BuildAndSaveFlowExecutionRecord(
	ctx context.Context,
	flowStatus, flowId, flowNodeId, entityId, entityType, currentStep string,
	eventData any,
) (*postgres_entity.FlowExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.BuildAndSaveFlowExecutionRecord")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

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
		tracing.TraceErr(span, err)
		return nil, err
	}

	return flowExecutionRecord, nil
}

func (w *workflowService) SaveFlowExecutionRecord(ctx context.Context, flowExecutionRecord postgres_entity.FlowExecution) (*postgres_entity.FlowExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.SaveFlowExecutionRecord")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if flowExecutionRecord.FlowID == "" {
		err := errors.New("Tenant, FlowID, or EntityID missing")
		span.LogFields(log.Object("executionRecord", flowExecutionRecord))
		tracing.TraceErr(span, err)
		return nil, err
	}

	if flowExecutionRecord.ID == "" {
		return w.postgres.FlowExecutionRepository.Create(ctx, flowExecutionRecord)
	}

	return w.postgres.FlowExecutionRepository.Update(ctx, flowExecutionRecord)
}

// Flow Agent Execution

func (w *workflowService) SaveFlowAgentExecutionRecord(ctx context.Context, flowAgentExecutionRecord postgres_entity.AgentExecution) (*postgres_entity.AgentExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.SaveFlowExecutionAgentRecord")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set in context")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if flowAgentExecutionRecord.ID == "" {
		if *flowAgentExecutionRecord.AgentID == "" || flowAgentExecutionRecord.ID == "" {
			span.LogFields(log.Object("flowAgentExecutionRecord", flowAgentExecutionRecord))
			err := errors.New("Agent or FlowExecutionID missing")
			tracing.TraceErr(span, err)
			return nil, err
		}
		return w.postgres.AgentExecutionRepository.Create(ctx, flowAgentExecutionRecord)
	}

	return w.postgres.AgentExecutionRepository.Update(ctx, flowAgentExecutionRecord)
}
