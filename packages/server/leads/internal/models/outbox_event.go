package models

import (
	"time"

	"github.com/customeros/customeros/packages/server/leads/internal/enum"
)

type OutboxEvent struct {
	ID           string            `gorm:"column:id;primaryKey;type:varchar(50);"`
	EventType    enum.Events       `gorm:"column:event_type;type:varchar(100);index;not null"`
	EntityID     string            `gorm:"column:entity;type:varchar(50)" json:"entityId"`
	Publisher    enum.LeadsService `gorm:"column:publisher;type:varchar(50);index;not null" json:"publisher"`
	Tenant       string            `gorm:"column:tenant;type:varchar(50);not null" json:"tenant"`
	SessionID    string            `gorm:"column:session_id;type:varchar(50)" json:"sessionId"`
	Payload      []byte            `gorm:"column:payload;type:bytea;not null"`
	Status       enum.OutboxStatus `gorm:"column:status;type:varchar(20);index;not null"`
	CreatedAt    time.Time         `gorm:"column:created_at;not null;index"`      // When the event was created
	ProcessedAt  *time.Time        `gorm:"column:processed_at;"`                  // When the event was processed
	RetryCount   int               `gorm:"column:retry_count;type:int;default:0"` // For error handling
	ErrorMessage string            `gorm:"column:error_message;type:text;"`       // Last error message if failed
	LockUntil    *time.Time        `gorm:"column:lock_until;index"`               // For distributed processing
}

func (OutboxEvent) TableName() string {
	return "outbox_events"
}
