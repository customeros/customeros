// @openapi 3.0.0
package mailstack

import (
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
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
		statusCode, errMsg, result, err := h.services.CommonServices.MailstackService.RegisterMailbox(ctx, tenant, domain, interfaces.CreateMailboxRequest{
			Username:        mailboxRequest.Username,
			Password:        mailboxRequest.Password,
			ForwardingTo:    mailboxRequest.ForwardingTo,
			WebmailEnabled:  mailboxRequest.WebmailEnabled,
			LinkedUserEmail: mailboxRequest.LinkedUser,
		})
		if err != nil {
			h.responseHandler.HandleError(c, statusCode, &errMsg)
			tracing.TraceErr(span, err)
			return
		}

		if statusCode != http.StatusOK {
			h.responseHandler.HandleError(c, statusCode, &errMsg)
			return
		}

		h.responseHandler.HandleSuccess(c, MailboxResponse{
			Mailbox: MailboxRecord{
				Email:             result.Email,
				Password:          result.Password,
				ForwardingTo:      result.ForwardingTo,
				ForwardingEnabled: result.ForwardingEnabled,
				WebmailEnabled:    result.WebmailEnabled,
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
		tracing.SetDefaultRestSpanTags(ctx, span)

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

		// Get mailboxes using service
		statusCode, errMsg, mailboxes, err := h.services.CommonServices.MailstackService.GetMailboxes(ctx, tenant, domain, "")
		if err != nil {
			message := "Internal server error"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			tracing.TraceErr(span, err)
			return
		}

		// Handle non-200 responses
		if statusCode != http.StatusOK {
			h.responseHandler.HandleError(c, statusCode, &errMsg)
			return
		}

		// Convert service response to API response
		var response MailboxesResponse
		for _, mailbox := range mailboxes {
			response.Mailboxes = append(response.Mailboxes, MailboxRecord{
				Email:             mailbox.Email,
				ForwardingTo:      mailbox.ForwardingTo,
				ForwardingEnabled: mailbox.ForwardingEnabled,
				WebmailEnabled:    mailbox.WebmailEnabled,
			})
		}

		h.responseHandler.HandleSuccess(c, response)
	}
}
