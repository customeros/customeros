package enum

import "fmt"

type FlowAction string

const (
	ActionContactCreate       FlowAction = "contact.create"
	ActionEmailCreate         FlowAction = "email.create"
	ActionEmailAddressUpdate  FlowAction = "email_address.update"
	ActionTimelineEventCreate FlowAction = "timeline_event.create"
	ActionOrganizationCreate  FlowAction = "organization.create"
	ActionNotSet              FlowAction = ""
)

func (t FlowAction) String() string {
	return string(t)
}

func GetFlowAction(s string) (FlowAction, error) {
	switch FlowAction(s) {
	case
		ActionContactCreate,
		ActionEmailCreate,
		ActionEmailAddressUpdate,
		ActionOrganizationCreate,
		ActionTimelineEventCreate:

		return FlowAction(s), nil

	default:
		return "", fmt.Errorf("invalid FlowAction: %s", s)
	}
}
