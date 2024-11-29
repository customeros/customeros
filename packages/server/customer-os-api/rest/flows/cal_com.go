package flows

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func CalDotCom(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := commontracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "CalDotCom", c.Request.Header)
		defer span.Finish()
		commontracing.TagComponentRest(span)

		if !strings.HasPrefix(c.ContentType(), "application/json") {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrUnsupportedContentType)
			return
		}

		// Read the raw body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Unable to read message body"))
			return
		}
		// Important: Restore the body for later use
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Get the signature from header
		signature := c.GetHeader("X-Cal-Signature-256")
		if signature == "" {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Missing signature header"))
			return
		}

		// determine tenant, lookup api key
		secretKey := "CAL_WEBHOOK_SECRET"

		valid, err := VerifyCalWebhookSignature(body, signature, secretKey)
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Unable to verify message payloar"))
			return
		}

		if !valid {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrUnauthorized)
			return
		}

		httpContext := rest.HTTPContext{
			GinContext:     c,
			ServiceContext: &ctx,
			Span:           span,
			Services:       services,
		}

		handleCalDotComEvent(httpContext)
	}
}

func VerifyCalWebhookSignature(payload []byte, signature string, secretKey string) (bool, error) {
	if signature == "" || secretKey == "" {
		return false, fmt.Errorf("signature or secret key is empty")
	}

	h := hmac.New(sha256.New, []byte(secretKey))
	_, err := h.Write(payload)
	if err != nil {
		return false, fmt.Errorf("failed to write payload to HMAC: %w", err)
	}

	calculatedHash := hex.EncodeToString(h.Sum(nil))
	expectedSignature := "sha256=" + calculatedHash
	return hmac.Equal([]byte(expectedSignature), []byte(signature)), nil
}

func handleCalDotComEvent(ctx rest.HTTPContext) {
	var webhook CalDotComPayload
	err := ctx.GinContext.BindJSON(&webhook)
	if err != nil {
		rest.SendError(ctx.GinContext, ctx.Span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Unable to parse cal.com payload"))
		tracing.TraceErr(ctx.Span, errors.Wrap(err, "Unable to parse cal.com payload"))
		return

	}

	switch webhook.TriggerEvent {
	case "BOOKING_CREATED":
		err := processBookingCreatedEvent(ctx, webhook)
		// create contacts
	case "BOOKING_RESCHEDULED":
		// todo -- Handle reschedule
	case "BOOKING_CANCELLED":
		// todo -- Handle cancellation
	}

	return
}

func processBookingCreatedEvent(ctx rest.HTTPContext, payload CalDotComPayload) error {
	return nil
}
