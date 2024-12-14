package flows

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func DeleteFlowNode(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.DeleteFlowNode", c.Request.Header)
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

		nodeId := c.Param("nodeId")

		switch nodeId {
		case "":
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("nodeID not provided"))
			return
		default:
			deleteRecord := entity.FlowNode{
				ID:     nodeId,
				FlowID: flowId,
				Active: false,
			}
			deletedNode, err := s.Repositories.PostgresRepositories.FlowNodeRepository.Update(ctx, deleteRecord)
			if err != nil {
				rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
				return
			}
			if deletedNode.Active != false {
				rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("node deletion failed"))
				return
			}
			c.JSON(http.StatusNoContent, enum.BuildBaseResponse(enum.StatusSuccess))
		}
	}
}
