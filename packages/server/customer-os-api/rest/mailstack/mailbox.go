// @openapi 3.0.0
package mailstack

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// @title MailStack API
// @version 1.0
// @description API for managing mailboxes and domain email configuration
// @BasePath /mailstack/v1

// RegisterNewMailbox creates a new mailbox for a domain
// @Summary Create mailbox
// @Description Creates a new mailbox for the specified domain with optional forwarding and webmail access
// @Tags Mailboxes
// @Accept json
// @Produce json
// @Param domain path string true "Domain name" example(example.com)
// @Param body body MailboxRequest true "Mailbox configuration"
// @Success 200 {object} MailboxResponse "Mailbox created successfully"
// @Success 200 {object} MailboxResponse "Mailbox created successfully with generated password"
// @Failure 400 {object} handlers.ErrorResponse "Invalid request - Missing or invalid parameters"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 404 {object} handlers.ErrorResponse "Domain not found"
// @Failure 409 {object} handlers.ErrorResponse "Mailbox already exists"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /domains/{domain}/mailboxes [post]
// @Security ApiKeyAuth
func (h *MailstackHandler) RegisterNewMailbox() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RegisterNewMailbox", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		// get domain from path
		domain := c.Param("domain")
		if domain == "" {
			message := "Missing parameter: domain"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}
		span.LogKV("request.domain", domain)

		// get tenant from context
		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// Parse and validate request body
		var mailboxRequest MailboxRequest
		if err := c.ShouldBindJSON(&mailboxRequest); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Invalid request body"))
			// log body
			body, _ := c.GetRawData()
			span.LogFields(tracingLog.String("request.body", string(body)))
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		// Create mailbox using service
		result, err := h.services.CommonServices.MailstackService.RegisterMailbox(ctx, tenant, domain, interfaces.CreateMailboxRequest{
			Username:        mailboxRequest.Username,
			Password:        mailboxRequest.Password,
			ForwardingTo:    mailboxRequest.ForwardingTo,
			WebmailEnabled:  mailboxRequest.WebmailEnabled,
			LinkedUserEmail: mailboxRequest.LinkedUser,
		})
		if err != nil {
			message := "Internal server error"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			tracing.TraceErr(span, err)
			return
		}

		// Handle non-201 responses
		if result.StatusCode != http.StatusCreated {
			if result.StatusCode == http.StatusInternalServerError || result.StatusCode == http.StatusUnauthorized {
				message := "Internal server error"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}
			h.responseHandler.HandleError(c, result.StatusCode, &result.ErrorMsg)
			return
		}

		h.responseHandler.HandleSuccess(c, MailboxResponse{
			Mailbox: MailboxRecord{
				Email:             result.Mailbox.Email,
				Password:          result.Mailbox.Password,
				ForwardingTo:      result.Mailbox.ForwardingTo,
				ForwardingEnabled: result.Mailbox.ForwardingEnabled,
				WebmailEnabled:    result.Mailbox.WebmailEnabled,
			},
		})
	}
}

// GetMailboxes retrieves all mailboxes for a domain
// @Summary List mailboxes
// @Description Retrieves all mailboxes configured for the specified domain
// @Tags Mailboxes
// @Accept json
// @Produce json
// @Param domain path string true "Domain name" example(example.com)
// @Success 200 {object} MailboxesResponse "Successfully retrieved mailboxes"
// @Failure 400 {object} handlers.ErrorResponse "Missing domain parameter"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 404 {object} handlers.ErrorResponse "Domain not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /domains/{domain}/mailboxes [get]
// @Security ApiKeyAuth
func (h *MailstackHandler) GetMailboxes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetMailboxes", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		// get domain from path
		domain := c.Param("domain")
		if domain == "" {
			message := "Missing parameter: domain"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}
		span.LogKV("request.domain", domain)

		// get tenant from context
		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// Create request to Mailstack API
		url := h.services.Cfg.Common.Internal.MailstackApiConfig.ApiUrl + "/v1/mailboxes"
		if domain != "" {
			url += "?domain=" + domain
		}
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			message := "Unable to create request to Mailstack API"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		// Add required headers
		req.Header.Set("X-CUSTOMER-OS-API-KEY", h.services.Cfg.Common.Internal.MailstackApiConfig.ApiKey)
		req.Header.Set("tenant", tenant)

		// Forward Jaeger trace context
		carrier := opentracing.HTTPHeadersCarrier(req.Header)
		err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
		if err != nil {
			span.LogFields(tracingLog.Error(err))
		}

		// Create HTTP client with default transport and timeout
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		// Make request to Mailstack API
		resp, err := client.Do(req)
		if err != nil {
			message := "Unable to connect to Mailstack API"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}
		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode != http.StatusOK {
			// Read error response body
			var errorResponse struct {
				Error string `json:"error"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
				errorResponse.Error = "Unknown error occurred"
			}
			tracing.TraceErr(span, errors.New(errorResponse.Error))

			// For 500 errors, use a generic message
			if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
				message := "Internal server error"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}

			// For other errors, propagate the status code and message from Mailstack
			h.responseHandler.HandleError(c, resp.StatusCode, &errorResponse.Error)
			return
		}

		// Parse response
		var response MailboxesResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			message := "Unable to parse Mailstack API response"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		h.responseHandler.HandleSuccess(c, response)
	}
}
