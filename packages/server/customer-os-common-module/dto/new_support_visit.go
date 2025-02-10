package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type NewSupportVisit struct{}

func (e NewSupportVisit) Name() enum.AgentListenerEvent {
	return enum.EventNewSupportVisit
}
