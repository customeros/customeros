package flows

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type FlowAgentResponse struct {
	enum.BaseResponse
	Agents []FlowAgentRecord `json:"actions"`
}

type FlowAgentSingleResponse struct {
	enum.BaseResponse
	Agent FlowAgentRecord `json:"action"`
}

type FlowAgentRecord struct {
	Agent       string `json:"action"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetAgents(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.Agents", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		actions, err := s.Repositories.PostgresRepositories.FlowAgentRegistryRepository.FindAll(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}

		results := make([]FlowAgentRecord, len(*actions))

		for i, action := range *actions {
			record := FlowAgentRecord{
				Agent:       action.Agent,
				Name:        action.FriendlyName,
				Description: action.Description,
			}
			results[i] = record
		}

		if len(*actions) == 1 {
			c.JSON(http.StatusOK, FlowAgentSingleResponse{
				BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
				Agent:        results[0],
			})
			return
		}

		c.JSON(http.StatusOK, FlowAgentResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Agents:       results,
		})
	}
}
