package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

var (
	_ interfaces.AgentEvents = CompanyIdentified{}
)

type CompanyIdentified struct {
	AgentExecutionId string `json:"agentExecutionId"`
}
