package entity

import (
	"time"
)

type SlackSyncRunStatus struct {
	ID             uint      `gorm:"primarykey"`
	RunId          string    `gorm:"run_id;not null"`
	Tenant         string    `gorm:"type:varchar(255)" binding:"required"`
	SlackChannelId string    `gorm:"type:varchar(50)" binding:"required"`
	OrganizationId string    `gorm:"type:varchar(50)" binding:"required"`
	StartAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	EndAt          time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Failed         bool      `gorm:"column:failed"`
}

func (SlackSyncRunStatus) TableName() string {
	return "slack_sync_run_status"
}
