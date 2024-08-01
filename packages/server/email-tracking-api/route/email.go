package route

import (
	"context"
	"github.com/gin-gonic/gin"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/email-tracking-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/email-tracking-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/email-tracking-api/tracing"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func AddEmailTrackRoute(ctx context.Context, route *gin.Engine, log logger.Logger, commonServices *commonservice.Services) {
	route.GET("/v1/l",
		handler.TracingEnhancer(ctx, "/v1/l"),
		trackLinkRequest(commonServices, log))
	route.GET("/v1/s",
		handler.TracingEnhancer(ctx, "/v1/s"),
		trackOpenRequest(commonServices, log))
	route.GET("/v1/u",
		handler.TracingEnhancer(ctx, "/v1/u"),
		trackUnsubscribeRequest(commonServices, log))
}

func trackLinkRequest(commonServices *commonservice.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "trackLinkRequest", c.Request.Header)
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
		tracing.LogObjectAsJson(span, "emailLookup", emailLookup)

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

		if emailLookup.TrackClicks {
			// Get IP address and email address
			ipAddress := c.ClientIP()

			// Store click data
			_, err = commonServices.PostgresRepositories.EmailTrackingRepository.Register(ctx, postgresentity.EmailTracking{
				Tenant:      emailLookup.Tenant,
				MessageId:   emailLookup.MessageId,
				LinkId:      emailLookup.LinkId,
				EventType:   postgresentity.EmailTrackingEventTypeLinkClick,
				IP:          ipAddress,
				RecipientId: emailLookup.RecipientId,
				Campaign:    emailLookup.Campaign,
			})
			if err != nil {
				log.Error(ctx, "Error storing click data", err)
			}
		}

		// Redirect to the specified URL
		c.Redirect(http.StatusFound, ensureAbsoluteURL(emailLookup.RedirectUrl))
	}
}

func ensureAbsoluteURL(url string) string {
	// Check if the URL already has a scheme
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}

	// If not, prepend "https://"
	return "https://" + url
}

func trackOpenRequest(commonServices *commonservice.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "trackOpenRequest", c.Request.Header)
		defer span.Finish()

		// Extract the 'c' query parameter
		emailLookupId := c.Query("c")
		span.LogFields(tracingLog.String("emailLookupId", emailLookupId))
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
		tracing.LogObjectAsJson(span, "emailLookup", emailLookup)

		if emailLookup == nil {
			return
		}
		// if email lookup is not of expected type, return 400
		if emailLookup.Type != postgresentity.EmailLookupTypeSpyPixel {
			log.Error(ctx, "Email lookup is not of expected type")
			return
		}
		if emailLookup.TrackOpens {
			// Get IP address and email address
			ipAddress := c.ClientIP()

			// Store click data
			_, err = commonServices.PostgresRepositories.EmailTrackingRepository.Register(ctx, postgresentity.EmailTracking{
				Tenant:      emailLookup.Tenant,
				MessageId:   emailLookup.MessageId,
				RecipientId: emailLookup.RecipientId,
				Campaign:    emailLookup.Campaign,
				EventType:   postgresentity.EmailTrackingEventTypeOpen,
				IP:          ipAddress,
			})
			if err != nil {
				log.Error(ctx, "Error storing email open data", err)
			}
		}
	}
}

func trackUnsubscribeRequest(commonServices *commonservice.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "trackUnsubscribeRequest", c.Request.Header)
		defer span.Finish()

		// Extract the 'c' query parameter
		emailLookupId := c.Query("u")
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
		tracing.LogObjectAsJson(span, "emailLookup", emailLookup)

		// if id not found, return 404
		if emailLookup == nil {
			c.String(http.StatusNotFound, "Not found")
			return
		}
		// if email lookup is not of expected type, return 400
		if emailLookup.Type != postgresentity.EmailLookupTypeUnsubscribe {
			log.Error(ctx, "Email lookup is not of expected type")
			c.String(http.StatusNotFound, "Not found")
			return
		}

		// Get IP address and email address
		ipAddress := c.ClientIP()

		// Store click data
		_, err = commonServices.PostgresRepositories.EmailTrackingRepository.Register(ctx, postgresentity.EmailTracking{
			Tenant:      emailLookup.Tenant,
			MessageId:   emailLookup.MessageId,
			LinkId:      emailLookup.LinkId,
			EventType:   postgresentity.EmailTrackingEventTypeUnsubscribe,
			IP:          ipAddress,
			RecipientId: emailLookup.RecipientId,
			Campaign:    emailLookup.Campaign,
		})
		if err != nil {
			log.Error(ctx, "Error storing unsubscribe data", err)
		}

		// Redirect to the specified URL
		c.Redirect(http.StatusFound, ensureAbsoluteURL(emailLookup.UnsubscribeUrl))
	}
}
