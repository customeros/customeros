package entity

import (
	"time"
)

type AutomationExecution struct {
	ID             string     `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	AutomationID   string     `gorm:"primaryKey;type:varchar(50)" json:"automationId"`
	TriggerEventID string     `gorm:"column:trigger_event_id;type:varchar(50)" json:"triggerEventId"`
	AgentID        string     `gorm:"column:agent_id;type:varchar(255)" json:"agentId"`
	FlowID         string     `gorm:"column:flow_id;type:varchar(255);not null;index" json:"flowId"`
	Status         string     `gorm:"column:status;type:varchar(255);not null" json:"status"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	StartedAt      time.Time  `gorm:"column:started_at;autoCreateTime" json:"startedAt"`
	CompletedAt    time.Time  `gorm:"column:completed_at;autoCreateTime" json:"completedAt"`
	UpdatedAt      *time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	ErrorMessage   *string    `gorm:"column:error_message;type:text" json:"errorMessage"`
}

func (AutomationExecution) TableName() string {
	return "automation_execution"
}
