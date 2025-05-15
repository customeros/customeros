package company_icp_fit

import (
	"context"

	nats_common "github.com/customeros/customeros/packages/server/customer-os-common-module/nats"

	"github.com/customeros/customeros/packages/server/ai/interfaces"
	"github.com/customeros/customeros/packages/server/ai/internal/telemetry"
)

type icpFitAnalyzer struct{}

func NewICPFitAnalyzer(natsConn *nats_common.NATSConnections) interfaces.NatsService {
	return &icpFitAnalyzer{}
}

func (s *icpFitAnalyzer) Start(ctx context.Context) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "icpFitAnalyzer.Start")
	defer span.Finish()

	return nil
}

func (s *icpFitAnalyzer) Stop() {
}
