package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/gen"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/repository"
	"log"
)

type VisitorService interface {
	getVisitorsForSync(context.Context) (*gen.Visitor, error)
}

type visitorService struct {
	trackedVisitorRepository *repository.TrackedVisitorRepository
}

func NewVisitorService(trackedVisitorRepository *repository.TrackedVisitorRepository) VisitorService {
	return &visitorService{
		trackedVisitorRepository: trackedVisitorRepository,
	}
}

func (v *visitorService) getVisitorsForSync(ctx context.Context) (*gen.Visitor, error) {
	visitors, err := (*v.trackedVisitorRepository).GetVisitorsWithSyncedFalse(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range visitors {
		log.Printf("Visitor: %v", v.VisitorID)
	}
	return nil, nil
}
