package enum

import "fmt"

type FlowStatus string

const (
	FlowStatusActive   FlowStatus = "ACTIVE"
	FlowStatusArchived FlowStatus = "ARCHIVED"
	FlowStatusInactive FlowStatus = "INACTIVE"
)

func (t FlowStatus) String() string {
	return string(t)
}

func GetFlowStatus(s string) (FlowStatus, error) {
	switch FlowStatus(s) {
	case
		FlowStatusActive,
		FlowStatusArchived,
		FlowStatusInactive:

		return FlowStatus(s), nil

	default:
		return "", fmt.Errorf("invalid FlowStatus: %s", s)
	}
}
