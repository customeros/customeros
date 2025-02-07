package enummapper

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

var agentGoalNameMap = map[enum.AgentGoal]string{
	enum.AgentGoalEvaluateICPFit:     "A company is qualified",
	enum.AgentGoalIdentifyWebVisitor: "A website visitor is identified",

	enum.AgentGoalSpotHelpNeeded: "Spot Help Needed",
}

func MapAgentGoalName(input enum.AgentGoal) string {
	return agentGoalNameMap[input]
}
