package integrations

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
)

func (h *IntegrationHandler) CalDotCom(c *gin.Context, tenant string) {
	spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "IntegrationHandler.CalDotCom")
	defer spans.Finish()

	if !strings.HasPrefix(c.ContentType(), "application/json") {
		h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
		return
	}

	// Read the raw body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		message := "Unable to read message body"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		return
	}
	// Important: Restore the body for later use
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// Get the signature from header
	signature := c.GetHeader("X-Cal-Signature-256")
	if signature == "" {
		message := "Missing signature header"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		return
	}

	// update context with tenant, pass this where tenant is needed
	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    tenant,
		AppSource: constants.AppSourceCustomerOsApiRest,
	})

	// lookup secret
	webhook, err := h.services.Repositories.PostgresRepositories.WebhooksRepository.Find(ctx, postgres_entity.Webhooks{WebhookPath: c.Request.URL.Path})
	if err != nil {
		spans.TraceError(err)
		h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
		return
	}

	valid, err := h.verifyCalWebhookSignature(body, signature, webhook.Secret)
	if err != nil {
		spans.TraceError(err)
		message := "Unable to verify message payload"
		h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
		return
	}

	if !valid {
		h.responseHandler.HandleError(c, http.StatusUnauthorized, nil)
		return
	}

	h.handleCalDotComEvent(c)
}

func (h *IntegrationHandler) verifyCalWebhookSignature(payload []byte, signature string, secretKey string) (bool, error) {
	if signature == "" || secretKey == "" {
		return false, fmt.Errorf("signature or secret key is empty")
	}

	hm := hmac.New(sha256.New, []byte(secretKey))
	_, err := hm.Write(payload)
	if err != nil {
		return false, fmt.Errorf("failed to write payload to HMAC: %w", err)
	}

	calculatedHash := hex.EncodeToString(hm.Sum(nil))
	expectedSignature := "sha256=" + calculatedHash
	return hmac.Equal([]byte(expectedSignature), []byte(signature)), nil
}

func (h *IntegrationHandler) handleCalDotComEvent(c *gin.Context) {
	spans, _ := telemetry.StartRestSpan(c.Request.Context(), "IntegrationHandler.handleCalDotComEvent")
	defer spans.Finish()

	var webhook CalDotComPayload
	err := c.BindJSON(&webhook)
	if err != nil {
		message := "Unable to parse cal.com payload"
		h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
		spans.TraceError(errors.Wrap(err, "Unable to parse cal.com payload"))
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
}
