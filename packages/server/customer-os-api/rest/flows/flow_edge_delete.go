package flows

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func DeleteFlowEdge(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.DeleteFlowEdge", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		// validate tenant owns the flow specified in the path
		flowId, belongsToTenant, err := validateFlowBelongsToTenant(c, s)
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}
		if !belongsToTenant {
			rest.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("unable to locate flow"))
			return
		}

		edgeId := c.Param("edgeId")
		edgePrefix := strings.HasPrefix(strings.ToLower(edgeId), "edge_")
		if !edgePrefix {
			edgeId = fmt.Sprintf("edge_%s", edgeId)
		}

		switch edgeId {
		case "":
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("missing edgeID"))
			return
		default:
			query := entity.FlowEdge{
				ID:     edgeId,
				FlowID: flowId,
				Active: false,
			}
			deletedEdge, err := s.Repositories.PostgresRepositories.FlowEdgeRepository.Update(ctx, query)
			if err != nil {
				rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
				return
			}
			if deletedEdge.Active != false {
				rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("edge deletion failed"))
				return
			}
			c.JSON(http.StatusNoContent, enum.BuildBaseResponse(enum.StatusSuccess))
		}
	}
}
