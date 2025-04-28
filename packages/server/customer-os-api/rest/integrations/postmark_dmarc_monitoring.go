package integrations

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (h *IntegrationHandler) PostmarkDMARCMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "IntegrationHandler.PostmarkDMARCMonitor")
		defer spans.Finish()

		// Validate Postmark User-Agent
		if c.Request.UserAgent() == "" || !strings.EqualFold(c.Request.UserAgent(), "Postmark") {
			spans.TraceError(fmt.Errorf("invalid user agent %s", c.Request.UserAgent()))
			h.responseHandler.HandleError(c, http.StatusForbidden, nil)
			return
		}

		// Get raw request body
		rawBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to read request body"))
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		// Process DMARC report via mailstack service
		err = h.services.CommonServices.MailstackService.ProcessDMARCMonitoringReport(ctx, rawBody)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "failed to process DMARC report"))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		// Return accepted response
		h.responseHandler.HandleAccepted(c)
	}
}
