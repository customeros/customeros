package service

import (
	"errors"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type WorkflowService interface {
	// Flow
	SaveFlow(ctx context.Context, flowRecord entity.Flows) (*entity.Flows, error)
	GetFlowsByTrigger(ctx context.Context, listenerEvent enum.FlowListenerEvent) ([]entity.Flows, error)
	GetNextStepInFlow(ctx context.Context, flowId string, fromNodeId *string) (*FlowNextStep, error)
	GetFlowsForListenerEvent(ctx context.Context, sourceEvent enum.FlowListenerEvent, eventType string, eventData any) ([]entity.Flows, error)

	// Flow Execution
	BuildAndSaveFlowExecutionRecord(ctx context.Context, flowStatus, flowId, flowNodeId, entityId, entityType, currentStep string, eventData any) (*entity.FlowExecution, error)
	SaveFlowExecutionRecord(ctx context.Context, flowExecutionRecord entity.FlowExecution) (*entity.FlowExecution, error)
	GetFlowExecutionRecordById(ctx context.Context, id string) (*entity.FlowExecution, error)

	// FlowAgent Execution
	SaveFlowAgentExecutionRecord(ctx context.Context, AgentExecutionRecord entity.AgentExecution) (*entity.AgentExecution, error)

	// Dead Flow Events
	// SendToDeadEvents(ctx context.Context, sourceEvent enum.FlowListenerEvent, eventType string, eventData any) error

	// Validation
	ValidateEventType(ctx context.Context, nodeType enum.FlowNodeType, event string) bool
	ValidateFlowBelongsToTenant(ctx context.Context, flowId string) (bool, error)
	// ValidateListener(ctx context.Context, listenerEvent enum.FlowListenerEvent) (bool, error)
	ValidateNodeType(ctx context.Context, nodeType string) (bool, *enum.FlowNodeType)
	ValidateTransition(ctx context.Context, fromNodeId string, toNodeId string) (bool, error)
}

type workflowService struct {
	services *Services
}

func NewWorkflowService(services *Services) WorkflowService {
	return &workflowService{
		services: services,
	}
}

// Flow
func (w *workflowService) SaveFlow(ctx context.Context, flowRecord entity.Flows) (*entity.Flows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.SaveFlow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if flowRecord.ID == "" {
		return w.services.PostgresRepositories.FlowsRepository.Create(ctx, flowRecord)
	}

	return w.services.PostgresRepositories.FlowsRepository.Update(ctx, flowRecord)
}

func (w *workflowService) GetFlowsByTrigger(ctx context.Context, listenerEvent enum.FlowListenerEvent) ([]entity.Flows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.GetFlowsByTrigger")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set in context")
		tracing.TraceErr(span, err)
		return nil, err

	}

	query := entity.Flows{
		AgentID: "",
	}

	allFlows, err := w.services.PostgresRepositories.FlowsRepository.FindAll(ctx, query)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return allFlows, nil
}

type FlowNextStep struct {
	FlowID      string
	FromNodeID  string
	ToNodeID    string
	ToNodeType  enum.FlowNodeType
	ToNodeAgent enum.FlowAgent
}

func (w *workflowService) GetNextStepInFlow(ctx context.Context, flowId string, fromNodeId *string) (*FlowNextStep, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.getNextStepInFlow")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	if fromNodeId == nil {
		query := entity.Flows{
			ID: flowId,
		}

		flow, err := w.services.PostgresRepositories.FlowsRepository.Find(ctx, query)
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

func (w *workflowService) GetFlowsForListenerEvent(ctx context.Context, sourceEvent enum.FlowListenerEvent, eventType string, eventData any) ([]entity.Flows, error) {
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

func (w *workflowService) nextStepInFlow(ctx context.Context, flowId, fromNodeId string) (*FlowNextStep, error) {
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

	query := entity.FlowEdge{
		FlowID:     flowId,
		FromNodeID: fromNodeId,
	}

	edges, err := w.services.PostgresRepositories.FlowEdgeRepository.FindAll(ctx, query)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if len(edges) == 0 { // mark flow as completed
		return &FlowNextStep{
			FlowID:     flowId,
			FromNodeID: fromNodeId,
			ToNodeType: enum.NodeFlowEnd,
		}, nil
	}

	// take the first result
	// todo - handle conditions
	edge := (edges)[0]

	nextNodeQuery := entity.FlowNode{
		ID:     edge.ToNodeID,
		FlowID: edge.FlowID,
	}

	nextNodeDetails, err := w.services.PostgresRepositories.FlowNodeRepository.Find(ctx, nextNodeQuery)
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

	nextStep := FlowNextStep{
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

func (w *workflowService) GetFlowExecutionRecordById(ctx context.Context, id string) (*entity.FlowExecution, error) {
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

	searchParams := entity.FlowExecution{
		ID: id,
	}

	result, err := w.services.PostgresRepositories.FlowExecutionRepository.Find(ctx, searchParams)
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
) (*entity.FlowExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.BuildAndSaveFlowExecutionRecord")
	defer span.Finish()
	tracing.SetDefaultListenerSpanTags(ctx, span)

	record := entity.FlowExecution{
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

func (w *workflowService) SaveFlowExecutionRecord(ctx context.Context, flowExecutionRecord entity.FlowExecution) (*entity.FlowExecution, error) {
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
		return w.services.PostgresRepositories.FlowExecutionRepository.Create(ctx, flowExecutionRecord)
	}

	return w.services.PostgresRepositories.FlowExecutionRepository.Update(ctx, flowExecutionRecord)
}

// Flow Agent Execution

func (w *workflowService) SaveFlowAgentExecutionRecord(ctx context.Context, flowAgentExecutionRecord entity.AgentExecution) (*entity.AgentExecution, error) {
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
		return w.services.PostgresRepositories.AgentExecutionRepository.Create(ctx, flowAgentExecutionRecord)
	}

	return w.services.PostgresRepositories.AgentExecutionRepository.Update(ctx, flowAgentExecutionRecord)
}

// Dead Flow Events

// func (w *workflowService) SendToDeadEvents(ctx context.Context, sourceEvent enum.FlowListenerEvent, eventType string, eventData any) error {
// 	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.SendToDeadEvents")
// 	defer span.Finish()
// 	tracing.SetDefaultListenerSpanTags(ctx, span)
//
// 	deadEvent, err := w.buildDeadEvent(ctx, sourceEvent, eventType, eventData)
// 	if err != nil {
// 		return err
// 	}
//
// 	_, err = w.services.PostgresRepositories.FlowDeadEventsRepository.Create(ctx, deadEvent)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return err
// 	}
// 	return nil
// }
//
// func (w *workflowService) buildDeadEvent(ctx context.Context, sourceEvent enum.FlowListenerEvent, eventType string, eventData any) (entity.FlowDeadEvents, error) {
// 	span, ctx := opentracing.StartSpanFromContext(ctx, "EventHandlers.buildDeadEventFromMeetingSummary")
// 	defer span.Finish()
// 	tracing.SetDefaultListenerSpanTags(ctx, span)
//
// 	data, err := utils.ObjectToString(eventData)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return entity.FlowDeadEvents{}, err
// 	}
//
// 	return entity.FlowDeadEvents{
// 		Tenant:    common.GetTenantFromContext(ctx),
// 		NodeType:  enum.NodeFlowListenerEvent.String(),
// 		EventType: eventType,
// 		Event:     sourceEvent.String(),
// 		CreatedAt: utils.Now(),
// 		Data:      &data,
// 	}, nil
// }
