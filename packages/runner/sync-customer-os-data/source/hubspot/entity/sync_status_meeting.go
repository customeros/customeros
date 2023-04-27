package entity

import "time"

type SyncStatusMeeting struct {
	Id                    string    `gorm:"column:id;primaryKey"`
	AirbyteAbId           string    `gorm:"column:_airbyte_ab_id;primaryKey"`
	AirbyteMeetingsHashid string    `gorm:"column:_airbyte_engagements_meetings_hashid;primaryKey"`
	SyncedToCustomerOs    bool      `gorm:"column:synced_to_customer_os"`
	SyncAttempt           int       `gorm:"column:synced_to_customer_os_attempt"`
	SyncedAt              time.Time `gorm:"column:synced_to_customer_os_at"`
	RunId                 string    `gorm:"column:run_id"`
}

func (SyncStatusMeeting) TableName() string {
	return "openline_sync_status_meetings"
}
