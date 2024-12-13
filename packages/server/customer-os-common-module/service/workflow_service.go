package service

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type WorkflowService interface {
	// Flow
	SaveFlow(ctx context.Context, flowRecord entity.Flow) (*entity.Flow, error)
	GetFlowsByTrigger(ctx context.Context, listenerEvent enum.FlowListenerEvent) (*[]entity.Flow, error)
	GetNextStepInFlow(ctx context.Context, flowId, fromNodeId string) (*FlowNextStep, error)

	// Flow Execution
	SaveFlowExecutionRecord(ctx context.Context, flowExecutionRecord entity.FlowExecution) (*entity.FlowExecution, error)
	GetFlowExecutionRecordById(ctx context.Context, id string) (*entity.FlowExecution, error)

	// FlowAction Execution
	SaveFlowActionExecutionRecord(ctx context.Context, flowActionExecutionRecord entity.FlowActionExecution) (*entity.FlowActionExecution, error)

	// Validation
	ValidateEventType(ctx context.Context, nodeType enum.FlowNodeType, event string) bool
	ValidateFlowBelongsToTenant(ctx context.Context, flowId string) bool
	ValidateListener(ctx context.Context, listenerEvent enum.FlowListenerEvent) (bool, error)
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
func (w *workflowService) SaveFlow(ctx context.Context, flowRecord entity.Flow) (*entity.Flow, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.SaveFlow")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if flowRecord.Tenant == "" {
		flowRecord.Tenant = common.GetTenantFromContext(ctx)
		if flowRecord.Tenant == "" {
			err := errors.New("tenant not set in context")
			tracing.TraceErr(span, err)
			return nil, err
		}
	}

	if flowRecord.ID == "" {
		return w.services.PostgresRepositories.FlowRepository.Create(ctx, flowRecord)
	}

	return w.services.PostgresRepositories.FlowRepository.Update(ctx, flowRecord)

}

func (w *workflowService) GetFlowsByTrigger(ctx context.Context, listenerEvent enum.FlowListenerEvent) (*[]entity.Flow, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.GetFlowsByTrigger")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set in context")
		tracing.TraceErr(span, err)
		return nil, err

	}

	query := entity.Flow{
		Tenant:    tenant,
		TriggerOn: listenerEvent.String(),
	}

	allFlows, err := w.services.PostgresRepositories.FlowRepository.FindAll(ctx, query)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return allFlows, nil
}

type FlowNextStep struct {
	FlowID       string
	FromNodeID   string
	ToNodeID     string
	ToNodeType   enum.FlowNodeType
	ToNodeAction enum.FlowAction
}

func (w *workflowService) GetNextStepInFlow(ctx context.Context, flowId, fromNodeId string) (*FlowNextStep, error) {
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

	if len(*edges) == 0 {
		return nil, nil
	}

	// take the first result
	// todo - handle conditions
	edge := (*edges)[0]

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

	if toNodeType == enum.NodeFlowAction {
		action, err := enum.GetFlowAction(*nextNodeDetails.Event)
		if err != nil {
			tracing.TraceErr(span, err)
		}
		nextStep.ToNodeAction = action
	}
	return &nextStep, nil
}

// Flow Execution

func (w *workflowService) GetFlowExecutionRecordById(ctx context.Context, id string) (*entity.FlowExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.GetFlowExecutionRecordById")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if id == "" {
		err := errors.New("ID is empty")
		tracing.TraceErr(span, err)
		return nil, err
	}

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant is not set")
		tracing.TraceErr(span, err)
		return nil, err
	}

	searchParams := entity.FlowExecution{
		ID:     id,
		Tenant: tenant,
	}

	result, err := w.services.PostgresRepositories.FlowExecutionRepository.FindRecord(ctx, searchParams)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return result, nil
}

func (w *workflowService) SaveFlowExecutionRecord(ctx context.Context, flowExecutionRecord entity.FlowExecution) (*entity.FlowExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.SaveFlowExecutionRecord")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	if flowExecutionRecord.Tenant == "" || flowExecutionRecord.FlowID == "" || flowExecutionRecord.EntityID == "" {
		err := errors.New("Tenant, FlowID, or EntityID missing")
		span.LogFields(log.Object("executionRecord", flowExecutionRecord))
		tracing.TraceErr(span, err)
		return nil, err
	}

	if flowExecutionRecord.ID == "" {
		id, err := w.services.PostgresRepositories.FlowExecutionRepository.Create(ctx, flowExecutionRecord)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		flowExecutionRecord.ID = id
		return &flowExecutionRecord, nil
	}

	updatedRecord, err := w.services.PostgresRepositories.FlowExecutionRepository.Update(ctx, flowExecutionRecord)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return updatedRecord, nil

}

// Flow Action Execution

func (w *workflowService) SaveFlowActionExecutionRecord(ctx context.Context, flowActionExecutionRecord entity.FlowActionExecution) (*entity.FlowActionExecution, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WorkflowService.SaveFlowExecutionActionRecord")
	defer span.Finish()
	tracing.TagComponentPostgresRepository(span)

	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("Tenant not set in context")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if flowActionExecutionRecord.ID == "" {
		if flowActionExecutionRecord.Action == "" || flowActionExecutionRecord.FlowExecutionID == "" {
			span.LogFields(log.Object("flowActionExecutionRecord", flowActionExecutionRecord))
			err := errors.New("Action or FlowExecutionID missing")
			tracing.TraceErr(span, err)
			return nil, err
		}
		return w.services.PostgresRepositories.FlowActionExecutionRepository.Create(ctx, flowActionExecutionRecord)
	}

	return w.services.PostgresRepositories.FlowActionExecutionRepository.Update(ctx, flowActionExecutionRecord)

}
