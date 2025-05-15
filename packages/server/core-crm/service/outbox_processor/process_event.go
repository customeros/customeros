package outbox_processor

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

func (s *OutboxProcessor) processEvent(ctx context.Context, outboxEvent *postgres_entity.OutboxEvent) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "OutboxProcessor.processEvent")
	defer span.Finish()
	span.TagEventType(outboxEvent.EventType.String())
	span.TagEntity(outboxEvent.ID)

	switch {
	//case strings.HasPrefix(event.EventType.String(), "webtracker"):
	//	return s.processWebTrackerEvent(ctx, event)

	default:
		err := errors.New("event type not implemented")
		span.TraceError(err)
		return err
	}
}

//
//func (s *OutboxProcessor) processWebscraperEvent(ctx context.Context, event *models.OutboxEvent) error {
//	span, ctx := telemetry.StartServiceSpan(ctx, "OutboxProcessor.processWebscraperEvent")
//	defer span.Finish()
//
//	scrapedEvent := &pb.WebpageScraped{}
//	err := proto.Unmarshal(event.Payload, scrapedEvent)
//	if err != nil {
//		span.TraceError(err)
//		return err
//	}
//
//	eventLog := &models.ScraperEvent{
//		ID:        event.ID,
//		Timestamp: event.CreatedAt,
//		Event:     event.EventType,
//		Publisher: event.Publisher,
//		Domain:    scrapedEvent.Url,
//		Url:       scrapedEvent.Url,
//		Payload:   event.Payload,
//	}
//
//	err = s.repositories.ScraperEvent.Create(ctx, eventLog)
//	if err != nil {
//		span.TraceError(err)
//		return err
//	}
//
//	err = s.publishEvent(ctx, event)
//	if err != nil {
//		span.TraceError(err)
//		return err
//	}
//	return nil
//}
//
//func (s *OutboxProcessor) processWebTrackerEvent(ctx context.Context, event *models.OutboxEvent) error {
//	span, ctx := telemetry.StartServiceSpan(ctx, "OutboxProcessor.processWebTrackerEvent")
//	defer span.Finish()
//	span.TagEventType(event.EventType.String())
//	span.TagEntity(event.ID)
//
//	// write event to data warehouse
//	eventLog := &models.WebTrackerEvent{
//		ID:        event.ID,
//		Event:     event.EventType,
//		Publisher: event.Publisher,
//		Timestamp: event.CreatedAt,
//		Tenant:    event.Tenant,
//		TrackerID: event.EntityID,
//		SessionID: event.SessionID,
//		Payload:   event.Payload,
//	}
//
//	err := s.repositories.WebTrackerEvent.Create(ctx, eventLog)
//	if err != nil {
//		span.TraceError(err)
//		return err
//	}
//
//	// determine if event needs to be published
//	if !strings.HasPrefix(event.EventType.String(), "webtracker.event") {
//		err = s.publishEvent(ctx, event)
//		if err != nil {
//			span.TraceError(err)
//			return err
//		}
//	}
//
//	return nil
//}
//
//func (s *OutboxProcessor) publishEvent(ctx context.Context, event *models.OutboxEvent) error {
//	span, ctx := telemetry.StartServiceSpan(ctx, "OutboxProcessor.publishEvent")
//	defer span.Finish()
//	span.TagEventType(event.EventType.String())
//	span.TagEntity(event.ID)
//
//	// Create message with headers
//	msg := nats.NewMsg(event.EventType.String())
//	msg.Data = event.Payload
//	msg.Header.Set(nats_internal.HEADER_TENANT, event.Tenant)
//
//	// Publish to the stored subject
//	_, err := s.natsConn.JS.PublishMsg(msg)
//	if err != nil {
//		span.TraceError(err)
//		return fmt.Errorf("failed to publish outbox event: %w", err)
//	}
//
//	return nil
//}
