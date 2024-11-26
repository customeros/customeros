package webhooks

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

func PostmarkInboundEmail(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RedirectToPayInvoice", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		body := c.Request.Body
		requestBody, err := io.ReadAll(body)
		if err != nil {
			tracing.LogObjectAsJson(span, "body", body)
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest)
			return
		}
		return
	}
}
