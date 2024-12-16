package flows

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func CreateFlow(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.CreateFlow", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		request, err := getFlowRequestPayload(c)
		if err != nil {
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("unable to parse payload"))
			return
		}

		// validate trigger event
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

		// create flow in database
		flowRecord := buildFlowEntity(ctx, request)
		result, err := s.Repositories.PostgresRepositories.FlowRepository.Create(ctx, flowRecord)
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("could not create flow"))
			return
		}

		// create first node in database
		triggerNode := entity.FlowNode{
			FlowID: result.ID,
			Type:   commonEnum.NodeFlowListenerEvent.String(),
			Event:  &result.TriggerOn,
		}
		nodeResult, err := s.Repositories.PostgresRepositories.FlowNodeRepository.Create(ctx, triggerNode)
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("could not create trigger node"))
		}

		// update Flow with trigger node ID
		result.TriggerNodeID = nodeResult.ID
		result, err = s.Repositories.PostgresRepositories.FlowRepository.Update(ctx, *result)
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("could not set trigger nodeID"))
			return
		}

		c.JSON(http.StatusOK, FlowResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Flow: FlowRecord{
				ID:            result.ID,
				Name:          result.Name,
				Description:   *result.Description,
				Trigger:       result.TriggerOn,
				TriggerNodeID: nodeResult.ID,
				VisibleUI:     result.VisibleInUI,
				Status:        result.Status,
				CreatedAt:     result.CreatedAt,
			},
		})
	}
}

func buildFlowEntity(ctx context.Context, request FlowRecord) entity.Flow {
	return entity.Flow{
		Tenant:        common.GetTenantFromContext(ctx),
		Name:          request.Name,
		Description:   &request.Description,
		TriggerOn:     request.Trigger,
		TriggerNodeID: request.TriggerNodeID,
		VisibleInUI:   request.VisibleUI,
		Status:        commonEnum.FlowStatusInactive.String(),
	}
}

func getFlowRequestPayload(c *gin.Context) (FlowRecord, error) {
	span, _ := tracing.StartTracerSpan(c.Request.Context(), "Flows.getFlowRequestPayload")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var req FlowRecord
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return req, err
	}

	// make Visible = fasle if not provided
	visibleDefault := false
	if req.VisibleUI == nil {
		req.VisibleUI = &visibleDefault
	}

	return req, nil
}
