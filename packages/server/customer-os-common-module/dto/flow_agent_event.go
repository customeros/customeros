package dto

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
)

type FlowAgentEvent struct {
	FlowExecutionId  string
	Tenant           string
	ExternalSystemId enum.Source
	SourceEvent      enum.FlowListenerEvent
	Name             enum.FlowAgent
	DataType         string
	Data             any
}

type FlowAgentExecutionResultEvent struct {
	FlowExecutionID      string
	FlowAgentExecutionID string
	Tenant               string
	Status               enum.FlowAgentExecutionStatus
	ErrorMessage         *string
	Data                 any
	DataType             string
}
