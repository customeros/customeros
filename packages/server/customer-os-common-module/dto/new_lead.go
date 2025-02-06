package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

var _ interfaces.AgentEvents = NewLead{}

type NewLead struct{}

func (e NewLead) Name() enum.AgentListenerEvent {
	return enum.EventNewLead
}
