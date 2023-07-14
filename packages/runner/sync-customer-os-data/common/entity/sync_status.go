package entity

import "time"

type SyncStatus struct {
	Entity             string    `gorm:"column:entity;primaryKey"`
	AirbyteAbId        string    `gorm:"column:_airbyte_ab_id;primaryKey"`
	SyncedToCustomerOs bool      `gorm:"column:synced_to_customer_os"`
	SyncAttempt        int       `gorm:"column:synced_to_customer_os_attempt"`
	SyncedAt           time.Time `gorm:"column:synced_to_customer_os_at"`
	RunId              string    `gorm:"column:run_id"`
	ExternalId         string    `gorm:"column:external_id"`
	Reason             string    `gorm:"column:reason"`
}

func (SyncStatus) TableName() string {
	return "openline_sync_status"
}
