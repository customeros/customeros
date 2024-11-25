package events

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func Fathom(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := commontracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Fathom", c.Request.Header)
		defer span.Finish()
		commontracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			return
		}

		if !strings.HasPrefix(c.ContentType(), "application/json") {
			rest.SendError(c, http.StatusBadRequest, "Unsupported Content-Type")
		}

		httpContext := rest.HTTPContext{
			GinContext:     c,
			ServiceContext: &ctx,
			Span:           span,
			Services:       services,
			Tenant:         tenant,
		}

		if c.Request.UserAgent() == "" {
			rest.SendError(c, http.StatusForbidden, "User-Agent header is required")
		}

		if !strings.EqualFold(c.Request.UserAgent(), "Zapier") {
			rest.SendError(c, http.StatusForbidden, "User-Agent not authorized")
		}

		handleFathomAISummaryZapier(httpContext)
	}
}

func handleFathomAISummaryZapier(ctx rest.HTTPContext) {
	var aiSummaryData FathomAISummaryZapier
	if err := ctx.GinContext.BindJSON(&aiSummaryData); err == nil && aiSummaryData.AISummaryHTMLFormatted != "" {
		rest.SendSuccess(ctx.GinContext, http.StatusAccepted, "Processing Fathom AI Summary")

		go func() {
			if err := processFathomAISummaryZapier(ctx, aiSummaryData); err != nil {
				tracing.TraceErr(ctx.Span, errors.Wrap(err, "failed to process Fathom AI summary from zapier"))
			}
		}()
		return
	}
	return
}

func processFathomAISummaryZapier(ctx rest.HTTPContext, aiSummaryData FathomAISummaryZapier) error {
	// extract AI summary
	// format as a log entry?
	// write to database
	return nil
}
