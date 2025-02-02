package agent_capability

import "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"

type NoConfig struct{}

type CapabilityOutput struct {
	ExecutionValidated bool `json:"executionValidated"`
	Completed          bool `json:"completed"`
}

func GetCapabilityConfigStruct(capabilityType enum.AgentCapability) any {
	switch capabilityType {
	case enum.CapabilitySendSlackNotification:
		return &SendSlackNotificationConfig{}
	case enum.CapabilitySendWebVisitorSlackNotification:
		return &SendWebVisitorSlackNotificationConfig{
			ChannelID: SlackChannelIdConfig{
				Value: "",
			},
			CooldownHours: SlackCooldownHoursConfig{
				Value: 12,
			},
		}
	case enum.CapabilityIdentifyWebVisitor:
		return &IdentifyWebsiteVisitorConfig{
			Websites: WebsitesConfig{
				Value: []string{},
			},
		}
	case enum.CapabilityApplyTag:
		return &ApplyTagConfig{}
	case enum.CapabilityEvaluateCompanyICPFit:
		return &ICPQualificationConfig{}
	default:
		return nil
	}
}
