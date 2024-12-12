package flows

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type CreateFlowEdgeRequest struct {
}

type CreateFlowEdgeRecord struct {
}

type CreateFlowEdgeResponse struct {
	rest.BaseResponse
	Edge CreateFlowEdgeRecord `json:"edge"`
}

func CreateFlowEdge(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.CreateFlowEdge", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrInvalidAPIKey)
			return
		}

		// validate tenant owns the flow specified in the path
		flowId := c.Param("flowId")
		isValidFlow := ValidateFlowIsTenant(ctx, s, flowId)
		if !isValidFlow {
			err := errors.New("Flow does not belong to tenant")
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusNotFound, rest.ErrNotFound.WithMessage("unable to locate flolw"))
			return
		}

		request, err := getEdgeCreateRequest(c)
		if err != nil {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("unable to parse payload"))
			return
		}

		// create flow edge in database
		// do this via LinkFlowNodes in workflow service

		c.JSON(http.StatusOK, response)
	}
}
