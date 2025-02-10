package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

var agentGoalNameMap = map[enum.AgentGoal]string{
	enum.AgentGoalEvaluateICPFit:     "A company is qualified",
	enum.AgentGoalIdentifyWebVisitor: "Identify website visitors",

	enum.AgentGoalSpotHelpNeeded: "Spot Help Needed",
}

func MapAgentGoalName(input enum.AgentGoal) string {
	return agentGoalNameMap[input]
}
