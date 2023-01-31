package entity

import (
	"time"
)

type SyncRun struct {
	ID              uint      `gorm:"primarykey"`
	RunId           string    `gorm:"run_id;not null"`
	StarAt          time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	EndAt           time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	SyncedPageViews int       `gorm:"column:synced_page_views"`
}

func (SyncRun) TableName() string {
	return "tracked_data_sync_run"
}
