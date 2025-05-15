package outbox_processor

import (
	"context"
	"github.com/customeros/customeros/packages/server/core-crm/internal/utils"
	nats_common "github.com/customeros/customeros/packages/server/customer-os-common-module/nats"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"time"
)

type OutboxProcessor struct {
	natsConn *nats_common.NATSConnections
	postgres *postgres_repository.Repositories
}

func NewOutboxProcessor(
	natsConn *nats_common.NATSConnections,
	postgres *postgres_repository.Repositories,
) *OutboxProcessor {
	return &OutboxProcessor{
		natsConn: natsConn,
		postgres: postgres,
	}
}

const (
	MAX_MESSAGES_PER_BATCH = 50
	EVENT_LOCK_TIMEOUT     = 3 * time.Minute
)

func (p *OutboxProcessor) ProcessBatch(ctx context.Context) error {
	span, ctx := telemetry.StartCronSpan(ctx, "OutboxProcessor.ProcessBatch")
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
		// Lock the event
		err = p.postgres.OutboxRepository.MarkAsProcessing(innerCtx, event.ID, EVENT_LOCK_TIMEOUT)
		if err != nil {
			// Another worker might have picked it up
			continue
		}

		// Process the event
		err = p.processEvent(innerCtx, event)
		if err != nil {
			// Mark as failed and increment retry count
			_ = p.postgres.OutboxRepository.MarkAsFailed(innerCtx, event.ID, err.Error())
			_ = p.postgres.OutboxRepository.IncrementRetryCount(innerCtx, event.ID)
			continue
		}

		// Mark as completed
		_ = p.postgres.OutboxRepository.MarkAsCompleted(innerCtx, event.ID)
	}

	return nil
}
