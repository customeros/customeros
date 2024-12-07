package enum

import (
	"fmt"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

type FlowListenerEvent string

const (
	EventFathomMeetingSummaryCreated FlowListenerEvent = "fathom.meeting_summary.created"
	EventFlowContactAdded            FlowListenerEvent = "flow.contact.added"
	EventGrainMeetingSummaryCreated  FlowListenerEvent = "grain.meeting_summary.created"
	NotSet                           FlowListenerEvent = ""
)

func (e FlowListenerEvent) String() string {
	return string(e)
}

func (e FlowListenerEvent) Parse() (system, resource, action string, err error) {
	parts := strings.Split(e.String(), ".")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid event name format: %s", e.String())
	}

	return parts[0], parts[1], parts[2], nil
}

func (e FlowListenerEvent) ExternalSystem() (system enum.ExternalSystemId, err error) {
	systemId, _, _, err := e.Parse()
	if err != nil {
		return enum.ExternalSystemId(NotSet), fmt.Errorf("invalid event name")
	}

	return enum.DecodeExternalSystemId(systemId), nil
}

func GetFlowEvent(s string) (FlowListenerEvent, error) {
	switch FlowListenerEvent(s) {
	case EventFathomMeetingSummaryCreated, EventGrainMeetingSummaryCreated:
		return FlowListenerEvent(s), nil
	default:
		return "", fmt.Errorf("invalid FlowListenerEvent: %s", s)
	}
}
