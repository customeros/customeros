package me

import (
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"

	rest_handlers "github.com/customeros/customeros/packages/server/customer-os-api/rest"
)

type MeResponse struct {
	Tenant string `json:"tenant"`
}

func AuthorizeMe(h *rest_handlers.RestHandlers) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "AuthorizeMe", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.Response.HandleError(c, http.StatusNotFound, nil)
		}

		h.Response.HandleSuccess(c, MeResponse{
			Tenant: tenant,
		})
		return
	}
}
