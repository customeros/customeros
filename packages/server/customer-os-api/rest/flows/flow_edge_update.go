package flows

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func UpdateFlowEdge(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.UpdateFlowEdge", c.Request.Header)
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

		edgeId := c.Param("edgeId")
		edgePrefix := strings.HasPrefix(strings.ToLower(edgeId), "edge_")
		if !edgePrefix {
			edgeId = fmt.Sprintf("edge_%s", edgeId)
		}

		switch edgeId {
		case "":
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("missing nodeID"))
			return
		default:
			updateRecord, err := getEdgeRequestPayload(c)
			if err != nil {
				rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("unable to parse payload"))
				return
			}
			updateRecord.ID = edgeId
			updateRecord.FlowID = flowId
			updateFlowEdge(c, s, updateRecord)
		}
	}
}

func updateFlowEdge(c *gin.Context, s *service.Services, updateRecord FlowEdgeRecord) {
	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "Flows.updateFlowEdge")
	defer span.Finish()
	tracing.TagComponentRest(span)

	// validate flow exists and is active before updating
	query := entity.FlowEdge{
		ID:     updateRecord.ID,
		FlowID: updateRecord.FlowID,
		Status: commonEnum.FlowNodeEdgeStatusActive.String(),
	}

	edge, err := s.Repositories.PostgresRepositories.FlowEdgeRepository.Find(ctx, query)
	if err != nil {
		rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
		return
	}
	if edge == nil {
		rest.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("edge not found"))
		return
	}

	query.FromNodeID = updateRecord.FromNodeID
	query.ToNodeID = updateRecord.ToNodeID
	query.Condition = updateRecord.Condition

	if updateRecord.Data != nil {
		eventData, err := utils.AnyToJSONB(*updateRecord.Data)
		if err != nil {
			tracing.TraceErr(span, err)
		}
		query.Data = &eventData
	}

	updatedEdge, err := s.Repositories.PostgresRepositories.FlowEdgeRepository.Update(ctx, query)
	if err != nil {
		rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
		return
	}

	response := FlowEdgeResponse{
		BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
		Edge: FlowEdgeRecord{
			ID:         updatedEdge.ID,
			FlowID:     updatedEdge.FlowID,
			FromNodeID: updatedEdge.FromNodeID,
			ToNodeID:   updatedEdge.ToNodeID,
			Condition:  updatedEdge.Condition,
			CreatedAt:  updatedEdge.CreatedAt,
			UpdatedAt:  updatedEdge.UpdatedAt,
		},
	}

	if updatedEdge.Data != nil {
		data, err := utils.JSONBToAny(*updatedEdge.Data)
		if err != nil {
			tracing.TraceErr(span, err)
		}
		response.Edge.Data = &data
	}

	c.JSON(http.StatusOK, response)
}

func getEdgeRequestPayload(c *gin.Context) (FlowEdgeRecord, error) {
	span, _ := tracing.StartTracerSpan(c.Request.Context(), "Flows.getEdgeRequestPayload")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var req FlowEdgeRecord
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return req, err
	}

	return req, nil
}
