package source

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common"
)

type ProcessingEntity struct {
	ExternalId  string
	Entity      string
	TableSuffix string
}

type SourceDataService interface {
	Init()
	Close()
	SourceId() string
	GetDataForSync(ctx context.Context, dataType common.SyncedEntityType, batchSize int, runId string) []interface{}
	MarkProcessed(ctx context.Context, syncId, runId string, synced, skipped bool, reason string) error
}
