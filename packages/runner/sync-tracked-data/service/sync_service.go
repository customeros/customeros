package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-tracked-data/repository"
	"log"
)

type SyncService interface {
	Sync(runId string)
}

type syncService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewSyncService(repositories *repository.Repositories, serviceContainer *Services) SyncService {
	return &syncService{
		repositories: repositories,
		services:     serviceContainer,
	}
}

func (s *syncService) Sync(runId string) {
	ctx := context.Background()
	_, err := s.services.VisitorService.getVisitorsForSync(ctx)
	if err != nil {
		log.Printf("ERROR run id: %s failed to sync tracked data. error fetching visitors: %v", runId, err.Error())
	}
}
