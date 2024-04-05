package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

var sourceByModel = map[model.DataSource]neo4jentity.DataSource{
	model.DataSourceNa:             neo4jentity.DataSourceNA,
	model.DataSourceOpenline:       neo4jentity.DataSourceOpenline,
	model.DataSourceHubspot:        neo4jentity.DataSourceHubspot,
	model.DataSourceZendeskSupport: neo4jentity.DataSourceZendeskSupport,
	model.DataSourcePipedrive:      neo4jentity.DataSourcePipedrive,
	model.DataSourceSLACk:          neo4jentity.DataSourceSlack,
	model.DataSourceWebscrape:      neo4jentity.DataSourceWebscrape,
	model.DataSourceIntercom:       neo4jentity.DataSourceIntercom,
	model.DataSourceSalesforce:     neo4jentity.DataSourceSalesforce,
	model.DataSourceStripe:         neo4jentity.DataSourceStripe,
	model.DataSourceMixpanel:       neo4jentity.DataSourceMixpanel,
	model.DataSourceClose:          neo4jentity.DataSourceClose,
	model.DataSourceOutlook:        neo4jentity.DataSourceOutlook,
	model.DataSourceUnthread:       neo4jentity.DataSourceUnthread,
}

var sourceByValue = utils.ReverseMap(sourceByModel)

func MapDataSourceFromModel(input model.DataSource) neo4jentity.DataSource {
	return sourceByModel[input]
}

func MapDataSourceToModel(input neo4jentity.DataSource) model.DataSource {
	return sourceByValue[input]
}
