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

func UpdateFlowNode(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.UpdateFlowNode", c.Request.Header)
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

		nodeId := c.Param("nodeId")
		nodePrefix := strings.HasPrefix(strings.ToLower(nodeId), "node_")
		if !nodePrefix {
			nodeId = fmt.Sprintf("node_%s", nodeId)
		}

		switch nodeId {
		case "":
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("nodeID not provided"))
			return
		default:
			updateRecord, err := getNodeUpdatePayload(c, s)
			if err != nil {
				tracing.TraceErr(span, err)
				handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("unable to parse requst"))
				return
			}
			updateRecord.ID = nodeId
			updateRecord.FlowID = flowId
			updateFlowNode(c, s, updateRecord)
		}
	}
}

func updateFlowNode(c *gin.Context, s *service.Services, updateRecord FlowNodeRecord) {
	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "Flows.updateFlowNode")
	defer span.Finish()
	tracing.TagComponentRest(span)

	// validate flow exists and is active before updating
	query := entity.FlowNode{
		ID:     updateRecord.ID,
		FlowID: updateRecord.FlowID,
		Status: commonEnum.FlowNodeEdgeStatusActive.String(),
	}

	node, err := s.Repositories.PostgresRepositories.FlowNodeRepository.Find(ctx, query)
	if err != nil {
		handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
		return
	}
	if node == nil {
		handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("node not found"))
		return
	}

	query.Type = updateRecord.Type
	query.Event = updateRecord.Event

	if updateRecord.EventData != nil {
		eventData, err := utils.AnyToJSONB(*updateRecord.EventData)
		if err != nil {
			tracing.TraceErr(span, err)
		}
		query.EventData = &eventData
	}

	updatedNode, err := s.Repositories.PostgresRepositories.FlowNodeRepository.Update(ctx, query)
	if err != nil {
		handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
		return
	}

	response := FlowNodeResponse{
		BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
		Node: FlowNodeRecord{
			ID:        updatedNode.ID,
			FlowID:    updatedNode.FlowID,
			Type:      updatedNode.Type,
			Event:     updatedNode.Event,
			CreatedAt: updatedNode.CreatedAt,
			UpdatedAt: updatedNode.UpdatedAt,
		},
	}

	if updatedNode.EventData != nil {
		eventData, err := utils.JSONBToAny(*updatedNode.EventData)
		if err != nil {
			tracing.TraceErr(span, err)
		}
		response.Node.EventData = &eventData
	}

	c.JSON(http.StatusOK, response)
}

func getNodeUpdatePayload(c *gin.Context, s *service.Services) (FlowNodeRecord, error) {
	span, _ := tracing.StartTracerSpan(c.Request.Context(), "Flows.getNodeUpdatePayload")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var req FlowNodeRecord
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return req, err
	}

	return req, nil
}
