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

const maxAttempts = 10

func GetAirbyteUnprocessedRawRecords(ctx context.Context, db *gorm.DB, limit int, runId, syncedEntity, tableSuffix, tenant, sourceId string) (entity.AirbyteRaws, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetAirbyteUnprocessedRawRecords")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.Int("limit", limit), log.String("syncedEntity", syncedEntity), log.String("tableSuffix", tableSuffix))

	var airbyteRecords entity.AirbyteRaws
	rawTableName := fmt.Sprintf(`%s_%s_raw__stream_%s`, sourceId, tenant, tableSuffix)

	err := db.
		Raw(fmt.Sprintf(`SELECT a.*
FROM %s a
LEFT JOIN airbyte_sync_status s ON a._airbyte_raw_id = s._airbyte_raw_id and s.entity = ? and s.table_suffix = ? and s.tenant = ?
WHERE (s.synced_to_customer_os IS NULL OR s.synced_to_customer_os = FALSE)
  AND (s.synced_to_customer_os_attempt IS NULL OR s.synced_to_customer_os_attempt < ?)
  AND (s.run_id IS NULL OR s.run_id <> ?)
ORDER BY a._airbyte_extracted_at ASC
LIMIT ?`, rawTableName, syncedEntity, tableSuffix, tenant, maxAttempts, runId, limit)).
		Find(&airbyteRecords).Error

	if err != nil {
		return nil, err
	}
	return airbyteRecords, nil
}

func GetOpenlineUnprocessedRawRecords(ctx context.Context, db *gorm.DB, limit int, runId, syncedEntity, tableSuffix string) (entity.OpenlineRaws, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GetOpenlineUnprocessedRawRecords")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.Int("limit", limit), log.String("syncedEntity", syncedEntity), log.String("tableSuffix", tableSuffix))

	var rawRecords entity.OpenlineRaws

	err := db.
		Raw(fmt.Sprintf(`SELECT o.*
FROM _openline_raw_%s o
LEFT JOIN openline_sync_status s ON o.raw_id = s.raw_id and s.entity = ? and s.table_suffix = ?
WHERE (s.synced_to_customer_os IS NULL OR s.synced_to_customer_os = FALSE)
  AND (s.synced_to_customer_os_attempt IS NULL OR s.synced_to_customer_os_attempt < ?)
  AND (s.run_id IS NULL OR s.run_id <> ?)
ORDER BY o.emitted_at ASC
LIMIT ?`, tableSuffix), syncedEntity, tableSuffix, maxAttempts, runId, limit).
		Find(&rawRecords).Error

	if err != nil {
		return nil, err
	}
	return rawRecords, nil
}

func MarkAirbyteRawRecordProcessed(ctx context.Context, db *gorm.DB, tenant, syncedEntity, tableSuffix, airbyteRawId string, synced, skipped bool, runId, externalId, reason string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MarkAirbyteRawRecordProcessed")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.String("syncedEntity", syncedEntity), log.String("tableSuffix", tableSuffix))

	syncStatus := entity.SyncStatusForAirbyte{
		Entity:       syncedEntity,
		TableSuffix:  tableSuffix,
		AirbyteRawId: airbyteRawId,
		Tenant:       tenant,
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
		Where(&entity.SyncStatusForAirbyte{AirbyteRawId: airbyteRawId, Entity: syncedEntity, TableSuffix: tableSuffix}).
		Save(&syncStatus).Error
}

func MarkOpenlineRawRecordProcessed(ctx context.Context, db *gorm.DB, syncedEntity, tableSuffix, rawId string, synced, skipped bool, runId, externalId, reason string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MarkOpenlineRawRecordProcessed")
	defer span.Finish()
	tracing.SetDefaultPostgresRepositorySpanTags(ctx, span)
	span.LogFields(log.String("syncedEntity", syncedEntity), log.String("tableSuffix", tableSuffix))

	syncStatus := entity.SyncStatusForOpenline{
		Entity:      syncedEntity,
		TableSuffix: tableSuffix,
		RawId:       rawId,
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
		Where(&entity.SyncStatusForOpenline{RawId: rawId, Entity: syncedEntity, TableSuffix: tableSuffix}).
		Save(&syncStatus).Error
}
