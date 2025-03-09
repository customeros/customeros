package dto

import "time"

type AgentExecutionObservability struct {
	SkipPublishingObservability bool

	// Basic request metadata
	CapabilityExecutionID string `json:"capability_execution_id"`
	ExecutionID           string `json:"agent_execution_id"`
	AgentID               string `json:"agent_id"`
	AgentType             string `json:"agent_type"`
	AgentScope            string `json:"agent_scope"`
	Tenant                string `json:"tenant,omitempty"`
	UserID                string `json:"user_id,omitempty"`
	Capability            string `json:"capability"`
	TriggerEvent          string `json:"trigger_event"`

	// Execution context
	InputData  any `json:"input_data"`
	OutputData any `json:"output_data"`

	// Status information
	Status       string     `json:"status,omitempty"`
	Attempt      int        `json:"attempt,omitempty"`
	StartedAt    time.Time  `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	Success      bool       `json:"success"`
	ErrorMessage string     `json:"error_message,omitempty"`
	Retry        bool       `json:"retry"`
	RetryAt      *time.Time `json:"retry_at,omitempty"`

	// Tracing context
	TraceID string `json:"trace_id,omitempty"` // For Jaeger integration
}
