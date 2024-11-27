package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type MeResponse struct {
	BaseResponse
	Tenant string `json:"tenant"`
}

func AuthorizeMe(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "CreateContact", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := ValidateTenant(c, ctx, span)
		if tenant == "" {
			return
		}

		c.JSON(http.StatusOK, MeResponse{
			BaseResponse: BuildBaseResponse(StatusSuccess),
			Tenant:       tenant,
		})
		return
	}
}
