package postgres_entity

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type AgentExecution struct {
	ID           string                    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	AgentID      *string                   `gorm:"column:agent_id;type:varchar(50);not null" json:"agent_id" binding:"required"`
	Tenant       string                    `gorm:"column:tenant;type:varchar(255)" json:"tenant"`
	TriggerEvent string                    `gorm:"column:trigger_event;type:varchar(50)" json:"triggerEvent"`
	FlowID       *string                   `gorm:"column:flow_id;type:varchar(255);" json:"flowId"`
	Status       enum.AgentExecutionStatus `gorm:"column:status;type:varchar(255);not null;default:'PENDING'" json:"status"`
	CreatedAt    time.Time                 `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	StartedAt    *time.Time                `gorm:"column:started_at" json:"startedAt"`
	UpdatedAt    *time.Time                `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	CompletedAt  *time.Time                `gorm:"column:completed_at" json:"completedAt"`
	ErrorMessage *string                   `gorm:"column:error_message;type:text" json:"errorMessage"`
	GoalAchieved *bool                     `gorm:"column:goal_achieved;type:boolean" json:"goalAchieved"`
	TraceId      string                    `gorm:"column:trace_id;type:varchar(255)" json:"traceId"`

	// Retry related fields
	RetryCount  int        `gorm:"column:retry_count;type:int;default:0" json:"retryCount"`
	MaxRetries  int        `gorm:"column:max_retries;type:int;default:3" json:"maxRetries"`
	NextRetryAt *time.Time `gorm:"column:next_retry_at" json:"nextRetryAt"`

	// Async state management
	StateData   map[string]any `gorm:"column:state_data;type:jsonb" json:"stateData"`            // Stores execution state for resume
	CurrentStep string         `gorm:"column:current_step;type:varchar(255)" json:"currentStep"` // Current capability being executed
	Checkpoints map[string]any `gorm:"column:checkpoints;type:jsonb" json:"checkpoints"`         // Stores completion state of each step
}

func (AgentExecution) TableName() string {
	return "agent_execution"
}
