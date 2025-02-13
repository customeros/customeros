package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type NewEmail struct{}

func (e NewEmail) Name() enum.AgentListenerEvent {
	return enum.EventNewEmail
}
