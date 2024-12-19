package enum

import "fmt"

type FlowNodeEdgeStatus string

const (
	FlowNodeEdgeStatusActive   FlowNodeEdgeStatus = "active"
	FlowNodeEdgeStatusArchived FlowNodeEdgeStatus = "archived"
)

func (t FlowNodeEdgeStatus) String() string {
	return string(t)
}

func GetFlowNodeEdgeStatus(s string) (FlowStatus, error) {
	switch FlowNodeEdgeStatus(s) {
	case
		FlowNodeEdgeStatusActive,
		FlowNodeEdgeStatusArchived:

		return FlowStatus(s), nil

	default:
		return "", fmt.Errorf("invalid FlowNodeEdgeStatus: %s", s)
	}
}
