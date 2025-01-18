// @openapi 3.0.0
package mailstack

import (
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"

	"github.com/customeros/customeros/packages/server/customer-os-api/enum"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
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
// @Failure 400 {object} handlers.ErrorResponse "Invalid request - Missing required fields or invalid format"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 406 {object} handlers.ErrorResponse "Not Acceptable - Domain TLD not supported or premium domain"
// @Failure 409 {object} handlers.ErrorResponse "Conflict - Domain already registered"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error - Configuration failed or service unavailable"
// @Router /mailstack/v1/domains [post]
// @Security ApiKeyAuth
func RegisterNewDomain(services *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RegisterNewDomain", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// Parse and validate request body
		var req RegisterNewDomainRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest)
			span.LogFields(tracingLog.String("result", "Invalid request body"))
			return
		}

		// Check for missing domain
		if req.Domain == "" {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Missing required field: domain"))
			return
		} else if req.Website == "" {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Missing required field: website"))
			return
		}

		registerNewDomainResponse, err := registerDomain(ctx, tenant, req.Domain, req.Website, services)
		if err != nil {
			if errors.Is(err, coserrors.ErrNotSupported) {
				handlers.SendError(c, span, http.StatusNotAcceptable, enum.ErrBadRequest.WithMessage("Domain TLD not supported"))
				return
			} else if errors.Is(err, coserrors.ErrDomainUnavailable) {
				handlers.SendError(c, span, http.StatusConflict, enum.ErrConflict.WithMessage("Domain already registered"))
				return
			} else if errors.Is(err, coserrors.ErrDomainPremium) {
				handlers.SendError(c, span, http.StatusNotAcceptable, enum.ErrBadRequest.WithMessage("Premium domain names are not supported"))
				return
			} else if errors.Is(err, coserrors.ErrDomainPriceExceeded) {
				handlers.SendError(c, span, http.StatusNotAcceptable, enum.ErrBadRequest.WithMessage("Unauthorized to purchase domain, please contact support"))
				return
			} else if errors.Is(err, coserrors.ErrDomainConfigurationFailed) {
				handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Unable to configure domain"))
				return
			} else if errors.Is(err, coserrors.ErrConnectionTimeout) {
				handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Connection timeout, please retry"))
				return
			} else {
				handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Domain registratin failed, please contact support"))
				return
			}
		}

		c.JSON(http.StatusCreated, DomainResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Domain:       registerNewDomainResponse,
		})
	}
}

func registerDomain(ctx context.Context, tenant, domain, website string, services *cosapi_services.Services) (DomainRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "registerDomain")
	defer span.Finish()

	registerNewDomainResponse := DomainRecord{}

	var err error

	// check if domain tld is supported
	// Extract the TLD from the domain (e.g., "com" from "example.com")
	tld := strings.Split(domain, ".")[1]
	tldSupported := false
	for _, supportedTld := range services.Cfg.Common.Internal.MailstackConfig.SupportedTlds {
		if tld == supportedTld {
			tldSupported = true
			break
		}
	}
	if !tldSupported {
		return registerNewDomainResponse, coserrors.ErrNotSupported
	}

	// step 1 - check domain availability
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
	if domainPrice > services.Cfg.Common.External.NamecheapConfig.MaxPrice {
		return registerNewDomainResponse, coserrors.ErrDomainPriceExceeded
	}

	// step 3 - register domain
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
// @Failure 400 {object} handlers.ErrorResponse "Invalid request - Missing required fields or invalid format"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 404 {object} handlers.ErrorResponse "Not Found - Domain not found or not owned by tenant"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error - Configuration failed or service unavailable"
// @Router /mailstack/v1/domains/configure [post]
// @Security ApiKeyAuth
func ConfigureDomain(services *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "ConfigureDomain", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// Parse and validate request body
		var req ConfigureDomainRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest)
			return
		}

		// Check for missing domain
		if req.Domain == "" {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Missing required field: domain"))
			return
		} else if req.Website == "" {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Missing required field: website"))
			return
		}

		domainResponse, err := configureDomain(ctx, tenant, req.Domain, req.Website, services)
		if err != nil {
			if errors.Is(err, coserrors.ErrDomainNotFound) {
				handlers.SendError(c, span, http.StatusNotFound, enum.ErrNotFound)
				return
			} else if errors.Is(err, coserrors.ErrDomainConfigurationFailed) {
				handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Unable to configure domain"))
				return
			} else {
				handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Domain registration failed"))
				return
			}
		}

		c.JSON(http.StatusCreated, DomainResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Domain:       domainResponse,
		})
	}
}

func configureDomain(ctx context.Context, tenant, domain, redirectWebsite string, services *cosapi_services.Services) (DomainRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "configureDomain")
	defer span.Finish()

	domainResponse := DomainRecord{}
	domainResponse.Domain = domain

	var err error

	domainBelongsToTenant, err := services.Repositories.PostgresRepositories.MailStackDomainRepository.CheckDomainOwnership(ctx, tenant, domain)
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
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error - Unable to retrieve domains"
// @Router /mailstack/v1/domains [get]
// @Security ApiKeyAuth
func GetDomains(services *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetDomains", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// get all active domains from postgres
		activeDomainRecords, err := services.Repositories.PostgresRepositories.MailStackDomainRepository.GetActiveDomains(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error retrieving domains"))
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Unable to retrieve domains"))
			return
		}

		response := DomainsResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Domains:      make([]DomainRecord, 0, len(activeDomainRecords)),
		}

		for _, domainRecord := range activeDomainRecords {
			domain, err := services.CommonServices.NamecheapService.GetDomainInfo(ctx, tenant, domainRecord.Domain)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error getting domain info"))
				span.LogFields(tracingLog.String("result", "Error getting domain info"))
				handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Unable to retrieve domain info"))
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

func RecommendDomain(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RecommendDomain", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// get root domain
		baseName, exists := c.GetQuery("baseName")
		if !exists {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("must provide baseName parameter"))
			return
		}

		// get domain recommendations
		recommendations := s.CommonServices.MailboxService.RecommendOutboundDomains(ctx, baseName, 20)

		c.JSON(http.StatusOK, DomainRecommendationResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Domains:      recommendations,
		})
	}
}
