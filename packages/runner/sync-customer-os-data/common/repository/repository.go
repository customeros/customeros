package repository

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/common/entity"
	"gorm.io/gorm"
	"time"
)

const (
	ContactEntity      = "contacts"
	CompanyEntity      = "companies"
	OwnerEntity        = "owners"
	NoteEntity         = "engagements_notes"
	MeetingEntity      = "engagements_meetings"
	EmailMessageEntity = "engagements_emails"
)

func GetAirbyteUnprocessedRecords(db *gorm.DB, limit int, runId, tableSuffix string) (entity.AirbyteRaws, error) {
	var airbyteRecords entity.AirbyteRaws

	err := db.
		Raw(fmt.Sprintf(`SELECT a.*
FROM _airbyte_raw_%s a
LEFT JOIN openline_sync_status s ON a._airbyte_ab_id = s._airbyte_ab_id and s.entity = ?
WHERE (s.synced_to_customer_os IS NULL OR s.synced_to_customer_os = FALSE)
  AND (s.synced_to_customer_os_attempt IS NULL OR s.synced_to_customer_os_attempt < ?)
  AND (s.run_id IS NULL OR s.run_id <> ?)
ORDER BY a._airbyte_emitted_at ASC
LIMIT ?`, tableSuffix), tableSuffix, 10, runId, limit).
		Find(&airbyteRecords).Error

	if err != nil {
		return nil, err
	}
	return airbyteRecords, nil
}

func MarkProcessed(db *gorm.DB, syncedEntity, airbyteAbId string, synced bool, runId, externalSyncId string) error {
	syncStatus := entity.SyncStatus{
		Entity:         syncedEntity,
		AirbyteAbId:    airbyteAbId,
		ExternalSyncId: externalSyncId,
	}
	db.FirstOrCreate(&syncStatus, syncStatus)

	return db.Model(&syncStatus).
		Where(&entity.SyncStatus{AirbyteAbId: airbyteAbId, Entity: syncedEntity}).
		Updates(entity.SyncStatus{
			SyncedToCustomerOs: synced,
			SyncedAt:           time.Now(),
			SyncAttempt:        syncStatus.SyncAttempt + 1,
			RunId:              runId,
		}).Error
}
