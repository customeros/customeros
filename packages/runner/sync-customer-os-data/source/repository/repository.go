package repository

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/source/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gorm.io/gorm"
)

func GetAirbyteUnprocessedRecords(ctx context.Context, db *gorm.DB, limit int, runId, syncedEntity, tableSuffix string) (entity.AirbyteRaws, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetAirbyteUnprocessedRecords")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.Int("limit", limit), log.String("syncedEntity", syncedEntity), log.String("tableSuffix", tableSuffix))

	var airbyteRecords entity.AirbyteRaws

	err := db.
		Raw(fmt.Sprintf(`SELECT a.*
FROM _airbyte_raw_%s a
LEFT JOIN openline_sync_status s ON a._airbyte_ab_id = s._airbyte_ab_id and s.entity = ? and s.table_suffix = ?
WHERE (s.synced_to_customer_os IS NULL OR s.synced_to_customer_os = FALSE)
  AND (s.synced_to_customer_os_attempt IS NULL OR s.synced_to_customer_os_attempt < ?)
  AND (s.run_id IS NULL OR s.run_id <> ?)
ORDER BY a._airbyte_emitted_at ASC
LIMIT ?`, tableSuffix), syncedEntity, tableSuffix, 10, runId, limit).
		Find(&airbyteRecords).Error

	if err != nil {
		return nil, err
	}
	return airbyteRecords, nil
}

func MarkProcessed(ctx context.Context, db *gorm.DB, syncedEntity, tableSuffix, airbyteAbId string, synced, skipped bool, runId, externalId, reason string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MarkProcessed")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.String("syncedEntity", syncedEntity), log.String("tableSuffix", tableSuffix))

	syncStatus := entity.SyncStatus{
		Entity:      syncedEntity,
		TableSuffix: tableSuffix,
		AirbyteAbId: airbyteAbId,
	}
	db.FirstOrCreate(&syncStatus, syncStatus)
	syncStatus.Reason = reason
	syncStatus.Skipped = skipped
	syncStatus.SyncedToCustomerOs = synced
	syncStatus.SyncedAt = utils.Now()
	syncStatus.RunId = runId
	syncStatus.ExternalId = externalId
	syncStatus.SyncAttempt = syncStatus.SyncAttempt + 1

	return db.Model(&syncStatus).
		Where(&entity.SyncStatus{AirbyteAbId: airbyteAbId, Entity: syncedEntity, TableSuffix: tableSuffix}).
		Save(&syncStatus).Error
}
