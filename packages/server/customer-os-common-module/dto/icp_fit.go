package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

var (
	_ interfaces.AgentEvents = IcpFit{}
	_ interfaces.AgentEvents = IcpNotAFit{}
)

type IcpFit struct {
	AgentExecutionId string `json:"agentExecutionId"`
}

func (e IcpFit) Name() enum.AgentListenerEvent {
	return enum.EventICPFit
}

type IcpNotAFit struct {
	AgentExecutionId string `json:"agentExecutionId"`
}

func (e IcpNotAFit) Name() enum.AgentListenerEvent {
	return enum.EventICPNotAFit
}
