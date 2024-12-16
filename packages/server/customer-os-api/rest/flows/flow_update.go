package flows

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func UpdateFlow(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.UpdateFlow", c.Request.Header)
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

		request, err := getFlowRequestPayload(c)
		if err != nil {
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("unable to parse payload"))
			return
		}

		// validate trigger event if present
		if request.Trigger != "" {
			validTrigger, err := s.CommonServices.WorkflowService.ValidateListener(ctx, commonEnum.FlowListenerEvent(request.Trigger))
			if err != nil {
				tracing.TraceErr(span, err)
				rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("unable to validate flow trigger event"))
				return
			}

			if !validTrigger {
				err := errors.New("Invalid trigger")
				tracing.TraceErr(span, err)
				rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("invalid trigger event"))
				return
			}

		}

		// update flow in database
		flowRecord := buildFlowEntity(ctx, request)
		flowRecord.ID = flowId
		result, err := s.Repositories.PostgresRepositories.FlowRepository.Update(ctx, flowRecord)
		if err != nil {
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest)
			return
		}
		if result == nil {
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("unable to update flow"))
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
