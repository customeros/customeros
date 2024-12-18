package flows

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

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
		flowId, belongsToTenant, err := validateFlowBelongsToTenant(c, s)
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}
		if !belongsToTenant {
			rest.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("unable to locate flow"))
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

func buildCreateEdgeResponse(c *gin.Context, flowEdge *entity.FlowEdge) FlowEdgeResponse {
	response := FlowEdgeRecord{
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

	return FlowEdgeResponse{
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
		Status:     commonEnum.FlowNodeEdgeStatusActive.String(),
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
