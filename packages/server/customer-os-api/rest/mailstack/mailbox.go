// @openapi 3.0.0
package restmailstack

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	service2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/coserrors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
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
// @Failure 400 {object} rest.ErrorResponse "Invalid request - Missing or invalid parameters"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 404 {object} rest.ErrorResponse "Domain not found"
// @Failure 409 {object} rest.ErrorResponse "Mailbox already exists"
// @Failure 500 {object} rest.ErrorResponse "Internal server error"
// @Router /domains/{domain}/mailboxes [post]
// @Security ApiKeyAuth
func RegisterNewMailbox(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RegisterNewMailbox", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		// get domain from path
		domain := c.Param("domain")
		if domain == "" {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Missing parameter: domain"))
			return
		}
		span.LogKV("request.domain", domain)

		// get tenant from context
		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrUnauthorized)
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
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest)
			return
		}

		username := strings.TrimSpace(mailboxRequest.Username)
		if username == "" {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Missing parameter: username"))
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
		if err := validateMailboxUsername(username); err != nil {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("username is invalid"))
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

		err := services.CommonServices.MailboxService.CreateMailbox(ctx, nil, service2.CreateMailboxRequest{
			Domain:          domain,
			Username:        username,
			Password:        password,
			LinkedUserEmail: mailboxRequest.LinkedUser,
			WebmailEnabled:  mailboxRequest.WebmailEnabled,
			ForwardingTo:    forwardingTo,
		})

		if err != nil {
			if errors.Is(err, coserrors.ErrDomainNotFound) {
				rest.SendError(c, span, http.StatusNotFound, rest.ErrNotFound.WithMessage("Domain not found"))
				return
			} else if errors.Is(err, coserrors.ErrMailboxExists) {
				rest.SendError(c, span, http.StatusConflict, rest.ErrConflict.WithMessage("Username already exists"))
				return
			} else {
				rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Mailbox setup failed, please contact support"))
				return
			}
		}

		mailbox, err := services.CommonServices.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, username+"@"+domain)
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Error retrieving mailbox"))
			tracing.TraceErr(span, errors.Wrap(err, "Error retrieving mailbox"))
			return
		}

		err = services.CommonServices.RabbitMQService.PublishEvent(ctx, mailbox.ID, model.MAILBOX, dto.MailstackProvisionMailbox{})
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Error provisioning mailbox"))
			tracing.TraceErr(span, errors.Wrap(err, "Error provisioning mailbox"))
			return
		}

		if passwordGenerated {
			response.Password = password
		}
		c.JSON(http.StatusOK, MailboxResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Mailbox:      response,
		})
	}
}

func validateMailboxUsername(username string) error {
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
// @Failure 400 {object} rest.ErrorResponse "Missing domain parameter"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 404 {object} rest.ErrorResponse "Domain not found"
// @Failure 500 {object} rest.ErrorResponse "Internal server error"
// @Router /domains/{domain}/mailboxes [get]
// @Security ApiKeyAuth
func GetMailboxes(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetMailboxes", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		// get domain from path
		domain := c.Param("domain")
		if domain == "" {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Missing parameter: domain"))
			return
		}
		span.LogKV("request.domain", domain)

		// get tenant from context
		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrUnauthorized)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// get mailboxes for domain from postgres
		mailboxRecords, err := services.CommonServices.PostgresRepositories.TenantSettingsMailboxRepository.GetAllByDomain(ctx, domain)
		if err != nil {
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Error retrieving mailboxes"))
			tracing.TraceErr(span, errors.Wrap(err, "Error retrieving mailboxes"))
			return
		}

		response := MailboxesResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Mailboxes:    make([]MailboxRecord, 0, len(mailboxRecords)),
		}
		for _, mailboxRecord := range mailboxRecords {
			mailboxDetails, err := services.CommonServices.OpenSrsService.GetMailboxDetails(ctx, mailboxRecord.MailboxUsername)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error getting mailbox details"))
				rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Error getting mailbox details"))
				return
			}
			response.Mailboxes = append(response.Mailboxes, MailboxRecord{
				Email:             mailboxRecord.MailboxUsername,
				ForwardingEnabled: mailboxDetails.ForwardingEnabled,
				ForwardingTo:      mailboxDetails.ForwardingTo,
				WebmailEnabled:    mailboxDetails.WebmailEnabled,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}
