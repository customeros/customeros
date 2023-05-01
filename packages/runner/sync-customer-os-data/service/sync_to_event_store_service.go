package service

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

const syncToEventStoreBatchSize = 2

type SyncToEventStoreService interface {
	SyncEmails(ctx context.Context)
}

type syncToEventStoreService struct {
	repositories *repository.Repositories
	services     *Services
	grpcClients  *grpc_client.Clients
	batchSize    int
}

func NewSyncToEventStoreService(repositories *repository.Repositories, services *Services, grpcClients *grpc_client.Clients) SyncToEventStoreService {
	return &syncToEventStoreService{
		repositories: repositories,
		services:     services,
		grpcClients:  grpcClients,
		batchSize:    syncToEventStoreBatchSize,
	}
}

func (s *syncToEventStoreService) SyncEmails(ctx context.Context) {
	logrus.Infof("start sync emails to eventstore at %v", time.Now().UTC())
	syncedCount := 0
	logrus.Infof("completed sync %v emails to eventstore at %v", syncedCount, time.Now().UTC())
}
