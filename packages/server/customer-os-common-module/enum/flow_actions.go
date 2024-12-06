package enum

import "fmt"

type FlowAction string

const (
	ActionContactCreate       FlowAction = "contact.create"
	ActionOrganizationCreate  FlowAction = "organization.create"
	ActionTimelineEventCreate FlowAction = "timeline_event.create"
)

func (t FlowAction) String() string {
	return string(t)
}

func GetFlowAction(s string) (FlowAction, error) {
	switch FlowAction(s) {
	case
		ActionContactCreate,
		ActionOrganizationCreate,
		ActionTimelineEventCreate:
		return FlowAction(s), nil

	default:
		return "", fmt.Errorf("invalid FlowAction: %s", s)
	}
}
