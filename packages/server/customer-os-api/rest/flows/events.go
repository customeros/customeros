package flows

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type FlowEventsResponse struct {
	rest.BaseResponse
	Events []FlowEventRecord `json:"events"`
}

type FlowEventSingleResponse struct {
	rest.BaseResponse
	Event FlowEventRecord `json:"event"`
}

type FlowEventRecord struct {
	System      string `json:"system"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetEvents(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.GetEvents", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrInvalidAPIKey)
			return
		}

		events, err := s.Repositories.PostgresRepositories.FlowEventsRepository.GetAllFlowEvents(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer)
			return
		}

		results := make([]FlowEventRecord, len(events))

		for i, event := range events {
			record := FlowEventRecord{
				System:      event.ExternalSystem,
				Resource:    event.Resource,
				Action:      event.Action,
				Name:        event.EventName,
				Description: event.Description,
			}
			results[i] = record
		}

		if len(events) == 1 {
			c.JSON(http.StatusOK, FlowEventSingleResponse{
				BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
				Event:        results[0],
			})
			return
		}

		c.JSON(http.StatusOK, FlowEventsResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Events:       results,
		})
	}
}
