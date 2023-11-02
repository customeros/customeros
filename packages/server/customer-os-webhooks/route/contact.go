package route

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	commonService "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	cosHandler "github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"io"
	"net/http"
	"time"
)

func AddContactRoutes(ctx context.Context, route *gin.Engine, services *service.Services, log logger.Logger) {
	route.POST("/sync/contacts",
		cosHandler.TracingEnhancer(ctx, "/sync/contacts"),
		commonService.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonService.CUSTOMER_OS_WEBHOOKS),
		syncContactsHandler(services, log))
	route.POST("/sync/contact",
		cosHandler.TracingEnhancer(ctx, "/sync/contact"),
		commonService.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonService.CUSTOMER_OS_WEBHOOKS),
		syncContactHandler(services, log))
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
		timeout := time.Duration(len(contacts)) * utils.LongDuration
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
		timeout := utils.LongDuration
		ctx, cancel := context.WithTimeout(ctx, timeout)
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
