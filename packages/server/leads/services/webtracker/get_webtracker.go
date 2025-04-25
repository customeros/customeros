package webtracker

import (
	"context"

	"github.com/customeros/customeros/packages/server/leads/internal/models"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
)

func (s *webtrackerService) GetWebtrackerByOrigin(ctx context.Context, origin string) (*models.WebTracker, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "webtrackerService.GetWebtrackerByOrigin")
	defer span.Finish()

	return s.repositories.WebTracker.GetByDomain(ctx, origin)
}
