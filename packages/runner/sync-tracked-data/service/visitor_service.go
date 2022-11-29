package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/repository"
	"log"
)

type VisitorService interface {
	getVisitorsForSync(context.Context) (*entity.Visitor, error)
}

type visitorService struct {
	trackedVisitorRepository *repository.TrackedVisitorRepository
}

func NewVisitorService(trackedVisitorRepository *repository.TrackedVisitorRepository) VisitorService {
	return &visitorService{
		trackedVisitorRepository: trackedVisitorRepository,
	}
}

func (v *visitorService) getVisitorsForSync(ctx context.Context) (*entity.Visitor, error) {
	visitors, err := (*v.trackedVisitorRepository).GetVisitorsWithSyncedFalse(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range visitors {
		log.Printf("Visitor: %v", v.VisitorId)
	}
	return nil, nil
}
