package entity

import "time"

type SyncStatusTicketComment struct {
	Id                          int64     `gorm:"column:id;primaryKey"`
	AirbyteAbId                 string    `gorm:"column:_airbyte_ab_id;primaryKey"`
	AirbyteTicketCommentsHashid string    `gorm:"column:_airbyte_ticket_comments_hashid;primaryKey"`
	SyncedToCustomerOs          bool      `gorm:"column:synced_to_customer_os"`
	SyncAttempt                 int       `gorm:"column:synced_to_customer_os_attempt"`
	SyncedAt                    time.Time `gorm:"column:synced_to_customer_os_at"`
	RunId                       string    `gorm:"column:run_id"`
}

func (SyncStatusTicketComment) TableName() string {
	return "openline_sync_status_ticket_comments"
}
