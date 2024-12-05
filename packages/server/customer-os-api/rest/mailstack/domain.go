// @openapi 3.0.0
package restmailstack

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/coserrors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

// @title MailStack API
// @version 1.0
// @description API for managing domain registration and configuration for mail services
// @BasePath /mailstack/v1

// RegisterNewDomain registers a new domain for mail service
// @Summary Register new domain
// @Description Registers and configures a new domain for mail services, including DNS setup
// @Tags Domains
// @Accept json
// @Produce json
// @Param body body RegisterNewDomainRequest true "Domain registration details"
// @Success 201 {object} DomainResponse "Domain registered and configured successfully"
// @Failure 400 {object} rest.ErrorResponse "Invalid request - Missing required fields or invalid format"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 406 {object} rest.ErrorResponse "Not Acceptable - Domain TLD not supported or premium domain"
// @Failure 409 {object} rest.ErrorResponse "Conflict - Domain already registered"
// @Failure 500 {object} rest.ErrorResponse "Internal server error - Configuration failed or service unavailable"
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
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrUnauthorized)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// Parse and validate request body
		var req RegisterNewDomainRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest)
			span.LogFields(tracingLog.String("result", "Invalid request body"))
			return
		}

		// Check for missing domain
		if req.Domain == "" {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Missing required field: domain"))
			return
		} else if req.Website == "" {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Missing required field: website"))
			return
		}

		registerNewDomainResponse, err := registerDomain(ctx, tenant, req.Domain, req.Website, services)
		if err != nil {
			if errors.Is(err, coserrors.ErrNotSupported) {
				rest.SendError(c, span, http.StatusNotAcceptable, rest.ErrBadRequest.WithMessage("Domain TLD not supported"))
				return
			} else if errors.Is(err, coserrors.ErrDomainUnavailable) {
				rest.SendError(c, span, http.StatusConflict, rest.ErrConflict.WithMessage("Domain already registered"))
				return
			} else if errors.Is(err, coserrors.ErrDomainPremium) {
				rest.SendError(c, span, http.StatusNotAcceptable, rest.ErrBadRequest.WithMessage("Premium domain names are not supported"))
				return
			} else if errors.Is(err, coserrors.ErrDomainPriceExceeded) {
				rest.SendError(c, span, http.StatusNotAcceptable, rest.ErrBadRequest.WithMessage("Unauthorized to purchase domain, please contact support"))
				return
			} else if errors.Is(err, coserrors.ErrDomainConfigurationFailed) {
				rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Unable to configure domain"))
				return
			} else if errors.Is(err, coserrors.ErrConnectionTimeout) {
				rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Connection timeout, please retry"))
				return
			} else {
				rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Domain registratin failed, please contact support"))
				return
			}
		}

		c.JSON(http.StatusCreated, DomainResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Domain:       registerNewDomainResponse,
		})
	}
}

func registerDomain(ctx context.Context, tenant, domain, website string, services *service.Services) (DomainRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "registerDomain")
	defer span.Finish()

	registerNewDomainResponse := DomainRecord{}

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
	isAvailable, isPremium, err := services.CommonServices.NamecheapService.CheckDomainAvailability(ctx, domain)
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
	domainPrice, err := services.CommonServices.NamecheapService.GetDomainPrice(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error getting domain price"))
		return registerNewDomainResponse, err
	}
	if domainPrice > services.Cfg.ExternalServices.NamecheapConfig.MaxPrice {
		return registerNewDomainResponse, coserrors.ErrDomainPriceExceeded
	}

	//step 3 - register domain
	err = services.CommonServices.NamecheapService.PurchaseDomain(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error purchasing domain"))
		return registerNewDomainResponse, err
	}

	// step 4 - configure domain
	return configureDomain(ctx, tenant, domain, website, services)
}

// ConfigureDomain configures DNS for an existing domain
// @Summary Configure domain DNS
// @Description Sets up DNS records and mail services for an existing domain
// @Tags Domains
// @Accept json
// @Produce json
// @Param body body ConfigureDomainRequest true "Domain configuration details"
// @Success 201 {object} DomainResponse "Domain configured successfully"
// @Failure 400 {object} rest.ErrorResponse "Invalid request - Missing required fields or invalid format"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 404 {object} rest.ErrorResponse "Not Found - Domain not found or not owned by tenant"
// @Failure 500 {object} rest.ErrorResponse "Internal server error - Configuration failed or service unavailable"
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
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrUnauthorized)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// Parse and validate request body
		var req ConfigureDomainRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest)
			return
		}

		// Check for missing domain
		if req.Domain == "" {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Missing required field: domain"))
			return
		} else if req.Website == "" {
			rest.SendError(c, span, http.StatusBadRequest, rest.ErrBadRequest.WithMessage("Missing required field: website"))
			return
		}

		domainResponse, err := configureDomain(ctx, tenant, req.Domain, req.Website, services)
		if err != nil {
			if errors.Is(err, coserrors.ErrDomainNotFound) {
				rest.SendError(c, span, http.StatusNotFound, rest.ErrNotFound)
				return
			} else if errors.Is(err, coserrors.ErrDomainConfigurationFailed) {
				rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Unable to configure domain"))
				return
			} else {
				rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Domain registration failed"))
				return
			}
		}

		c.JSON(http.StatusCreated, DomainResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Domain:       domainResponse,
		})
	}
}

func configureDomain(ctx context.Context, tenant, domain, redirectWebsite string, services *service.Services) (DomainRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "configureDomain")
	defer span.Finish()

	domainResponse := DomainRecord{}
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

	err = services.CommonServices.MailstackService.ConfigureMailstackDomain(ctx, domain, redirectWebsite)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error configuring domain"))
		return domainResponse, coserrors.ErrDomainConfigurationFailed
	}

	// get domain details
	domainInfo, err := services.CommonServices.NamecheapService.GetDomainInfo(ctx, tenant, domain)
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

// GetDomains retrieves all active domains
// @Summary List active domains
// @Description Retrieves all active domains for the authenticated tenant
// @Tags Domains
// @Accept json
// @Produce json
// @Success 200 {object} DomainsResponse "Successfully retrieved domain list"
// @Failure 401 {object} rest.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 500 {object} rest.ErrorResponse "Internal server error - Unable to retrieve domains"
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
			rest.SendError(c, span, http.StatusUnauthorized, rest.ErrUnauthorized)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// get all active domains from postgres
		activeDomainRecords, err := services.CommonServices.PostgresRepositories.MailStackDomainRepository.GetActiveDomains(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error retrieving domains"))
			rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Unable to retrieve domains"))
			return
		}

		response := DomainsResponse{
			BaseResponse: rest.BuildBaseResponse(rest.StatusSuccess),
			Domains:      make([]DomainRecord, 0, len(activeDomainRecords)),
		}

		for _, domainRecord := range activeDomainRecords {
			domain, err := services.CommonServices.NamecheapService.GetDomainInfo(ctx, tenant, domainRecord.Domain)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error getting domain info"))
				span.LogFields(tracingLog.String("result", "Error getting domain info"))
				rest.SendError(c, span, http.StatusInternalServerError, rest.ErrInternalServer.WithMessage("Unable to retrieve domain info"))
				return
			}
			response.Domains = append(response.Domains, DomainRecord{
				Domain:      domain.DomainName,
				CreatedDate: domain.CreatedDate,
				ExpiredDate: domain.ExpiredDate,
				Nameservers: domain.Nameservers,
			})
		}

		c.JSON(http.StatusOK, response)
	}
}
