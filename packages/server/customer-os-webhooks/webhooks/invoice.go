package webhook

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
)

type InvoiceHandler struct {
	BaseHandler
}

func NewInvoiceHandler(services *service.Services, logger logger.Logger) *InvoiceHandler {
	return &InvoiceHandler{
		BaseHandler: BaseHandler{
			services: services,
			logger:   logger,
		},
	}
}

func (h *InvoiceHandler) GetPath() string {
	return "/sync/invoice"
}

func (h *InvoiceHandler) Handle(c *gin.Context) {
	// Start tracing
	ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncInvoice", c.Request.Header)
	defer span.Finish()

	// Read and validate tenant header
	tenant := c.GetHeader("tenant")
	if tenant == "" {
		h.handleError(c, span, http.StatusBadRequest, "Missing or empty tenant header", nil)
		return
	}

	// Set context with tenant
	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    tenant,
		AppSource: constants.AppSourceCustomerOsWebhooks,
	})

	// Limit request body size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeCommon)
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		h.logger.Errorf("(SyncInvoice) error reading request body: %s", err.Error())
		h.handleError(c, span, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Parse JSON request body
	var invoice model.InvoiceData
	if err = json.Unmarshal(requestBody, &invoice); err != nil {
		tracing.TraceErr(span, err)
		h.logger.Errorf("(SyncInvoice) Failed unmarshalling body request: %s", err.Error())
		h.handleError(c, span, http.StatusUnprocessableEntity, "Cannot unmarshal request body", err)
		return
	}

	// Set timeout context
	ctx, cancel := context.WithTimeout(ctx, constants.Duration1Min)
	defer cancel()

	// Process invoice
	syncResult, err := h.services.InvoiceService.SyncInvoices(ctx, []model.InvoiceData{invoice})
	if err != nil {
		tracing.TraceErr(span, err)
		h.logger.Errorf("(SyncInvoice) error in sync invoice: %s", err.Error())
		status := http.StatusInternalServerError
		message := "Failed processing invoice"
		if errors.IsBadRequest(err) {
			status = http.StatusBadRequest
			message = err.Error()
		}
		h.handleError(c, span, status, message, err)
		return
	}

	c.JSON(http.StatusOK, syncResult)
}
