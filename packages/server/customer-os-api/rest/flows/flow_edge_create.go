package flows

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type CreateFlowEdgeRequest struct {
	FromNodeID string  `json:"fromNodeId"`
	ToNodeID   string  `json:"toNodeId"`
	Condition  *string `json:"condition"`
	Data       *any    `json:"data"`
}

type CreateFlowEdgeRecord struct {
	ID         string    `json:"id"`
	FlowID     string    `json:"flowId"`
	FromNodeID string    `json:"fromNodeId"`
	ToNodeID   string    `json:"toNodeId"`
	Condition  *string   `json:"condition,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	Data       *any      `json:"data,omitempty"`
}

type CreateFlowEdgeResponse struct {
	enum.BaseResponse
	Edge CreateFlowEdgeRecord `json:"edge"`
}

func CreateFlowEdge(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.CreateFlowEdge", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		// validate tenant owns the flow specified in the path
		flowId := c.Param("flowId")
		isValidFlow, err := s.CommonServices.WorkflowService.ValidateFlowBelongsToTenant(ctx, flowId)
		if err != nil {
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
		}
		if !isValidFlow {
			err := errors.New("Flow does not belong to tenant")
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("unable to locate flolw"))
			return
		}

		request, err := getEdgeCreateRequest(c, s)
		if err != nil {
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("unable to parse payload"))
			return
		}

		flowEdge, err := createFlowEdge(c, s, request, flowId)

		respPayload := buildCreateEdgeResponse(c, flowEdge)

		c.JSON(http.StatusOK, respPayload)
	}
}

func buildCreateEdgeResponse(c *gin.Context, flowEdge *entity.FlowEdge) CreateFlowEdgeResponse {
	response := CreateFlowEdgeRecord{
		ID:         flowEdge.ID,
		FlowID:     flowEdge.FlowID,
		FromNodeID: flowEdge.FromNodeID,
		ToNodeID:   flowEdge.ToNodeID,
		Condition:  flowEdge.Condition,
		CreatedAt:  flowEdge.CreatedAt,
	}

	if flowEdge.Data != nil {
		data, err := utils.JSONBToAny(*flowEdge.Data)
		if err == nil {
			response.Data = &data
		}
	}

	return CreateFlowEdgeResponse{
		BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
		Edge:         response,
	}
}

func createFlowEdge(c *gin.Context, s *service.Services, request CreateFlowEdgeRequest, flowId string) (*entity.FlowEdge, error) {
	span, ctx := tracing.StartTracerSpan(c.Request.Context(), "Flows.createFlowEdge")
	defer span.Finish()
	tracing.TagComponentRest(span)

	flowEdgeRecord := entity.FlowEdge{
		FlowID:     flowId,
		FromNodeID: request.FromNodeID,
		ToNodeID:   request.ToNodeID,
		Condition:  request.Condition,
	}

	if request.Data != nil {
		data, err := utils.AnyToJSONB(request.Data)
		if err != nil {
			tracing.TraceErr(span, err)
			return &flowEdgeRecord, err
		}
		flowEdgeRecord.Data = &data
	}

	return s.Repositories.PostgresRepositories.FlowEdgeRepository.Create(ctx, flowEdgeRecord)
}

func getEdgeCreateRequest(c *gin.Context, s *service.Services) (CreateFlowEdgeRequest, error) {
	span, ctx := tracing.StartTracerSpan(c.Request.Context(), "Flows.getEdgeCreateRequest")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var req CreateFlowEdgeRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return req, err
	}

	// validate transition
	ok, err := s.CommonServices.WorkflowService.ValidateTransition(ctx, req.FromNodeID, req.ToNodeID)
	if err != nil {
		tracing.TraceErr(span, err)
		return req, err
	}

	if !ok {
		err = errors.New("not a valid edge between from and to nodes")
		tracing.TraceErr(span, err)
		return req, err
	}

	return req, nil
}
