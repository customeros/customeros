package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

var (
	_ interfaces.AgentEvents = WebVisitorIdentified{}
	_ interfaces.AgentEvents = WebVisitorNotIdentified{}
)

type WebVisitorIdentified struct{}

func (e WebVisitorIdentified) Name() enum.AgentListenerEvent {
	return enum.EventWebVisitorIdentified
}

type WebVisitorNotIdentified struct{}

func (e WebVisitorNotIdentified) Name() enum.AgentListenerEvent {
	return enum.EventWebVisitorNotIdentified
}
