package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/leads/dto"
	"github.com/customeros/customeros/packages/server/leads/internal/caches"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
	"github.com/customeros/customeros/packages/server/leads/services"
)

type WebsiteEventsHandler struct {
	cache    *caches.OriginTenantCache
	services *services.Services
}

func NewWebsiteEventsHandler(services *services.Services) *WebsiteEventsHandler {
	return &WebsiteEventsHandler{
		cache:    caches.NewOriginTenantCache(),
		services: services,
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

		tenant, webtrackerID, err := h.getTenantAndTrackerID(ctx, c.GetHeader("Origin"))
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

		c.JSON(http.StatusAccepted, gin.H{"accepted": "true"})
		h.services.WebEventProcessor.Process(ctx, trackerData, webtrackerID)
		return
	}
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

func (h *WebsiteEventsHandler) getTenantAndTrackerID(ctx context.Context, origin string) (string, string, error) {
	span, _ := telemetry.StartRestSpan(ctx, "WebsiteEventsHandler.getTenantAndTrackerID")
	defer span.Finish()
	span.LogKV("origin", origin)

	cleanedOrigin := utils.StripUrlToBasePath(origin)

	tenant, webtrackerID, err := h.cache.GetDataForOrigin(cleanedOrigin)
	if tenant != "" {
		span.LogKV("result.tenant.cached", tenant)
		return tenant, webtrackerID, nil
	}

	tenant, webtrackerID, err = h.findTrackerIDForOrigin(ctx, cleanedOrigin)
	if err != nil {
		span.TraceError(err)
		return "", "", err
	}

	return tenant, webtrackerID, err
}

func (h *WebsiteEventsHandler) findTrackerIDForOrigin(ctx context.Context, cleanedOrigin string) (string, string, error) {
	span, _ := telemetry.StartRestSpan(ctx, "WebsiteEventsHandler.findTrackerIDForOrigin")
	defer span.Finish()

	webtracker, err := h.services.WebtrackerService.GetWebtrackerByOrigin(ctx, cleanedOrigin)
	if err != nil {
		span.TraceError(err)
		return "", "", err
	}

	err = h.cache.SetDataForOrigin(cleanedOrigin, webtracker.Tenant, webtracker.ID)
	if err != nil {
		span.TraceError(err)
	}

	return webtracker.Tenant, webtracker.ID, nil
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
