package entity

import (
	"time"
)

type FlowDeadListenerEvents struct {
	ID            string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant        string    `gorm:"column:tenant;type:varchar(255);not null" json:"tenant" binding:"required"`
	ListenerEvent string    `gorm:"column:listener_event;type:varchar(255);not null" json:"listenerEvent" binding:"required"`
	EventType     string    `gorm:"column:event_type;type:varchar(255);not null" json:"eventType" binding:"required"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	Data          *string   `gorm:"column:data;type:text" json:"data"`
}

func (FlowDeadListenerEvents) TableName() string {
	return "flow_dead_listener_events"
}
