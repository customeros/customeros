package event_logger

import (
	"context"

	"github.com/customeros/mailstack/proto/pb"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"

	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
)

func (s *LeadEventLoggerService) processErrorMessage(ctx context.Context, msg *nats.Msg) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "LeadEventLoggerService.processErrorMessage")
	defer spans.Finish()

	message := &pb.ErrorEvent{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		spans.TraceError(err)
		return
	}

	return
}
