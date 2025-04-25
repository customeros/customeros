package web_event_processor

import (
	"context"

	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

func (s *webEventProcessor) processIdentifyEvent(ctx context.Context, message *pb.WebTrackerEvent) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebEventProcessor.processIdentifyEvent")
	defer spans.Finish()

	return nil
}
