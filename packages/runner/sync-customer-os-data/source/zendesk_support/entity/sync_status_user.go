package entity

import "time"

type SyncStatusUser struct {
	Id                 int64     `gorm:"column:id;primaryKey"`
	AirbyteAbId        string    `gorm:"column:_airbyte_ab_id;primaryKey"`
	AirbyteUsersHashid string    `gorm:"column:_airbyte_users_hashid;primaryKey"`
	SyncedToCustomerOs bool      `gorm:"column:synced_to_customer_os"`
	SyncAttempt        int       `gorm:"column:synced_to_customer_os_attempt"`
	SyncedAt           time.Time `gorm:"column:synced_to_customer_os_at"`
	RunId              string    `gorm:"column:run_id"`
}

func (SyncStatusUser) TableName() string {
	return "openline_sync_status_users"
}
