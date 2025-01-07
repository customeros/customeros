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
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func GetFlowEdges(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.GetFlowEdges", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := handlers.ValidateTenant(c, ctx, span)
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		// validate tenant owns the flow specified in the path
		flowId, belongsToTenant, err := validateFlowBelongsToTenant(c, s)
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}
		if !belongsToTenant {
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("unable to locate flow"))
			return
		}

		edgeId := c.Param("edgeId")
		edgePrefix := strings.HasPrefix(strings.ToLower(edgeId), "edge_")
		if !edgePrefix {
			edgeId = fmt.Sprintf("edge_%s", edgeId)
		}

		query := entity.FlowEdge{
			ID:     edgeId,
			FlowID: flowId,
			Status: commonEnum.FlowNodeEdgeStatusActive.String(),
		}

		switch edgeId {
		case "":
			getFlowEdge(c, s, query)
		default:
			getFlowEdges(c, s, query)
		}
	}
}

func getFlowEdge(c *gin.Context, s *service.Services, query entity.FlowEdge) {
	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "Flows.getFlowEdge")
	defer span.Finish()
	tracing.TagComponentRest(span)

	edge, err := s.Repositories.PostgresRepositories.FlowEdgeRepository.Find(ctx, query)
	if err != nil {
		tracing.TraceErr(span, err)
		handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
		return
	}
	if edge == nil {
		handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("edge does not exist"))
		return
	}

	response := buildEdgeResponse(edge)
	c.JSON(http.StatusOK, FlowEdgeResponse{
		BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
		Edge:         response,
	})
	return
}

func getFlowEdges(c *gin.Context, s *service.Services, query entity.FlowEdge) {
	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "Flows.getFlowEdges")
	defer span.Finish()
	tracing.TagComponentRest(span)

	edges, err := s.Repositories.PostgresRepositories.FlowEdgeRepository.FindAll(ctx, query)
	if err != nil {
		tracing.TraceErr(span, err)
		handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
		return
	}
	if len(*edges) == 0 {
		handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("no edges exist for this flow"))
		return
	}

	records := make([]FlowEdgeRecord, len(*edges))
	for i, edge := range *edges {
		records[i] = buildEdgeResponse(&edge)
	}

	c.JSON(http.StatusOK, FlowEdgesResponse{
		BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
		Edges:        records,
	})
}

func buildEdgeResponse(flowEdge *entity.FlowEdge) FlowEdgeRecord {
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

	return response
}
