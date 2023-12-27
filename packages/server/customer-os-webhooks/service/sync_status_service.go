package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type SyncStatus struct {
	FailedSync bool
	Skipped    bool
	Reason     string
}

type SyncResult struct {
	Skipped   int `json:"skipped"`
	Failed    int `json:"failed"`
	Completed int `json:"completed"`
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
	SaveSyncResults(ctx context.Context, tenant, externalSystem, appSource, entityType string, syncDate time.Time, statuses []SyncStatus)
	PrepareSyncResult(statuses []SyncStatus) SyncResult
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

func (s *syncStatusService) SaveSyncResults(ctx context.Context, tenant, externalSystem, appSource, entityType string, syncDate time.Time, statuses []SyncStatus) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "SyncStatusService.SaveSyncResults")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "statuses", statuses)
	span.LogFields(log.String("externalSystem", externalSystem), log.String("appSource", appSource), log.String("entityType", entityType))

	completed, failed, skipped := 0, 0, 0
	reason := ""
	for _, status := range statuses {
		if status.FailedSync {
			if status.Reason != "" {
				failed++
				if reason != "" {
					reason += "\n"
				}
				reason = status.Reason
			}
		} else if status.Skipped {
			skipped++
			if status.Reason != "" {
				if reason != "" {
					reason += "\n"
				}
				reason = status.Reason
			}
		} else {
			completed++
		}
	}
	s.repositories.SyncRunWebhookRepository.Save(ctx, postgresentity.SyncRunWebhook{
		Tenant:         tenant,
		ExternalSystem: externalSystem,
		AppSource:      appSource,
		StartAt:        syncDate,
		EndAt:          time.Now(),
		Entity:         entityType,
		Total:          completed + failed + skipped,
		Completed:      completed,
		Failed:         failed,
		Skipped:        skipped,
		Reason:         reason,
	})
}

func (s *syncStatusService) PrepareSyncResult(statuses []SyncStatus) SyncResult {
	completed, failed, skipped := 0, 0, 0
	for _, status := range statuses {
		if status.FailedSync {
			failed++
		} else if status.Skipped {
			skipped++
		} else {
			completed++
		}
	}
	return SyncResult{
		Completed: completed,
		Failed:    failed,
		Skipped:   skipped,
	}
}
