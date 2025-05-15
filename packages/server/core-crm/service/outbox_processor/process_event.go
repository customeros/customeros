package outbox_processor

import (
	"context"
	"errors"
	"fmt"
	nats_common "github.com/customeros/customeros/packages/server/customer-os-common-module/nats"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/customeros/packages/server/enums"
	"github.com/nats-io/nats.go"
)

func (s *OutboxProcessor) processEvent(ctx context.Context, outboxEvent *postgres_entity.OutboxEvent) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "OutboxProcessor.processEvent")
	defer span.Finish()
	span.TagEventType(outboxEvent.EventType.String())
	span.TagEntity(outboxEvent.ID)

	switch {
	case outboxEvent.EventType == enums.EventTenantCreated:
		return s.processTenantCreatedEvent(ctx, outboxEvent)
	case outboxEvent.EventType == enums.EventOrganizationCreated:
		return s.processOrganizationCreatedEvent(ctx, outboxEvent)

	default:
		err := errors.New("event type not implemented")
		span.TraceError(err)
		return err
	}
}

func (s *OutboxProcessor) processTenantCreatedEvent(ctx context.Context, event *postgres_entity.OutboxEvent) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "OutboxProcessor.processTenantCreatedEvent")
	defer span.Finish()
	span.TagEventType(event.EventType.String())
	span.TagEntity(event.ID)

	// TODO write event to data warehouse

	return s.publishEvent(ctx, event, enums.StreamTenant)
}

func (s *OutboxProcessor) processOrganizationCreatedEvent(ctx context.Context, event *postgres_entity.OutboxEvent) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "OutboxProcessor.processOrganizationCreatedEvent")
	defer span.Finish()
	span.TagEventType(event.EventType.String())
	span.TagEntity(event.ID)

	// TODO write event to data warehouse

	return s.publishEvent(ctx, event, enums.StreamOrganization)
}

func (s *OutboxProcessor) publishEvent(ctx context.Context, event *postgres_entity.OutboxEvent, stream enums.NatsStream) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "OutboxProcessor.publishEvent")
	defer span.Finish()
	span.TagEventType(event.EventType.String())
	span.TagEntity(event.ID)

	// Create message with headers
	msg := nats.NewMsg(event.EventType.String())
	msg.Data = event.Payload
	msg.Header.Set(string(nats_common.NATS_HEADER_TENANT), event.Tenant)

	natsConn, err := s.natsConns.GetNatsConnection(stream)
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to get NATS connection: %w", err)
	}

	// Publish to the stored subject
	_, err = natsConn.JS.PublishMsg(msg)
	if err != nil {
		span.TraceError(err)
		return fmt.Errorf("failed to publish outbox event: %w", err)
	}

	return nil
}
