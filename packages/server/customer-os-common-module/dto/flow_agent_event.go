package dto

import (
	neoEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
)

type FlowAgentEvent struct {
	FlowExecutionId  string
	ExternalSystemId neoEnum.ExternalSystemId
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
