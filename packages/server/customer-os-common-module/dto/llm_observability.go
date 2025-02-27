package dto

import "time"

type LLMObservability struct {
	// Basic request metadata
	RequestID string    `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
	UserID    string    `json:"user_id,omitempty"`
	Tenant    string    `json:"tenant,omitempty"`

	// Model information
	Model       string `json:"model"`
	Temperature string `json:"temperature"`

	// Request details
	PromptTokens   int `json:"prompt_tokens"`
	ResponseTokens int `json:"response_tokens"`
	TotalTokens    int `json:"total_tokens"`

	// Cost tracking
	CostUSD float64 `json:"cost_usd,omitempty"`

	// Status information
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error_message,omitempty"`

	// Tracing context
	TraceID      string `json:"trace_id,omitempty"` // For Jaeger integration
	SpanID       string `json:"span_id,omitempty"`
	ParentSpanID string `json:"parent_span_id,omitempty"`

	// Content analysis
	SystemPrompt string `json:"systemPrompt,omitempty"`
	Prompt       string `json:"prompt,omitempty"`
	Response     string `json:"response,omitempty"`
}
