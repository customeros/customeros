package interfaces

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type AgentListenerUntyped interface {
	Type() enum.AgentListenerEvent
	Name() string
	DefaultConfig() any
	SubscribedAgents() []enum.AgentType
}
