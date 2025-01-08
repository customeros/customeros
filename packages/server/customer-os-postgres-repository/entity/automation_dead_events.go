package entity

import (
	"time"
)

type AutomationDeadEvents struct {
	ID           string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	AutomationID string    `gorm:"primaryKey;type:varchar(50)" json:"automationId"`
	Tenant       string    `gorm:"column:tenant;type:varchar(255);not null" json:"tenant" binding:"required"`
	Event        string    `gorm:"column:event;type:varchar(255);not null" json:"event" binding:"required"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	Data         *string   `gorm:"column:data;type:text" json:"data"`
}

func (AutomationDeadEvents) TableName() string {
	return "automation_dead_events"
}
