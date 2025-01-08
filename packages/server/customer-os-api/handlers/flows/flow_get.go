package flows

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func GetFlows(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.GetFlows", c.Request.Header)
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
		if flowId != "" && !belongsToTenant {
			handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("unable to locate flow"))
			return
		}

		query := entity.Flow{
			ID:     c.Param("flowId"),
			Tenant: tenant,
		}

		switch query.ID {
		case "":
			getAllFlows(c, s, query)
		default:
			getFlow(c, s, query)
		}
	}
}

func getFlow(c *gin.Context, s *service.Services, query entity.Flow) {
	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "Flows.getFlow")
	defer span.Finish()
	tracing.TagComponentRest(span)

	flowRecord, err := s.Repositories.PostgresRepositories.FlowRepository.Find(ctx, query)
	if err != nil {
		handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest)
		return
	}
	if flowRecord == nil {
		handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("flow does not exist"))
		return
	}

	c.JSON(http.StatusOK, FlowResponse{
		BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
		Flow:         buildFlowRecord(ctx, *flowRecord),
	})
}

func getAllFlows(c *gin.Context, s *service.Services, query entity.Flow) {
	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "Flows.getAllFlows")
	defer span.Finish()
	tracing.TagComponentRest(span)

	flowRecords, err := s.Repositories.PostgresRepositories.FlowRepository.FindAll(ctx, query)
	if err != nil {
		handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest)
		return
	}

	response := make([]FlowRecord, len(*flowRecords))
	for i, record := range *flowRecords {
		response[i] = buildFlowRecord(ctx, record)
	}

	c.JSON(http.StatusOK, FlowsResponse{
		BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
		Flows:        response,
	})
}

func buildFlowRecord(ctx context.Context, record entity.Flow) FlowRecord {
	return FlowRecord{
		ID:            record.ID,
		Name:          record.Name,
		Description:   *record.Description,
		Trigger:       record.TriggerOn,
		TriggerNodeID: record.TriggerNodeID,
		VisibleUI:     record.VisibleInUI,
		Status:        record.Status,
		CreatedAt:     record.CreatedAt,
		UpdatedAt:     record.UpdatedAt,
	}
}
