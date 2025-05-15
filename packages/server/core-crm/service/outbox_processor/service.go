package outbox_processor

import (
	"context"
	"github.com/customeros/customeros/packages/server/core-crm/internal/utils"
	nats_common "github.com/customeros/customeros/packages/server/customer-os-common-module/nats"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"time"
)

type OutboxProcessor struct {
	natsConns *nats_common.NATSConnections
	postgres  *postgres_repository.Repositories
}

func NewOutboxProcessor(
	natsConns *nats_common.NATSConnections,
	postgres *postgres_repository.Repositories,
) *OutboxProcessor {
	return &OutboxProcessor{
		natsConns: natsConns,
		postgres:  postgres,
	}
}

const (
	MAX_MESSAGES_PER_BATCH = 50
	EVENT_LOCK_TIMEOUT     = 3 * time.Minute
)

func (p *OutboxProcessor) ProcessBatch(ctx context.Context) error {
	span, ctx := telemetry.StartCronSpan(ctx, "OutboxProcessor.ProcessBatch", telemetry.WithNewRoot())
	defer span.Finish()

	// Get pending events
	events, err := p.postgres.OutboxRepository.GetPendingEvents(ctx, MAX_MESSAGES_PER_BATCH)
	if err != nil {
		span.TraceError(err)
		return err
	}
	span.LogKV("events.count", len(events))

	for _, event := range events {
		innerCtx := utils.WithCustomContext(ctx, &utils.CustomContext{
			Tenant: event.Tenant,
		})
		p.ProcessOutboxEvent(innerCtx, event)
	}

	return nil
}

func (p *OutboxProcessor) ProcessOutboxEvent(ctx context.Context, event *postgres_entity.OutboxEvent) {
	span, ctx := telemetry.StartServiceSpan(ctx, "OutboxProcessor.processOutboxEvent")
	defer span.Finish()
	span.TagEventType(event.EventType.String())
	span.TagEntity(event.ID)

	// Lock the event
	err := p.postgres.OutboxRepository.MarkAsProcessing(ctx, event.ID, EVENT_LOCK_TIMEOUT)
	if err != nil {
		// Another worker might have picked it up
		return
	}

	// Process the event
	err = p.processEvent(ctx, event)
	if err != nil {
		// Mark as failed and increment retry count
		_ = p.postgres.OutboxRepository.MarkAsFailed(ctx, event.ID, err.Error())
		_ = p.postgres.OutboxRepository.IncrementRetryCount(ctx, event.ID)
		return
	}

	// Mark as completed
	_ = p.postgres.OutboxRepository.MarkAsCompleted(ctx, event.ID)
}
