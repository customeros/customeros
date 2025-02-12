package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

var agentGoalNameMap = map[enum.AgentGoal]string{
	enum.AgentGoalEvaluateICPFit:         "Qualify companies",
	enum.AgentGoalIdentifyWebVisitor:     "Identify companies that visit my website",
	enum.AgentGoalSpotHelpNeeded:         "Spot companies that may need help",
	enum.AgentGoalCaptureExternalMeeting: "Capture and share external meetings",
	enum.AgentGoalReceiveReply:           "Receive a reply to our outreach",
	enum.AgentGoalGetPaid:                "Get paid",
}

func MapAgentGoalName(input enum.AgentGoal) string {
	return agentGoalNameMap[input]
}
