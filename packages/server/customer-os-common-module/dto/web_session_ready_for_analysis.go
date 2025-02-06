package dto

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

var _ interfaces.AgentEvents = WebSessionReadyForAnalysis{}

type WebSessionReadyForAnalysis struct{}

func (e WebSessionReadyForAnalysis) Name() enum.AgentListenerEvent {
	return enum.EventWebSessionReadyForAnalysis
}
