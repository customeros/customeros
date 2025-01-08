package entity

import (
	"time"
)

type AgentExecution struct {
	ID                    string     `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	AutomationExecutionID string     `gorm:"column:automation_execution_id;type:uuid" json:"automationExecutionId"`
	AgentID               string     `gorm:"column:agent_id;type:varchar(50);not null" json:"agent_id" binding:"required"`
	FlowNodeID            string     `gorm:"column:flow_node_id;type:varchar(50)" json:"flowNodeId"`
	Status                string     `gorm:"column:status;type:varchar(255);not null;default:'pending'" json:"status"`
	InputPayload          string     `gorm:"column:input_payload;type:text" json:"input"`
	OutputPayload         string     `gorm:"column:output_payload;type:text" json:"output"`
	StartedAt             *time.Time `gorm:"column:started_at" json:"startedAt"`
	CompletedAt           *time.Time `gorm:"column:completed_at" json:"completedAt"`
	ErrorMessage          *string    `gorm:"column:error_message;type:text" json:"errorMessage"`
	CreatedAt             time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	RetryCount            int        `gorm:"column:retry_count" json:"retryCount"`
}

func (AgentExecution) TableName() string {
	return "agent_execution"
}
