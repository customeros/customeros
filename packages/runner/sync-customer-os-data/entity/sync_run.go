package entity

import (
	"time"
)

type SyncRun struct {
	ID                         uint      `gorm:"primarykey"`
	RunId                      string    `gorm:"run_id;not null"`
	StartAt                    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	EndAt                      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	TenantSyncSettingsId       uint
	TenantSyncSettings         TenantSyncSettings
	TotalCompletedEntities     int `gorm:"column:total_synced_entities"`
	TotalFailedEntities        int `gorm:"column:total_failed_entities"`
	CompletedContacts          int `gorm:"column:synced_contacts"`
	FailedContacts             int `gorm:"column:failed_contacts"`
	CompletedUsers             int `gorm:"column:synced_users"`
	FailedUsers                int `gorm:"column:failed_users"`
	CompletedOrganizations     int `gorm:"column:synced_organizations"`
	FailedOrganizations        int `gorm:"column:failed_organizations"`
	CompletedNotes             int `gorm:"column:synced_notes"`
	FailedNotes                int `gorm:"column:failed_notes"`
	CompletedEmailMessages     int `gorm:"column:synced_email_messages"`
	FailedEmailMessages        int `gorm:"column:failed_email_messages"`
	CompletedIssues            int `gorm:"column:synced_issues"`
	FailedIssues               int `gorm:"column:failed_issues"`
	CompletedMeetings          int `gorm:"column:synced_meetings"`
	FailedMeetings             int `gorm:"column:failed_meetings"`
	CompletedInteractionEvents int `gorm:"column:synced_interaction_events"`
	FailedInteractionEvents    int `gorm:"column:failed_interaction_events"`
}

func (SyncRun) TableName() string {
	return "sync_run"
}
