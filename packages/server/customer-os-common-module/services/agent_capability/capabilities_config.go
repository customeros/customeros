package agent_capability

import "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

type NoConfig struct{}

func GetCapabilityConfigStruct(capabilityType enum.AgentCapabilityType) any {
	switch capabilityType {
	case enum.CapabilitySendSlackNotification:
		return &SendSlackNotificationConfig{}
	default:
		return &NoConfig{}
	}
}
