package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

var agentGoalNameMap = map[enum.AgentGoal]string{
	enum.AgentGoalEvaluateICPFit:         "Evaluate company ICP fit",
	enum.AgentGoalIdentifyWebVisitor:     "Identify website visitors",
	enum.AgentGoalSpotHelpNeeded:         "Spot compaines that may need help",
	enum.AgentGoalCaptureExternalMeeting: "Capture and share external meetings",
	enum.AgentGoalReceiveReply:           "Receive a reply to our outreach",
}

func MapAgentGoalName(input enum.AgentGoal) string {
	return agentGoalNameMap[input]
}
