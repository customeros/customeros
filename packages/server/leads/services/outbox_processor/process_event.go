package outbox_processor

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"

	"github.com/customeros/customeros/packages/server/leads/internal/models"
	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
)

func (s *OutboxProcessor) processEvent(ctx context.Context, event *models.OutboxEvent) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "OutboxProcessor.processEvent")
	defer span.Finish()

	switch {
	case strings.HasPrefix(event.EventType.String(), "webtracker"):
		return s.processWebTrackerEvent(ctx, event)

	case strings.HasPrefix(event.EventType.String(), "proxy"):
		// TODO
		return nil

	case strings.HasPrefix(event.EventType.String(), "lead"):
		// TODO
		return nil

	default:
		err := errors.New("event type not implemented")
		span.TraceError(err)
		return err
	}
}

func (s *OutboxProcessor) processWebTrackerEvent(ctx context.Context, event *models.OutboxEvent) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "OutboxProcessor.processOutboxEvent")
	defer span.Finish()

	// write event to data warehouse
	eventLog := &models.WebTrackerEvent{
		ID:        event.ID,
		Event:     event.EventType,
		Publisher: event.Publisher,
		Timestamp: event.CreatedAt,
		Tenant:    event.Tenant,
		TrackerID: event.EntityID,
		SessionID: event.SessionID,
		Payload:   event.Payload,
	}

	err := s.repositories.WebTrackerEvent.Create(ctx, eventLog)
	if err != nil {
		span.TraceError(err)
		return err
	}

	// determine if event needs to be published
	if !strings.HasPrefix(event.EventType.String(), "webtracker.event") {
		err := s.publishEvent(ctx, event)
		if err != nil {
			span.TraceError(err)
			return err
		}
	}

	return nil
}

func (s *OutboxProcessor) publishEvent(ctx context.Context, event *models.OutboxEvent) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "OutboxProcessor.publishEvent")
	defer span.Finish()

	// Create message with headers
	msg := nats.NewMsg(event.EventType.String())
	msg.Data = event.Payload
	msg.Header.Set(nats_internal.HEADER_TENANT, utils.GetTenantFromContext(ctx))
	msg.Header.Set(nats_internal.HEADER_USERID, utils.GetUserIdFromContext(ctx))

	// Publish to the stored subject
	_, err := s.natsConn.JS.PublishMsg(msg)
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to publish stored email: %w", err)
	}

	return nil
}
