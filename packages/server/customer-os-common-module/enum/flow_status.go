package enum

import "fmt"

type FlowStatus string

const (
	FlowStatusArchived FlowStatus = "archived"
	FlowStatusOff      FlowStatus = "off"
	FlowStatusOn       FlowStatus = "on"
)

func (t FlowStatus) String() string {
	return string(t)
}

func GetFlowStatus(s string) (FlowStatus, error) {
	switch FlowStatus(s) {
	case
		FlowStatusArchived,
		FlowStatusOff,
		FlowStatusOn:

		return FlowStatus(s), nil

	default:
		return "", fmt.Errorf("invalid FlowStatus: %s", s)
	}
}
