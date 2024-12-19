package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type MeResponse struct {
	enum.BaseResponse
	Tenant string `json:"tenant"`
}

func AuthorizeMe(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "AuthorizeMe", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := ValidateTenant(c, ctx, span)

		c.JSON(http.StatusOK, MeResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Tenant:       tenant,
		})
		return
	}
}
