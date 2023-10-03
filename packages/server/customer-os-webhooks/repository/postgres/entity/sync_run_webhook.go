package postgresentity

import (
	"github.com/google/uuid"
	"time"
)

type SyncRunWebhook struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Tenant         string    `gorm:"column:tenant;size:50"`
	ExternalSystem string    `gorm:"column:external_system;size:50"`
	AppSource      string    `gorm:"column:app_source;size:50"`
	Entity         string    `gorm:"column:entity;size:50"`
	StartAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	EndAt          time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Reason         string    `gorm:"column:reason"`
	Total          int       `gorm:"column:total"`
	Completed      int       `gorm:"column:completed"`
	Skipped        int       `gorm:"column:skipped"`
	Failed         int       `gorm:"column:failed"`
}

func (SyncRunWebhook) TableName() string {
	return "sync_run_webhook"
}
