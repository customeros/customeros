package webpage_classification

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/customeros/customeros/packages/server/ai/internal/telemetry"
)

func (s *webpageClassification) handleWebpageClassification(ctx context.Context, msg *nats.Msg) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "webpageClassification.handleWebpageClassification")
	defer span.Finish()

	return nil
}
