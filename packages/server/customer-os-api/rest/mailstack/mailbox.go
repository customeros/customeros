package restmailstack

import (
	"fmt"
	"github.com/gin-gonic/gin"
	coserrors "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"regexp"
	"strings"
)

// RegisterNewMailbox registers a new mailbox for the given domain
// @Summary Register a new mailbox
// @Description Registers a new mailbox for the specified domain
// @Tags MailStack API
// @Accept json
// @Produce json
// @Param domain path string true "Domain for which to register the mailbox"
// @Param body body MailboxRequest true "Mailbox registration payload"
// @Success 200 {object} MailboxResponse "Mailbox setup successful"
// @Failure 400 {object} rest.ErrorResponse "Invalid request body, missing input fields, or invalid username format"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized access - API key invalid or expired"
// @Failure 404 {object} rest.ErrorResponse "Domain not found"
// @Failure 409 {object} rest.ErrorResponse "Mailbox already exists"
// @Failure 500 {object} rest.ErrorResponse "Internal server error"
// @Router /mailstack/v1/domains/{domain}/mailboxes [post]
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
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing domain",
				})
			return
		}
		span.LogKV("request.domain", domain)

		// get tenant from context
		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			c.JSON(http.StatusUnauthorized,
				rest.ErrorResponse{
					Status:  "error",
					Message: "API key invalid or expired",
				})
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
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Invalid request body or missing input fields",
				})
			span.LogFields(tracingLog.String("result", "Invalid request body"))
			return
		}

		username := strings.TrimSpace(mailboxRequest.Username)
		if username == "" {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing username",
				})
			span.LogFields(tracingLog.String("result", "Missing username"))
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
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: err.Error(),
				})
			span.LogFields(tracingLog.String("result", "Invalid username format"))
			return
		}

		// add mailbox
		forwardingEnabled := true
		forwardingTo := mailboxRequest.ForwardingTo
		additionalForwardingTo := fmt.Sprintf("bcc@%s.customeros.ai", strings.ToLower(tenant))
		forwardingTo = append(forwardingTo, additionalForwardingTo)

		response, err := addMailbox(ctx, tenant, domain, username, password, forwardingEnabled, mailboxRequest.WebmailEnabled, forwardingTo, services)
		if err != nil {
			if errors.Is(err, coserrors.ErrDomainNotFound) {
				c.JSON(http.StatusNotFound,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Domain not found",
					})
				span.LogFields(tracingLog.String("result", "Domain not found"))
				return
			} else if errors.Is(err, coserrors.ErrMailboxExists) {
				c.JSON(http.StatusConflict,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Username already exists",
					})
				span.LogFields(tracingLog.String("result", "Mailbox already exists"))
				return
			} else {
				c.JSON(http.StatusInternalServerError,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Mailbox setup failed, please contact support",
					})
				span.LogFields(tracingLog.String("result", "Internal server error"))
				return
			}
		}

		response.Status = "success"
		response.Message = "Mailbox setup successful"
		if passwordGenerated {
			response.Password = password
		}
		c.JSON(http.StatusOK, response)
	}
}

func addMailbox(ctx context.Context, tenant, domain string, username, password string, forwardingEnabled, webmailEnabled bool, forwardingTo []string, services *service.Services) (MailboxResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "addMailbox")
	defer span.Finish()

	mailboxResponse := MailboxResponse{
		Email: username + "@" + domain,
	}

	// Check domain belongs to tenant
	domainBelongsToTenant, err := services.CommonServices.PostgresRepositories.MailStackDomainRepository.CheckDomainOwnership(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error checking domain"))
		return mailboxResponse, err
	}
	if !domainBelongsToTenant {
		return mailboxResponse, coserrors.ErrDomainNotFound
	}

	// Check mailbox doesn't already exist
	mailboxRecord, err := services.CommonServices.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, tenant, mailboxResponse.Email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error checking mailbox"))
		return mailboxResponse, err
	}
	if mailboxRecord != nil {
		return mailboxResponse, coserrors.ErrMailboxExists
	}

	err = services.OpensrsService.SetMailbox(ctx, tenant, domain, username, password, forwardingEnabled, forwardingTo, webmailEnabled)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error setting mailbox"))
		return mailboxResponse, err
	}

	// Save mailbox details in postgres
	err = services.CommonServices.PostgresRepositories.TenantSettingsMailboxRepository.Merge(ctx, tenant, &entity.TenantSettingsMailbox{
		Domain: domain, MailboxUsername: username + "@" + domain, Tenant: tenant, MailboxPassword: password, MinMinutesBetweenEmails: 5, MaxMinutesBetweenEmails: 10,
	})

	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error saving mailbox"))
		return mailboxResponse, err
	}

	mailboxResponse.WebmailEnabled = webmailEnabled
	mailboxResponse.ForwardingEnabled = forwardingEnabled
	mailboxResponse.ForwardingTo = forwardingTo

	return mailboxResponse, nil
}

func validateMailboxUsername(username string) error {
	// Regular expression for a valid username (allows alphanumeric, dots, underscores, hyphens)
	var re = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	if !re.MatchString(username) {
		return errors.New("invalid username format: only alphanumeric characters, dots, underscores, and hyphens are allowed")
	}
	// Additional checks (length, etc.) can be added if necessary
	return nil
}

// GetMailboxes retrieves all mailboxes for a specified domain
// @Summary Get all mailboxes
// @Description Retrieves a list of all mailboxes associated with a specified domain
// @Tags MailStack API
// @Accept json
// @Produce json
// @Param domain path string true "Domain for which to retrieve mailboxes"
// @Success 200 {object} MailboxesResponse "Successfully retrieved mailboxes"
// @Failure 400 {object} rest.ErrorResponse "Missing domain"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized access - API key invalid or expired"
// @Failure 500 {object} rest.ErrorResponse "Error retrieving mailboxes"
// @Router /mailstack/v1/domains/{domain}/mailboxes [get]
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
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing domain",
				})
			return
		}
		span.LogKV("request.domain", domain)

		// get tenant from context
		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			c.JSON(http.StatusUnauthorized,
				rest.ErrorResponse{
					Status:  "error",
					Message: "API key invalid or expired",
				})
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// get mailboxes for domain from postgres
		mailboxRecords, err := services.CommonServices.PostgresRepositories.TenantSettingsMailboxRepository.GetAllByDomain(ctx, tenant, domain)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error retrieving mailboxes"))
			span.LogFields(tracingLog.String("result", "Error retrieving mailboxes"))
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Error retrieving mailboxes",
				})
			return
		}

		response := MailboxesResponse{
			Status: "success",
		}
		for _, mailboxRecord := range mailboxRecords {
			mailboxDetails, err := services.OpensrsService.GetMailboxDetails(ctx, mailboxRecord.MailboxUsername)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error getting mailbox details"))
				span.LogFields(tracingLog.String("result", "Error getting mailbox details"))
				c.JSON(http.StatusInternalServerError,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Error getting mailbox details",
					})
				return
			}
			response.Mailboxes = append(response.Mailboxes, MailboxResponse{
				Email:             mailboxRecord.MailboxUsername,
				ForwardingEnabled: mailboxDetails.ForwardingEnabled,
				ForwardingTo:      mailboxDetails.ForwardingTo,
				WebmailEnabled:    mailboxDetails.WebmailEnabled,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}
