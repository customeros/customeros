package route

import (
	"context"
	"encoding/json"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
)

func AddSyncEmailRoutes(ctx context.Context, route *gin.Engine, services *service.Services, log logger.Logger, cache *commoncaches.Cache) {
	//route.POST("/sync/emails",
	//	handler.TracingEnhancer(ctx, "/sync/emails"),
	//	security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_WEBHOOKS, security.WithCache(cache)),
	//	syncEmailsHandler(services, log))
	route.POST("/sync/email",
		tracing.TracingEnhancer(ctx, "/sync/email"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_WEBHOOKS, security.WithCache(cache)),
		syncEmailHandler(services, log))
}

//func syncEmailsHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncEmails", c.Request.Header)
//		defer span.Finish()
//
//		// Read the tenant header
//		tenant := c.GetHeader("tenant")
//		if tenant == "" {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
//			return
//		}
//		ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: tenant})
//
//		// Limit the size of the request body
//		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeMessages)
//		requestBody, err := io.ReadAll(c.Request.Body)
//		if err != nil {
//			tracing.TraceErr(span, err)
//			log.Errorf("(SyncEmails) error reading request body: %s", err.Error())
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
//			return
//		}
//
//		// Parse the JSON request body
//		var emails []model.EmailData
//		if err = json.Unmarshal(requestBody, &emails); err != nil {
//			tracing.TraceErr(span, err)
//			log.Errorf("(SyncEmails) Failed unmarshalling body request: %s", err.Error())
//			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
//			return
//		}
//
//		if len(emails) == 0 {
//			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing emails in request"})
//			return
//		}
//
//		// Context timeout, allocate per email
//		timeout := time.Duration(len(emails)) * common.Min1Duration
//		if timeout > constants.RequestMaxTimeout {
//			timeout = constants.RequestMaxTimeout
//		}
//		ctx, cancel := context.WithTimeout(ctx, timeout)
//		defer cancel()
//
//		syncResult, err := services.SyncEmailService.SyncEmail(ctx, emails)
//		if err != nil {
//			tracing.TraceErr(span, err)
//			log.Errorf("(SyncEmails) error in sync emails: %s", err.Error())
//			if errors.IsBadRequest(err) {
//				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//			} else {
//				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing log entries"})
//			}
//		} else {
//			c.JSON(http.StatusOK, syncResult)
//		}
//	}
//}

func syncEmailHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncEmail", c.Request.Header)
		defer span.Finish()

		// Read the tenant header
		tenant := c.GetHeader("tenant")
		if tenant == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or empty tenant header"})
			return
		}
		ctx = common.WithCustomContext(ctx, &common.CustomContext{Tenant: tenant})

		// Limit the size of the request body
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, constants.RequestMaxBodySizeMessages)
		requestBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncEmail) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var email model.EmailData
		if err = json.Unmarshal(requestBody, &email); err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncEmail) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		// Context timeout, allocate per email
		ctx, cancel := context.WithTimeout(ctx, common.Min1Duration)
		defer cancel()

		organizationSync, interactionEventSync, contactSync, err := services.SyncEmailService.SyncEmail(ctx, email)
		//combine the results into a single response
		combinedResults := service.CombinedResults{
			OrganizationSync:     organizationSync,
			InteractionEventSync: interactionEventSync,
			ContactSync:          contactSync,
		}
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncEmail) error in sync email: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing log entry"})
			}
		} else {
			c.JSON(http.StatusOK, combinedResults)
		}
	}
}
