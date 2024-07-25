package route

import (
	"context"
	"github.com/gin-gonic/gin"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/email-tracking-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/email-tracking-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/email-tracking-api/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"net/http"
)

func AddEmailTrackRoute(ctx context.Context, route *gin.Engine, log logger.Logger, commonServices *commonservice.Services) {
	route.GET("/v1/l",
		handler.TracingEnhancer(ctx, "/v1/l"),
		handleLinkRequest(ctx, commonServices, log))
	route.GET("/v1/s",
		handler.TracingEnhancer(ctx, "/v1/s"),
		handleTrackRequest(ctx, commonServices, log))
}

func handleLinkRequest(ctx context.Context, commonServices *commonservice.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(ctx, "HandleLinkRequest")
		defer span.Finish()

		// Extract the 'c' query parameter
		emailLookupId := c.Query("c")
		if emailLookupId == "" {
			c.String(http.StatusBadRequest, "Missing required parameter")
			return
		}

		// Check email lookup id
		emailLookup, err := commonServices.PostgresRepositories.EmailLookupRepository.GetById(ctx, emailLookupId)
		if err != nil {
			log.Error(ctx, "Error retrieving email lookup", err)
			c.String(http.StatusInternalServerError, "An error occurred")
			return
		}
		// if id not found, return 404
		if emailLookup == nil {
			c.String(http.StatusNotFound, "Not found")
			return
		}
		// if email lookup is not of expected type, return 400
		if emailLookup.Type != postgresentity.EmailLookupTypeLink {
			log.Error(ctx, "Email lookup is not of expected type")
			c.String(http.StatusNotFound, "Not found")
			return
		}

		// Get IP address and email address
		ipAddress := c.ClientIP()

		// Store click data
		_, err = commonServices.PostgresRepositories.EmailTrackingRepository.Register(ctx, postgresentity.EmailTracking{
			Tenant:    emailLookup.Tenant,
			MessageId: emailLookup.MessageId,
			LinkId:    emailLookup.LinkId,
			EventType: postgresentity.EmailTrackingEventTypeLinkClick,
			IP:        ipAddress,
		})
		if err != nil {
			log.Error(ctx, "Error storing click data", err)
		}

		// Redirect to the specified URL
		c.Redirect(http.StatusFound, emailLookup.RedirectUrl)
	}
}

func handleTrackRequest(ctx context.Context, commonServices *commonservice.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(ctx, "HandleTrackRequest")
		defer span.Finish()

		// Extract the 'c' query parameter
		emailLookupId := c.Query("c")
		if emailLookupId == "" {
			err := errors.New("Missing required parameter")
			tracing.TraceErr(span, err)
			log.Error(ctx, err.Error())
			return
		}

		// Check email lookup id
		emailLookup, err := commonServices.PostgresRepositories.EmailLookupRepository.GetById(ctx, emailLookupId)
		if err != nil {
			log.Error(ctx, "Error retrieving email lookup", err)
			return
		}
		if emailLookup == nil {
			return
		}
		// if email lookup is not of expected type, return 400
		if emailLookup.Type != postgresentity.EmailLookupTypeLink {
			log.Error(ctx, "Email lookup is not of expected type")
			return
		}

		// Get IP address and email address
		ipAddress := c.ClientIP()

		// Store click data
		_, err = commonServices.PostgresRepositories.EmailTrackingRepository.Register(ctx, postgresentity.EmailTracking{
			Tenant:    emailLookup.Tenant,
			MessageId: emailLookup.MessageId,
			EventType: postgresentity.EmailTrackingEventTypeOpen,
			IP:        ipAddress,
		})
		if err != nil {
			log.Error(ctx, "Error storing email open data", err)
		}
	}
}
