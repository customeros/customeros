package outbox_processor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
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

func (p *OutboxProcessor) ProcessBatch(ctx context.Context) error {
	// Get pending events
	events, err := p.outboxRepo.GetPendingEvents(ctx, 10)
	if err != nil {
		return err
	}

	for _, event := range events {
		// Lock the event
		err := p.outboxRepo.MarkAsProcessing(ctx, event.ID, 5*time.Minute)
		if err != nil {
			// Another worker might have picked it up
			continue
		}

		// Process the event
		err = p.processEvent(ctx, event)
		if err != nil {
			// Mark as failed and increment retry count
			p.outboxRepo.MarkAsFailed(ctx, event.ID, err.Error())
			p.outboxRepo.IncrementRetryCount(ctx, event.ID)
			continue
		}

		// Mark as processed
		p.outboxRepo.MarkAsProcessed(ctx, event.ID)
	}

	return nil
}
