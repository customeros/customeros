package restmailstack

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
	"net/http"
)

// RegisterNewDomain registers a new domain for the mail service
// @Summary Register a new domain
// @Description Registers a new domain with a list of mailboxes in the MailStack system.
// @Tags MailStack API
// @Accept  json
// @Produce  json
// @Param   body   body    RegisterNewDomainRequest  true  "Domain registration payload"
// @Success 201 {object} RegisterNewDomainResponse "Domain registered successfully"
// @Failure 400  "Invalid request body or missing input fields"
// @Failure 401  "Unauthorized access - API key invalid or expired"
// @Failure 500  "Internal server error"
// @Router /mailstack/v1/domains [post]
// @Security ApiKeyAuth
func RegisterNewDomain(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "CreateOrganization", c.Request.Header)
		defer span.Finish()
		tracing.SetDefaultRestSpanTags(ctx, span)

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
		span.SetTag(tracing.SpanTagTenant, tenant)

		// Parse and validate request body
		var req RegisterNewDomainRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Invalid request body or missing input fields",
				})
			span.LogFields(tracingLog.String("result", "Invalid request body"))
			return
		}

		// Check for missing domain
		if req.Domain == "" {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing required field: domain",
				})
			span.LogFields(tracingLog.String("result", "Missing domain"))
			return
		}

		// Check if mailboxes are provided and validate each mailbox
		if len(req.Mailboxes) == 0 {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing required field: mailboxes",
				})
			span.LogFields(tracingLog.String("result", "Missing mailboxes"))
			return
		}

		for _, mailbox := range req.Mailboxes {
			if mailbox.Username == "" || mailbox.Password == "" {
				c.JSON(http.StatusBadRequest,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Each mailbox must have a username and password",
					})
				span.LogFields(tracingLog.String("result", "Invalid mailbox configuration"))
				return
			}
		}

		callDomainRegistrationWithMailBoxes(ctx, req)

		// Placeholder for response logic (to be implemented)
		c.JSON(http.StatusCreated,
			RegisterNewDomainResponse{
				Status:  "success",
				Message: "Domain registered successfully",
				// More response details will be added here later, including:
				// "domain": req.Domain,
				// "mailboxes": req.Mailboxes,
			})
	}
}

func callDomainRegistrationWithMailBoxes(ctx context.Context, req RegisterNewDomainRequest) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "callDomainRegistrationWithMailBoxes")
	defer span.Finish()

	// step 1 - register domain

	// step 2 - configure domain with cloudflare

	// step 3 - configure mailboxes

	// step 4 - warm mailboxes
}
