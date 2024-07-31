package entity

import "time"

type EmailTrackingEventType string

const (
	EmailTrackingEventTypeOpen      EmailTrackingEventType = "email_open"
	EmailTrackingEventTypeLinkClick EmailTrackingEventType = "email_link_click"
)

type EmailTracking struct {
	ID        string                 `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant    string                 `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	CreatedAt time.Time              `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time              `gorm:"column:updated_at;type:timestamp;DEFAULT:current_timestamp" json:"updatedAt"`
	Timestamp time.Time              `gorm:"column:timestamp;type:timestamp;DEFAULT:current_timestamp" json:"timestamp"`
	MessageId string                 `gorm:"column:message_id;type:varchar(64);NOT NULL" json:"messageId"`
	LinkId    string                 `gorm:"column:link_id;type:varchar(64);" json:"linkId"`
	EventType EmailTrackingEventType `gorm:"column:event_type;type:varchar(255);NOT NULL" json:"eventType"`
	IP        string                 `gorm:"column:ip;type:varchar(255);" json:"ip" `
}

func (EmailTracking) TableName() string {
	return "email_tracking"
}
