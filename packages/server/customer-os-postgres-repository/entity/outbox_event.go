package postgres_entity

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"gorm.io/gorm"
	"time"

	"github.com/customeros/customeros/packages/server/enums"
)

type OutboxStatus string

const (
	OutboxPending    OutboxStatus = "pending"
	OutboxProcessing OutboxStatus = "processing"
	OutboxCompleted  OutboxStatus = "completed"
	OutboxFailed     OutboxStatus = "failed"
)

type OutboxEvent struct {
	ID           string              `gorm:"column:id;primaryKey;type:varchar(50);"`
	EventType    enums.NatsEventType `gorm:"column:event_type;type:varchar(100);index;not null"`
	EntityID     string              `gorm:"column:entity;type:varchar(50)" json:"entityId"`
	Publisher    string              `gorm:"column:publisher;type:varchar(50);index;not null" json:"publisher"`
	Tenant       string              `gorm:"column:tenant;type:varchar(50);not null" json:"tenant"`
	SessionID    string              `gorm:"column:session_id;type:varchar(50)" json:"sessionId"`
	Payload      []byte              `gorm:"column:payload;type:bytea;not null"`
	Status       OutboxStatus        `gorm:"column:status;type:varchar(20);index;not null"`
	CreatedAt    time.Time           `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	ProcessedAt  *time.Time          `gorm:"column:processed_at;"`                  // When the event was processed
	RetryCount   int                 `gorm:"column:retry_count;type:int;default:0"` // For error handling
	ErrorMessage string              `gorm:"column:error_message;type:text;"`       // Last error message if failed
	LockUntil    *time.Time          `gorm:"column:lock_until;index"`               // For distributed processing
}

func (OutboxEvent) TableName() string {
	return "outbox_events"
}

func (o *OutboxEvent) BeforeCreate(tx *gorm.DB) error {
	if o.ID == "" {
		o.ID = utils.GenerateNanoIdWithPrefix("outbox", 21)
	}
	return nil
}
