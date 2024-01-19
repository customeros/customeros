package entity

import (
	"time"
)

type EventBuffer struct {
	Tenant          string    `gorm:"column:tenant;size:50"`
	UUID            string    `gorm:"column:uuid;size:250;primaryKey"`
	ExpiryTimestamp time.Time `gorm:"column:expiry_timestamp;size:100"`
	CreatedDate     time.Time `gorm:"default:current_timestamp"`
	// event store Event fields
	EventType          string    `gorm:"column:event_type;size:40"`
	EventData          []byte    `gorm:"column:event_data;"`
	EventID            string    `gorm:"column:event_id;size:50"`
	EventTimestamp     time.Time `gorm:"column:event_timestamp;size:100"`
	EventAggregateType string    `gorm:"column:event_aggregate_type;size:40"`
	EventAggregateID   string    `gorm:"column:event_aggregate_id;size:50"`
	EventVersion       int64     `gorm:"column:event_version;"`
	EventMetadata      []byte    `gorm:"column:event_metadata;"`
}

func (EventBuffer) TableName() string {
	return "event_buffer"
}
