package flows

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func CreateFlowNode(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.CreateFlowNode", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := handlers.ValidateTenant(c, ctx, span)
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		flowId, belongsToTenant, err := validateFlowBelongsToTenant(c, s)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}
		if !belongsToTenant {
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("unable to locate flow"))
			return
		}

		request, err := getNodeCreateRequest(c, s)
		if err != nil {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("unable to parse payload"))
			return
		}

		// create flow node in database
		flowNode, err := createFlowNodeRecord(ctx, request, flowId)
		if err != nil {
			tracing.TraceErr(span, err)
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}

		newFlowNode, err := s.Repositories.PostgresRepositories.FlowNodeRepository.Create(ctx, flowNode)
		if err != nil {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest)
			return
		}

		response := buildFlowNodeCreateResponse(ctx, newFlowNode)

		c.JSON(http.StatusOK, response)
	}
}

func buildFlowNodeCreateResponse(ctx context.Context, flowNode *entity.FlowNode) FlowNodeResponse {
	span, ctx := tracing.StartTracerSpan(ctx, "Flows.buildFlowNodeCreateResponse")
	defer span.Finish()
	tracing.TagComponentRest(span)

	payload := FlowNodeRecord{
		ID:        flowNode.ID,
		FlowID:    flowNode.FlowID,
		Type:      flowNode.Type,
		Event:     flowNode.Event,
		CreatedAt: flowNode.CreatedAt,
	}

	response := FlowNodeResponse{
		BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
		Node:         payload,
	}

	if flowNode.EventData == nil {
		return response
	}

	data, err := utils.JSONBToAny(*flowNode.EventData)
	if err != nil {
		tracing.TraceErr(span, err)
		return response
	}

	response.Node.EventData = &data
	return response
}

func createFlowNodeRecord(ctx context.Context, request CreateFlowNodeRequest, flowId string) (entity.FlowNode, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "Flows.createFlowNodeRecord")
	defer span.Finish()
	tracing.TagComponentRest(span)

	flowNode := entity.FlowNode{
		FlowID: flowId,
		Type:   request.Type,
		Event:  request.Event,
		Status: commonEnum.FlowNodeEdgeStatusActive.String(),
	}

	if request.EventData == nil {
		return flowNode, nil
	}

	data, err := utils.AnyToJSONB(request.EventData)
	if err != nil {
		return flowNode, err
	}

	flowNode.EventData = &data
	return flowNode, nil
}

func getNodeCreateRequest(c *gin.Context, s *service.Services) (CreateFlowNodeRequest, error) {
	span, ctx := tracing.StartTracerSpan(c.Request.Context(), "Flows.getNodeCreateRequest")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var req CreateFlowNodeRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return req, err
	}

	nodeOk, nodeType := s.CommonServices.WorkflowService.ValidateNodeType(ctx, req.Type)
	if !nodeOk {
		err = errors.New("node type is invalid")
		tracing.TraceErr(span, err)
		return req, err
	}

	eventOk := s.CommonServices.WorkflowService.ValidateEventType(ctx, *nodeType, *req.Event)
	if !eventOk {
		err = errors.New("node event is invalid")
		tracing.TraceErr(span, err)
		return req, err
	}

	return req, nil
}
