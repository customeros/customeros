package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/customeros/customeros/packages/server/leads/dto"
	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
	"github.com/customeros/customeros/packages/server/leads/proto/mappers"
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

const REQUEST_TIMEOUT = 60 * time.Second

func (h *WebsiteEventsHandler) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "WebsiteTrackerEventsHandler.RevealWebsiteVisitors")
		defer spans.Finish()

		if err := h.validateHeaders(c); err != nil {
			spans.TraceError(err)
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		tenant, err := h.validateTrackingAllowed(ctx, c.GetHeader("Origin"))
		if err != nil {
			spans.TraceError(err)
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		ctx = utils.SetTenantInContext(ctx, tenant)

		trackerData := h.parsePayload(c, tenant)
		if trackerData == nil {
			err = fmt.Errorf("unable to build tracking record")
			spans.TraceError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		// check if bot, return early if not trusted IP
		if !h.isTrustedIP(ctx, trackerData.IP) {
			c.JSON(http.StatusAccepted, gin.H{"accepted": "true"})
			return
		}

		// attach to session
		sessionID, err := h.attachToSession(ctx, trackerData)
		if err != nil {
			err = fmt.Errorf("unable to build tracking record")
			spans.TraceError(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		// publish event
		switch trackerData.EventType {
		case enum.WebTrackerPageView:
			err := h.publishPageViewEvent(ctx, trackerData, sessionID)
			if err != nil {
				spans.TraceError(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to process page view event"})
				return
			}
		case enum.WebTrackerPageExit:
			err := h.publishPageExitEvent(ctx, trackerData, sessionID)
			if err != nil {
				spans.TraceError(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to process page exit event"})
				return
			}
		case enum.WebTrackerClick:
			// do nothing for now
		case enum.WebTrackerIdentify:
			err := h.publishIdentifyEvent(ctx, trackerData, sessionID)
			if err != nil {
				spans.TraceError(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to process identify event"})
				return
			}
		default:
			err = errors.New("Unsupported web tracker event")
			spans.TraceError(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"accepted": "true"})
		return
	}
}

func (h *WebsiteEventsHandler) isTrustedIP(ctx context.Context, ipAddress string) (bool, error) {
	span, ctx := telemetry.StartRestSpan(ctx, "WebsiteEventsHandler.isTrustedIP")
	defer span.Finish()

	reqData, err := proto.Marshal(request)
	if err != nil {
		span.TraceError(err)
		return nil, err
	}

	// Send request to service
	msg := nats.NewMsg(enum.EventAskIPData.String())
	msg.Header = nats.Header{
		enum.TENANT_HEADER:  []string{utils.GetTenantFromContext(ctx)},
		enum.USER_ID_HEADER: []string{utils.GetUserIdFromContext(ctx)},
	}
	msg.Data = reqData

	resp, err := s.natsConn.Conn.RequestMsg(msg, REQUEST_TIMEOUT)
	if err != nil {
		span.TraceError(err)
		return false, err
	}

	// Unmarshal response
	response := &pb.EmailClassificationResponse{}
	if err := proto.Unmarshal(resp.Data, response); err != nil {
		span.TraceError(err)
		return false, err
	}

	return !ipThreats.IsThreat
}

func (h *WebsiteEventsHandler) attachToSession(ctx context.Context, event *dto.WebTrackerEvent) (string, error) {
	span, ctx := telemetry.StartRestSpan(ctx, "WebsiteEventsHandler.attachToSession")
	defer span.Finish()

	// check to see if sessionID exists in Nats KV
	sessionID, err := h.natsConn.SessionCache.Get(ctx, event.VisitorID)
	if err != nil {
		span.TraceError(err)
		return "", err
	}
	if sessionID != "" {
		return sessionID, nil
	}

	// generate new ID
	sessionID = generateSessionId()
	err = h.natsConn.SessionCache.Set(ctx, event.VisitorID, sessionID)
	if err != nil {
		span.TraceError(err)
		return "", err
	}
	return sessionID, nil
}

func (h *WebsiteEventsHandler) publishPageViewEvent(ctx context.Context, event *dto.WebTrackerEvent, sessionID string) error {
	span, ctx := telemetry.StartRestSpan(ctx, "WebsiteEventsHandler.publishPageViewEvent")
	defer span.Finish()

	subject := fmt.Sprintf("webtracker.%s.event.page_view", utils.GetTenantFromContext(ctx))
	return h.publishEvent(ctx, event, sessionID, subject)
}

func (h *WebsiteEventsHandler) publishPageExitEvent(ctx context.Context, event *dto.WebTrackerEvent, sessionID string) error {
	span, ctx := telemetry.StartRestSpan(ctx, "WebsiteEventsHandler.publishPageExitEvent")
	defer span.Finish()

	subject := fmt.Sprintf("webtracker.%s.event.page_exit", utils.GetTenantFromContext(ctx))
	return h.publishEvent(ctx, event, sessionID, subject)
}

func (h *WebsiteEventsHandler) publishIdentifyEvent(ctx context.Context, event *dto.WebTrackerEvent, sessionID string) error {
	span, ctx := telemetry.StartRestSpan(ctx, "WebsiteEventsHandler.publishIdentifyEvent")
	defer span.Finish()

	subject := fmt.Sprintf("webtracker.%s.event.identify", utils.GetTenantFromContext(ctx))
	return h.publishEvent(ctx, event, sessionID, subject)
}

func (h *WebsiteEventsHandler) publishEvent(ctx context.Context, event *dto.WebTrackerEvent, sessionID, subject string) error {
	span, ctx := telemetry.StartRestSpan(ctx, "WebsiteEventsHandler.publishEvent")
	defer span.Finish()

	mappedEvent := mappers.ConvertToProtoWebTrackerEvent(event)
	mappedEvent.SessionId = sessionID

	data, err := proto.Marshal(mappedEvent)
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to marshal stored email: %w", err)
	}

	// Create message with headers
	msg := nats.NewMsg(subject)
	msg.Data = data
	msg.Header.Set("X-Tenant", utils.GetTenantFromContext(ctx))
	msg.Header.Set("X-UserId", utils.GetUserIdFromContext(ctx))

	// Publish to the stored subject
	_, err = h.natsConn.JS.PublishMsg(msg)
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to publish page view event: %w", err)
	}

	return nil
}

func (h *WebsiteEventsHandler) validateHeaders(c *gin.Context) error {
	span, _ := telemetry.StartRestSpan(c.Request.Context(), "WebsiteEventsHandler.validateHeaders")
	defer span.Finish()

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
	span, _ := telemetry.StartRestSpan(ctx, "WebsiteEventsHandler.validateTrackingAllowed")
	defer span.Finish()
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
	span, _ := telemetry.StartRestSpan(c.Request.Context(), "WebsiteEventsHandler.checkAgentForOrigin")
	defer span.Finish()

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

func (h *WebsiteEventsHandler) parsePayload(c *gin.Context, tenant string) *dto.WebTrackerEvent {
	span, _ := telemetry.StartRestSpan(c.Request.Context(), "WebsiteEventsHandler.parsePayload")
	defer span.Finish()

	tracking := dto.WebTrackerEvent{}

	rawJSON, err := c.GetRawData()
	if err != nil {
		span.TraceError(err)
		return nil
	}
	span.LogFields(log.String("rawJSON", string(rawJSON)))

	var inputMap map[string]any
	if err := json.Unmarshal(rawJSON, &inputMap); err != nil {
		span.TraceError(err)
		return nil
	}

	if err := utils.Decode(inputMap, &tracking); err != nil {
		span.TraceError(err)
		return nil
	}
	tracking.UserAgent = utils.SanitizeUTF8(tracking.UserAgent)
	tracking.Referrer = utils.SanitizeUTF8(tracking.Referrer)
	tracking.Href = utils.SanitizeUTF8(tracking.Href)
	tracking.VisitorID = utils.SanitizeUTF8(tracking.VisitorID)

	return &tracking
}

func generateSessionId() string {
	return utils.GenerateNanoIDWithPrefix("sid", 21)
}
