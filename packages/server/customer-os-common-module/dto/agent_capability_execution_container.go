package dto

import "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

type CapabilityExecutionContainer struct {
	AgentID               string
	AgentExecutionID      string
	Capability            enum.AgentCapability
	CapabilityExecutionID string
	InputData             any
	OutputData            any
}
