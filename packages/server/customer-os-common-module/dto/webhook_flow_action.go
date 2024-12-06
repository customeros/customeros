package dto

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"

	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
)

type FlowActionEvent struct {
	ExternalSystemId enum.ExternalSystemId
	SourceEvent      commonenum.FlowListenerEvent
	Name             commonenum.FlowAction
	DataType         string
	Data             any
}
