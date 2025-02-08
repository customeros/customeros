package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

var (
	_ interfaces.AgentEvents = WebVisitorIdentified{}
	_ interfaces.AgentEvents = WebVisitorNotIdentified{}
)

type WebVisitorIdentified struct {
	AgentExecutionId string `json:"agentExecutionId"`
}

type WebVisitorNotIdentified struct {
	AgentExecutionId string `json:"agentExecutionId"`
}
