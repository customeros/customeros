package entity

type DataSource string

const (
	DataSourceNA             DataSource = ""
	DataSourceOpenline       DataSource = "openline"
	DataSourceGmail          DataSource = "gmail"
	DataSourceHubspot        DataSource = "hubspot"
	DataSourceZendeskSupport DataSource = "zendesk_support"
	DataSourcePipedrive      DataSource = "pipedrive"
	DataSourceSlack          DataSource = "slack"
	DataSourceWebscrape      DataSource = "webscrape"
	DataSourceIntercom       DataSource = "intercom"
	DataSourceSalesforce     DataSource = "salesforce"
	DataSourceStripe         DataSource = "stripe"
	DataSourceMixpanel       DataSource = "mixpanel"
)

var AllDataSource = []DataSource{
	DataSourceOpenline,
	DataSourceHubspot,
	DataSourceZendeskSupport,
	DataSourcePipedrive,
	DataSourceSlack,
	DataSourceWebscrape,
	DataSourceIntercom,
	DataSourceSalesforce,
	DataSourceStripe,
	DataSourceMixpanel,
}

func (ds DataSource) String() string {
	return string(ds)
}

func GetDataSource(s string) DataSource {
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
