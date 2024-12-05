package enum

import (
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type ExternalSystemId string

const (
	Attio          ExternalSystemId = "attio"
	CalCom         ExternalSystemId = "calcom"
	Close          ExternalSystemId = "close"
	Fathom         ExternalSystemId = "fathom"
	GCal           ExternalSystemId = "gcal"
	GMail          ExternalSystemId = "gmail"
	Grain          ExternalSystemId = "grain"
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
	NotSet         ExternalSystemId = ""
)

var validExternalSystems = func() map[string]ExternalSystemId {
	systems := []ExternalSystemId{
		Attio, CalCom, Close, Fathom, GCal, Grain, GMail,
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

func (e ExternalSystemId) IntegrationID(rotationCount int) string {
	input := fmt.Sprintf("%s:%d", e.String(), rotationCount)
	return utils.GenerateHashId(input, 12)
}
