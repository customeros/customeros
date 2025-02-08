package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

var (
	_ interfaces.AgentEvents = IcpFit{}
	_ interfaces.AgentEvents = IcpNotAFit{}
)

type IcpFit struct {
	AgentExecutionId string `json:"agentExecutionId"`
}

type IcpNotAFit struct {
	AgentExecutionId string `json:"agentExecutionId"`
}
