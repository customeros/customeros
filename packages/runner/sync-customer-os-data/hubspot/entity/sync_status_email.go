package entity

import "time"

type SyncStatusEmail struct {
	Id                  string    `gorm:"column:id;primaryKey"`
	AirbyteAbId         string    `gorm:"column:_airbyte_ab_id;primaryKey"`
	AirbyteEmailsHashid string    `gorm:"column:_airbyte_engagements_emails_hashid;primaryKey"`
	SyncedToCustomerOs  bool      `gorm:"column:synced_to_customer_os"`
	SyncAttempt         int       `gorm:"column:synced_to_customer_os_attempt"`
	SyncedAt            time.Time `gorm:"column:synced_to_customer_os_at"`
	RunId               string    `gorm:"column:run_id"`
}

func (SyncStatusEmail) TableName() string {
	return "openline_sync_status_emails"
}
