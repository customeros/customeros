package public

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_listeners"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-api/caches"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type WebsiteTrackerEventsHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
	cache           *caches.OriginTenantCache
}

func NewWebsiteTrackerEventsHandler(services *cosapi_services.Services, responseHandler *response.Response) *WebsiteTrackerEventsHandler {
	return &WebsiteTrackerEventsHandler{
		services: services,
		cache:    caches.NewOriginTenantCache(),
	}
}

type ReferrerQueryParams struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (h *WebsiteTrackerEventsHandler) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "WebsiteTrackerEventsHandler.RevealWebsiteVisitors", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		if err := h.validateHeaders(c); err != nil {
			h.responseHandler.HandleError(c, http.StatusForbidden, nil)
			return
		}

		tenant, err := h.validateTrackingAllowed(ctx, c.GetHeader("Origin"))
		if err != nil {
			h.responseHandler.HandleError(c, http.StatusForbidden, nil)
			return
		}
		span.SetTag(tracing.SpanTagTenant, tenant)
		ctx = common.SetTenantInContext(ctx, tenant)

		trackerData := h.buildTrackerDbData(c, tenant)
		if trackerData == nil {
			err = fmt.Errorf("unable to build tracking record")
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		// return early, continue processing
		h.responseHandler.HandleAccepted(c)

		if !h.isTrustedIP(ctx, trackerData.IP) {
			return
		}

		// check if the event type is known
		if !enum.IsValidWebTrackerEvent(trackerData.EventType) {
			err = fmt.Errorf("unsupported web-tracker event type: %s", trackerData.EventType)
			tracing.TraceErr(span, err)
			return
		}

		if err := h.assignEventToSession(ctx, trackerData); err != nil {
			if trackerData.EventType != enum.WebTrackerPageExit.String() {
				tracing.TraceErr(span, err)
			}
			return
		}

		if _, err := h.services.Repositories.PostgresRepositories.WebTrackerEventsRepository.Create(ctx, *trackerData); err != nil {
			tracing.TraceErr(span, err)
			return
		}
	}
}

func (h *WebsiteTrackerEventsHandler) validateHeaders(c *gin.Context) error {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "WebsiteTrackerEventsHandler.validateHeaders")
	defer span.Finish()
	tracing.TagComponentRest(span)

	origin := c.GetHeader("Origin")
	referer := c.GetHeader("Referer")
	userAgent := c.GetHeader("User-Agent")

	span.LogKV("origin", origin)
	span.LogKV("referer", referer)
	span.LogKV("userAgent", userAgent)

	switch {
	case origin == "":
		err := errors.New("missing origin")
		span.LogFields(log.String("result.error", err.Error()))
		return err
	case referer == "":
		err := errors.New("missing referer")
		span.LogFields(log.String("result.error", err.Error()))
		return err
	case userAgent == "":
		err := errors.New("missing userAgent")
		span.LogFields(log.String("result.error", err.Error()))
		return err
	default:
		return nil
	}
}

func (h *WebsiteTrackerEventsHandler) validateTrackingAllowed(ctx context.Context, origin string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.validateTrackingAllowed")
	defer span.Finish()
	tracing.TagComponentRest(span)
	span.LogKV("origin", origin)

	cleanedOrigin := utils.StripUrlToBasePath(origin)

	tenant := h.cache.GetTenantForOrigin(cleanedOrigin)
	if tenant != "" {
		span.LogKV("result.tenant.cached", tenant)
		return tenant, nil
	}

	agents, err := h.services.Repositories.PostgresRepositories.AgentRepository.GetActiveConfiguredAgentsByTypesCrossTenant(ctx, []enum.AgentType{enum.AgentWebVisitorIdentifier})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get agents"))
		return "", err
	}

	// check if agent has intent to identify visitor
	tenant = h.findTenantByOrigin(ctx, agents, cleanedOrigin)

	if tenant == "" {
		err = fmt.Errorf("tenant not found for origin: %s", origin)
		span.LogFields(log.Bool("result.tenant.found", false))
		return "", err
	}
	span.LogKV("result.tenant", tenant)
	return tenant, nil
}

func (h *WebsiteTrackerEventsHandler) findTenantByOrigin(ctx context.Context, agents []postgres_entity.Agent, cleanedOrigin string) string {
	for _, agent := range agents {
		if tenant := h.checkAgentForOrigin(ctx, agent, cleanedOrigin); tenant != "" {
			h.cache.SetTenantForOrigin(cleanedOrigin, tenant)
			return tenant
		}
	}
	return ""
}

func (h *WebsiteTrackerEventsHandler) checkAgentForOrigin(ctx context.Context, agent postgres_entity.Agent, cleanedOrigin string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.checkAgentForOrigin")
	defer span.Finish()
	tracing.TagComponentRest(span)

	for _, listener := range agent.Listeners {
		if listener.Type != enum.EventNewWebSession {
			continue
		}

		var config agent_listeners.IdentifyWebsiteVisitorConfig
		err := listener.GetConfig(&config)
		if err != nil {
			tracing.TraceErr(span, err)
			return ""
		}

		for _, website := range config.Websites.Value {
			if utils.StripUrlToBasePath(website) == cleanedOrigin {
				return agent.Tenant
			}
		}
	}
	return ""
}

func (h *WebsiteTrackerEventsHandler) assignEventToSession(ctx context.Context, trackerData *postgres_entity.WebTrackerEvents) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.assignEventsToSession")
	defer span.Finish()
	tracing.TagComponentRest(span)

	span.LogKV(
		"ip", trackerData.IP,
		"hostname", trackerData.Hostname,
		"visitor_id", trackerData.VisitorID,
		"event_type", trackerData.EventType,
	)

	query := postgres_entity.WebSession{
		Tenant:    trackerData.Tenant,
		IP:        trackerData.IP,
		Hostname:  trackerData.Hostname,
		VisitorID: trackerData.VisitorID,
		IsActive:  true,
	}

	session, err := h.services.Repositories.PostgresRepositories.WebSessionRepository.FindSession(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if session == nil && (trackerData.EventType == enum.WebTrackerPageView.String() || trackerData.EventType == enum.WebTrackerClick.String()) {
		session, err = h.createWebSession(ctx, trackerData)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	if session == nil {
		err = errors.New("session not found and not created")
		if trackerData.EventType != enum.WebTrackerPageExit.String() {
			tracing.TraceErr(span, err)
		}
		return err
	}

	trackerData.SessionID = session.ID
	err = h.updateSessionLastActivity(ctx, trackerData.SessionID, trackerData.EventType)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (h *WebsiteTrackerEventsHandler) updateSessionLastActivity(ctx context.Context, sessionID, eventType string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.updateSessionLastActivityTimestamp")
	defer span.Finish()
	tracing.TagComponentRest(span)
	span.LogKV("sessionID", sessionID, "eventType", eventType)

	_, err := h.services.Repositories.PostgresRepositories.WebSessionRepository.UpdateLastActivity(ctx, sessionID, eventType)
	if err != nil {
		err = errors.New("unable to update web session")
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (h *WebsiteTrackerEventsHandler) createWebSession(ctx context.Context, trackerData *postgres_entity.WebTrackerEvents) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.createWebSession")
	defer span.Finish()
	tracing.TagComponentRest(span)

	query := postgres_entity.WebSession{
		Tenant:        trackerData.Tenant,
		IP:            trackerData.IP,
		VisitorID:     trackerData.VisitorID,
		Hostname:      trackerData.Hostname,
		Referrer:      &trackerData.Referrer,
		StartTime:     utils.Now(),
		LastEventType: trackerData.EventType,
		LastActivity:  utils.Now(),
		IsActive:      true,
	}

	params := h.parseURLParams(trackerData.Search)

	if params != nil {
		paramString, err := h.setReferrerQueryParams(ctx, params)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if paramString != nil {
			query.QueryParams = paramString
		}
	}

	newSession, err := h.services.Repositories.PostgresRepositories.WebSessionRepository.Create(ctx, query)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if newSession == nil {
		err = errors.New("unable to create new web session")
		tracing.TraceErr(span, err)
		return nil, err
	}
	return newSession, nil
}

func (h *WebsiteTrackerEventsHandler) setReferrerQueryParams(ctx context.Context, queryParams []ReferrerQueryParams) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.SetReferrerQueryParams")
	defer span.Finish()

	bytes, err := json.Marshal(queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal QueryParam: %w", err)
	}
	results := string(bytes)
	return &results, nil
}

func (h *WebsiteTrackerEventsHandler) buildTrackerDbData(c *gin.Context, tenant string) *postgres_entity.WebTrackerEvents {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "WebsiteTrackerEventsHandler.buildTrackerEventData")
	defer span.Finish()
	tracing.TagTenant(span, tenant)

	tracking := postgres_entity.WebTrackerEvents{}

	rawJSON, err := c.GetRawData()
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get raw data"))
		return nil
	}
	span.LogFields(log.String("rawJSON", string(rawJSON)))

	var inputMap map[string]any
	if err := json.Unmarshal(rawJSON, &inputMap); err != nil {
		tracing.TraceErr(span, err)
		return nil
	}

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
	tracking.VisitorID = utils.SanitizeUTF8(tracking.VisitorID)

	return &tracking
}

func (h *WebsiteTrackerEventsHandler) isTrustedIP(ctx context.Context, ipAddress string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.isTrustedIp")
	defer span.Finish()
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	span.LogKV("ipAddress", ipAddress)

	ipThreats, err := h.services.CommonServices.VerifyService.Threats(ctx, ipAddress)
	if err != nil || ipThreats == nil {
		tracing.TraceErr(span, err)
		return true
	}

	return !ipThreats.IsThreat
}

func (h *WebsiteTrackerEventsHandler) parseURLParams(queryString string) []ReferrerQueryParams {
	queryString = strings.TrimPrefix(queryString, "?")
	pairs := strings.Split(queryString, "&")
	params := make([]ReferrerQueryParams, 0, len(pairs))

	for _, pair := range pairs {
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
