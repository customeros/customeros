package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/customeros/mailsherpa/domaincheck"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/user-admin-api/service"
)

type NewOrRepeatVisitor string

const (
	VisitorNew     NewOrRepeatVisitor = "new"
	VisitorRepeat  NewOrRepeatVisitor = "repeat"
	VisitorUnknown NewOrRepeatVisitor = ""
)

type RawTrackerEvent string

const (
	EventPageExit RawTrackerEvent = "page_exit"
	EventPageView RawTrackerEvent = "page_view"
	EventClick    RawTrackerEvent = "click"
)

func addTrackingRoutes(rg *gin.RouterGroup, services *service.Services, config *config.Config) {
	rg.POST("", func(ginContext *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(context.Background(), "Tracking.track", ginContext.Request.Header)
		defer span.Finish()

		origin := ginContext.GetHeader("Origin")
		referer := ginContext.GetHeader("Referer")
		userAgent := ginContext.GetHeader("User-Agent")

		span.LogFields(tracingLog.String("origin", origin))
		span.LogFields(tracingLog.String("referer", referer))
		span.LogFields(tracingLog.String("userAgent", userAgent))

		if origin == "" || referer == "" || userAgent == "" {
			err := fmt.Errorf("missing required headers")
			tracing.TraceErr(span, err)
			ginContext.JSON(http.StatusForbidden, gin.H{})
			return
		}

		// tenant validation
		tenant, err := services.CommonServices.PostgresRepositories.TrackingAllowedOriginRepository.GetTenantForOrigin(ctx, origin)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get tenant for origin"))
			ginContext.JSON(http.StatusForbidden, gin.H{})
			return
		}

		if tenant == nil || *tenant == "" {
			err = fmt.Errorf("tenant not found for origin")
			tracing.TraceErr(span, err)
			span.LogFields(tracingLog.String("result.info", "tenant not found for origin"))
			ginContext.JSON(http.StatusForbidden, gin.H{})
			return
		}

		span.SetTag(tracing.SpanTagTenant, *tenant)

		trackerData := buildTrackerEventData(ginContext, tenant)
		if trackerData == nil {
			err = fmt.Errorf("unable to build tracking record")
			tracing.TraceErr(span, err)
			ginContext.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("%v", err.Error()),
			})
			return
		}

		// return early, continue processing
		ginContext.JSON(http.StatusAccepted, gin.H{})

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

		// try to deanonymize the IP address
		domain, linkedinSlug, err := identifyIP(ctx, services, trackerData.IP)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		if domain == nil {
			return
		}

		// update all IP events w/ ID data
		query := entity.TrackerEvents{
			IP:           trackerData.IP,
			Domain:       domain,
			LinkedinSlug: linkedinSlug,
		}

		_, err = services.CommonServices.PostgresRepositories.TrackerEventsRepository.Update(ctx, query)

		// update all visitorIDs w/ ID data
		query = entity.TrackerEvents{
			VisitorId:    trackerData.VisitorId,
			Domain:       domain,
			LinkedinSlug: linkedinSlug,
		}

		_, err = services.CommonServices.PostgresRepositories.TrackerEventsRepository.UpdateWhereNoCompanyID(ctx, query)

		// create and publish event
		eventData := data_fields.WebsiteVisitEvent{
			ID:       savedRecord.ID,
			Tenant:   *tenant,
			Domain:   *domain,
			Referrer: trackerData.Referrer,
		}

		webhookEvent := dto.WebhookEvent{
			DataType: eventData.Type(),
			Data:     eventData,
		}

		visitorType, err := newOrRepeatVisitor(ctx, services, tenant, domain)
		if err != nil {
			tracing.TraceErr(span, err)
		}
		switch visitorType {
		case VisitorNew:
			webhookEvent.Name = enum.EventRevealWebsiteVisitNew
		case VisitorRepeat:
			webhookEvent.Name = enum.EventRevealWebsiteVisitRepeat
		case VisitorUnknown:
			return
		default:
			err = errors.New("Invalid visitor type")
			tracing.TraceErr(span, err)
			return
		}

		err = services.CommonServices.RabbitMQService.PublishWebhookEvent(ctx, webhookEvent)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to publish event"))
		}

		return
	})
}

func identifyIP(ctx context.Context, s *service.Services, ipAddress string) (domain, linkedinSlug *string, err error) {
	span, ctx := tracing.StartTracerSpan(ctx, "Tracking.identifyIP")
	defer span.Finish()

	snitcherData, err := s.CommonServices.EnrichmentService.Snitcher(ctx, ipAddress)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, nil, err
	}

	if snitcherData == nil || snitcherData.Data == nil {
		return nil, nil, nil
	}

	_, primaryDomain := domaincheck.PrimaryDomainCheck(snitcherData.Data.Domain)

	return &primaryDomain, &snitcherData.Data.Profiles.LinkedIn.Handle, nil
}

func newOrRepeatVisitor(ctx context.Context, s *service.Services, tenant, domain *string) (NewOrRepeatVisitor, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "Tracking.newOrRepeatVisitor")
	defer span.Finish()

	query := entity.TrackerEvents{
		Tenant:    *tenant,
		Domain:    domain,
		EventType: string(EventPageExit),
	}

	results, err := s.CommonServices.PostgresRepositories.TrackerEventsRepository.FindAll(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return VisitorUnknown, err
	}

	if len(*results) == 0 {
		return VisitorNew, nil
	}
	return VisitorRepeat, nil
}

func buildTrackerEventData(c *gin.Context, tenant *string) *entity.TrackerEvents {
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
