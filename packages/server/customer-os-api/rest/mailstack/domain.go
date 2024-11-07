package restmailstack

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/coserrors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

// RegisterNewDomain registers a new domain for the mail service
// @Summary Register a new domain
// @Description Registers a new domain for the mail service
// @Tags MailStack API
// @Accept json
// @Produce json
// @Param body body RegisterNewDomainRequest true "Domain registration payload"
// @Success 201 {object} DomainResponse "Domain registered successfully"
// @Failure 400 {object} rest.ErrorResponse "Invalid request body or missing input fields"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized access - API key invalid or expired"
// @Failure 409 {object} rest.ErrorResponse "Domain is already registered"
// @Failure 406 {object} rest.ErrorResponse "Restrictions on domain purchase"
// @Failure 500 {object} rest.ErrorResponse "Internal server error"
// @Router /mailstack/v1/domains [post]
// @Security ApiKeyAuth
func RegisterNewDomain(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RegisterNewDomain", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

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
		} else if req.Website == "" {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing required field: website",
				})
			span.LogFields(tracingLog.String("result", "Missing website"))
			return
		}

		registerNewDomainResponse, err := registerDomain(ctx, tenant, req.Domain, req.Website, services)
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
				return
			} else if errors.Is(err, coserrors.ErrDomainPriceExceeded) {
				c.JSON(http.StatusNotAcceptable,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Domain price exceeds the maximum allowed price, please contact support",
					})
				return
			} else if errors.Is(err, coserrors.ErrDomainConfigurationFailed) {
				c.JSON(http.StatusInternalServerError,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Error configuring domain, please contact support",
					})
				return
			} else if errors.Is(err, coserrors.ErrConnectionTimeout) {
				c.JSON(http.StatusInternalServerError,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Connection timeout, please try again",
					})
				return
			} else {
				c.JSON(http.StatusInternalServerError,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Domain registration failed, please contact support",
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

func registerDomain(ctx context.Context, tenant, domain, website string, services *service.Services) (DomainResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "registerDomain")
	defer span.Finish()

	var registerNewDomainResponse = DomainResponse{}
	registerNewDomainResponse.Domain = domain

	var err error

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

	//step 1 - check domain availability
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

	//step 3 - register domain
	err = services.NamecheapService.PurchaseDomain(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error purchasing domain"))
		return registerNewDomainResponse, err
	}

	//step 4 - configure domain
	return configureDomain(ctx, tenant, domain, website, services)
}

// ConfigureDomain configures the given domain for the mail service
// @Summary Configure domain
// @Description Configures the DNS records for the given domain
// @Tags MailStack API
// @Accept json
// @Produce json
// @Param body body ConfigureDomainRequest true "Domain payload"
// @Success 201 {object} DomainResponse "Domain configured successfully"
// @Failure 400 {object} rest.ErrorResponse "Invalid request body or missing input fields"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized access - API key invalid or expired"
// @Failure 404 {object} rest.ErrorResponse "Domain not found"
// @Failure 500 {object} rest.ErrorResponse "Internal server error"
// @Router /mailstack/v1/domains/configure [post]
// @Security ApiKeyAuth
func ConfigureDomain(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "ConfigureDomain", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

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
		var req ConfigureDomainRequest
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
		} else if req.Website == "" {
			c.JSON(http.StatusBadRequest,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Missing required field: website",
				})
			span.LogFields(tracingLog.String("result", "Missing website"))
			return
		}

		domainResponse, err := configureDomain(ctx, tenant, req.Domain, req.Website, services)
		if err != nil {
			if errors.Is(err, coserrors.ErrDomainNotFound) {
				c.JSON(http.StatusNotFound,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Domain not found",
					})
				span.LogFields(tracingLog.String("result", "Domain not found"))
				return
			} else if errors.Is(err, coserrors.ErrDomainConfigurationFailed) {
				c.JSON(http.StatusInternalServerError,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Error configuring domain, please contact support",
					})
				return
			} else {
				c.JSON(http.StatusInternalServerError,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Domain registration failed, please contact support",
					})
				span.LogFields(tracingLog.String("result", "Internal server error"))
				return
			}
		}

		domainResponse.Status = "success"
		domainResponse.Message = "Domain configured successfully"

		// Placeholder for response logic (to be implemented)
		c.JSON(http.StatusCreated, domainResponse)
	}
}

func configureDomain(ctx context.Context, tenant, domain, website string, services *service.Services) (DomainResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "configureDomain")
	defer span.Finish()

	var domainResponse = DomainResponse{}
	domainResponse.Domain = domain

	var err error

	domainBelongsToTenant, err := services.CommonServices.PostgresRepositories.MailStackDomainRepository.CheckDomainOwnership(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error checking domain"))
		return domainResponse, err
	}
	if !domainBelongsToTenant {
		return domainResponse, coserrors.ErrDomainNotFound
	}

	// setup domain in cloudflare
	nameservers, err := services.CloudflareService.SetupDomainForMailStack(ctx, tenant, domain, website)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error setting up domain in Cloudflare"))
		return domainResponse, coserrors.ErrDomainConfigurationFailed
	}

	// setup domain in openSRS
	err = services.CommonServices.OpenSrsService.SetupDomainForMailStack(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error setting up domain in OpenSRS"))
		return domainResponse, coserrors.ErrDomainConfigurationFailed
	}

	// replace nameservers in namecheap
	err = services.NamecheapService.UpdateNameservers(ctx, tenant, domain, nameservers)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error updating nameservers"))
		return domainResponse, coserrors.ErrDomainConfigurationFailed
	}

	// mark domain as configured
	err = services.CommonServices.PostgresRepositories.MailStackDomainRepository.MarkConfigured(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error setting domain as configured"))
	}

	// get domain details
	domainInfo, err := services.NamecheapService.GetDomainInfo(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error getting domain info"))
		return domainResponse, err
	}
	domainResponse.CreatedDate = domainInfo.CreatedDate
	domainResponse.ExpiredDate = domainInfo.ExpiredDate
	domainResponse.Nameservers = domainInfo.Nameservers
	domainResponse.Domain = domainInfo.DomainName

	return domainResponse, nil
}

// GetDomains retrieves all active domains for the tenant
// @Summary Get active domains
// @Description Retrieves a list of all active domains associated with the tenant
// @Tags MailStack API
// @Accept json
// @Produce json
// @Success 200 {object} DomainsResponse "Successfully retrieved domains"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized access - API key invalid or expired"
// @Failure 500 {object} rest.ErrorResponse "Internal server error"
// @Router /mailstack/v1/domains [get]
// @Security ApiKeyAuth
func GetDomains(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetDomains", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

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

		// get all active domains from postgres
		activeDomainRecords, err := services.CommonServices.PostgresRepositories.MailStackDomainRepository.GetActiveDomains(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error retrieving domains"))
			span.LogFields(tracingLog.String("result", "Error retrieving domains"))
			c.JSON(http.StatusInternalServerError,
				rest.ErrorResponse{
					Status:  "error",
					Message: "Error retrieving domains",
				})
			return
		}

		response := DomainsResponse{
			Status: "success",
		}
		for _, domainRecord := range activeDomainRecords {
			domain, err := services.NamecheapService.GetDomainInfo(ctx, tenant, domainRecord.Domain)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error getting domain info"))
				span.LogFields(tracingLog.String("result", "Error getting domain info"))
				c.JSON(http.StatusInternalServerError,
					rest.ErrorResponse{
						Status:  "error",
						Message: "Error getting domain info",
					})
				return
			}
			response.Domains = append(response.Domains, DomainResponse{
				Domain:      domain.DomainName,
				CreatedDate: domain.CreatedDate,
				ExpiredDate: domain.ExpiredDate,
				Nameservers: domain.Nameservers,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}
