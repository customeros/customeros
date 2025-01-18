package service

import (
	"context"
	"github.com/customeros/customeros/packages/runner/sync-customer-os-data/source"
	"time"
)

type result struct {
	completed int
	failed    int
	skipped   int
}

type SyncService interface {
	Sync(ctx context.Context, sourceService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int)
}
