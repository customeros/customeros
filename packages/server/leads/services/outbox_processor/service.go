package outbox_processor

import (
	"context"
	"time"

	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
)

type OutboxProcessor struct {
	natsConn     *nats_internal.NATSConnections
	repositories *repository.Repositories
}

func NewOutboxProcessor(
	natsConn *nats_internal.NATSConnections,
	repos *repository.Repositories,
) *OutboxProcessor {
	return &OutboxProcessor{
		natsConn:     natsConn,
		repositories: repos,
	}
}

const (
	MAX_MESSAGES_PER_BATCH = 50
	EVENT_LOCK_TIMEOUT     = 3 * time.Minute
)

func (p *OutboxProcessor) ProcessBatch(ctx context.Context) error {
	// Get pending events
	events, err := p.repositories.Outbox.GetPendingEvents(ctx, MAX_MESSAGES_PER_BATCH)
	if err != nil {
		return err
	}

	for _, event := range events {
		// Lock the event
		err := p.repositories.Outbox.MarkAsProcessing(ctx, event.ID, EVENT_LOCK_TIMEOUT)
		if err != nil {
			// Another worker might have picked it up
			continue
		}

		// Process the event
		err = p.processEvent(ctx, event)
		if err != nil {
			// Mark as failed and increment retry count
			p.repositories.Outbox.MarkAsFailed(ctx, event.ID, err.Error())
			p.repositories.Outbox.IncrementRetryCount(ctx, event.ID)
			continue
		}

		// Mark as processed
		p.repositories.Outbox.MarkAsCompleted(ctx, event.ID)
	}

	return nil
}
