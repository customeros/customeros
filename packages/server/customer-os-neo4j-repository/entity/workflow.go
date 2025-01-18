package entity

import (
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

// Main workflow definition
type WorkflowStatus string

const (
	WorkflowStatusInactive WorkflowStatus = "INACTIVE"
	WorkflowStatusActive   WorkflowStatus = "ACTIVE"
	WorkflowStatusArchived WorkflowStatus = "ARCHIVED"
)

func GetWorkflowStatus(s string) WorkflowStatus {
	return WorkflowStatus(s)
}

func (w WorkflowStatus) String() string {
	return string(w)

}

type Workflow struct {
	ID             string                   `json:"id"`
	TenantID       string                   `json:"tenantId"`
	Name           string                   `json:"name"`
	Description    string                   `json:"description"`
	Status         WorkflowStatus           `json:"status"`
	ListenerEvents []enum.FlowListenerEvent `json:"listenerEvent"`
	Nodes          []WorkflowNode           `json:"nodes"`
	Edges          []WorkflowEdge           `json:"edges"`
	CreatedAt      time.Time                `json:"createdAt"`
	UpdatedAt      time.Time                `json:"updatedAt"`
	CreatedBy      string                   `json:"createdBy"`
	LastModifiedBy string                   `json:"lastModifiedBy"`
}

// Node definition
type WorkflowNode struct {
	ID             string            `json:"id"`
	Type           enum.FlowNodeType `json:"flowNodeType"`
	Event          string            `json:"event"`
	Position       Position          `json:"position"`
	Width          float64           `json:"width"`
	Height         float64           `json:"height"`
	SourcePosition string            `json:"sourcePosition"`
	TargetPosition string            `json:"targetPosition"`
	InternalID     string            `json:"internalId,omitempty"`
}

// Node position
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Edge definition
type WorkflowEdge struct {
	ID        string    `json:"id"`
	Source    string    `json:"source"`
	Target    string    `json:"target"`
	Type      string    `json:"type"`
	MarkerEnd MarkerEnd `json:"markerEnd"`
	Condition string    `json:"condition,omitempty"`
}

type MarkerEnd struct {
	Type   string `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Wait node config

type DurationUnit string

const (
	DurationMinutes DurationUnit = "minutes"
	DurationHours   DurationUnit = "hours"
	DurationDays    DurationUnit = "days"
)

type WaitNodeConfig struct {
	Duration     int    `json:"waitDuration"`
	DurationUnit string `json:"waitDurationUnit"`
	NextNodeID   string `json:"nextNodeId"`
}
