package entity

import (
	"time"
)

type SyncRun struct {
	ID                   uint      `gorm:"primarykey"`
	StarAt               time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	EndAt                time.Time
	TenantSyncSettingsId int
	TenantSyncSettings   TenantSyncSettings
	Status               string `gorm:"column:status"`
}

func (SyncRun) TableName() string {
	return "sync_run"
}
