package flows

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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

func CalDotCom(c *gin.Context, s *service.Services) {
	span, ctx := commontracing.StartTracerSpan(c.Request.Context(), "Flows.CalDotCom")
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

	// determine tenant
	tenant, err := s.Repositories.PostgresRepositories.TenantRepository.GetTenant(ctx, c.Param("tenantId"))
	if err != nil {
		tracing.TraceErr(span, err)
		rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Unable to identify tenant"))
		return
	}

	// lookup secret
	webhook, err := s.Repositories.PostgresRepositories.FlowWebhooksRepository.FindWebhookByPath(ctx, tenant, c.Request.URL.Path)
	if err != nil {
		tracing.TraceErr(span, err)
		rest.SendError(c, span, http.StatusInternalServerError, rest.ErrNotFound)
		return
	}

	valid, err := VerifyCalWebhookSignature(body, signature, webhook.Secret)
	if err != nil {
		tracing.TraceErr(span, err)
		rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Unable to verify message payloar"))
		return
	}

	if !valid {
		rest.SendError(c, span, http.StatusUnauthorized, rest.ErrUnauthorized)
		return
	}

	handleCalDotComEvent(c)
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

func handleCalDotComEvent(c *gin.Context) {
	span, _ := commontracing.StartTracerSpan(c.Request.Context(), "Flows.handleCalDotComEvent")
	defer span.Finish()
	commontracing.TagComponentRest(span)

	var webhook CalDotComPayload
	err := c.BindJSON(&webhook)
	if err != nil {
		rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Unable to parse cal.com payload"))
		tracing.TraceErr(span, errors.Wrap(err, "Unable to parse cal.com payload"))
		return

	}

	switch webhook.TriggerEvent {
	case "BOOKING_CREATED":
		// err := processBookingCreatedEvent(ctx, webhook)
		// create contacts
	case "BOOKING_RESCHEDULED":
		// todo -- Handle reschedule
	case "BOOKING_CANCELLED":
		// todo -- Handle cancellation
	}

	return
}

func processBookingCreatedEvent(c *gin.Context, payload *CalDotComPayload) error {
	return nil
}
