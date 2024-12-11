package entity

import (
	"time"
)

type ActionExecution struct {
	ID              string     `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	FlowExecutionID string     `gorm:"column:flow_execution_id;type:uuid;default:gen_random_uuid()" json:"flowExecutionId"`
	FlowNodeID      string     `gorm:"column:flow_node_id;type:varchar(255)" json:"flowNodeId"`
	Action          string     `gorm:"column:action;type:varchar(255);not null" json:"action" binding:"required"`
	Status          string     `gorm:"column:status;type:varchar(255);not null;default:'pending'" json:"status"`
	ScheduledFor    *time.Time `gorm:"column:scheduled_for" json:"scheduledFor"`
	StartedAt       *time.Time `gorm:"column:started_at" json:"startedAt"`
	CompletedAt     *time.Time `gorm:"column:completed_at" json:"completedAt"`
	ErrorMessage    *string    `gorm:"column:error_message;type:text" json:"errorMessage"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	Result          *string    `gorm:"column:result;type:text" json:"result"`
}

func (ActionExecution) TableName() string {
	return "flow_action_execution"
}
