package flows

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type FlowListenerEventsResponse struct {
	enum.BaseResponse
	Events []FlowListenerEventRecord `json:"events"`
}

type FlowListenerEventSingleResponse struct {
	enum.BaseResponse
	Event FlowListenerEventRecord `json:"event"`
}

type FlowListenerEventRecord struct {
	System      string `json:"system"`
	Event       string `json:"event"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetListeners(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.GetListeners", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := handlers.ValidateTenant(c, ctx, span)
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		events, err := s.Repositories.PostgresRepositories.FlowListenerRegistryRepository.FindAll(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}

		results := make([]FlowListenerEventRecord, len(*events))

		for i, event := range *events {
			record := FlowListenerEventRecord{
				System:      event.ExternalSystem,
				Event:       event.ListenerEvent,
				Name:        event.FriendlyName,
				Description: event.Description,
			}
			results[i] = record
		}

		if len(*events) == 1 {
			c.JSON(http.StatusOK, FlowListenerEventSingleResponse{
				BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
				Event:        results[0],
			})
			return
		}

		c.JSON(http.StatusOK, FlowListenerEventsResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Events:       results,
		})
	}
}
