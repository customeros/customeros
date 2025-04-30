package web_event_processor

import (
	"context"

	"gorm.io/gorm"

	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
	"github.com/customeros/customeros/packages/server/leads/proto/pb"
)

type WebEventProcessor interface {
	Process(ctx context.Context, webtrackerID string, event *pb.WebTrackerEvent)
}

type webEventProcessor struct {
	natsConn     *nats_internal.NATSConnections
	leadsWriteDB *gorm.DB
	repositories *repository.Repositories
}

func NewWebEventProcessor(
	natsConn *nats_internal.NATSConnections,
	leadsWriteDB *gorm.DB,
	repositories *repository.Repositories,
) WebEventProcessor {
	return &webEventProcessor{
		natsConn:     natsConn,
		leadsWriteDB: leadsWriteDB,
		repositories: repositories,
	}
}
