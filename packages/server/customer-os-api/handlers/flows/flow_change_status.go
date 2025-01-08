package flows

import (
	"net/http"

	"github.com/gin-gonic/gin"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func TurnOn(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.TurnOn", c.Request.Header)
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

		// update flow status in database
		result, err := s.Repositories.PostgresRepositories.FlowRepository.Update(ctx, entity.Flow{
			ID:     flowId,
			Status: commonEnum.FlowStatusOn.String(),
		})
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("could not create flow"))
			return
		}

		c.JSON(http.StatusOK, FlowResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Flow: FlowRecord{
				ID:            result.ID,
				Name:          result.Name,
				Description:   *result.Description,
				Trigger:       result.TriggerOn,
				TriggerNodeID: result.TriggerNodeID,
				VisibleUI:     result.VisibleInUI,
				Status:        result.Status,
				CreatedAt:     result.CreatedAt,
			},
		})
	}
}

func TurnOff(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.TurnOff", c.Request.Header)
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

		// update flow status in database
		result, err := s.Repositories.PostgresRepositories.FlowRepository.Update(ctx, entity.Flow{
			ID:     flowId,
			Status: commonEnum.FlowStatusOff.String(),
		})
		if err != nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("could not create flow"))
			return
		}

		c.JSON(http.StatusOK, FlowResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Flow: FlowRecord{
				ID:            result.ID,
				Name:          result.Name,
				Description:   *result.Description,
				Trigger:       result.TriggerOn,
				TriggerNodeID: result.TriggerNodeID,
				VisibleUI:     result.VisibleInUI,
				Status:        result.Status,
				CreatedAt:     result.CreatedAt,
			},
		})
	}
}
