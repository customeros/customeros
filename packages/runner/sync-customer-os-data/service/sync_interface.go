package service

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source"
	"golang.org/x/net/context"
	"time"
)

type SyncService interface {
	Sync(ctx context.Context, sourceService source.SourceDataService, syncDate time.Time, tenant, runId string, batchSize int) (int, int, int)
}
