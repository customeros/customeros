package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

var (
	_ interfaces.AgentEvents = (*IcpFit)(nil)
	_ interfaces.AgentEvents = (*IcpNotAFit)(nil)
)

type IcpFit struct{}

func (e *IcpFit) Name() enum.AgentListenerEvent {
	return enum.EventICPFit
}

type IcpNotAFit struct{}

func (e *IcpNotAFit) Name() enum.AgentListenerEvent {
	return enum.EventICPNotAFit
}
