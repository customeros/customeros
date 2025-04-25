package models

import (
	"time"

	"github.com/customeros/customeros/packages/server/leads/internal/enum"
)

type OutboxEvent struct {
	ID           string            `gorm:"column:id;primaryKey;type:varchar(50);"`
	EventType    enum.Events       `gorm:"column:event_type;type:varchar(100);index;not null"`
	Payload      []byte            `gorm:"column:payload;type:jsonb;not null"`
	Status       enum.OutboxStatus `gorm:"column:status;type:varchar(20);index;not null"`
	CreatedAt    time.Time         `gorm:"column:created_at;not null;index"`      // When the event was created
	ProcessedAt  *time.Time        `gorm:"column:processed_at;"`                  // When the event was processed
	RetryCount   int               `gorm:"column:retry_count;type:int;default:0"` // For error handling
	ErrorMessage string            `gorm:"column:error_message;type:text;"`       // Last error message if failed
	LockUntil    *time.Time        `gorm:"column:lock_until;index"`               // For distributed processing
	Version      int               `gorm:"column:version;type:int;default:1"`     // For optimistic locking
}

func (OutboxEvent) TableName() string {
	return "outbox_events"
}
