// @openapi 3.0.0
package mailstack

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
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
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "MailstackHandler.RegisterNewDomain", c.Request.Header)
		defer span.Finish()
		tracing.SetDefaultRestSpanTags(ctx, span)

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
			message := "Invalid request body"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		// Create request to Mailstack API
		jsonBody, err := json.Marshal(req)
		if err != nil {
			message := "Unable to marshal request body"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		// Create request to Mailstack API
		mailstackReq, err := http.NewRequestWithContext(ctx, "POST", h.services.Cfg.Common.Internal.MailstackApiConfig.ApiUrl+"/v1/domains", bytes.NewBuffer(jsonBody))
		if err != nil {
			message := "Unable to create request to Mailstack API"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		// Add required headers
		mailstackReq.Header.Set("Content-Type", "application/json")
		mailstackReq.Header.Set("X-CUSTOMER-OS-API-KEY", h.services.Cfg.Common.Internal.MailstackApiConfig.ApiKey)
		mailstackReq.Header.Set("tenant", tenant)

		// Forward Jaeger trace context
		carrier := opentracing.HTTPHeadersCarrier(mailstackReq.Header)
		err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
		if err != nil {
			span.LogFields(tracingLog.Error(err))
		}

		// Create HTTP client with default transport
		client := &http.Client{}

		// Make request to Mailstack API
		resp, err := client.Do(mailstackReq)
		if err != nil {
			message := "Unable to connect to Mailstack API"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}
		defer resp.Body.Close()

		// Check response status
		if resp.StatusCode != http.StatusCreated {
			// Read error response body
			var errorResponse struct {
				Error string `json:"error"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
				errorResponse.Error = "Unknown error occurred"
			}
			tracing.TraceErr(span, errors.New(errorResponse.Error))

			// For 500 errors, use a generic message
			if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
				message := "Internal server error"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}

			// For other errors, propagate the status code and message from Mailstack
			h.responseHandler.HandleError(c, resp.StatusCode, &errorResponse.Error)
			return
		}

		// Parse response
		var response DomainResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			message := "Unable to parse Mailstack API response"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		h.responseHandler.HandleSuccess(c, response)
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
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "MailstackHandler.ConfigureDomain", c.Request.Header)
		defer span.Finish()
		tracing.SetDefaultRestSpanTags(ctx, span)

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
			message := "Invalid request body"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		// Create request to Mailstack API
		jsonBody, err := json.Marshal(req)
		if err != nil {
			message := "Unable to marshal request body"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		// Create request to Mailstack API
		mailstackReq, err := http.NewRequestWithContext(ctx, "POST", h.services.Cfg.Common.Internal.MailstackApiConfig.ApiUrl+"/v1/domains/configure", bytes.NewBuffer(jsonBody))
		if err != nil {
			message := "Unable to create request to Mailstack API"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		// Add required headers
		mailstackReq.Header.Set("Content-Type", "application/json")
		mailstackReq.Header.Set("X-CUSTOMER-OS-API-KEY", h.services.Cfg.Common.Internal.MailstackApiConfig.ApiKey)
		mailstackReq.Header.Set("tenant", tenant)

		// Forward Jaeger trace context
		carrier := opentracing.HTTPHeadersCarrier(mailstackReq.Header)
		err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
		if err != nil {
			span.LogFields(tracingLog.Error(err))
		}

		// Create HTTP client with default transport
		client := &http.Client{}

		// Make request to Mailstack API
		resp, err := client.Do(mailstackReq)
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
			tracing.TraceErr(span, errors.New(errorResponse.Error))

			// For 500 errors, use a generic message
			if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
				message := "Internal server error"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}

			// For other errors, propagate the status code and message from Mailstack
			h.responseHandler.HandleError(c, resp.StatusCode, &errorResponse.Error)
			return
		}

		// Parse response
		var response DomainResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			message := "Unable to parse Mailstack API response"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		h.responseHandler.HandleSuccess(c, response)
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
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "MailstackHandler.GetDomains", c.Request.Header)
		defer span.Finish()
		tracing.SetDefaultRestSpanTags(ctx, span)

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
			tracing.TraceErr(span, errors.New(errorResponse.Error))

			// For 500 errors, use a generic message
			if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
				message := "Internal server error"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}

			// For other errors, propagate the status code and message from Mailstack
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
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "MailstackHandler.RecommendDomain", c.Request.Header)
		defer span.Finish()
		tracing.SetDefaultRestSpanTags(ctx, span)

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
			tracing.TraceErr(span, errors.New(errorResponse.Error))

			// For 500 errors, use a generic message
			if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusUnauthorized {
				message := "Internal server error"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}

			// For other errors, propagate the status code and message from Mailstack
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
