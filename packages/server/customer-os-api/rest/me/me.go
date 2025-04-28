package me

import (
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/gin-gonic/gin"

	rest_handlers "github.com/customeros/customeros/packages/server/customer-os-api/rest"
)

type MeResponse struct {
	Tenant string `json:"tenant"`
}

func AuthorizeMe(h *rest_handlers.RestHandlers) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "AuthorizeMe")
		defer spans.Finish()

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.Response.HandleError(c, http.StatusNotFound, nil)
			return
		}

		h.Response.HandleSuccess(c, MeResponse{
			Tenant: tenant,
		})
		return
	}
}
