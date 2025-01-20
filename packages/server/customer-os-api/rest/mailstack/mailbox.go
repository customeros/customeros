// @openapi 3.0.0
package mailstack

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
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

		username := strings.TrimSpace(mailboxRequest.Username)
		if username == "" {
			message := "Missing parameter: username"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}
		span.LogKV("request.username", username)

		password := strings.TrimSpace(mailboxRequest.Password)
		passwordGenerated := false
		if password == "" {
			passwordGenerated = true
			password = utils.GenerateLowerAlpha(1) + utils.GenerateKey(11, false)
		}

		// validate username format
		if err := h.validateMailboxUsername(username); err != nil {
			message := "username is invalid"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		// add mailbox
		forwardingTo := mailboxRequest.ForwardingTo
		additionalForwardingTo := fmt.Sprintf("bcc@%s.customeros.ai", strings.ToLower(tenant))
		forwardingTo = append(forwardingTo, additionalForwardingTo)

		response := MailboxRecord{
			Email:             username + "@" + domain,
			WebmailEnabled:    mailboxRequest.WebmailEnabled,
			ForwardingEnabled: true,
			ForwardingTo:      forwardingTo,
		}

		err := h.services.CommonServices.MailboxService.CreateMailbox(ctx, nil, interfaces.CreateMailboxRequest{
			Domain:          domain,
			Username:        username,
			Password:        password,
			LinkedUserEmail: mailboxRequest.LinkedUser,
			WebmailEnabled:  mailboxRequest.WebmailEnabled,
			ForwardingTo:    forwardingTo,
		})
		if err != nil {
			if errors.Is(err, coserrors.ErrDomainNotFound) {
				message := "domain not found"
				h.responseHandler.HandleError(c, http.StatusNotFound, &message)
				return
			} else if errors.Is(err, coserrors.ErrMailboxExists) {
				message := "username already exists"
				h.responseHandler.HandleError(c, http.StatusConflict, &message)
				return
			} else {
				message := "Mailbox setup failed"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}
		}

		mailbox, err := h.services.Repositories.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, username+"@"+domain)
		if err != nil {
			message := "Error retrieving mailbox"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			tracing.TraceErr(span, errors.Wrap(err, "Error retrieving mailbox"))
			return
		}

		err = h.services.CommonServices.Events.Publisher.PublishEvent(ctx, mailbox.ID, model.MAILBOX, dto.MailstackProvisionMailbox{})
		if err != nil {
			message := "Error provisioning mailbox"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			tracing.TraceErr(span, errors.Wrap(err, "Error provisioning mailbox"))
			return
		}

		if passwordGenerated {
			response.Password = password
		}
		h.responseHandler.HandleSuccess(c, MailboxResponse{
			Mailbox: response,
		})
	}
}

func (h *MailstackHandler) validateMailboxUsername(username string) error {
	// Regular expression for a valid username (allows alphanumeric, dots, underscores, hyphens)
	re := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !re.MatchString(username) {
		return errors.New("invalid username format: only alphanumeric characters, dots, underscores, and hyphens are allowed")
	}
	// Additional checks (length, etc.) can be added if necessary
	return nil
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

		// get mailboxes for domain from postgres
		mailboxRecords, err := h.services.Repositories.PostgresRepositories.TenantSettingsMailboxRepository.GetAllByDomain(ctx, domain)
		if err != nil {
			message := "Error retrieving mailboxes"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			tracing.TraceErr(span, errors.Wrap(err, "Error retrieving mailboxes"))
			return
		}

		response := MailboxesResponse{
			Mailboxes: make([]MailboxRecord, 0, len(mailboxRecords)),
		}
		for _, mailboxRecord := range mailboxRecords {
			mailboxDetails, err := h.services.CommonServices.OpenSRSService.GetMailboxDetails(ctx, mailboxRecord.MailboxUsername)
			if err != nil {
				message := "Could not get mailbox details"
				tracing.TraceErr(span, errors.Wrap(err, message))
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}
			response.Mailboxes = append(response.Mailboxes, MailboxRecord{
				Email:             mailboxRecord.MailboxUsername,
				ForwardingEnabled: mailboxDetails.ForwardingEnabled,
				ForwardingTo:      mailboxDetails.ForwardingTo,
				WebmailEnabled:    mailboxDetails.WebmailEnabled,
			})
		}

		h.responseHandler.HandleSuccess(c, response)
	}
}
