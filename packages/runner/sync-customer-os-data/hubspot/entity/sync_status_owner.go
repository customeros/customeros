package entity

import "time"

type SyncStatusOwner struct {
	Id                  string    `gorm:"column:id;primaryKey"`
	AirbyteAbId         string    `gorm:"column:_airbyte_ab_id;primaryKey"`
	AirbyteOwnersHashid string    `gorm:"column:_airbyte_owners_hashid;primaryKey"`
	SyncedToCustomerOs  bool      `gorm:"column:synced_to_customer_os"`
	SyncAttempt         int       `gorm:"column:synced_to_customer_os_attempt"`
	SyncedAt            time.Time `gorm:"column:synced_to_customer_os_at"`
	RunId               string    `gorm:"column:run_id"`
}

func (SyncStatusOwner) TableName() string {
	return "openline_sync_status_owners"
}
