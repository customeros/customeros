package public

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	cosapi_services "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ReferrerQueryParams struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func RevealWebsiteEvents(services *cosapi_services.Services) gin.HandlerFunc {
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
		tenant, err := services.Repositories.PostgresRepositories.TrackingAllowedOriginRepository.GetTenantForOrigin(ctx, origin)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to get tenant for origin"))
			c.JSON(http.StatusForbidden, gin.H{})
			return
		}

		if tenant == "" {
			err = fmt.Errorf("tenant not found for origin")
			tracing.TraceErr(span, err)
			span.LogFields(log.String("result.info", "tenant not found for origin"))
			c.JSON(http.StatusForbidden, gin.H{})
			return
		}

		span.SetTag(tracing.SpanTagTenant, tenant)

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

		// assign event to session (new or existing)
		err = assignEventToSession(ctx, services, trackerData)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		// write event to tracking table
		_, err = services.Repositories.PostgresRepositories.WebTrackerEventsRepository.Create(ctx, *trackerData)
		if err != nil {
			tracing.TraceErr(span, err)
			return
		}

		return
	}

}

func assignEventToSession(ctx context.Context, s *cosapi_services.Services, trackerData *entity.WebTrackerEvents) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "Tracking.assignEventsToSession")
	defer span.Finish()

	// find active session for visitor
	query := entity.WebSession{
		Tenant:    trackerData.Tenant,
		VisitorID: trackerData.VisitorID,
		IsActive:  true,
	}
	session, err := s.Repositories.PostgresRepositories.WebSessionRepository.FindSession(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// if no active session found, create one if event is page_view
	if session == nil && trackerData.EventType == enum.WebTrackerPageView.String() {

		query := entity.WebSession{
			Tenant:        trackerData.Tenant,
			VisitorID:     trackerData.VisitorID,
			IP:            trackerData.IP,
			Referrer:      &trackerData.Referrer,
			StartTime:     utils.Now(),
			LastEventType: trackerData.EventType,
			LastActivity:  utils.Now(),
			IsActive:      true,
		}

		params := parseURLParams(trackerData.Search)
		if params != nil {
			paramString, err := setReferrerQueryParams(ctx, params)
			if err != nil {
				tracing.TraceErr(span, err)
				return err
			}
			if paramString != nil {
				query.QueryParams = paramString
			}
		}

		newSession, err := s.Repositories.PostgresRepositories.WebSessionRepository.Create(ctx, query)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		if newSession == nil {
			err = errors.New("unable to create new web session")
			tracing.TraceErr(span, err)
			return err
		}
		trackerData.SessionID = newSession.ID
		return nil
	}

	trackerData.SessionID = session.ID

	// update existing session last activity
	updateQuery := entity.WebSession{
		ID:           session.ID,
		LastActivity: utils.Now(),
	}
	_, err = s.Repositories.PostgresRepositories.WebSessionRepository.Update(ctx, updateQuery)
	if err != nil {
		err = errors.New("unable to update web session")
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func setReferrerQueryParams(ctx context.Context, queryParams []ReferrerQueryParams) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Tracking.SetReferrerQueryParams")
	defer span.Finish()

	bytes, err := json.Marshal(queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal QueryParam: %w", err)
	}
	results := string(bytes)
	return &results, nil

}

func buildTrackerDbData(c *gin.Context, tenant string) *entity.WebTrackerEvents {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "Tracking.buildTrackerEventData")
	defer span.Finish()
	tracing.TagTenant(span, tenant)

	tracking := entity.WebTrackerEvents{}

	// 1 Get the raw request body
	rawJSON, err := c.GetRawData()
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get raw data"))
		return nil
	}
	span.LogFields(log.String("rawJSON", string(rawJSON)))

	// 2 Unmarshal into a map
	var inputMap map[string]any
	if err := json.Unmarshal(rawJSON, &inputMap); err != nil {
		tracing.TraceErr(span, err)
		return nil
	}

	// 3 Decode using mapstructure-based decode function
	if err := utils.Decode(inputMap, &tracking); err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	tracking.Tenant = tenant
	tracking.UserAgent = utils.SanitizeUTF8(tracking.UserAgent)
	tracking.Referrer = utils.SanitizeUTF8(tracking.Referrer)
	tracking.Origin = utils.SanitizeUTF8(tracking.Origin)
	tracking.Href = utils.SanitizeUTF8(tracking.Href)
	tracking.Search = utils.SanitizeUTF8(tracking.Search)
	tracking.Hostname = utils.SanitizeUTF8(tracking.Hostname)
	tracking.Pathname = utils.SanitizeUTF8(tracking.Pathname)
	return &tracking
}

func isTrustedIP(ctx context.Context, s *cosapi_services.Services, ipAddress string) bool {
	span, ctx := tracing.StartTracerSpan(ctx, "Tracking.isTrustedIp")
	defer span.Finish()

	ipThreats, err := s.CommonServices.VerifyService.Threats(ctx, ipAddress)
	if err != nil || ipThreats == nil {
		tracing.TraceErr(span, err)
		return true
	}

	return !ipThreats.IsThreat
}

func parseURLParams(queryString string) []ReferrerQueryParams {
	// Remove leading ? if present
	queryString = strings.TrimPrefix(queryString, "?")

	// Split the string by & to get individual param-value pairs
	pairs := strings.Split(queryString, "&")

	// Create slice to hold results
	params := make([]ReferrerQueryParams, 0, len(pairs))

	// Parse each pair into the struct
	for _, pair := range pairs {
		// Split pair by = to separate param and value
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			params = append(params, ReferrerQueryParams{
				Name:  parts[0],
				Value: parts[1],
			})
		}
	}

	return params
}
