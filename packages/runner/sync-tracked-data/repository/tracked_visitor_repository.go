package repository

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/gen"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/gen/visitor"
)

type TrackedVisitorRepository interface {
	GetVisitorsWithSyncedFalse(context.Context) ([]*gen.Visitor, error)
}

type trackedVisitorRepository struct {
	client *gen.Client
}

func NewTrackedVisitorRepository(client *gen.Client) TrackedVisitorRepository {
	return &trackedVisitorRepository{
		client: client,
	}
}

func (t *trackedVisitorRepository) GetVisitorsWithSyncedFalse(ctx context.Context) ([]*gen.Visitor, error) {
	return t.client.Visitor.Query().
		Where(visitor.SyncedToCustomerOsEQ(false),
			visitor.VisitorIDNotNil()).
		All(ctx)
}
