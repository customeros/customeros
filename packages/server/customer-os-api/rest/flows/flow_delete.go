package flows

import (
	"net/http"

	"github.com/gin-gonic/gin"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func DeleteFlow(s *service.Services) gin.HandlerFunc {
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

		// update flow in database
		flowRecord := entity.Flow{
			ID:     flowId,
			Status: commonEnum.FlowStatusArchived.String(),
		}
		result, err := s.Repositories.PostgresRepositories.FlowRepository.Update(ctx, flowRecord)
		if err != nil {
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest)
			return
		}
		if result == nil {
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("flow deletion failed"))
		}
		c.JSON(http.StatusNoContent, enum.BuildBaseResponse(enum.StatusSuccess))
	}
}
