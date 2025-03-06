package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type NewLead struct{}

func (e NewLead) ListenerEvent() enum.AgentListenerEvent {
	return enum.EventNewLead
}
