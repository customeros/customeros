package flows

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type FlowActionResponse struct {
	enum.BaseResponse
	Actions []FlowActionRecord `json:"actions"`
}

type FlowActionSingleResponse struct {
	enum.BaseResponse
	Action FlowActionRecord `json:"action"`
}

type FlowActionRecord struct {
	Action      string `json:"action"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetActions(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.Actions", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		actions, err := s.Repositories.PostgresRepositories.FlowActionRegistryRepository.FindAll(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}

		results := make([]FlowActionRecord, len(*actions))

		for i, action := range *actions {
			record := FlowActionRecord{
				Action:      action.Action,
				Name:        action.FriendlyName,
				Description: action.Description,
			}
			results[i] = record
		}

		if len(*actions) == 1 {
			c.JSON(http.StatusOK, FlowActionSingleResponse{
				BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
				Action:       results[0],
			})
			return
		}

		c.JSON(http.StatusOK, FlowActionResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Actions:      results,
		})
	}
}
