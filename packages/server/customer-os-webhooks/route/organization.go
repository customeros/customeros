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

func AddOrganizationRoutes(ctx context.Context, route *gin.Engine, services *service.Services, log logger.Logger) {
	route.POST("/sync/organizations",
		cosHandler.TracingEnhancer(ctx, "/sync/organizations"),
		commonService.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonService.CUSTOMER_OS_WEBHOOKS),
		syncOrganizationsHandler(services, log))
	route.POST("/sync/organization",
		cosHandler.TracingEnhancer(ctx, "/sync/organization"),
		commonService.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonService.CUSTOMER_OS_WEBHOOKS),
		syncOrganizationHandler(services, log))
}

func syncOrganizationsHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncOrganizations", c.Request.Header)
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
			log.Errorf("(SyncOrganizations) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var organizations []model.OrganizationData
		if err = json.Unmarshal(requestBody, &organizations); err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncOrganizations) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		if len(organizations) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing organizations in request"})
			return
		}

		// Context timeout, allocate per organization
		timeout := time.Duration(len(organizations)) * utils.LongDuration
		if timeout > constants.RequestMaxTimeout {
			timeout = constants.RequestMaxTimeout
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		syncResult, err := services.OrganizationService.SyncOrganizations(ctx, organizations)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncOrganizations) error in sync organizations: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing organizations"})
			}
		} else {
			c.JSON(http.StatusOK, syncResult)
		}
	}
}

func syncOrganizationHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncOrganization", c.Request.Header)
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
			log.Errorf("(SyncOrganization) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var organization model.OrganizationData
		if err = json.Unmarshal(requestBody, &organization); err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncOrganization) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		// Context timeout, allocate per organization
		timeout := utils.LongDuration
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		syncResult, err := services.OrganizationService.SyncOrganizations(ctx, []model.OrganizationData{organization})
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncOrganization) error in sync organization: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing organization"})
			}
		} else {
			c.JSON(http.StatusOK, syncResult)
		}
	}
}
