package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type FlowAgentEvent struct {
	FlowExecutionId  string
	Tenant           string
	ExternalSystemId enum.Source
	SourceEvent      enum.AgentListenerEvent
	Name             enum.FlowAgent
	DataType         string
	Data             any
}

type FlowAgentExecutionResultEvent struct {
	FlowExecutionID      string
	FlowAgentExecutionID string
	Tenant               string
	Status               string
	ErrorMessage         *string
	Data                 any
	DataType             string
}
