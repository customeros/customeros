package route

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/config"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
)

func AddContactRoutes(ctx context.Context, route *gin.Engine, cfg *config.Config, services *service.Services, log logger.Logger, cache *commoncaches.Cache) {
	route.POST("/sync/contacts",
		handler.TracingEnhancer(ctx, "/sync/contacts"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_WEBHOOKS, security.WithCache(cache)),
		syncContactsHandler(services, log))
	route.POST("/sync/contact",
		handler.TracingEnhancer(ctx, "/sync/contact"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_WEBHOOKS, security.WithCache(cache)),
		syncContactHandler(services, log))
	route.POST("/sync/better-contact",
		handler.TracingEnhancer(ctx, "/sync/better-contact"),
		syncBetterContactResponse(cfg, services, log))
}

func syncContactsHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncContacts", c.Request.Header)
		defer span.Finish()

		// Read the tenant header
		tenant := c.GetHeader("tenant")
		if tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
			return
		}
		ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: tenant})

		// Limit the size of the request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeCommon)
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncContacts) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var contacts []model.ContactData
		if err = json.Unmarshal(requestBody, &contacts); err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncContacts) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		if len(contacts) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing contacts in request"})
			return
		}

		// Context timeout, allocate per contact
		timeout := time.Duration(len(contacts)) * common.Min1Duration
		if timeout > constants.RequestMaxTimeout {
			timeout = constants.RequestMaxTimeout
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		syncResult, err := services.ContactService.SyncContacts(ctx, contacts)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncContacts) error in sync contacts: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing contacts"})
			}
		} else {
			c.JSON(http.StatusOK, syncResult)
		}
	}
}

func syncContactHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncContact", c.Request.Header)
		defer span.Finish()

		// Read the tenant header
		tenant := c.GetHeader("tenant")
		if tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
			return
		}
		ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: tenant})

		// Limit the size of the request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeCommon)
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncContact) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var contact model.ContactData
		if err = json.Unmarshal(requestBody, &contact); err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncContact) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		// Context timeout, allocate per contact
		ctx, cancel := context.WithTimeout(ctx, common.Min1Duration)
		defer cancel()

		syncResult, err := services.ContactService.SyncContacts(ctx, []model.ContactData{contact})
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncContact) error in sync contact: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing contact"})
			}
		} else {
			c.JSON(http.StatusOK, syncResult)
		}
	}
}

func syncBetterContactResponse(cfg *config.Config, services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncBetterContact", c.Request.Header)
		defer span.Finish()

		// Read the tenant header
		apiKeyHeader := c.Query("apiKey")
		if apiKeyHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
			return
		}

		if apiKeyHeader != cfg.BetterContactCallbackApiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}

		// Limit the size of the request body
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var betterContactResponse entity.BetterContactResponseBody
		if err = json.Unmarshal(requestBody, &betterContactResponse); err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		ctx, cancel := context.WithTimeout(ctx, common.Min1Duration)
		defer cancel()

		err = services.CommonServices.PostgresRepositories.EnrichDetailsBetterContactRepository.AddResponse(ctx, betterContactResponse.Id, string(requestBody))
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing better contact response"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	}
}
