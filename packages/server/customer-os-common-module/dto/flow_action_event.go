package dto

import (
	neoEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
)

type FlowActionEvent struct {
	FlowExecutionId  string
	ExternalSystemId neoEnum.ExternalSystemId
	SourceEvent      enum.FlowListenerEvent
	Name             enum.FlowAction
	DataType         string
	Data             any
}

type FlowActionExecutionResultEvent struct {
	FlowExecutionID       string
	FlowActionExecutionID string
	Tenant                string
	Status                enum.FlowActionExecutionStatus
	ErrorMessage          *string
	Data                  any
	DataType              string
}
