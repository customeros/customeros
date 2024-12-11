package entity

import (
	"time"
)

type FlowExecution struct {
	ID             string     `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant         string     `gorm:"column:tenant;type:varchar(255);not null;index" json:"tenant" binding:"required"`
	FlowID         string     `gorm:"column:flow_id;type:varchar(255);not null;index" json:"flowId" binding:"required"`
	EntityID       string     `gorm:"column:entity_id;type:varchar(255);not null;index" json:"entityId" binding:"required"`
	EntityType     string     `gorm:"column:entity_type;type:varchar(255);not null" json:"entityType" binding:"required"`
	Status         string     `gorm:"column:status;type:varchar(255);not null" json:"status"`
	ScheduledFor   *time.Time `gorm:"column:scheduled_for" json:"scheduledFor"`
	StartedAt      *time.Time `gorm:"column:started_at;not null" json:"startedAt"`
	CompletedAt    *time.Time `gorm:"column:completed_at" json:"completedAt"`
	ErrorMessage   *string    `gorm:"column:error_message;type:text" json:"errorMessage"`
	BlockedReason  *string    `gorm:"column:blocked_reason;type:text" json:"blockedReason"`
	NextStep       string     `gorm:"column:nextStep;type:varchar(255);index" json:"nextStep"`
	NextStepNodeId string     `gorm:"column:next_step_node_id;type:varchar(255);index" json:"nextStepNodeId"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Context        *string    `gorm:"column:context;type:text" json:"context"`
}

func (FlowExecution) TableName() string {
	return "flow_execution"
}
