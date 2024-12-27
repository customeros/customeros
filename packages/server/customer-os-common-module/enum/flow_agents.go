package enum

import (
	"fmt"
	"strings"
)

type FlowAgent string

const (
	AgentContactCreate       FlowAgent = "contact.create"
	AgentEmailSendNew        FlowAgent = "email.send_new"
	AgentEmailSendReply      FlowAgent = "email.send_reply"
	AgentLinkedinConnect     FlowAgent = "linkedin.connect"
	AgentLinkedinMessage     FlowAgent = "linkedin.message"
	AgentOrganizationCreate  FlowAgent = "organization.create"
	AgentSlackNotify         FlowAgent = "slack.notify"
	AgentTimelineEventCreate FlowAgent = "timeline_event.create"
)

func (t FlowAgent) String() string {
	return string(t)
}

func (t FlowAgent) Agent() string {
	agent, _, _ := strings.Cut(t.String(), ".")
	return agent
}

func GetFlowAgent(s string) (FlowAgent, error) {
	switch FlowAgent(s) {
	case
		AgentContactCreate,
		AgentEmailSendNew,
		AgentEmailSendReply,
		AgentLinkedinConnect,
		AgentLinkedinMessage,
		AgentOrganizationCreate,
		AgentSlackNotify,
		AgentTimelineEventCreate:
		return FlowAgent(s), nil

	default:
		return "", fmt.Errorf("invalid FlowAgent: %s", s)
	}
}
