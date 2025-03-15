// @openapi 3.0.0
package mailstack

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	common_srv "github.com/customeros/customeros/packages/server/customer-os-common-module/services/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

// @title MailStack API
// @version 1.0
// @description API for managing mailboxes and domain email configuration
// @BasePath /mailstack/v1

// RegisterNewMailbox creates a new mailbox for a domain
// @Summary Create mailbox
// @Description Creates a new mailbox for the specified domain with optional forwarding and webmail access
// @Tags Mailboxes
// @Accept json
// @Produce json
// @Param domain path string true "Domain name" example(example.com)
// @Param body body MailboxRequest true "Mailbox configuration"
// @Success 200 {object} MailboxResponse "Mailbox created successfully"
// @Success 200 {object} MailboxResponse "Mailbox created successfully with generated password"
// @Failure 400 {object} handlers.ErrorResponse "Invalid request - Missing or invalid parameters"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 404 {object} handlers.ErrorResponse "Domain not found"
// @Failure 409 {object} handlers.ErrorResponse "Mailbox already exists"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /domains/{domain}/mailboxes [post]
// @Security ApiKeyAuth
func (h *MailstackHandler) RegisterNewMailbox() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RegisterNewMailbox", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		// get domain from path
		domain := c.Param("domain")
		if domain == "" {
			message := "Missing parameter: domain"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}
		span.LogKV("request.domain", domain)

		// get tenant from context
		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
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
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		// Get user ID from linked email
		var userId string
		if mailboxRequest.LinkedUser != "" {
			userDbNode, err := h.services.Repositories.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, mailboxRequest.LinkedUser)
			if err != nil {
				message := "Error finding linked user"
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				tracing.TraceErr(span, errors.Wrap(err, "Error finding linked user"))
				return
			}
			if userDbNode != nil {
				userEntity := neo4jmapper.MapDbNodeToUserEntity(userDbNode)
				userId = userEntity.Id
			}
		}

		// Create request body for Mailstack API
		reqBody := struct {
			Username       string   `json:"username"`
			Password       string   `json:"password"`
			Domain         string   `json:"domain"`
			ForwardingTo   []string `json:"forwardingTo"`
			WebmailEnabled bool     `json:"webmailEnabled"`
			UserId         string   `json:"userId"`
		}{
			Username:       mailboxRequest.Username,
			Password:       mailboxRequest.Password,
			Domain:         domain,
			ForwardingTo:   mailboxRequest.ForwardingTo,
			WebmailEnabled: mailboxRequest.WebmailEnabled,
			UserId:         userId,
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			message := "Unable to marshal request body"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		// Create request to Mailstack API
		mailstackReq, err := http.NewRequestWithContext(ctx, "POST", h.services.Cfg.Common.Internal.MailstackApiConfig.ApiUrl+"/v1/mailboxes", bytes.NewBuffer(jsonBody))
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

		// Create HTTP client with default transport and timeout
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

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
		var apiResponse MailboxRecord
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			message := "Unable to parse Mailstack API response"
			tracing.TraceErr(span, errors.Wrap(err, message))
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		// create email node
		emailFields := interfaces.EmailFields{
			Email: apiResponse.Email,
		}
		var linkWith *common_srv.LinkWith
		if userId != "" {
			linkWith = &common_srv.LinkWith{
				Type: model.USER,
				Id:   userId,
			}
		}
		_, err = h.services.CommonServices.EmailService.Merge(ctx, nil, tenant, emailFields, linkWith)
		if err != nil {
			message := "Internal server error"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			tracing.TraceErr(span, errors.Wrap(err, message))
			return
		}

		err = h.services.CommonServices.Events.Publisher.PublishFanoutEvent(ctx, apiResponse.ID, model.MAILBOX, dto.MailstackProvisionMailbox{})
		if err != nil {
			message := "Error provisioning mailbox"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			tracing.TraceErr(span, errors.Wrap(err, "Error provisioning mailbox"))
			return
		}

		mailboxResponse := MailboxResponse{
			Mailbox: apiResponse,
		}

		h.responseHandler.HandleSuccess(c, mailboxResponse)
	}
}

// GetMailboxes retrieves all mailboxes for a domain
// @Summary List mailboxes
// @Description Retrieves all mailboxes configured for the specified domain
// @Tags Mailboxes
// @Accept json
// @Produce json
// @Param domain path string true "Domain name" example(example.com)
// @Success 200 {object} MailboxesResponse "Successfully retrieved mailboxes"
// @Failure 400 {object} handlers.ErrorResponse "Missing domain parameter"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 404 {object} handlers.ErrorResponse "Domain not found"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /domains/{domain}/mailboxes [get]
// @Security ApiKeyAuth
func (h *MailstackHandler) GetMailboxes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GetMailboxes", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		// get domain from path
		domain := c.Param("domain")
		if domain == "" {
			message := "Missing parameter: domain"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}
		span.LogKV("request.domain", domain)

		// get tenant from context
		tenant := common.GetTenantFromContext(ctx)
		// if tenant missing return auth error
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			span.LogFields(tracingLog.String("result", "Missing tenant in context"))
			return
		}

		// get mailboxes for domain from postgres
		mailboxRecords, err := h.services.Repositories.PostgresRepositories.TenantSettingsMailboxRepository.GetAllByDomain(ctx, domain)
		if err != nil {
			message := "Error retrieving mailboxes"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			tracing.TraceErr(span, errors.Wrap(err, "Error retrieving mailboxes"))
			return
		}

		response := MailboxesResponse{
			Mailboxes: make([]MailboxRecord, 0, len(mailboxRecords)),
		}
		for _, mailboxRecord := range mailboxRecords {
			mailboxDetails, err := h.services.CommonServices.OpenSRSService.GetMailboxDetails(ctx, mailboxRecord.MailboxUsername)
			if err != nil {
				message := "Could not get mailbox details"
				tracing.TraceErr(span, errors.Wrap(err, message))
				h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
				return
			}
			response.Mailboxes = append(response.Mailboxes, MailboxRecord{
				Email:             mailboxRecord.MailboxUsername,
				ForwardingEnabled: mailboxDetails.ForwardingEnabled,
				ForwardingTo:      mailboxDetails.ForwardingTo,
				WebmailEnabled:    mailboxDetails.WebmailEnabled,
			})
		}

		h.responseHandler.HandleSuccess(c, response)
	}
}
