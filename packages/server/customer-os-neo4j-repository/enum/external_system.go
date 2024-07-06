package enum

type ExternalSystemId string

const (
	Hubspot        ExternalSystemId = "hubspot"
	ZendeskSupport ExternalSystemId = "zendesk_support"
	CalCom         ExternalSystemId = "calcom"
	Pipedrive      ExternalSystemId = "pipedrive"
	Slack          ExternalSystemId = "slack"
	Intercom       ExternalSystemId = "intercom"
	Salesforce     ExternalSystemId = "salesforce"
	Stripe         ExternalSystemId = "stripe"
	Mixpanel       ExternalSystemId = "mixpanel"
	Close          ExternalSystemId = "close"
	Unthread       ExternalSystemId = "unthread"
	Outlook        ExternalSystemId = "outlook"
	Attio          ExternalSystemId = "attio"
	WeConnect      ExternalSystemId = "weconnect"
)

func (e ExternalSystemId) String() string {
	return string(e)
}

func DecodeExternalSystemId(value string) ExternalSystemId {
	switch value {
	case "hubspot":
		return Hubspot
	case "zendesk_support":
		return ZendeskSupport
	case "calcom":
		return CalCom
	case "pipedrive":
		return Pipedrive
	case "slack":
		return Slack
	case "intercom":
		return Intercom
	case "salesforce":
		return Salesforce
	case "stripe":
		return Stripe
	case "mixpanel":
		return Mixpanel
	case "close":
		return Close
	case "unthread":
		return Unthread
	case "outlook":
		return Outlook
	case "attio":
		return Attio
	case "weconnect":
		return WeConnect
	}
	return ""
}
