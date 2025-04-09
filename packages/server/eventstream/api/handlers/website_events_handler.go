package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/customeros/customeros/packages/server/eventstream/internal/enum"
	nats_internal "github.com/customeros/customeros/packages/server/eventstream/internal/nats"
	"github.com/customeros/customeros/packages/server/eventstream/internal/telemetry"
	"github.com/customeros/customeros/packages/server/eventstream/internal/utils"
	"github.com/customeros/customeros/packages/server/eventstream/proto/mappers"
)

type WebsiteEventsHandler struct {
	natsConn *nats_internal.NATSConnections
	cache    *caches.OriginTenantCache
}

func NewWebsiteEventsHandler(natsConn *nats_internal.NATSConnections) *WebsiteTrackerEventsHandler {
	return &WebsiteEventsHandler{
		natsConn: natsConn,
		cache:    caches.NewOriginTenantCache(),
	}
}

type WebTrackerEvent struct {
	ID               string                `json:"id"`
	VisitorID        string                `json:"visitorId"`
	IP               string                `json:"ip" `
	EventType        enum.EventStreamEvent `json:"eventType"`
	EventData        string                `json:"eventData"`
	Timestamp        time.Time             `json:"timestamp"`
	Href             string                `json:"href"`
	Referrer         string                `json:"referrer"`
	UserAgent        string                `json:"userAgent"`
	Language         string                `json:"language"`
	CookiesEnabled   bool                  `json:"cookiesEnabled"`
	ScreenResolution string                `json:"screenResolution"`
}

func (h *WebsiteEventsHandler) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := telemetry.StartRestSpan(c.Request.Context(), "WebsiteTrackerEventsHandler.RevealWebsiteVisitors", c.Request.Header)
		defer span.Close()

		if err := h.validateHeaders(c); err != nil {
			h.responseHandler.HandleError(c, http.StatusForbidden, nil)
			return
		}

		tenant, err := h.validateTrackingAllowed(ctx, c.GetHeader("Origin"))
		if err != nil {
			h.responseHandler.HandleError(c, http.StatusForbidden, nil)
			return
		}
		ctx = common.SetTenantInContext(ctx, tenant)

		trackerData := h.parsePayload(c, tenant)
		if trackerData == nil {
			err = fmt.Errorf("unable to build tracking record")
			span.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		// publish event
		switch trackerData.EventType {
		case enum.WebTrackerPageView:
			err := h.publishPageViewEvent(ctx, trackerData)
		//
		case enum.WebTrackerPageExit:
		//
		case enum.WebTrackerClick:
		//
		case enum.WebTrackerIdentify:
		//
		default:
			// error unsupported
		}

		// return early, continue processing
		h.responseHandler.HandleAccepted(c)
	}
}

func (h *WebsiteEventsHandler) publishPageViewEvent(ctx context.Context, event *WebTrackerEvent) error {
	span, ctx := telemetry.StartRestSpan(ctx, "WebsiteEventsHandler.publishPageViewEvent")
	defer span.Finish()

	data, err := proto.Marshal(mappers.ConvertToProtoWebTrackerEvent(event))
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to marshal stored email: %w", err)
	}

	// Create message with headers
	msg := nats.NewMsg(fmt.Sprintf("eventstream.%s.event.page_view", utils.GetTenantFromContext(ctx)))
	msg.Data = data
	msg.Header.Set("X-Tenant", utils.GetTenantFromContext(ctx))
	msg.Header.Set("X-UserId", utils.GetUserIdFromContext(ctx))

	// Publish to the stored subject
	_, err = h.natsConn.JS.PublishMsg(msg)
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to publish stored email: %w", err)
	}

	return nil
}

func (h *WebsiteEventsHandler) validateHeaders(c *gin.Context) error {
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

func (h *WebsiteEventsHandler) validateTrackingAllowed(ctx context.Context, origin string) (string, error) {
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

func (h *WebsiteEventsHandler) findTenantByOrigin(ctx context.Context, agents []postgres_entity.Agent, cleanedOrigin string) string {
	for _, agent := range agents {
		if tenant := h.checkAgentForOrigin(ctx, agent, cleanedOrigin); tenant != "" {
			h.cache.SetTenantForOrigin(cleanedOrigin, tenant)
			return tenant
		}
	}
	return ""
}

func (h *WebsiteEventsHandler) checkAgentForOrigin(ctx context.Context, agent postgres_entity.Agent, cleanedOrigin string) string {
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

func (h *WebsiteEventsHandler) parsePayload(c *gin.Context, tenant string) *WebTrackerEvent {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "WebsiteTrackerEventsHandler.buildTrackerEventData")
	defer span.Finish()
	tracing.TagTenant(span, tenant)

	tracking := WebTrackerEvent{}

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
