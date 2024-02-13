package entity

import "time"

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
	Outlook        ExternalSystemId = "outlook"
)

const (
	ExternalNodeContact string = "Contact"
	ExternalNodeMeeting string = "Meeting"
)

type ExternalSystemEntity struct {
	ExternalSystemId ExternalSystemId
	Relationship     struct {
		ExternalId     string
		SyncDate       *time.Time
		ExternalUrl    *string
		ExternalSource *string
	}
	DataloaderKey string
}

type ExternalSystemEntities []ExternalSystemEntity

func ExternalSystemTypeFromString(input string) ExternalSystemId {
	for _, v := range []ExternalSystemId{
		Hubspot, ZendeskSupport, CalCom, Pipedrive, Slack, Intercom, Salesforce, Stripe, Mixpanel, Close, Outlook,
	} {
		if string(v) == input {
			return v
		}
	}
	// Return a default value or handle the case when the input string doesn't match any ExternalSystemId
	return ""
}
