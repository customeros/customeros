package outbox_processor

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"time"
)

const (
	LIMIT    = 1000
	LOOKBACK = 72 * time.Hour
)

func (s *OutboxProcessor) Cleanup(ctx context.Context) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "OutboxProcessor.Cleanup")
	defer span.Finish()

	_, err := s.postgres.OutboxRepository.DeleteProcessedEvents(ctx, LOOKBACK, LIMIT)
	if err != nil {
		span.TraceError(err)
		return err
	}

	return nil
}
