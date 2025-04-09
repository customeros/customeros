package web_event_processor

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/customeros/customeros/packages/server/eventstream/internal/telemetry"
)

func (s *WebEventProcessor) handlePageExit(ctx context.Context, msg *nats.Msg) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebEventProcessor.handlePageExit")
	defer spans.Finish()
}

func (s *WebEventProcessor) handlePageView(ctx context.Context, msg *nats.Msg) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebEventProcessor.handlePageView")
	defer spans.Finish()
}

func (s *WebEventProcessor) handleClick(ctx context.Context, msg *nats.Msg) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebEventProcessor.handleClick")
	defer spans.Finish()
}

func (s *WebEventProcessor) handleIdentify(ctx context.Context, msg *nats.Msg) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "WebEventProcessor.handleIdentify")
	defer spans.Finish()
}
