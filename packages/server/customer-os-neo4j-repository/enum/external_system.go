package enum

type ExternalSystemId string

const (
	Attio          ExternalSystemId = "attio"
	CalCom         ExternalSystemId = "calcom"
	Close          ExternalSystemId = "close"
	Fathom         ExternalSystemId = "fathom"
	GCal           ExternalSystemId = "gcal"
	Gain           ExternalSystemId = "grain"
	GMail          ExternalSystemId = "gmail"
	Hubspot        ExternalSystemId = "hubspot"
	Intercom       ExternalSystemId = "intercom"
	Mailstack      ExternalSystemId = "mailstack"
	Mixpanel       ExternalSystemId = "mixpanel"
	Outlook        ExternalSystemId = "outlook"
	Pipedrive      ExternalSystemId = "pipedrive"
	Postmark       ExternalSystemId = "postmark"
	Salesforce     ExternalSystemId = "salesforce"
	Slack          ExternalSystemId = "slack"
	Stripe         ExternalSystemId = "stripe"
	Unthread       ExternalSystemId = "unthread"
	WeConnect      ExternalSystemId = "weconnect"
	ZendeskSell    ExternalSystemId = "zendesk-sell"
	ZendeskSupport ExternalSystemId = "zendesk_support"
)

var validExternalSystems = func() map[string]ExternalSystemId {
	systems := []ExternalSystemId{
		Attio, CalCom, Close, Fathom, GCal, Gain, GMail,
		Hubspot, Intercom, Mailstack, Mixpanel, Outlook,
		Pipedrive, Postmark, Salesforce, Slack, Stripe,
		Unthread, WeConnect, ZendeskSell, ZendeskSupport,
	}

	m := make(map[string]ExternalSystemId)
	for _, s := range systems {
		m[string(s)] = s
	}
	return m
}()

func (e ExternalSystemId) String() string {
	return string(e)
}

func DecodeExternalSystemId(value string) ExternalSystemId {
	if system, ok := validExternalSystems[value]; ok {
		return system
	}
	return ""
}
