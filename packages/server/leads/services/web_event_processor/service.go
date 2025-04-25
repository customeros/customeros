package web_event_processor

import (
	"context"

	"github.com/customeros/customeros/packages/server/leads/dto"
	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
)

type WebEventProcessor interface {
	Process(ctx context.Context, event *dto.WebTrackerEvent, webtrackerID string)
}

type webEventProcessor struct {
	natsConn     *nats_internal.NATSConnections
	repositories *repository.Repositories
}

func NewWebEventProcessor(
	natsConn *nats_internal.NATSConnections,
	repositories *repository.Repositories,
) WebEventProcessor {
	return &webEventProcessor{
		natsConn:     natsConn,
		repositories: repositories,
	}
}
