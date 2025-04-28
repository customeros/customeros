package web_event_processor

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"

	leads_errors "github.com/customeros/customeros/packages/server/leads/errors"
	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/internal/models"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
	proto_mappers "github.com/customeros/customeros/packages/server/leads/proto/mappers"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

const IP_DATA_LOOKBACK_DAYS = -90 // days

func (s *webEventProcessor) Process(ctx context.Context, webtrackerID string, event *pb.WebTrackerEvent) {
	span, ctx := telemetry.StartServiceSpan(ctx, "webEventProcessor.Process")
	defer span.Finish()

	tenant := utils.GetTenantFromContext(ctx)
	if tenant == "" {
		span.TraceError(leads_errors.ErrTenantMissing)
		return
	}

	// get SessionID
	sessionID, err := s.natsConn.SessionCache.Get(ctx, event.VisitorId)
	if err != nil {
		span.TraceError(err)
		return
	}
	if sessionID == "" {
		sessionID, err = s.newSession(ctx, webtrackerID, event)
		if err != nil {
			span.TraceError(err)
			return
		}
	}

	payload, err := proto.Marshal(event)
	if err != nil {
		span.TraceError(err)
		return
	}

	// update lastEventAt on webtracker && write webtracker event
	err = s.logWebtrackerEvent(ctx, webtrackerID, &models.WebTrackerEvent{
		ID:        utils.GenerateNanoIDWithPrefix("wevt", 16),
		Event:     proto_mappers.ConvertFromProtoEventType(event.EventType),
		Publisher: enum.WebEventProcessor,
		Timestamp: event.Timestamp.AsTime(),
		Tenant:    tenant,
		SessionID: sessionID,
		Payload:   payload,
	})

	return
}

func (s *webEventProcessor) newSession(ctx context.Context, webtrackerID string, event *pb.WebTrackerEvent) (string, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "webEventProcessor.newSession")
	defer span.Finish()

	sessionID, err := s.createNewSession(ctx, webtrackerID, event)
	if err != nil {
		span.TraceError(err)
		return "", err
	}

	webSessionEvent := &pb.WebtrackerSessionNew{
		SessionId: sessionID,
		TrackerId: webtrackerID,
		VisitorId: event.VisitorId,
		Ip:        event.Ip,
		EventType: event.EventType,
		EventData: event.EventData,
		Timestamp: event.Timestamp,
		Href:      event.Href,
		Referrer:  event.Referrer,
		UserAgent: event.UserAgent,
		Language:  event.Language,
	}

	payload, err := proto.Marshal(webSessionEvent)
	if err != nil {
		span.TraceError(err)
		return "", err
	}

	// Create outbox entry
	outboxEvent := &models.OutboxEvent{
		ID:        utils.GenerateEventID(),
		EventType: enum.EventWebtrackerSessionNew,
		Payload:   payload,
		Status:    enum.OutboxPending,
		CreatedAt: time.Now(),
	}

	err = s.repositories.Outbox.Create(ctx, outboxEvent)
	if err != nil {
		span.TraceError(err)
		return "", err
	}
	return sessionID, nil
}

func (s *webEventProcessor) logWebtrackerEvent(ctx context.Context, webtrackerID string, event *models.WebTrackerEvent) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "webEventProcessor.logWebtrackerEvent")
	defer span.Finish()

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Update the WebTracker's last event timestamp
		err := s.repositories.WebTracker.UpdateLastEventAtWithTxn(ctx, tx, webtrackerID, event.Timestamp)
		if err != nil {
			span.TraceError(err)
			return err
		}

		// Create an outbox entry to log the event in TimescaleDB
		outboxEvent := &models.OutboxEvent{
			ID:        utils.GenerateEventID(),
			EventType: mapWebtrackerToLeadEvent(event.Event),
			Payload:   event.Payload,
			Status:    enum.OutboxPending,
			CreatedAt: time.Now(),
		}

		err = s.repositories.Outbox.CreateWithTxn(ctx, tx, outboxEvent)
		if err != nil {
			span.TraceError(err)
			return err
		}

		return nil
	})
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to update tracker and create outbox entry: %w", err)
	}

	return nil
}

func (s *webEventProcessor) createNewSession(ctx context.Context, webtrackerID string, event *pb.WebTrackerEvent) (string, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "webEventProcessor.attachToSession")
	defer span.Finish()

	// generate new ID
	sessionID := utils.GenerateNanoIDWithPrefix("sess", 21)
	err := s.natsConn.SessionCache.Set(ctx, event.VisitorId, sessionID)
	if err != nil {
		span.TraceError(err)
		return "", err
	}

	return sessionID, nil
}

func mapWebtrackerToLeadEvent(e enum.WebTrackerEvent) enum.Events {
	switch {
	case e == enum.WebTrackerPageView:
		return enum.EventWebtrackerPageView
	case e == enum.WebTrackerPageExit:
		return enum.EventWebtrackerPageExit
	case e == enum.WebTrackerClick:
		return enum.EventWebtrackerClick
	case e == enum.WebTrackerIdentify:
		return enum.EventWebtrackerVisitorIdentified
	default:
		return ""
	}
}
