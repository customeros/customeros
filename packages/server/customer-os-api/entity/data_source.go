package entity

type DataSource string

const (
	DataSourceNA       DataSource = ""
	DataSourceOpenline DataSource = "openline"
	DataSourceHubspot  DataSource = "hubspot"
	DataSourceZendesk  DataSource = "zendesk"
)

var AllDataSource = []DataSource{
	DataSourceOpenline,
	DataSourceHubspot,
	DataSourceZendesk,
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
