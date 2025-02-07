package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

var _ interfaces.AgentEvents = NewSupportVisit{}

type NewSupportVisit struct{}

func (e NewSupportVisit) Name() enum.AgentListenerEvent {
	return enum.EventNewSupportVisit
}
