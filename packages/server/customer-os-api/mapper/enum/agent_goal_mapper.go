package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

var agentGoalNameMap = map[enum.AgentGoal]string{
	enum.AgentGoalEvaluateICPFit:         "Evaluate the Ideal Customer Profile (ICP) fit of a company",
	enum.AgentGoalIdentifyWebVisitor:     "Identify companies that visit my website",
	enum.AgentGoalSpotHelpNeeded:         "Flag companies that may need help",
	enum.AgentGoalCaptureExternalMeeting: "Capture and share external meetings with the team",
	enum.AgentGoalReceiveReply:           "Receive a reply to our outreach",
	enum.AgentGoalGetPaid:                "Get paid",
}

func MapAgentGoalName(input enum.AgentGoal) string {
	return agentGoalNameMap[input]
}
