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
	TotalCompleted             int `gorm:"column:total_synced_entities"`
	TotalFailed                int `gorm:"column:total_failed_entities"`
	TotalSkipped               int `gorm:"column:total_skipped_entities"`
	CompletedContacts          int `gorm:"column:synced_contacts"`
	FailedContacts             int `gorm:"column:failed_contacts"`
	SkippedContacts            int `gorm:"column:skipped_contacts"`
	CompletedUsers             int `gorm:"column:synced_users"`
	FailedUsers                int `gorm:"column:failed_users"`
	SkippedUsers               int `gorm:"column:skipped_users"`
	CompletedOrganizations     int `gorm:"column:synced_organizations"`
	FailedOrganizations        int `gorm:"column:failed_organizations"`
	SkippedOrganizations       int `gorm:"column:skipped_organizations"`
	CompletedNotes             int `gorm:"column:synced_notes"`
	FailedNotes                int `gorm:"column:failed_notes"`
	SkippedNotes               int `gorm:"column:skipped_notes"`
	CompletedEmailMessages     int `gorm:"column:synced_email_messages"`
	FailedEmailMessages        int `gorm:"column:failed_email_messages"`
	SkippedEmailMessages       int `gorm:"column:skipped_email_messages"`
	CompletedIssues            int `gorm:"column:synced_issues"`
	FailedIssues               int `gorm:"column:failed_issues"`
	SkippedIssues              int `gorm:"column:skipped_issues"`
	CompletedMeetings          int `gorm:"column:synced_meetings"`
	FailedMeetings             int `gorm:"column:failed_meetings"`
	SkippedMeetings            int `gorm:"column:skipped_meetings"`
	CompletedInteractionEvents int `gorm:"column:synced_interaction_events"`
	FailedInteractionEvents    int `gorm:"column:failed_interaction_events"`
	SkippedInteractionEvents   int `gorm:"column:skipped_interaction_events"`
	CompletedLogEntries        int `gorm:"column:synced_log_entries"`
	FailedLogEntries           int `gorm:"column:failed_log_entries"`
	SkippedLogEntries          int `gorm:"column:skipped_log_entries"`
}

func (SyncRun) TableName() string {
	return "sync_run"
}

func (s *SyncRun) SumTotalCompleted() {
	s.TotalCompleted = s.CompletedContacts + s.CompletedUsers + s.CompletedOrganizations + s.CompletedNotes + s.CompletedEmailMessages +
		s.CompletedIssues + s.CompletedMeetings + s.CompletedInteractionEvents + s.CompletedLogEntries
}

func (s *SyncRun) SumTotalFailed() {
	s.TotalFailed = s.FailedContacts + s.FailedUsers + s.FailedOrganizations + s.FailedNotes + s.FailedEmailMessages +
		s.FailedIssues + s.FailedMeetings + s.FailedInteractionEvents + s.FailedLogEntries
}

func (s *SyncRun) SumTotalSkipped() {
	s.TotalSkipped = s.SkippedContacts + s.SkippedUsers + s.SkippedOrganizations + s.SkippedNotes + s.SkippedEmailMessages +
		s.SkippedIssues + s.SkippedMeetings + s.SkippedInteractionEvents + s.SkippedLogEntries
}
