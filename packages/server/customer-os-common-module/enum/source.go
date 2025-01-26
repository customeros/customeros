package enum

import (
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type Source string

const (
	SourceUnknown        Source = ""
	SourceAgent          Source = "agent"
	SourceAttio          Source = "attio"
	SourceCalCom         Source = "calcom"
	SourceClose          Source = "close"
	SourceCustomerOS     Source = "customer-os"
	SourceFathom         Source = "fathom"
	SourceGCal           Source = "gcal"
	SourceGmail          Source = "gmail"
	SourceGrain          Source = "grain"
	SourceHubspot        Source = "hubspot"
	SourceIntercom       Source = "intercom"
	SourceMailstack      Source = "mailstack"
	SourceMixpanel       Source = "mixpanel"
	SourceOutlook        Source = "outlook"
	SourcePipedrive      Source = "pipedrive"
	SourcePostmark       Source = "postmark"
	SourceSalesforce     Source = "salesforce"
	SourceShopify        Source = "shopify"
	SourceSlack          Source = "slack"
	SourceStripe         Source = "stripe"
	SourceUnthread       Source = "unthread"
	SourceWebscrape      Source = "webscrape"
	SourceWebtracker     Source = "webtracker"
	SourceZendeskSell    Source = "zendesk-sell"
	SourceZendeskSupport Source = "zendesk_support"
)

var AllSources = []Source{
	SourceUnknown,
	SourceAgent,
	SourceAttio,
	SourceCalCom,
	SourceClose,
	SourceCustomerOS,
	SourceFathom,
	SourceGCal,
	SourceGmail,
	SourceGrain,
	SourceHubspot,
	SourceIntercom,
	SourceMailstack,
	SourceMixpanel,
	SourceOutlook,
	SourcePipedrive,
	SourcePostmark,
	SourceSalesforce,
	SourceShopify,
	SourceSlack,
	SourceStripe,
	SourceUnthread,
	SourceWebscrape,
	SourceZendeskSell,
	SourceZendeskSupport,
}

func (ds Source) String() string {
	return string(ds)
}

func DecodeSource(s string) Source {
	if IsValidSource(s) {
		return Source(s)
	}
	return SourceUnknown
}

func IsValidSource(s string) bool {
	for _, ds := range AllSources {
		if ds == Source(s) {
			return true
		}
	}
	return false
}

func (ds Source) IntegrationID(rotationCount int) string {
	input := fmt.Sprintf("%s:%d", ds.String(), rotationCount)
	return utils.GenerateHashId(input, 12)
}
