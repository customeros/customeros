package flows

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type CreateFlowRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Trigger     string `json:"triggerOn" binding:"required"`
	VisibleUI   *bool  `json:"visible"`
}

type CreateFlowRecord struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Trigger     string    `json:"triggerOn,omitempty"`
	VisibleUI   bool      `json:"visible"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"UpadatedAt,omitempty"`
}

type CreateFlowResponse struct {
	rest.BaseResponse
	Flow CreateFlowRecord `json:"flow"`
}

func CreateFlow(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.CreateFlow", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrInvalidAPIKey)
			return
		}

		request, err := getFlowCreateRequest(c)
		if err != nil {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("unable to parse payload"))
			return
		}

		// validate trigger event
		validTrigger, err := s.CommonServices.WorkflowService.ValidateListener(ctx, enum.FlowListenerEvent(request.Trigger))
		if err != nil {
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("unable to validate flow trigger event"))
			return
		}

		if !validTrigger {
			err := errors.New("Invalid trigger")
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("invalid trigger event"))
			return
		}

		// create flow in database
		flowRecord := createFlowRecord(ctx, request)
		flowRecord, err = s.Repositories.PostgresRepositories.FlowRepository.CreateFlow(ctx, flowRecord)
		if err != nil {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest)
			return
		}

		c.JSON(http.StatusOK, CreateFlowResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Flow: CreateFlowRecord{
				ID:          flowRecord.ID,
				Name:        flowRecord.Name,
				Description: *flowRecord.Description,
				Trigger:     flowRecord.TriggerOn,
				VisibleUI:   *flowRecord.VisibleInUI,
				Status:      flowRecord.Status,
				CreatedAt:   flowRecord.CreatedAt,
			},
		})
	}
}

func createFlowRecord(ctx context.Context, request CreateFlowRequest) entity.Flow {
	return entity.Flow{
		Tenant:      common.GetTenantFromContext(ctx),
		Name:        request.Name,
		Description: &request.Description,
		TriggerOn:   request.Trigger,
		VisibleInUI: request.VisibleUI,
		Status:      enum.FlowStatusInactive.String(),
	}
}

func getFlowCreateRequest(c *gin.Context) (CreateFlowRequest, error) {
	span, _ := tracing.StartTracerSpan(c.Request.Context(), "Flows.getPostPayload")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var req CreateFlowRequest
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
