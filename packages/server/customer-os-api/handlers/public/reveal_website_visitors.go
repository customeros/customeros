package public

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

func RevealWebsiteVisitors(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "public.RevealWebsiteVisitors", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		origin := c.GetHeader("Origin")
		referer := c.GetHeader("Referer")
		userAgent := c.GetHeader("User-Agent")

		span.LogKV("origin", origin)
		span.LogKV("referer", referer)
		span.LogKV("userAgent", userAgent)

		if origin == "" || referer == "" || userAgent == "" {
			err := fmt.Errorf("missing required headers")
			tracing.TraceErr(span, err)
			c.JSON(http.StatusForbidden, gin.H{})
			return
		}

		// tenant validation
		tenant, err := services.CommonServices.PostgresRepositories.TrackingAllowedOriginRepository.GetTenantForOrigin(ctx, origin)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get tenant for origin"))
			c.JSON(http.StatusForbidden, gin.H{})
			return
		}

		if tenant == nil {
			err = fmt.Errorf("tenant not found for origin")
			tracing.TraceErr(span, err)
			span.LogFields(log.String("result.info", "tenant not found for origin"))
			c.JSON(http.StatusForbidden, gin.H{})
			return
		}

		span.SetTag(tracing.SpanTagTenant, *tenant)

		trackerData := buildTrackerDbData(c, tenant)
		if trackerData == nil {
			err = fmt.Errorf("unable to build tracking record")
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("%v", err.Error()),
			})
			return
		}

		// return early, continue processing
		c.JSON(http.StatusAccepted, gin.H{})

		// if bot, return
		if !isTrustedIP(ctx, services, trackerData.IP) {
			return
		}

		// write event to tracking table
		savedRecord, err := services.CommonServices.PostgresRepositories.TrackerEventsRepository.Create(ctx, *trackerData)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		// create and publish event
		eventData := data_fields.WebsiteVisitEvent{
			ID:          savedRecord.ID,
			Tenant:      *tenant,
			IPAddress:   &trackerData.IP,
			VisitorId:   trackerData.VisitorID,
			Website:     trackerData.Hostname,
			PageVisited: trackerData.Pathname,
			Referrer:    trackerData.Referrer,
			Params:      parseURLParams(trackerData.Search),
		}

		webhookEvent := dto.WebhookEvent{
			DataType: eventData.Type(),
			Data:     eventData,
		}

		err = services.CommonServices.RabbitMQService.PublishWebhookEvent(ctx, webhookEvent)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to publish event"))
		}

		return
	}

}

func buildTrackerDbData(c *gin.Context, tenant *string) *entity.TrackerEvents {
	span, _ := tracing.StartTracerSpan(c.Request.Context(), "Tracking.buildTrackerEventData")
	defer span.Finish()

	tracking := entity.TrackerEvents{}

	if err := c.BindJSON(&tracking); err != nil {
		tracing.TraceErr(span, err)
		return nil
	}

	tracking.Tenant = *tenant
	tracking.UserAgent = utils.SanitizeUTF8(tracking.UserAgent)
	tracking.Referrer = utils.SanitizeUTF8(tracking.Referrer)
	tracking.Origin = utils.SanitizeUTF8(tracking.Origin)
	tracking.Href = utils.SanitizeUTF8(tracking.Href)
	tracking.Search = utils.SanitizeUTF8(tracking.Search)
	tracking.Hostname = utils.SanitizeUTF8(tracking.Hostname)
	tracking.Pathname = utils.SanitizeUTF8(tracking.Pathname)
	return &tracking
}

func isTrustedIP(ctx context.Context, s *service.Services, ipAddress string) bool {
	span, ctx := tracing.StartTracerSpan(ctx, "Tracking.isTrustedIp")
	defer span.Finish()

	ipThreats, err := s.CommonServices.VerifyService.Threats(ctx, ipAddress)
	if err != nil || ipThreats == nil {
		tracing.TraceErr(span, err)
		return true
	}

	return !ipThreats.IsThreat
}

func parseURLParams(queryString string) []data_fields.URLParams {
	// Remove leading ? if present
	queryString = strings.TrimPrefix(queryString, "?")

	// Split the string by & to get individual param-value pairs
	pairs := strings.Split(queryString, "&")

	// Create slice to hold results
	params := make([]data_fields.URLParams, 0, len(pairs))

	// Parse each pair into the struct
	for _, pair := range pairs {
		// Split pair by = to separate param and value
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			params = append(params, data_fields.URLParams{
				Param: parts[0],
				Value: parts[1],
			})
		}
	}

	return params
}
