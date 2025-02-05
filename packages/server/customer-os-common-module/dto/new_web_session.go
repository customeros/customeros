package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

var _ interfaces.AgentEvents = (*NewWebSession)(nil)

type NewWebSession struct {
	WebSessionID string `json:"webSessionId"`
	IPAddress    string `json:"ipAddress"`
	VisitorID    string `json:"visitorId"`
	Hostname     string `json:"hostname"`
}

func (f *NewWebSession) Name() enum.AgentListenerEvent {
	return enum.EventNewWebSession
}
