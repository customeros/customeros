package enum

import "fmt"

type FlowExecutionStatus string

const (
	FlowExecutionBlocked   FlowExecutionStatus = "BLOCKED"
	FlowExecutionCompleted FlowExecutionStatus = "COMPLETED"
	FlowExecutionError     FlowExecutionStatus = "ERROR"
	FlowExecutionRunning   FlowExecutionStatus = "RUNNING"
)

func (t FlowExecutionStatus) String() string {
	return string(t)
}

func GetFlowExecutionStatus(s string) (FlowExecutionStatus, error) {
	switch FlowExecutionStatus(s) {
	case
		FlowExecutionBlocked,
		FlowExecutionCompleted,
		FlowExecutionError,
		FlowExecutionRunning:
		return FlowExecutionStatus(s), nil

	default:
		return "", fmt.Errorf("invalid FlowExecutionStatus: %s", s)
	}
}

type FlowBlockedReason string

const (
	FlowBlockedNotActive      FlowBlockedReason = "FLOW_INACTIVE"
	FlowBlockedArchived       FlowBlockedReason = "FLOW_ARCHIVED"
	FlowBlockedIncompleteData FlowBlockedReason = "INCOMPLETE_DATA"
)

func (t FlowBlockedReason) String() string {
	return string(t)
}

func GetFlowBlockedReason(s string) (FlowBlockedReason, error) {
	switch FlowBlockedReason(s) {
	case
		FlowBlockedArchived,
		FlowBlockedIncompleteData,
		FlowBlockedNotActive:
		return FlowBlockedReason(s), nil

	default:
		return "", fmt.Errorf("invalid FlowBlockedReason: %s", s)
	}
}
