package dto

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	neoEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

type FlowActionEvent struct {
	FlowExecutionId  string
	ExternalSystemId neoEnum.ExternalSystemId
	SourceEvent      enum.FlowListenerEvent
	Name             enum.FlowAction
	DataType         string
	Data             any
}
