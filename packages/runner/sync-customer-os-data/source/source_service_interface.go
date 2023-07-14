package source

import (
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
	GetDataForSync(dataType common.SyncedEntityType, batchSize int, runId string) []interface{}
	MarkProcessed(syncId, runId string, synced, skipped bool, reason string) error
}
