package mailstack

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type DNSRecord struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type DNSResponse struct {
	Records []DNSRecord `json:"dnsRecords"`
}

type DNSRecordResponse struct {
	Record DNSRecord `json:"dnsRecord"`
}

func (h *MailstackHandler) DNS() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "mailstack.DNS", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			span.LogFields(log.String("result", "Missing tenant in context"))
			return
		}

		domain := c.Param("domain")

		// Create request to Mailstack API
		req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/v1/domains/%s/dns", h.services.Cfg.Common.Internal.MailstackApiConfig.ApiUrl, domain), nil)
		if err != nil {
			message := "Unable to create request to Mailstack API"
			tracing.TraceErr(span, err)
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
			span.LogFields(log.Error(err))
		}

		// Create HTTP client with default transport
		client := &http.Client{}

		// Make request to Mailstack API
		resp, err := client.Do(req)
		if err != nil {
			message := "Unable to connect to Mailstack API"
			tracing.TraceErr(span, err)
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
				tracing.TraceErr(span, err)
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}

			// For other errors, propagate the status code and message from Mailstack
			tracing.TraceErr(span, errors.New(errorResponse.Error))
			h.responseHandler.HandleError(c, resp.StatusCode, &errorResponse.Error)
			return
		}

		// Parse response
		var response DNSResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			message := "Unable to parse Mailstack API response"
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		h.responseHandler.HandleSuccess(c, response)
	}
}

func (h *MailstackHandler) AddDNSRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "mailstack.AddDNSRecord", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			span.LogFields(log.String("result", "Missing tenant in context"))
			return
		}

		domain := c.Param("domain")

		// get dns record payload
		record, err := h.getDNSRequestPayload(c)
		if err != nil {
			message := "Unable to parse request"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		// Create request body
		requestBody, err := json.Marshal(record)
		if err != nil {
			message := "Unable to marshal request"
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		// Create request to Mailstack API
		req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/v1/domains/%s/dns", h.services.Cfg.Common.Internal.MailstackApiConfig.ApiUrl, domain), strings.NewReader(string(requestBody)))
		if err != nil {
			message := "Unable to create request to Mailstack API"
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		// Add required headers
		req.Header.Set("X-CUSTOMER-OS-API-KEY", h.services.Cfg.Common.Internal.MailstackApiConfig.ApiKey)
		req.Header.Set("tenant", tenant)
		req.Header.Set("Content-Type", "application/json")

		// Forward Jaeger trace context
		carrier := opentracing.HTTPHeadersCarrier(req.Header)
		err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
		if err != nil {
			span.LogFields(log.Error(err))
		}

		// Create HTTP client with default transport
		client := &http.Client{}

		// Make request to Mailstack API
		resp, err := client.Do(req)
		if err != nil {
			message := "Unable to connect to Mailstack API"
			tracing.TraceErr(span, err)
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
				tracing.TraceErr(span, err)
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}

			// For other errors, propagate the status code and message from Mailstack
			tracing.TraceErr(span, errors.New(errorResponse.Error))
			h.responseHandler.HandleError(c, resp.StatusCode, &errorResponse.Error)
			return
		}

		// Parse response
		var response DNSRecordResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			message := "Unable to parse Mailstack API response"
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		h.responseHandler.HandleSuccess(c, response)
	}
}

func (h *MailstackHandler) DeleteDNSRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "mailstack.DeleteDNSRecord", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			span.LogFields(log.String("result", "Missing tenant in context"))
			return
		}

		domain := c.Param("domain")
		dnsId := c.Param("dnsId")

		// Create request to Mailstack API
		req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/v1/domains/%s/dns/%s", h.services.Cfg.Common.Internal.MailstackApiConfig.ApiUrl, domain, dnsId), nil)
		if err != nil {
			message := "Unable to create request to Mailstack API"
			tracing.TraceErr(span, err)
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
			span.LogFields(log.Error(err))
		}

		// Create HTTP client with default transport
		client := &http.Client{}

		// Make request to Mailstack API
		resp, err := client.Do(req)
		if err != nil {
			message := "Unable to connect to Mailstack API"
			tracing.TraceErr(span, err)
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
				tracing.TraceErr(span, err)
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}

			// For other errors, propagate the status code and message from Mailstack
			tracing.TraceErr(span, errors.New(errorResponse.Error))
			h.responseHandler.HandleError(c, resp.StatusCode, &errorResponse.Error)
			return
		}

		h.responseHandler.HandleSuccess(c, nil)
	}
}

func (h *MailstackHandler) getDNSRequestPayload(c *gin.Context) (DNSRecord, error) {
	span, _ := tracing.StartTracerSpan(c.Request.Context(), "Flows.getDNSRequestPayload")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var req DNSRecord
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return req, err
	}

	return req, nil
}
