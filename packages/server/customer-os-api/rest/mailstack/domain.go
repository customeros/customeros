// @openapi 3.0.0
package mailstack

import (
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "MailstackHandler.RegisterNewDomain")
		defer spans.Finish()

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			spans.TraceError(errors.New("Missing tenant in context"))
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		// Parse and validate request body
		var req RegisterNewDomainRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			message := "Invalid request body"
			spans.TraceError(errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		// Call service to register domain
		statusCode, errorMsg, domainRecord, err := h.services.CommonServices.MailstackService.RegisterNewDomain(ctx, tenant, req.Domain, req.Website)
		if err != nil || statusCode != http.StatusOK || domainRecord == nil {
			if errorMsg != "" {
				h.responseHandler.HandleError(c, statusCode, &errorMsg)
			} else {
				message := "Internal server error"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			}
			return
		}

		h.responseHandler.HandleSuccess(c, DomainResponse{
			Domain: DomainRecord{
				Domain:      domainRecord.Domain,
				CreatedDate: domainRecord.CreatedDate,
				ExpiredDate: domainRecord.ExpiredDate,
				Nameservers: domainRecord.Nameservers,
			},
		})
	}
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
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "MailstackHandler.ConfigureDomain")
		defer spans.Finish()

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			spans.LogKV("result", "Missing tenant in context")
			return
		}

		// Parse and validate request body
		var req ConfigureDomainRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			message := "Invalid request body"
			spans.TraceError(errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		// Call service to configure domain
		statusCode, errorMsg, domainRecord, err := h.services.CommonServices.MailstackService.ConfigureDomain(ctx, tenant, req.Domain, req.Website)
		if err != nil || statusCode != http.StatusOK {
			if errorMsg != "" {
				h.responseHandler.HandleError(c, statusCode, &errorMsg)
			} else {
				message := "Internal server error"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			}
			return
		}

		h.responseHandler.HandleSuccess(c, DomainResponse{
			Domain: DomainRecord{
				Domain:      domainRecord.Domain,
				CreatedDate: domainRecord.CreatedDate,
				ExpiredDate: domainRecord.ExpiredDate,
				Nameservers: domainRecord.Nameservers,
			},
		})
	}
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
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "MailstackHandler.GetDomains")
		defer spans.Finish()

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			spans.LogKV("result", "Missing tenant in context")
			return
		}

		// Call service to get domains
		statusCode, errorMsg, domainRecords, err := h.services.CommonServices.MailstackService.GetDomains(ctx, tenant)
		if err != nil || statusCode != http.StatusOK {
			if errorMsg != "" {
				h.responseHandler.HandleError(c, statusCode, &errorMsg)
			} else {
				message := "Internal server error"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			}
			return
		}

		// Map domain records to response type
		var domains []DomainRecord
		for _, record := range domainRecords {
			domains = append(domains, DomainRecord{
				Domain:      record.Domain,
				CreatedDate: record.CreatedDate,
				ExpiredDate: record.ExpiredDate,
				Nameservers: record.Nameservers,
			})
		}

		h.responseHandler.HandleSuccess(c, DomainsResponse{
			Domains: domains,
		})
	}
}

// RecommendDomain suggests available domain names based on a base name
// @Summary Get domain name suggestions
// @Description Returns a list of available domain name suggestions based on the provided base name
// @Tags Domains
// @Accept json
// @Produce json
// @Param baseName query string true "Base name to generate domain suggestions from"
// @Success 200 {object} DomainRecommendationResponse "Successfully retrieved domain suggestions"
// @Failure 400 {object} handlers.ErrorResponse "Invalid request - Missing required query parameter"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error - Unable to generate suggestions"
// @Router /mailstack/v1/domains/recommendations [get]
// @Security ApiKeyAuth
func (h *MailstackHandler) RecommendDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "MailstackHandler.RecommendDomain")
		defer spans.Finish()

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			spans.LogKV("result", "Missing tenant in context")
			return
		}

		// get root domain
		baseName, exists := c.GetQuery("baseName")
		if !exists {
			message := "Must provide baseName"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		// Call service to get domain recommendations
		statusCode, errorMsg, recommendations, err := h.services.CommonServices.MailstackService.RecommendDomain(ctx, tenant, baseName)
		if err != nil || statusCode != http.StatusOK {
			if errorMsg != "" {
				h.responseHandler.HandleError(c, statusCode, &errorMsg)
			} else {
				message := "Internal server error"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			}
			return
		}

		h.responseHandler.HandleSuccess(c, DomainRecommendationResponse{
			Domains: recommendations,
		})
	}
}
