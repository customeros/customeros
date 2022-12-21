package entity

import "time"

type SyncStatusNote struct {
	Id                 string    `gorm:"column:id;primaryKey"`
	AirbyteAbId        string    `gorm:"column:_airbyte_ab_id;primaryKey"`
	AirbyteNotesHashid string    `gorm:"column:_airbyte_engagements_notes_hashid;primaryKey"`
	SyncedToCustomerOs bool      `gorm:"column:synced_to_customer_os"`
	SyncAttempt        int       `gorm:"column:synced_to_customer_os_attempt"`
	SyncedAt           time.Time `gorm:"column:synced_to_customer_os_at"`
}

func (SyncStatusNote) TableName() string {
	return "openline_sync_status_notes"
}
