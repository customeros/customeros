package outbox_processor

import (
	"context"

	"github.com/customeros/customeros/packages/server/leads/internal/models"
)

func (s *OutboxProcessor) processEvent(ctx context.Context, event *models.OutboxEvent) {}
