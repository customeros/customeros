package flows

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type CreateFlowNodeRequest struct {
	Type      string   `json:"type"`
	Event     *string  `json:"event"`
	PositionX *float64 `json:"positionX"`
	PositionY *float64 `json:"positionY"`
	Data      *any     `json:"data"`
}

type CreateFlowNodeRecord struct {
	ID        string    `json:"id"`
	FlowID    string    `json:"flowId"`
	Type      string    `json:"type"`
	Event     *string   `json:"event,omitempty"`
	PositionX *float64  `json:"positionX,omitempty"`
	PositionY *float64  `json:"positionY,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"UpadatedAt,omitempty"`
	Data      *any      `json:"data,omitempty"`
}

type CreateFlowNodeResponse struct {
	rest.BaseResponse
	Node CreateFlowNodeRecord `json:"node"`
}

func CreateFlowNode(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.CreateFlowNode", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrInvalidAPIKey)
			return
		}

		// validate tenant owns the flow specified in the path
		flowId := c.Param("flowId")
		isValidFlow := s.CommonServices.WorkflowService.ValidateFlowBelongsToTenant(ctx, flowId)
		if !isValidFlow {
			err := errors.New("Flow does not belong to tenant")
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusNotFound, rest.ErrNotFound.WithMessage("unable to locate flolw"))
			return
		}

		request, err := getNodeCreateRequest(c, s)
		if err != nil {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("unable to parse payload"))
			return
		}

		// create flow node in database
		flowNode, err := createFlowNodeRecord(ctx, request, flowId)
		if err != nil {
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer)
			return
		}

		flowNode, err = s.Repositories.PostgresRepositories.FlowNodeRepository.CreateFlowNode(ctx, flowNode)
		if err != nil {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest)
			return
		}

		response := buildFlowNodeCreateResponse(ctx, flowNode)

		c.JSON(http.StatusOK, response)
	}
}

func buildFlowNodeCreateResponse(ctx context.Context, flowNode entity.FlowNode) CreateFlowNodeResponse {
	span, ctx := tracing.StartTracerSpan(ctx, "Flows.buildFlowNodeCreateResponse")
	defer span.Finish()
	tracing.TagComponentRest(span)

	payload := CreateFlowNodeRecord{
		ID:        flowNode.ID,
		FlowID:    flowNode.FlowID,
		Type:      flowNode.Type,
		Event:     flowNode.Event,
		PositionX: flowNode.PositionX,
		PositionY: flowNode.PositionY,
		CreatedAt: flowNode.CreatedAt,
	}

	response := CreateFlowNodeResponse{
		BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
		Node:         payload,
	}

	if flowNode.Data == nil {
		return response
	}

	data, err := utils.JSONBToAny(*flowNode.Data)
	if err != nil {
		tracing.TraceErr(span, err)
		return response
	}

	response.Node.Data = &data
	return response
}

func createFlowNodeRecord(ctx context.Context, request CreateFlowNodeRequest, flowId string) (entity.FlowNode, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "Flows.createFlowNodeRecord")
	defer span.Finish()
	tracing.TagComponentRest(span)

	flowNode := entity.FlowNode{
		FlowID:    flowId,
		Type:      request.Type,
		Event:     request.Event,
		PositionX: request.PositionX,
		PositionY: request.PositionY,
	}

	if request.Data == nil {
		return flowNode, nil
	}

	data, err := utils.AnyToJSONB(request.Data)
	if err != nil {
		return flowNode, err
	}

	flowNode.Data = &data
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
