package visitor_identity

import (
	"context"
	"errors"

	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

func (s *VisitorIdentityService) identifyWebVisitor(ctx context.Context, req *pb.IdentifyVisitorRequest) *pb.IdentifyVisitorResponse {
	span, ctx := telemetry.StartServiceSpan(ctx, "VisitorIdentityService.identifyWebVisitor")

	// lookup IP in identity table
	span.TraceError(errors.New(""))

	// check to see if bot IP

	// attempt to identify IP

	return nil
}
