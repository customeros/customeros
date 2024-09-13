package restmailstack

import (
	"github.com/gin-gonic/gin"
	coserrors "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

// RegisterNewDomain registers a new domain for the mail service
// @Summary Register a new domain
// @Description Registers a new domain
// @Tags MailStack API
// @Accept  json
// @Produce  json
// @Param   body   body    RegisterNewDomainRequest  true  "Domain registration payload"
// @Success 201 {object} RegisterNewDomainResponse "Domain registered successfully"
// @Failure 400  "Invalid request body or missing input fields"
// @Failure 401  "Unauthorized access - API key invalid or expired"
// @Failure 409  "Domain is already registered"
// @Failure 406  "Restrictions on domain purchase"
// @Failure 500  "Internal server error"
// @Router /mailstack/v1/domains [post]
// @Security ApiKeyAuth
func RegisterNewDomain(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RegisterNewDomain", c.Request.Header)
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

		registerNewDomainResponse, err := registerDomain(ctx, tenant, req.Domain, services)
		if err != nil {
			if errors.Is(err, coserrors.ErrNotSupported) {
				c.JSON(http.StatusNotAcceptable,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Domain TLD not supported, please choose another domain",
					})
				span.LogFields(tracingLog.String("result", "Domain TLD not supported"))
				return
			} else if errors.Is(err, coserrors.ErrDomainUnavailable) {
				c.JSON(http.StatusConflict,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Domain is already registered, please choose another domain",
					})
				return
			} else if errors.Is(err, coserrors.ErrDomainPremium) {
				c.JSON(http.StatusNotAcceptable,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Not supporting premium domains, please choose another domain",
					})
			} else if errors.Is(err, coserrors.ErrDomainPriceExceeded) {
				c.JSON(http.StatusNotAcceptable,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Domain price exceeds the maximum allowed price, please contact support",
					})
			} else {
				c.JSON(http.StatusInternalServerError,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Domain registration failed",
					})
				span.LogFields(tracingLog.String("result", "Internal server error, please contact support"))
				return
			}
		}

		registerNewDomainResponse.Status = "success"
		registerNewDomainResponse.Message = "Domain registered successfully"

		// Placeholder for response logic (to be implemented)
		c.JSON(http.StatusCreated, registerNewDomainResponse)
	}
}

func registerDomain(ctx context.Context, tenant, domain string, services *service.Services) (RegisterNewDomainResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "registerDomain")
	defer span.Finish()

	var registerNewDomainResponse = RegisterNewDomainResponse{}
	registerNewDomainResponse.Domain = domain

	// check if domain tld is supported
	// Extract the TLD from the domain (e.g., "com" from "example.com")
	tld := strings.Split(domain, ".")[1]
	tldSupported := false
	for _, supportedTld := range services.Cfg.AppConfig.Mailstack.SupportedTlds {
		if tld == supportedTld {
			tldSupported = true
			break
		}
	}
	if !tldSupported {
		return registerNewDomainResponse, coserrors.ErrNotSupported
	}

	// step 1 - check domain availability
	isAvailable, isPremium, err := services.NamecheapService.CheckDomainAvailability(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error checking domain availability"))
		return registerNewDomainResponse, err
	}
	if !isAvailable {
		return registerNewDomainResponse, coserrors.ErrDomainUnavailable
	}
	if isPremium {
		return registerNewDomainResponse, coserrors.ErrDomainPremium
	}

	// step 2 - check pricing
	domainPrice, err := services.NamecheapService.GetDomainPrice(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error getting domain price"))
		return registerNewDomainResponse, err
	}
	if domainPrice > services.Cfg.ExternalServices.Namecheap.MaxPrice {
		return registerNewDomainResponse, coserrors.ErrDomainPriceExceeded
	}

	// step 3 - register domain
	err = services.NamecheapService.PurchaseDomain(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error purchasing domain"))
		return registerNewDomainResponse, err
	}

	// step X - configure domain with cloudflare

	return registerNewDomainResponse, nil
}
