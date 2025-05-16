package icp

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/customeros/customeros/packages/server/core-crm/internal/telemetry"
)

func (s *icpService) handleICPFitResponse(ctx context.Context, msg *nats.Msg) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "icpService.handleICPFitResponse")
	defer span.Finish()

	// TODO

	return nil
}
