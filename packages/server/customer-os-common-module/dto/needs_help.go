package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

var (
	_ interfaces.AgentEvents = NeedsHelp{}
	_ interfaces.AgentEvents = DoesNotNeedHelp{}
)

type NeedsHelp struct{}

func (e NeedsHelp) Name() enum.AgentListenerEvent {
	return enum.EventNeedsHelp
}

type DoesNotNeedHelp struct{}

func (e DoesNotNeedHelp) Name() enum.AgentListenerEvent {
	return enum.EventDoesNotNeedHelp
}
