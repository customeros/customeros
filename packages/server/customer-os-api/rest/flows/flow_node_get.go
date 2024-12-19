package flows

import (
	"context"
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

func GetFlowNodes(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.GetFlowNode", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		flowId, belongsToTenant, err := validateFlowBelongsToTenant(c, s)
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}
		if !belongsToTenant {
			rest.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("unable to locate flow"))
			return
		}

		nodeId := c.Param("nodeId")

		if nodeId != "" {
			nodePrefix := strings.HasPrefix(strings.ToLower(nodeId), "node_")
			if !nodePrefix {
				nodeId = fmt.Sprintf("node_%s", nodeId)
			}
		}

		switch nodeId {
		case "":
			getAllFlowNodes(c, s, flowId)
		default:
			getFlowNode(c, s, flowId, nodeId)
		}
	}
}

func getFlowNode(c *gin.Context, s *service.Services, flowId, nodeId string) {
	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "Flows.getFlowNode")
	defer span.Finish()
	tracing.TagComponentRest(span)

	node, err := s.Repositories.PostgresRepositories.FlowNodeRepository.Find(ctx, entity.FlowNode{
		ID:     nodeId,
		FlowID: flowId,
		Status: commonEnum.FlowNodeEdgeStatusActive.String(),
	})
	if err != nil {
		rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
		return
	}

	response := FlowNodeResponse{
		BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
		Node: FlowNodeRecord{
			ID:        node.ID,
			FlowID:    node.FlowID,
			Type:      node.Type,
			Event:     node.Event,
			CreatedAt: node.CreatedAt,
			UpdatedAt: node.UpdatedAt,
		},
	}

	if node.EventData != nil {
		eventData, err := utils.JSONBToAny(*node.EventData)
		if err != nil {
			tracing.TraceErr(span, err)
		}
		response.Node.EventData = &eventData
	}

	c.JSON(http.StatusOK, response)
}

func getAllFlowNodes(c *gin.Context, s *service.Services, flowId string) {
	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "Flows.getAllFlowNodes")
	defer span.Finish()
	tracing.TagComponentRest(span)

	nodes, err := s.Repositories.PostgresRepositories.FlowNodeRepository.FindAll(ctx, entity.FlowNode{
		FlowID: flowId,
		Status: commonEnum.FlowNodeEdgeStatusActive.String(),
	})
	if err != nil {
		rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
		return
	}

	response := FlowNodesResponse{
		BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
		Nodes:        buildFlowNodesResponse(ctx, nodes),
	}
	c.JSON(http.StatusOK, response)
}

func buildFlowNodesResponse(ctx context.Context, nodes *[]entity.FlowNode) []FlowNodeRecord {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Flows.getAllFlowNodes")
	defer span.Finish()
	tracing.TagComponentRest(span)

	records := make([]FlowNodeRecord, len(*nodes))
	for i, node := range *nodes {
		records[i] = FlowNodeRecord{
			ID:        node.ID,
			FlowID:    node.FlowID,
			Type:      node.Type,
			Event:     node.Event,
			CreatedAt: node.CreatedAt,
			UpdatedAt: node.UpdatedAt,
		}

		if node.EventData != nil {
			eventData, err := utils.JSONBToAny(*node.EventData)
			if err != nil {
				tracing.TraceErr(span, err)
			}
			records[i].EventData = &eventData
		}

	}
	return records
}
