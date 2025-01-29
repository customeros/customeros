package public

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-api/caches"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/agent_capability"
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

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
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusForbidden, nil)
			return
		}
		span.SetTag(tracing.SpanTagTenant, tenant)

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

		if err := h.assignEventToSession(ctx, trackerData); err != nil {
			tracing.TraceErr(span, err)
			return
		}

		if _, err := h.services.Repositories.PostgresRepositories.WebTrackerEventsRepository.Create(ctx, *trackerData); err != nil {
			tracing.TraceErr(span, err)
			return
		}
	}
}

func (h *WebsiteTrackerEventsHandler) validateHeaders(c *gin.Context) error {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "WebsiteTrackerEventsHandler.assignEventsToSession")
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
		tracing.TraceErr(span, err)
		return err
	case referer == "":
		err := errors.New("missing referer")
		tracing.TraceErr(span, err)
		return err
	case userAgent == "":
		err := errors.New("missing userAgent")
		tracing.TraceErr(span, err)
		return err
	default:
		return nil
	}
}

func (h *WebsiteTrackerEventsHandler) validateTrackingAllowed(ctx context.Context, origin string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.assignEventsToSession")
	defer span.Finish()
	tracing.TagComponentRest(span)
	span.LogKV("origin", origin)

	cleanedOrigin := utils.CleanUrlBasePath(origin)

	tenant := h.cache.GetTenantForOrigin(cleanedOrigin)
	if tenant != "" {
		span.LogKV("result.tenant.cached", tenant)
		return tenant, nil
	}

	agents, err := h.services.Repositories.PostgresRepositories.AgentsRepository.GetActiveAgentsByTypesCrossTenant(ctx, []enum.AgentType{enum.AgentVisitorID})
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get agents"))
		return "", err
	}

	// check if agent has intent to identify visitor
	for _, agent := range agents {
		for _, capability := range agent.CapabilitiesConfig.Capabilities {
			if capability.Type == enum.CapabilityIdentifyWebVisitor {
				// unmarshal capability config
				var config agent_capability.IdentifyWebsiteVisitorConfig
				if err = json.Unmarshal([]byte(capability.Config), &config); err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal capability config"))
					return "", err
				}
				for _, website := range config.Websites.Value {
					if utils.CleanUrlBasePath(website) == cleanedOrigin {
						tenant = agent.Tenant
						h.cache.SetTenantForOrigin(cleanedOrigin, tenant)
						break
					}
				}
			}
		}
	}

	if tenant == "" {
		err = fmt.Errorf("tenant not found for origin: %s", origin)
		tracing.TraceErr(span, err)
		return "", err
	}
	span.LogKV("result.tenant", tenant)
	return tenant, nil
}

func (h *WebsiteTrackerEventsHandler) assignEventToSession(ctx context.Context, trackerData *postgres_entity.WebTrackerEvents) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.assignEventsToSession")
	defer span.Finish()
	tracing.TagComponentRest(span)

	query := postgres_entity.WebSession{
		Tenant:    trackerData.Tenant,
		VisitorID: trackerData.VisitorID,
		IsActive:  true,
	}
	session, err := h.services.Repositories.PostgresRepositories.WebSessionRepository.FindSession(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if session == nil && trackerData.EventType == enum.WebTrackerPageView.String() {
		session, err = h.createWebSession(ctx, trackerData)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	trackerData.SessionID = session.ID
	return h.updateSessionLastActivity(ctx, trackerData.SessionID, trackerData.EventType)
}

func (h *WebsiteTrackerEventsHandler) updateSessionLastActivity(ctx context.Context, sessionID, eventType string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.updateSessionLastActivityTimestamp")
	defer span.Finish()
	tracing.TagComponentRest(span)

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
		VisitorID:     trackerData.VisitorID,
		IP:            trackerData.IP,
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
	return &tracking
}

func (h *WebsiteTrackerEventsHandler) isTrustedIP(ctx context.Context, ipAddress string) bool {
	span, ctx := tracing.StartTracerSpan(ctx, "WebsiteTrackerEventsHandler.isTrustedIp")
	defer span.Finish()

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
