package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	"time"
)

type SyncStatus struct {
	FailedSync bool
	Skipped    bool
	Reason     string
}

func NewFailedSyncStatus(reason string) SyncStatus {
	return SyncStatus{
		FailedSync: true,
		Reason:     reason,
	}
}

func NewSkippedSyncStatus(reason string) SyncStatus {
	return SyncStatus{
		Skipped: true,
		Reason:  reason,
	}
}

func NewSuccessfulSyncStatus() SyncStatus {
	return SyncStatus{}
}

type SyncStatusService interface {
	SaveSyncResults(ctx context.Context, tenant, externalSystem string, syncDate time.Time, statuses []SyncStatus)
}

type syncStatusService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewSyncStatusService(log logger.Logger, repositories *repository.Repositories) SyncStatusService {
	return &syncStatusService{
		log:          log,
		repositories: repositories,
	}
}

func (s syncStatusService) SaveSyncResults(ctx context.Context, tenant, externalSystem string, syncDate time.Time, statuses []SyncStatus) {
	//TODO implement me COS-364
}
