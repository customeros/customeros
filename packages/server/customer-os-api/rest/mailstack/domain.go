// @openapi 3.0.0
package mailstack

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
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
func (h *MailstackHandler) RegisterNewDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RegisterNewDomain", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			tracing.TraceErr(span, errors.New("Missing tenant in context"))
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		// Parse and validate request body
		var req RegisterNewDomainRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		// Check for missing domain
		if req.Domain == "" {
			message := "Missing required field: domain"
			tracing.TraceErr(span, errors.New(message))
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		} else if req.Website == "" {
			message := "Missing required field: website"
			tracing.TraceErr(span, errors.New(message))
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		registerNewDomainResponse, err := h.registerDomain(ctx, tenant, req.Domain, req.Website)
		if err != nil {
			if errors.Is(err, coserrors.ErrNotSupported) {
				message := "Domain TLD not supported"
				tracing.TraceErr(span, errors.New(message))
				h.responseHandler.HandleError(c, http.StatusNotAcceptable, &message)
				return
			} else if errors.Is(err, coserrors.ErrDomainUnavailable) {
				message := "Domain already registered"
				tracing.TraceErr(span, errors.New(message))
				h.responseHandler.HandleError(c, http.StatusConflict, &message)
				return
			} else if errors.Is(err, coserrors.ErrDomainPremium) {
				message := "Premium domain names are not supported"
				tracing.TraceErr(span, errors.New(message))
				h.responseHandler.HandleError(c, http.StatusNotAcceptable, &message)
				return
			} else if errors.Is(err, coserrors.ErrDomainPriceExceeded) {
				message := "Unauthorized to purchase domain"
				tracing.TraceErr(span, errors.New(message))
				h.responseHandler.HandleError(c, http.StatusNotAcceptable, &message)
				return
			} else if errors.Is(err, coserrors.ErrDomainConfigurationFailed) {
				message := "Unable to configure domain"
				tracing.TraceErr(span, errors.New(message))
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			} else if errors.Is(err, coserrors.ErrConnectionTimeout) {
				message := "Connection timeout, please retry"
				tracing.TraceErr(span, errors.New(message))
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			} else {
				message := "Domain registration failed, please contact support"
				tracing.TraceErr(span, errors.New(message))
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}
		}

		h.responseHandler.HandleSuccess(c, DomainResponse{
			Domain: registerNewDomainResponse,
		})
	}
}

func (h *MailstackHandler) registerDomain(ctx context.Context, tenant, domain, website string) (DomainRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "registerDomain")
	defer span.Finish()
	tracing.TagComponentRest(span)

	registerNewDomainResponse := DomainRecord{}

	var err error

	// check if domain tld is supported
	// Extract the TLD from the domain (e.g., "com" from "example.com")
	tld := strings.Split(domain, ".")[1]
	if !utils.IsStringInSlice(tld, h.services.Cfg.Common.Internal.MailstackApiConfig.SupportedTlds) {
		return registerNewDomainResponse, coserrors.ErrNotSupported
	}

	// step 1 - check domain availability
	isAvailable, isPremium, err := h.services.CommonServices.NamecheapService.CheckDomainAvailability(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error checking domain availability"))
		return registerNewDomainResponse, err
	}
	if !isAvailable {
		tracing.TraceErr(span, coserrors.ErrDomainUnavailable)
		return registerNewDomainResponse, coserrors.ErrDomainUnavailable
	}
	if isPremium {
		tracing.TraceErr(span, coserrors.ErrDomainPremium)
		return registerNewDomainResponse, coserrors.ErrDomainPremium
	}

	// step 2 - check pricing
	domainPrice, err := h.services.CommonServices.NamecheapService.GetDomainPrice(ctx, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error getting domain price"))
		return registerNewDomainResponse, err
	}
	if domainPrice > h.services.Cfg.Common.External.NamecheapConfig.MaxPrice {
		return registerNewDomainResponse, coserrors.ErrDomainPriceExceeded
	}

	// step 3 - register domain
	err = h.services.CommonServices.NamecheapService.PurchaseDomain(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error purchasing domain"))
		return registerNewDomainResponse, err
	}

	// step 4 - configure domain
	return h.configureDomain(ctx, tenant, domain, website)
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
func (h *MailstackHandler) ConfigureDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "ConfigureDomain", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// Parse and validate request body
		var req ConfigureDomainRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		// Check for missing domain
		if req.Domain == "" {
			message := "Missing required field: domain"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		} else if req.Website == "" {
			message := "Missing required field: website"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		domainResponse, err := h.configureDomain(ctx, tenant, req.Domain, req.Website)
		if err != nil {
			if errors.Is(err, coserrors.ErrDomainNotFound) {
				h.responseHandler.HandleError(c, http.StatusNotFound, nil)
				return
			} else if errors.Is(err, coserrors.ErrDomainConfigurationFailed) {
				message := "Unable to configure domain"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			} else {
				message := "Domain registration failed"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}
		}

		h.responseHandler.HandleSuccess(c, DomainResponse{
			Domain: domainResponse,
		})
	}
}

func (h *MailstackHandler) configureDomain(ctx context.Context, tenant, domain, redirectWebsite string) (DomainRecord, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "configureDomain")
	defer span.Finish()

	domainResponse := DomainRecord{}
	domainResponse.Domain = domain

	var err error

	domainBelongsToTenant, err := h.services.Repositories.PostgresRepositories.MailStackDomainRepository.CheckDomainOwnership(ctx, tenant, domain)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error checking domain"))
		return domainResponse, err
	}
	if !domainBelongsToTenant {
		return domainResponse, coserrors.ErrDomainNotFound
	}

	err = h.services.CommonServices.MailstackService.ConfigureMailstackDomain(ctx, domain, redirectWebsite)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Error configuring domain"))
		return domainResponse, coserrors.ErrDomainConfigurationFailed
	}

	// get domain details
	domainInfo, err := h.services.CommonServices.NamecheapService.GetDomainInfo(ctx, tenant, domain)
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
func (h *MailstackHandler) GetDomains() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetDomains", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// Create request to Mailstack API
		req, err := http.NewRequestWithContext(ctx, "GET", h.services.Cfg.Common.Internal.MailstackApiConfig.ApiUrl+"/v1/domains", nil)
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

		// Create HTTP client with default transport
		client := &http.Client{}

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

			// For 500 errors, use a generic message
			if resp.StatusCode == http.StatusInternalServerError {
				message := "Internal server error"
				tracing.TraceErr(span, errors.New(message))
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}

			// For other errors, propagate the status code and message from Mailstack
			tracing.TraceErr(span, errors.New(errorResponse.Error))
			h.responseHandler.HandleError(c, resp.StatusCode, &errorResponse.Error)
			return
		}

		// Parse response
		var response DomainsResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			message := "Unable to parse Mailstack API response"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		h.responseHandler.HandleSuccess(c, response)
	}
}

func (h *MailstackHandler) RecommendDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RecommendDomain", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// get root domain
		baseName, exists := c.GetQuery("baseName")
		if !exists {
			message := "Must provide baseName"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		// Create request to Mailstack API
		req, err := http.NewRequestWithContext(ctx, "GET", h.services.Cfg.Common.Internal.MailstackApiConfig.ApiUrl+"/v1/domains/recommendations?baseName="+baseName, nil)
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

		// Create HTTP client with default transport
		client := &http.Client{}

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

			// For 500 errors, use a generic message
			if resp.StatusCode == http.StatusInternalServerError {
				message := "Internal server error"
				tracing.TraceErr(span, errors.New(message))
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}

			// For other errors, propagate the status code and message from Mailstack
			tracing.TraceErr(span, errors.New(errorResponse.Error))
			h.responseHandler.HandleError(c, resp.StatusCode, &errorResponse.Error)
			return
		}

		// Parse response
		var recommendations []string
		if err := json.NewDecoder(resp.Body).Decode(&recommendations); err != nil {
			message := "Unable to parse Mailstack API response"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		h.responseHandler.HandleSuccess(c, DomainRecommendationResponse{
			Domains: recommendations,
		})
	}
}
