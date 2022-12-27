package entity

import (
	"time"
)

type SyncRun struct {
	ID                   uint      `gorm:"primarykey"`
	RunId                string    `gorm:"run_id;not null"`
	StarAt               time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	EndAt                time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	TenantSyncSettingsId uint
	TenantSyncSettings   TenantSyncSettings
	CompletedContacts    int `gorm:"column:synced_contacts"`
	FailedContacts       int `gorm:"column:failed_contacts"`
	CompletedUsers       int `gorm:"column:synced_users"`
	FailedUsers          int `gorm:"column:failed_users"`
	CompletedCompanies   int `gorm:"column:synced_companies"`
	FailedCompanies      int `gorm:"column:failed_companies"`
	CompletedNotes       int `gorm:"column:synced_notes"`
	FailedNotes          int `gorm:"column:failed_notes"`
}

func (SyncRun) TableName() string {
	return "sync_run"
}
