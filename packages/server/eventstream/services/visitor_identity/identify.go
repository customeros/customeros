package visitor_identity

import (
	"context"

	"github.com/customeros/customeros/packages/server/eventstream/internal/telemetry"
	"github.com/customeros/customeros/packages/server/eventstream/proto/pb"
)

func (s *VisitorIdentityService) identifyWebVisitor(ctx context.Context, req *pb.IdentifyVisitorRequest) *pb.IdentifyVisitorResponse {
	span, ctx := telemetry.StartServiceSpan(ctx, "VisitorIdentityService.identifyWebVisitor")

	// lookup IP in identity table

	// check to see if bot IP

	// attempt to identify IP

	return nil
}
