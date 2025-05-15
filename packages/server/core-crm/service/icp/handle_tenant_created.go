package icp

import (
	"context"

	"github.com/nats-io/nats.go"

	"github.com/customeros/customeros/packages/server/core-crm/internal/telemetry"
)

func (s *icpService) handleTenantCreated(ctx context.Context, msg *nats.Msg) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "icpService.handleTenantCreated")
	defer span.Finish()

	return nil
}
