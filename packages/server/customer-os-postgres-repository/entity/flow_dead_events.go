package entity

import (
	"time"
)

type FlowDeadEvents struct {
	ID        string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant    string    `gorm:"column:tenant;type:varchar(255);not null" json:"tenant" binding:"required"`
	NodeType  string    `gorm:"column:nodeType;type:varchar(255);not null" json:"nodeType" binding:"required"`
	EventType string    `gorm:"column:event_type;type:varchar(255);not null" json:"eventType" binding:"required"`
	Event     string    `gorm:"column:event;type:varchar(255);not null" json:"event" binding:"required"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	Data      *string   `gorm:"column:data;type:text" json:"data"`
}

func (FlowDeadEvents) TableName() string {
	return "flow_dead_events"
}
