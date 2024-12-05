package enum

import (
	"fmt"
	"strings"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

type FlowEvent string

const (
	EventFathomMeetingSummaryCreated FlowEvent = "fathom.meeting_summary.created"
	EventGrainMeetingSummaryCreated  FlowEvent = "grain.meeting_summary.created"
	EventMailstackEmailReceived      FlowEvent = "mailstack.email.received"
	EventNotSet                      FlowEvent = ""
)

func (e FlowEvent) String() string {
	return string(e)
}

func (e FlowEvent) Parse() (system, resource, action string, err error) {
	parts := strings.Split(e.String(), ".")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid event name format: %s", e.String())
	}

	return parts[0], parts[1], parts[2], nil
}

func (e FlowEvent) ExternalSystem() (system enum.ExternalSystemId, err error) {
	systemId, _, _, err := e.Parse()
	if err != nil {
		return enum.ExternalSystemId(EventNotSet), fmt.Errorf("invalid event name")
	}

	return enum.DecodeExternalSystemId(systemId), nil
}

func GetFlowEvent(s string) (FlowEvent, error) {
	switch FlowEvent(s) {
	case
		EventFathomMeetingSummaryCreated,
		EventGrainMeetingSummaryCreated,
		EventMailstackEmailReceived:

		return FlowEvent(s), nil
	default:
		return "", fmt.Errorf("invalid FlowEvent: %s", s)
	}
}
