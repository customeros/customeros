package dto

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"

	commonenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
)

type FlowActionEvent[T any] struct {
	ExternalSystemId enum.ExternalSystemId
	SourceEvent      commonenum.FlowEvent
	Name             commonenum.FlowAction
	Data             *T
}

func NewFlowActionEvent(eventName commonenum.FlowAction, externalSystem enum.ExternalSystemId, sourceEvent commonenum.FlowEvent, data any) FlowActionEvent[any] {
	return FlowActionEvent[any]{
		ExternalSystemId: externalSystem,
		SourceEvent:      sourceEvent,
		Data:             &data,
	}
}
