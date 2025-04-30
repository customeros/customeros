package web_event_processor

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/customeros/mailsherpa/mailvalidate"
	"google.golang.org/protobuf/proto"

	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/internal/models"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

type IdentifyData struct {
	Tag     string            `json:"tag"`
	ID      string            `json:"id"`
	Classes string            `json:"classes"`
	Email   string            `json:"email"`
	Dataset map[string]string `json:"dataset"`
}

func (s *webEventProcessor) handleIdentifyEvent(ctx context.Context, webtrackerID string, event *pb.WebTrackerEvent) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "webEventProcessor.handleIdentifyEvent")
	defer span.Finish()

	// parse email from event payload
	email, domain, err := s.parseEmailFromIdentifyEvent(ctx, event.EventData)
	if err != nil {
		span.TraceError(err)
		return err
	}
	if email == nil {
		return nil
	}

	// create event
	return s.createIdentifiedVisitorEvent(ctx, webtrackerID, *email, *domain, event)
}

func (s *webEventProcessor) createIdentifiedVisitorEvent(ctx context.Context, webtrackerID, email, domain string, event *pb.WebTrackerEvent) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "webEventProcessor.createIdentifiedVisitorEvent")
	defer span.Finish()

	message := &pb.WebtrackerVisitorIdentified{
		SessionId: event.SessionId,
		TrackerId: webtrackerID,
		VisitorId: event.VisitorId,
		Ip:        event.Ip,
		Domain:    domain,
		Email:     email,
	}

	payload, err := proto.Marshal(message)
	if err != nil {
		span.TraceError(err)
		return err
	}

	return s.repositories.Outbox.Create(ctx, &models.OutboxEvent{
		ID:        utils.GenerateEventID(),
		EventType: enum.EventWebtrackerVisitorIdentified,
		EntityID:  webtrackerID,
		Publisher: enum.WebEventProcessor,
		Tenant:    utils.GetTenantFromContext(ctx),
		SessionID: event.SessionId,
		Payload:   payload,
		Status:    enum.OutboxPending,
		CreatedAt: utils.Now(),
	})
}

func (s *webEventProcessor) parseEmailFromIdentifyEvent(ctx context.Context, eventData string) (*string, *string, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "webEventProcessor.parseEmailFromIdentifyEvent")
	defer span.Finish()

	if eventData == "" {
		return nil, nil, nil
	}

	cleanEvent := strings.TrimSpace(eventData)

	var payload IdentifyData
	err := json.Unmarshal([]byte(cleanEvent), &payload)
	if err != nil {
		span.TraceError(err)
		return nil, nil, err
	}

	if payload.Email == "" {
		return nil, nil, nil
	}

	email := strings.ToLower(strings.TrimSpace(payload.Email))

	emailValidation := mailvalidate.ValidateEmailSyntax(email)
	if !emailValidation.IsValid || emailValidation.IsSystemGenerated || emailValidation.IsRoleAccount {
		return nil, nil, nil
	}

	var domain *string

	if !emailValidation.IsFreeAccount {
		domain = &emailValidation.Domain
	}

	return &emailValidation.CleanEmail, domain, nil
}
