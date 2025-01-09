package entity

import (
	"time"
)

type AgentExecution struct {
	ID           string     `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	AgentID      *string    `gorm:"column:agent_id;type:varchar(50);not null" json:"agent_id" binding:"required"`
	TriggerEvent string     `gorm:"column:trigger_event;type:varchar(50)" json:"triggerEvent"`
	FlowID       *string    `gorm:"column:flow_id;type:varchar(255);" json:"flowId"`
	Status       string     `gorm:"column:status;type:varchar(255);not null;default:'pending'" json:"status"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	StartedAt    *time.Time `gorm:"column:started_at" json:"startedAt"`
	UpdatedAt    *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	CompletedAt  *time.Time `gorm:"column:completed_at" json:"completedAt"`
	ErrorMessage *string    `gorm:"column:error_message;type:text" json:"errorMessage"`
}

func (AgentExecution) TableName() string {
	return "agent_execution"
}
