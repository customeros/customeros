package integrations

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (h *IntegrationHandler) PostmarkDMARCMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Flows.PostmarkDMARCMonitor", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		// Validate Postmark User-Agent
		if c.Request.UserAgent() == "" || !strings.EqualFold(c.Request.UserAgent(), "Postmark") {
			tracing.TraceErr(span, fmt.Errorf("invalid user agent %s", c.Request.UserAgent()))
			h.responseHandler.HandleError(c, http.StatusForbidden, nil)
			return
		}

		// Get raw request body
		rawBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to read request body"))
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		// Process DMARC report via mailstack service
		err = h.services.CommonServices.MailstackService.ProcessDMARCMonitoringReport(ctx, rawBody)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to process DMARC report"))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		// Return accepted response
		h.responseHandler.HandleAccepted(c)
	}
}
