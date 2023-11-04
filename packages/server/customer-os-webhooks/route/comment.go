package route

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/errors"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"io"
	"net/http"
	"time"
)

func AddCommentRoutes(ctx context.Context, route *gin.Engine, services *service.Services, log logger.Logger, cache *commoncaches.Cache) {
	route.POST("/sync/comments",
		handler.TracingEnhancer(ctx, "/sync/comments"),
		commonservice.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonservice.CUSTOMER_OS_WEBHOOKS, commonservice.WithCache(cache)),
		syncCommentsHandler(services, log))
	route.POST("/sync/comment",
		handler.TracingEnhancer(ctx, "/sync/comment"),
		commonservice.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonservice.CUSTOMER_OS_WEBHOOKS, commonservice.WithCache(cache)),
		syncCommentHandler(services, log))
}

func syncCommentsHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncComments", c.Request.Header)
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
			log.Errorf("(SyncLogEntries) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var comments []model.CommentData
		if err = json.Unmarshal(requestBody, &comments); err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncLogEntries) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		if len(comments) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing comments in request"})
			return
		}

		// Context timeout, allocate per comment
		timeout := time.Duration(len(comments)) * utils.LongDuration
		if timeout > constants.RequestMaxTimeout {
			timeout = constants.RequestMaxTimeout
		}
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		syncResult, err := services.CommentService.SyncComments(ctx, comments)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncLogEntries) error in sync comments: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing log entries"})
			}
		} else {
			c.JSON(http.StatusOK, syncResult)
		}
	}
}

func syncCommentHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "SyncComment", c.Request.Header)
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
			log.Errorf("(SyncComment) error reading request body: %s", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Parse the JSON request body
		var comment model.CommentData
		if err = json.Unmarshal(requestBody, &comment); err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncLogEntries) Failed unmarshalling body request: %s", err.Error())
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Cannot unmarshal request body"})
			return
		}

		// Context timeout, allocate per comment
		timeout := utils.LongDuration
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		syncResult, err := services.CommentService.SyncComments(ctx, []model.CommentData{comment})
		if err != nil {
			tracing.TraceErr(span, err)
			log.Errorf("(SyncComment) error in sync comment: %s", err.Error())
			if errors.IsBadRequest(err) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed processing log entry"})
			}
		} else {
			c.JSON(http.StatusOK, syncResult)
		}
	}
}
