package dto

import "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

type CapabilityExecutionContainer struct {
	AgentID               string
	AgentExecutionID      string
	Capability            enum.AgentCapabilityType
	CapabilityExecutionID string
	InputData             any
	OutputData            any
}
