package entity

type DataSource string

const (
	DataSourceNA             DataSource = ""
	DataSourceAttio          DataSource = "attio"
	DataSourceClose          DataSource = "close"
	DataSourceFathom         DataSource = "fathom"
	DataSourceGmail          DataSource = "gmail"
	DataSourceGrain          DataSource = "grain"
	DataSourceHubspot        DataSource = "hubspot"
	DataSourceIntercom       DataSource = "intercom"
	DataSourceMailstack      DataSource = "mailstack"
	DataSourceMixpanel       DataSource = "mixpanel"
	DataSourceOpenline       DataSource = "openline"
	DataSourceOutlook        DataSource = "outlook"
	DataSourcePipedrive      DataSource = "pipedrive"
	DataSourceSalesforce     DataSource = "salesforce"
	DataSourceShopify        DataSource = "shopify"
	DataSourceSlack          DataSource = "slack"
	DataSourceStripe         DataSource = "stripe"
	DataSourceUnthread       DataSource = "unthread"
	DataSourceWebscrape      DataSource = "webscrape"
	DataSourceZendeskSell    DataSource = "zendesk-sell"
	DataSourceZendeskSupport DataSource = "zendesk_support"
)

var AllDataSource = []DataSource{
	DataSourceAttio,
	DataSourceClose,
	DataSourceFathom,
	DataSourceGmail,
	DataSourceGrain,
	DataSourceHubspot,
	DataSourceIntercom,
	DataSourceMailstack,
	DataSourceMixpanel,
	DataSourceOpenline,
	DataSourceOutlook,
	DataSourcePipedrive,
	DataSourceSalesforce,
	DataSourceShopify,
	DataSourceSlack,
	DataSourceStripe,
	DataSourceUnthread,
	DataSourceWebscrape,
	DataSourceZendeskSell,
	DataSourceZendeskSupport,
	DataSourceFathom,
	DataSourceGrain,
}

func (ds DataSource) String() string {
	return string(ds)
}

func DecodeDataSource(s string) DataSource {
	if IsValidDataSource(s) {
		return DataSource(s)
	}
	return DataSourceNA
}

func IsValidDataSource(s string) bool {
	for _, ds := range AllDataSource {
		if ds == DataSource(s) {
			return true
		}
	}
	return false
}
