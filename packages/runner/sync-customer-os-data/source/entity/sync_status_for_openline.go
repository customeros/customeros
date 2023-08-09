package entity

import "time"

type SyncStatusForOpenline struct {
	Entity             string    `gorm:"column:entity;primaryKey"`
	TableSuffix        string    `gorm:"column:table_suffix;primaryKey"`
	RawId              string    `gorm:"column:raw_id;primaryKey"`
	SyncedToCustomerOs bool      `gorm:"column:synced_to_customer_os"`
	Skipped            bool      `gorm:"column:skipped"`
	SyncAttempt        int       `gorm:"column:synced_to_customer_os_attempt"`
	SyncedAt           time.Time `gorm:"column:synced_to_customer_os_at"`
	RunId              string    `gorm:"column:run_id"`
	ExternalId         string    `gorm:"column:external_id"`
	Reason             string    `gorm:"column:reason"`
}

func (SyncStatusForOpenline) TableName() string {
	return "openline_sync_status"
}
