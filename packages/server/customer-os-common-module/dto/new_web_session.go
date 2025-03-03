package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type NewWebSession struct{}

func (f NewWebSession) Name() enum.AgentListenerEvent {
	return enum.EventNewWebSession
}
